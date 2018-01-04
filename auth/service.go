/*
	The auth service issues a JWT auth token for authentication inside the microservice architecture
 */

package auth

import "errors"

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
	
}