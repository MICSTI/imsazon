// This package contains the stock model

package product

import "errors"

// ProductId uniquely identifies a product
type ProductId string

func (p ProductId) String() string {
	return string(p)
}

type Product struct {
	Id				ProductId
	Name			string
	Description		string
	ImageUrl		string
	Quantity		int
}

func New(id ProductId, name string, description string, imageUrl string, quantity int) *Product {
	return &Product{
		Id:				id,
		Name:			name,
		Description:	description,
		ImageUrl:		imageUrl,
		Quantity:		quantity,
	}
}

// Repository interface provides access to an in-memory product store
type Repository interface {
	// add a product to the store
	// if the product id already exists, the name and description properties are updated and the quantity added
	// returns a new product object with the current stock status
	Add(product *Product) (Product, error)

	// withdraws a product from the store
	// returns a new product object with the current stock status
	Withdraw(product *Product) (Product, error)
}

var ErrProductUnknown = errors.New("Unknown product")