package order

import (
	"github.com/MICSTI/imsazon/models/product"
	"errors"
	"math/rand"
	"time"
)

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
	Status		OrderStatus
	Items		[]*product.SimpleProduct
}

func New(id OrderId, items []*product.SimpleProduct) *Order {
	return &Order{
		Id:			id,
		Items:		items,
	}
}

// Repository provides access to an order store
type Repository interface {
	Create(order *Order) error
	UpdateStatus(id OrderId, newStatus OrderStatus) error
	Find(id OrderId) (*Order, error)
	FindAll() []*Order
}

// ErrUnknown is used when an ordercould not be found.
var ErrUnknown = errors.New("Unknown order")

// create random string for OrderIds
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func getRandomOrderId() OrderId {
	return OrderId(RandStringBytesMaskImprSrc(8))
}