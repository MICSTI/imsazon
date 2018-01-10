
/*
	The stock service is responsible for keeping track of the inventory of IMSazon.
	It provides information about all stock items and their quantity in the store.
	It also provides methods to add and withdraw items from the store.
 */
package stock

import (
	"errors"
	productModel "github.com/MICSTI/imsazon/models/product"
)

// ErrInvalidArgument is returned when on or more arguments are invalid
var ErrInvalidArgument = errors.New("Invalid argument")

type Service interface {
	// GetItems returns an array of all stock products including their quantity
	GetItems() []*productModel.Product

	// Add adds an item with the specified quantity to the stock. Returns a new product object with the updated stock information.
	Add(productToAdd *productModel.Product) (*productModel.Product, error)

	// Withdraw removes the specified quantity from the stock. Returns a new product object with the updated stock information.
	Withdraw(productToWithdraw *productModel.Product) (*productModel.Product, error)
}

type service struct {
	products		productModel.Repository
}

func(s *service) GetItems() []*productModel.Product {
	p := s.products.FindAll()

	return p
}

func(s *service) Add(productToAdd *productModel.Product) (*productModel.Product, error) {
	if productToAdd.Id == "" {
		return &productModel.Product{}, ErrInvalidArgument
	}

	p, err := s.products.Add(productToAdd)

	if err != nil {
		return &productModel.Product{}, err
	}

	return p, nil
}

func(s *service) Withdraw(productToWithdraw *productModel.Product) (*productModel.Product, error) {
	if productToWithdraw.Id == "" {
		return &productModel.Product{}, ErrInvalidArgument
	}

	p, err := s.products.Withdraw(productToWithdraw)

	if err != nil {
		return &productModel.Product{}, err
	}

	return p, nil
}

func NewService(products productModel.Repository) Service {
	return &service{
		products: products,
		}
}