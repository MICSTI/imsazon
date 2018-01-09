/**
	The order service is responsible for storing all orders and the user they belong to.
 */
package order

import (
	"errors"
	orderModel "github.com/MICSTI/imsazon/models/order"
	"github.com/MICSTI/imsazon/models/user"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("Invalid argument")

// Service is the interface that provides order methods
type Service interface {
	// creates a new order
	Create(newOrder *orderModel.Order) (order *orderModel.Order, err error)

	// updates the status of an order
	UpdateStatus(id orderModel.OrderId, newStatus orderModel.OrderStatus) (order *orderModel.Order, err error)

	// returns an order by id
	GetById(id orderModel.OrderId) (*orderModel.Order, error)

	// returns all orders
	GetAll() []*orderModel.Order

	// returns all order for a specific user
	GetAllForUser(userId user.UserId) []*orderModel.Order
}

type service struct {
	orders			orderModel.Repository
}

func (s *service) Create(newOrder *orderModel.Order) (order *orderModel.Order, err error) {
	if newOrder.UserId == "" {
		return nil, ErrInvalidArgument
	}

	newOrder.Id = orderModel.GetRandomOrderId()

	return s.Create(newOrder)
}

func (s *service) UpdateStatus(id orderModel.OrderId, newStatus orderModel.OrderStatus) (order *orderModel.Order, err error) {
	if id == ""  {
		return nil, ErrInvalidArgument
	}

	return s.UpdateStatus(id, newStatus)
}

func (s *service) GetById(id orderModel.OrderId) (*orderModel.Order, error) {
	if id == "" {
		return nil, ErrInvalidArgument
	}

	return s.GetById(id)
}

func (s *service) GetAll() []*orderModel.Order {
	return s.GetAll()
}

func (s *service) GetAllForUser(userId user.UserId) []*orderModel.Order {
	if userId == "" {
		return nil
	}

	return s.GetAllForUser(userId)
}

// NewService returns an order service with necessary dependencies.
func NewService(orders orderModel.Repository) Service {
	return &service{
		orders:		orders,
	}
}