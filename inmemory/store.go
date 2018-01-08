/*
	This package implements the in-memory repository stores for all important data
 */

package inmemory

import (
	"sync"
	"github.com/MICSTI/imsazon/models/user"
	"github.com/MICSTI/imsazon/models/product"
	"github.com/MICSTI/imsazon/models/cart"
)

/* ---------- USER REPOSITORY ---------- */
type userRepository struct {
	mtx		sync.RWMutex
	users	map[user.UserId]*user.User
}

// adds a user to the repository store
func (r *userRepository) Add(u *user.User) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.users[u.Id] = u
	return nil
}

// attempts to find the user by UserId inside the repository store
func (r *userRepository) Find(id user.UserId) (*user.User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if val, ok := r.users[id]; ok {
		return val, nil
	}
	return nil, user.ErrUnknown
}

// returns alls users in an array
func (r *userRepository) FindAll() []*user.User {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	u := make([]*user.User, 0, len(r.users))
	for _, val := range r.users {
		u = append(u, val)
	}
	return u
}

// checks the login credentials of a user
func (r *userRepository) CheckLogin(username string, password string) (*user.User, error) {
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
	return nil, user.ErrUnknown
}

// returns an instance of a user repository
func NewUserRepository() user.Repository {
	r := &userRepository{
		users: make(map[user.UserId]*user.User),
	}

	r.users[user.U0001] = user.Rey
	r.users[user.U0002] = user.Kylo
	r.users[user.U0003] = user.Luke

	return r
}

/* ---------- PRODUCT REPOSITORY ---------- */
type productRepository struct {
	mtx		sync.RWMutex
	products	map[product.ProductId]*product.Product
}

func (r *productRepository) Store(p *product.Product) (*product.Product, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.products[p.Id] = p
	return p, nil
}

func (r *productRepository) Add(p *product.Product) (*product.Product, error) {
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

func (r *productRepository) Withdraw(p *product.Product) (*product.Product, error) {
	// first check if the product even exists
	stored, err := r.Find(p.Id)

	if err != nil {
		return nil, err
	}

	// check if there are enough items for withdrawing
	if stored.Quantity < p.Quantity {
		return nil, product.ErrNotEnoughItems
	}

	// update the properties of the stock item
	r.mtx.Lock()
	defer r.mtx.Unlock()

	stored.Quantity -= p.Quantity

	return stored, nil
}

func (r *productRepository) Find(id product.ProductId) (*product.Product, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if val, ok := r.products[id]; ok {
		return val, nil
	}
	return nil, product.ErrProductUnknown
}

func (r *productRepository) FindAll() []*product.Product {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	p := make([]*product.Product, 0, len(r.products))
	for _, val := range r.products {
		p = append(p, val)
	}
	return p
}

func NewProductRepository() product.Repository {
	r := &productRepository{
		products: make(map[product.ProductId]*product.Product),
	}

	r.products[product.P0001] = product.Lightsaber
	r.products[product.P0002] = product.MilleniumFalcon
	r.products[product.P0003] = product.BB8
	r.products[product.P0004] = product.Podracer
	r.products[product.P0005] = product.CarboniteFreezer

	return r
}

/* ---------- CART REPOSITORY ---------- */
type cartRepository struct {
	mtx			sync.RWMutex
	carts		map[user.UserId][]*product.SimpleProduct
}

func (r *cartRepository) Put(userId user.UserId, productId product.ProductId, quantity int) ([]*product.SimpleProduct, error) {
	return nil, nil
}

func (r *cartRepository) Remove(userId user.UserId, productId product.ProductId) ([]*product.SimpleProduct, error) {
	return nil, nil
}

func NewCartRepository() cart.Repository {
	return &cartRepository{
		carts: make(map[user.UserId][]*product.SimpleProduct)
	}
}