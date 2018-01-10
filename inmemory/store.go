/*
	This package implements the in-memory repository stores for all important data
 */

package inmemory

import (
	"sync"
	userModel "github.com/MICSTI/imsazon/models/user"
	productModel "github.com/MICSTI/imsazon/models/product"
	cartModel "github.com/MICSTI/imsazon/models/cart"
	orderModel "github.com/MICSTI/imsazon/models/order"
)

/* ---------- USER REPOSITORY ---------- */
type userRepository struct {
	mtx		sync.RWMutex
	users	map[userModel.UserId]*userModel.User
}

// adds a user to the repository store
func (r *userRepository) Add(u *userModel.User) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.users[u.Id] = u
	return nil
}

// attempts to find the user by UserId inside the repository store
func (r *userRepository) Find(id userModel.UserId) (*userModel.User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if val, ok := r.users[id]; ok {
		return val, nil
	}
	return nil, userModel.ErrUnknown
}

// returns alls users in an array
func (r *userRepository) FindAll() []*userModel.User {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	u := make([]*userModel.User, 0, len(r.users))
	for _, val := range r.users {
		u = append(u, val)
	}
	return u
}

// checks the login credentials of a user
func (r *userRepository) CheckLogin(username string, password string) (*userModel.User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	for _, val := range r.users {
		// check if the username and password match with a user in the store
		// -------------------- IMPORTANT NOTICE: -----------------
		// obviously it is unacceptable to compare plaintext passwords in a real-world application
		// however, since this is a prototype application we can overlook this flaw
		// admittedly though, it would be quite simpleto add a hashing method so only the hashes would be stored and compared
		if username == val.Username && password == val.Password {
			return val, nil
		}
	}
	return nil, userModel.ErrUnknown
}

// returns an instance of a user repository
func NewUserRepository() userModel.Repository {
	r := &userRepository{
		users: make(map[userModel.UserId]*userModel.User),
	}

	r.users[userModel.U0001] = userModel.Rey
	r.users[userModel.U0002] = userModel.Kylo
	r.users[userModel.U0003] = userModel.Luke

	return r
}

/* ---------- PRODUCT REPOSITORY ---------- */
type productRepository struct {
	mtx		sync.RWMutex
	products	map[productModel.ProductId]*productModel.Product
}

func (r *productRepository) Store(p *productModel.Product) (*productModel.Product, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.products[p.Id] = p
	return p, nil
}

func (r *productRepository) Add(p *productModel.Product) (*productModel.Product, error) {
	// first check if the product even exists
	stored, err := r.Find(p.Id)

	if err != nil {
		return r.Store(p)
	}

	// update the properties of the stock item
	r.mtx.Lock()
	defer r.mtx.Unlock()

	stored.Quantity += p.Quantity

	return stored, nil
}

func (r *productRepository) Withdraw(p *productModel.Product) (*productModel.Product, error) {
	// first check if the product even exists
	stored, err := r.Find(p.Id)

	if err != nil {
		return nil, err
	}

	// check if there are enough items for withdrawing
	if stored.Quantity < p.Quantity {
		return nil, productModel.ErrNotEnoughItems
	}

	// update the properties of the stock item
	r.mtx.Lock()
	defer r.mtx.Unlock()

	stored.Quantity -= p.Quantity

	return stored, nil
}

func (r *productRepository) Find(id productModel.ProductId) (*productModel.Product, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if val, ok := r.products[id]; ok {
		return val, nil
	}
	return nil, productModel.ErrProductUnknown
}

func (r *productRepository) FindAll() []*productModel.Product {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	p := make([]*productModel.Product, 0, len(r.products))
	for _, val := range r.products {
		p = append(p, val)
	}
	return p
}

func NewProductRepository() productModel.Repository {
	r := &productRepository{
		products: make(map[productModel.ProductId]*productModel.Product),
	}

	r.products[productModel.P0001] = productModel.Lightsaber
	r.products[productModel.P0002] = productModel.MilleniumFalcon
	r.products[productModel.P0003] = productModel.BB8
	r.products[productModel.P0004] = productModel.Podracer
	r.products[productModel.P0005] = productModel.CarboniteFreezer

	return r
}

/* ---------- CART REPOSITORY ---------- */
type cartRepository struct {
	mtx			sync.RWMutex
	carts		map[userModel.UserId][]*productModel.SimpleProduct
}

