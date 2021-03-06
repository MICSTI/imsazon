package cart

import (
	"github.com/MICSTI/imsazon/models/product"
	"github.com/MICSTI/imsazon/models/user"
)

// Repository interface provides access to an in-memory cart store
type Repository interface {
	// returns the shopping cart for a user containing all the items that are currently in it
	GetCart(user.UserId) ([]*product.SimpleProduct, error)

	// adds an item to a user's cart - if it already exists it will be updated
	Put(user.UserId, product.ProductId, int) ([]*product.SimpleProduct, error)

	// deletes an item from the user's cart
	Remove(id user.UserId, productId product.ProductId) ([]*product.SimpleProduct, error)
}