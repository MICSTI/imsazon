package stock

import (
	"github.com/MICSTI/imsazon/product"
	"github.com/go-kit/kit/endpoint"
	"context"
)

type getItemsRequest struct {

}

type getItemsResponse struct {
	Products	[]product.Product	`json:"products,omitempty"`
	Err			error				`json:"error,omitempty"`
}

func (r getItemsResponse) error() error { return r.Err }

func makeGetItemsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		products, err := s.GetItems()
		return getItemsResponse{Products: products, Err: err}, nil
	}
}

type addRequest struct {
	Product		product.Product
}

type addResponse struct {
	UpdatedProduct		product.Product	`json:"product,omitempty"`
	Err					error				`json:"error,omitempty"`
}

func (r addResponse) error() error { return r.Err }

func makeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addRequest)
		updatedProduct, err := s.Add(&req.Product)
		return addResponse{UpdatedProduct: updatedProduct, Err: err}, nil
	}
}

// TODO withdrawRequest, withdrawResponse, makeWithdrawEndpoint