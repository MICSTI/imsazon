/*
	The hello services has no other purpose than to greet the user.
	It is just a test service to get a feel for how a microservice should be implemented.
 */

package hello

import "errors"

// ErrInvalidArgument is returned when on or more arguments are invalid
var ErrInvalidArgument = errors.New("Invalid argument")

// Service is the interface that provides the hello method
type Service interface {
	// just sends out a standard hello message to everyone
	SayHello() (string)

	// sends out a specific hello message to somebody
	SayHelloTo(name string) (string, error)
}

type service struct {

}

func (s *service) SayHello() string {
	return "Hello there!"
}

func (s *service) SayHelloTo(name string) (string, error) {
	if name == "" {
		return "", ErrInvalidArgument
	}

	return "The warmest of welcomes to you, " + name, nil
}