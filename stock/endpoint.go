package stock

import (
	productModel "github.com/MICSTI/imsazon/models/product"
	"github.com/go-kit/kit/endpoint"
	"context"
)

type getItemsRequest struct {

}

type getItemsResponse struct {
	Products	[]*productModel.Product	`json:"products,omitempty"`
	Err			error				`json:"error,omitempty"`
}

func (r getItemsResponse) error() error { return r.Err }

func makeGetItemsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		products := s.GetItems()
		return getItemsResponse{Products: products, Err: nil}, nil
	}
}

type addRequest struct {
	Product		productModel.Product
}

type addResponse struct {
	UpdatedProduct		*productModel.Product		`json:"product,omitempty"`
	Err					error					`json:"error,omitempty"`
}

func (r addResponse) error() error { return r.Err }

func makeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addRequest)
		updatedProduct, err := s.Add(&req.Product)

		return addResponse{UpdatedProduct: updatedProduct, Err: err}, nil
	}
}

type withdrawRequest struct {
	Product		productModel.Product
}

type withdrawResponse struct {
	UpdatedProduct		*productModel.Product	`json:"product,omitempty"`
	Err					error				`json:"error,omitempty"`
}

func (r withdrawResponse) error() error { return r.Err }

func makeWithdrawEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(withdrawRequest)
		updatedProduct, err := s.Withdraw(&req.Product)
		return withdrawResponse{UpdatedProduct: updatedProduct, Err: err}, nil
	}
}