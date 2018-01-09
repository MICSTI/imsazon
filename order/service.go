/**
	The order service is responsible for storing all orders and the user they belong to.
 */
package order

import (
	"errors"
	"github.com/MICSTI/imsazon/models/order"
	"github.com/MICSTI/imsazon/models/user"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("Invalid argument")

// Service is the interface that provides order methods
type Service interface {
	// creates a new order
	Create(newOrder *order.Order) (order *order.Order, err error)

	// updates the status of an order
	UpdateStatus(id order.OrderId, newStatus order.OrderStatus) (order *order.Order, err error)

	// returns an order by id
	GetById(id order.OrderId) (*order.Order, error)

	// returns all orders
	GetAll() []*order.Order

	// returns all order for a specific user
	GetAllForUser(userId user.UserId) []*order.Order
}

type service struct {
	orders			order.Repository
}

func (s *service) Create(newOrder *order.Order) (order *order.Order, err error) {
	if newOrder.UserId == "" {
		return nil, ErrInvalidArgument
	}

	newOrder.Id = order.getRandomOrderId()
}