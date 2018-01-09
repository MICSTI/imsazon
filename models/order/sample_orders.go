package order

import "github.com/MICSTI/imsazon/models/product"

// Sample OrderIds
var (
	O0001 OrderId = "zdsklde4d"
	O0002 OrderId = "etd4cs5d3"
)

// sample orders
var (
	Order1 = &Order{
		Id:	O0001,
		Items: []*product.SimpleProduct{
			&product.SimpleProduct{
				Id: product.P0001,
				Quantity: 2,
			},
			&product.SimpleProduct{
				Id: product.P0003,
				Quantity: 1,
			},
		},
	}
	Order2 = &Order{
		Id: O0002,
		Items: []*product.SimpleProduct{
			&product.SimpleProduct{
				Id: product.P0002,
				Quantity: 1,
			},
		},
	}
)