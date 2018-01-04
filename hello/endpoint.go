package hello

import (
	"github.com/go-kit/kit/endpoint"
	"context"
)

type sayHelloRequest struct {

}

type sayHelloResponse struct {
	Greeting string `json:"greeting,omitempty"`
}

func makeSayHelloEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		greeting := s.SayHello()
		return sayHelloResponse{Greeting: greeting}, nil
	}
}

type sayHelloToRequest struct {
	Name string
}

type sayHelloToResponse struct {
	Greeting string `json:"greeting,omitempty"`
	Err error		`json:"error,omitempty"`
}

func (r sayHelloToResponse) error() error { return r.Err }

func makeSayHelloToEndpooint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(sayHelloToRequest)
		greeting, err := s.SayHelloTo(req.Name)
		return sayHelloToResponse{Greeting: greeting, Err: err}, nil
	}
}
