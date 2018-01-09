package order

import "github.com/MICSTI/imsazon/models/product"

// OrderId uniquely identifies an order
type OrderId string

func (s OrderId) String() string {
	return string(s)
}

// OrderStatus describes the status of an order
type OrderStatus int

// valid order statuses
const (
	Created	OrderStatus = iota
	PaymentSuccessful
	PaymentError
	Shipped
	ReturnRequested
	Returned
)

func (s OrderStatus) String() string {
	switch s {
	case Created:
		return "Created"
	case PaymentSuccessful:
		return "Payment Successful"
	case PaymentError:
		return "Payment Error"
	case Shipped:
		return "Shipped"
	case ReturnRequested:
		return "Return Requested"
	case Returned:
		return "Returned"
	}
	return "Unknown order status"
}

type Order struct {
	Id			OrderId
	Items		[]*product.SimpleProduct
}