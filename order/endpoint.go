package order

import (
	orderModel "github.com/MICSTI/imsazon/models/order"
	"github.com/go-kit/kit/endpoint"
	"context"
	"github.com/MICSTI/imsazon/models/user"
)

type createRequest struct {
	Order			*orderModel.Order
}

type createResponse struct {
	Order			*orderModel.Order		`json:"order,omitempty"`
	Err				error					`json:"error,omitempty"`
}

func (r createResponse) error() error { return r.Err }

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createRequest)
		createdOrder, err := s.Create(req.Order)
		return createResponse{Order: createdOrder, Err: err}, nil
	}
}

type updateStatusRequest struct {
	Id				orderModel.OrderId
	NewStatus		orderModel.OrderStatus
}

type updateStatusResponse struct {
	Order			*orderModel.Order		`json:"order,omitempty"`
	Err				error					`json:"error,omitempty"`
}

func (r updateStatusResponse) error() error { return r.Err }

func makeUpdateStatusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateStatusRequest)
		updatedOrder, err := s.UpdateStatus(req.Id, req.NewStatus)
		return updateStatusResponse{Order: updatedOrder, Err: err}, nil
	}
}

type getByIdRequest struct {
	Id				orderModel.OrderId
}

type getByIdResponse struct {
	Order			*orderModel.Order		`json:"order,omitempty"`
	Err				error					`json:"error,omitempty"`
}

func (r getByIdResponse) error() error { return r.Err }

func makeGetByIdEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getByIdRequest)
		o, err := s.GetById(req.Id)
		return getByIdResponse{Order: o, Err: err}, nil
	}
}

type getAllRequest struct {

}

type getAllResponse struct {
	Orders			[]*orderModel.Order		`json:"orders"`
}

func makeGetAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		orders := s.GetAll()
		return getAllResponse{Orders: orders}, nil
	}
}

type getAllForUserRequest struct {
	userId			user.UserId
}

type getAllForUserResponse struct {
	Orders			[]*orderModel.Order		`json:"orders"`
}

func makeGetAllForUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAllForUserRequest)
		orders := s.GetAllForUser(req.userId)
		return getAllForUserResponse{Orders: orders}, nil
	}
}