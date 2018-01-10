package shipping

import (
	orderModel "github.com/MICSTI/imsazon/models/order"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type shipRequest struct {
	OrderId			orderModel.OrderId
}

type shipResponse struct {
	Err				error		`json:"error,omitempty"`
}

func (r shipResponse) error() error { return r.Err }

func makeShipEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(shipRequest)
		err := s.Ship(req.OrderId)
		return shipResponse{Err: err}, nil
	}
}