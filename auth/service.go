/*
	The auth service issues a JWT auth token for authentication inside the microservice architecture
 */

package auth

import (
	"errors"
	"github.com/MICSTI/imsazon/user"
)

// ErrInvalidArgument is returned when one or more arguments are invalid
var ErrInvalidArgument = errors.New("Invalid argument")

// Service is the interface that provides the methods for obtaining an auth token
type Service interface {
	// Login checks the passed credentials and issues a JWT auth token in case they are valid
	Login(username string, password string) (string, error)

	// Check checks if the passed JWT auth token is valid
	Check(token string) (bool, error)
}

type service struct {
	users		user.Repository
}

func (s *service) Login(username string, password string) (string, error) {
	return "Hallo", nil
}

func (s *service) Check(token string) (bool, error) {
	return true, nil
}

// NewService returns a new instance of the auth service
func NewService(users user.Repository) Service {
	return &service{
		users:		users,
	}
}