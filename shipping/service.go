package shipping

import (
	"errors"
	orderModel "github.com/MICSTI/imsazon/models/order"
	"time"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("Invalid argument")

// Service is the interface that provides the shipping methods
type Service interface {
	// Ships the order from the physical store
	Ship(orderModel.OrderId) error
}

type service struct {

}

func (s *service) Ship(orderId orderModel.OrderId) (err error) {
	if orderId == "" {
		return ErrInvalidArgument
	}

	// we can't really do anything, so we just add a delay and trigger the sending of an email
	duration := time.Millisecond * 750
	time.Sleep(duration)

	// TODO call mail service


	return nil
}

// NewService creates a shipping service
func NewService() Service {
	return &service{

	}
}