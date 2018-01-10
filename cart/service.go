/**
	The cart service is responsible for maintaining the shopping cart of a user.
	Each user only has one cart.
 */
package cart

import (
	"errors"
	userModel "github.com/MICSTI/imsazon/models/user"
	productModel "github.com/MICSTI/imsazon/models/product"
	cartModel "github.com/MICSTI/imsazon/models/cart"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("Invalid argument")

// Service is the interface that provides the cart methods
type Service interface {
	// GetCart returns the cart for a user
	GetCart(userId userModel.UserId) ([]*productModel.SimpleProduct, error)

	// Put adds an item to a user's cart - if it already exists it will be updated
	Put(userId userModel.UserId, productId productModel.ProductId, quantity int) ([]*productModel.SimpleProduct, error)

	// Remove deletes an item from the user's cart
	Remove(userId userModel.UserId, productId productModel.ProductId) ([]*productModel.SimpleProduct, error)
}

type service struct {
	carts			cartModel.Repository
}

func (s *service) GetCart(userId userModel.UserId) ([]*productModel.SimpleProduct, error) {
	if userId == "" {
		return []*productModel.SimpleProduct{}, ErrInvalidArgument
	}

	return s.carts.GetCart(userId)
}

func (s *service) Put(userId userModel.UserId, productId productModel.ProductId, quantity int) (updatedCart []*productModel.SimpleProduct, err error) {
	if (userId == "" || productId == "" || quantity < 0) {
		return []*productModel.SimpleProduct{}, ErrInvalidArgument
	}

	return s.carts.Put(userId, productId, quantity)
}

func (s *service) Remove(userId userModel.UserId, productId productModel.ProductId) (updatedCart []*productModel.SimpleProduct, err error) {
	if (userId == "" || productId == "") {
		return []*productModel.SimpleProduct{}, ErrInvalidArgument
	}

	return s.carts.Remove(userId, productId)
}

// NewService creates a cart service with the necessary dependencies
func NewService(carts cartModel.Repository) Service {
	return &service{
		carts:		carts,
	}
}