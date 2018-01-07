package auth

import (
	"github.com/go-kit/kit/endpoint"
	"context"
	"github.com/MICSTI/imsazon/models/user"
)

type loginRequest struct {
	Username		string
	Password		string
}

type loginResponse struct {
	Token	string	`json:"token,omitempty"`
	Err		error	`json:"error,omitempty"`
}

func (r loginResponse) error() error { return r.Err }

func makeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loginRequest)
		token, err := s.Login(req.Username, req.Password)
		return loginResponse{Token: token, Err: err}, nil
	}
}

type checkRequest struct {
	Token	string
}

type checkResponse struct {
	UserId	user.UserId	`json:"userId,omitempty"`
	Err		error	`json:"error,omitempty"`
}

func (r checkResponse) error() error { return r.Err }

func makeCheckEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface {}) (interface{}, error) {
		req := request.(checkRequest)
		userId, err := s.Check(req.Token)
		return checkResponse{UserId: userId, Err: err}, nil
	}
}