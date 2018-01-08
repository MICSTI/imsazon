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
	// Add adds an item to a user's cart
	Add(userId user.UserId, productId product.ProductId, quantity int) error

	// Update updates the quantity of an item in a user's cart
	Update(id user.UserId, productId product.ProductId, quantity int) error

	// Delete deletes an item from the user's cart
	Delete(id user.UserId, productId product.ProductId)
}

type service struct {
	users			user.Repository
	products		product.Repository
}
