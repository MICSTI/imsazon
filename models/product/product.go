// This package contains the stock model

package product

import "errors"

// ProductId uniquely identifies a product
type ProductId string

func (p ProductId) String() string {
	return string(p)
}

type Product struct {
	Id				ProductId		`json:"id"`
	Name			string			`json:"name"`
	Description		string			`json:"description"`
	Category		string			`json:"category"`
	ImageUrl		string			`json:"imageUrl"`
	Price			float32			`json:"price"`
	Quantity		int				`json:"quantity"`
}

func New(id ProductId, name string, description string, category string, imageUrl string, price float32, quantity int) *Product {
	return &Product{
		Id:				id,
		Name:			name,
		Description:	description,
		Category:		category,
		ImageUrl:		imageUrl,
		Price:			price,
		Quantity:		quantity,
	}
}

type SimpleProduct struct {
	Id				ProductId		`json:"id"`
	Quantity		int				`json:"quantity"`
}

func NewSimpleProduct(id ProductId, quantity int) *SimpleProduct {
	return &SimpleProduct{
		Id:			id,
		Quantity:	quantity,
	}
}

// Repository interface provides access to an in-memory product store
type Repository interface {
	// directly stores a product in the store
	Store(product *Product) (*Product, error)

	// tries to find a product in the store by ProductId
	Find(id ProductId) (*Product, error)

	// returns an array of all products inside the store
	FindAll() []*Product

	// adds a product to the store
	// if the product id already exists, the name and description properties are updated and the quantity added
	// returns a new product object with the current stock status
	Add(product *Product) (*Product, error)

	// withdraws a product from the store
	// returns a new product object with the current stock status
	Withdraw(product *Product) (*Product, error)
}

var ErrProductUnknown = errors.New("Unknown product")
var ErrNotEnoughItems = errors.New("There are not enough items in the store for this operation")