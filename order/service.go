/**
	The order service is responsible for storing all orders and the user they belong to.
 */
package order

import "errors"

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("Invalid argument")

// Service is the interface that provides order methods
type Service interface {

}