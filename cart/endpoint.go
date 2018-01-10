package cart

import (
	userModel "github.com/MICSTI/imsazon/models/user"
	productModel "github.com/MICSTI/imsazon/models/product"
	"github.com/go-kit/kit/endpoint"
	"context"
)

type getCartRequest struct {
	UserId			userModel.UserId
}

type getCartResponse struct {
	UserId			userModel.UserId						`json:"userId,omitempty"`
	CartItems		[]*productModel.SimpleProduct		`json:"items"`
	Err				error							`json:"error,omitempty"`
}

func (r getCartResponse) error() error { return r.Err }

func makeGetCartEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getCartRequest)
		cartItems, err := s.GetCart(req.UserId)
		return getCartResponse{UserId: req.UserId, CartItems: cartItems, Err: err}, nil
	}
}

type putItemRequest struct {
	UserId			userModel.UserId
	ProductId		productModel.ProductId
	Quantity		int
}

type putItemResponse struct {
	UserId			userModel.UserId						`json:"userId,omitempty"`
	CartItems		[]*productModel.SimpleProduct		`json:"items"`
	Err				error							`json:"error,omitempty"`
}

func (r putItemResponse) error() error { return r.Err }

func makePutItemEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(putItemRequest)
		updatedCart, err := s.Put(req.UserId, req.ProductId, req.Quantity)
		return putItemResponse{UserId: req.UserId, CartItems: updatedCart, Err: err}, nil
	}
}

type removeItemRequest struct {
	UserId			userModel.UserId
	ProductId		productModel.ProductId
}

type removeItemResponse struct {
	UserId			userModel.UserId						`json:"userId,omitempty"`
	CartItems		[]*productModel.SimpleProduct		`json:"items"`
	Err				error							`json:"error,omitempty"`
}

func (r removeItemResponse) error() error { return r.Err }

func makeRemoveItemEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(removeItemRequest)
		updatedCart, err := s.Remove(req.UserId, req.ProductId)
		return removeItemResponse{UserId: req.UserId, CartItems: updatedCart, Err: err}, nil
	}
}