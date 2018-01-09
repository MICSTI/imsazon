/**
	The cart service is responsible for maintaining the shopping cart of a user.
	Each user only has one cart.
 */
package cart

import (
	"errors"
	"github.com/MICSTI/imsazon/models/user"
	"github.com/MICSTI/imsazon/models/product"
	"github.com/MICSTI/imsazon/models/cart"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("Invalid argument")

// Service is the interface that provides the cart methods
type Service interface {
	// GetCart returns the cart for a user
	GetCart(userId user.UserId) ([]*product.SimpleProduct, error)

	// Put adds an item to a user's cart - if it already exists it will be updated
	Put(userId user.UserId, productId product.ProductId, quantity int) ([]*product.SimpleProduct, error)

	// Remove deletes an item from the user's cart
	Remove(userId user.UserId, productId product.ProductId) ([]*product.SimpleProduct, error)
}

type service struct {
	carts			cart.Repository
}

func (s *service) GetCart(userId user.UserId) ([]*product.SimpleProduct, error) {
	if userId == "" {
		return []*product.SimpleProduct{}, ErrInvalidArgument
	}


}

func (s *service) Put(userId user.UserId, productId product.ProductId, quantity int) (updatedCart []*product.SimpleProduct, err error) {
	if (userId == "" || productId == "" || quantity < 0) {
		return []*product.SimpleProduct{}, ErrInvalidArgument
	}

	return s.carts.Put(userId, productId, quantity)
}

func (s *service) Remove(userId user.UserId, productId product.ProductId) (updatedCart []*product.SimpleProduct, err error) {
	if (userId == "" || productId == "") {
		return []*product.SimpleProduct{}, ErrInvalidArgument
	}

	return s.carts.Remove(userId, productId)
}

// NewService creates a cart service with the necessary dependencies
func NewService(carts cart.Repository) Service {
	return &service{
		carts:		carts,
	}
}