func (r *cartRepository) StoreUserCart(id userModel.UserId, cartItems []*productModel.SimpleProduct) []*productModel.SimpleProduct {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.carts[id] = cartItems
	return cartItems
}

// gets a user's cart by user id
func (r *cartRepository) FindUserCart(id userModel.UserId) ([]*productModel.SimpleProduct) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if val, ok := r.carts[id]; ok {
		return val
	} else {
		return nil
	}
}

func (r *cartRepository) FindItemInCart(userCart []*productModel.SimpleProduct, productToFind *productModel.SimpleProduct) (int, *productModel.SimpleProduct) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	for idx, val := range userCart {
		if val.Id == productToFind.Id {
			return idx, val
		}
	}
	return -1, nil
}

func (r *cartRepository) GetCart(id userModel.UserId) ([]*productModel.SimpleProduct, error) {
	userCart := r.FindUserCart(id)

	if userCart != nil {
		return userCart, nil
	} else {
		return r.StoreUserCart(id, []*productModel.SimpleProduct{}), nil
	}
}

func (r *cartRepository) Put(userId userModel.UserId, productId productModel.ProductId, quantity int) ([]*productModel.SimpleProduct, error) {
	// first create a SimpleProduct out of the parameters
	sp := productModel.NewSimpleProduct(
		productId,
		quantity,
	)

	// get the user's cart
	userCart := r.FindUserCart(userId)

	// try to find the item in the user's cart
	idx, itemInCart := r.FindItemInCart(userCart, sp)

	if idx < 0 {
		// item does not exist yet - we have to append it to the array
		updatedCart := append(userCart, sp)
		return r.StoreUserCart(userId, updatedCart), nil
	} else {
		// item does already exist - we have to update the properties
		r.mtx.Lock()
		defer r.mtx.Unlock()
		itemInCart.Quantity = sp.Quantity
		return userCart, nil
	}
}

func (r *cartRepository) Remove(userId userModel.UserId, productId productModel.ProductId) ([]*productModel.SimpleProduct, error) {
	// first create a SimpleProduct out of the parameters
	sp := productModel.NewSimpleProduct(
		productId,
		0,		// the quantity does not matter in this case
	)

	// get the user's cart
	userCart := r.FindUserCart(userId)

	// try to find the item in the user's cart
	idx, _ := r.FindItemInCart(userCart, sp)

	if idx >= 0 {
		// deletes the item on the position idx from the slice
		r.mtx.Lock()
		defer r.mtx.Unlock()
		userCart = append(userCart[:idx], userCart[idx + 1:]...)
	}

	// in case the item was not found we just don't do anything
	return userCart, nil
}

func NewCartRepository() cartModel.Repository {
	return &cartRepository{
		carts: make(map[userModel.UserId][]*productModel.SimpleProduct),
	}
}

/* ---------- ORDER REPOSITORY ---------- */
type orderRepository struct {
	mtx			sync.RWMutex
	orders		map[orderModel.OrderId]*orderModel.Order
}

func (r *orderRepository) Create(o *orderModel.Order) (order *orderModel.Order, err error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.orders[o.Id] = o
	return o, nil
}

func (r *orderRepository) UpdateStatus(id orderModel.OrderId, newStatus orderModel.OrderStatus) (order *orderModel.Order, err error) {
	o, err := r.Find(id)

	if err != nil {
		return nil, err
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()
	o.Status = newStatus
	return o, nil
}

func (r *orderRepository) Find(id orderModel.OrderId) (*orderModel.Order, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if val, ok := r.orders[id]; ok {
		return val, nil
	}
	return nil, orderModel.ErrUnknown
}

func (r *orderRepository) FindAll() []*orderModel.Order {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	o := make([]*orderModel.Order, 0, len(r.orders))
	for _, val := range r.orders {
		o = append(o, val)
	}
	return o
}

func (r *orderRepository) FindAllForUser(userId userModel.UserId) []*orderModel.Order {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	o := []*orderModel.Order{}
	for _, val := range r.orders {
		if userId == val.UserId {
			o = append(o, val)
		}
	}
	return o
}

func NewOrderRepository() orderModel.Repository {
	r := &orderRepository{
		orders: make(map[orderModel.OrderId]*orderModel.Order),
	}

	r.orders[orderModel.O0001] = orderModel.Order1
	r.orders[orderModel.O0002] = orderModel.Order2

	return r
}