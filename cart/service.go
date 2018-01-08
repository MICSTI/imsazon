/**
	The cart service is responsible for maintaining the shopping cart of a user.
	Each user only has one cart.
 */
package cart

import (
	"errors"
	"github.com/MICSTI/imsazon/models/user"
	"github.com/MICSTI/imsazon/models/product"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("Invalid argument")

// Service is the interface that provides the cart methods
type Service interface {
	// Put adds an item to a user's cart - if it already exists it will be updated
	Put(userId user.UserId, productId product.ProductId, quantity int) error

	// Remove deletes an item from the user's cart
	Remove(id user.UserId, productId product.ProductId) error
}

type service struct {
	users			user.Repository
	products		product.Repository
}
