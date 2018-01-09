package cart

import (
	"github.com/MICSTI/imsazon/models/user"
	"github.com/MICSTI/imsazon/models/product"
	"github.com/go-kit/kit/endpoint"
	"context"
)

type putItemRequest struct {
	UserId			user.UserId
	ProductId		product.ProductId
	Quantity		int
}

type putItemResponse struct {
	CartItems		[]*product.SimpleProduct
	Err				error
}

func (r putItemResponse) error() error { return r.Err }

func makePutItemEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(putItemRequest)
		updatedCart, err := s.Put(req.UserId, req.ProductId, req.Quantity)
		return putItemResponse{CartItems: updatedCart, Err: err}, nil
	}
}

type removeItemRequest struct {
	UserId			user.UserId
	ProductId		product.ProductId
}

type removeItemResponse struct {
	CartItems		[]*product.SimpleProduct
	Err				error
}

func (r removeItemResponse) error() error { return r.Err }

func makeRemoveItemEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(removeItemRequest)
		updatedCart, err := s.Remove(req.UserId, req.ProductId)
		return removeItemResponse{CartItems: updatedCart, Err: err}, nil
	}
}