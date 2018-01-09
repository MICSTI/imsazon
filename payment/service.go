/*
	The payment service is responsible for handling all payments.
	Payments happens when a user started the checkout process, before items are shipped.

				IMPORTANT NOTICE:
	This is only a fake service, which obviously does not issue any real charges to credit cards.
	It merely returns a success code in 80% of the cases and an error in the other 20% to simulate payment errors.
 */
package payment

import (
	"errors"
	"math/rand"
	"time"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("Invalid argument")

var ErrCard = errors.New(CardError.String())
var ErrValidation = errors.New(ValidationError.String())
var ErrNetwork = errors.New(NetworkError.String())
var ErrOther = errors.New(OtherError.String())

type CreditCardCharge struct {
	Id					string
	CardNumber			string
	Amount				float32
	Currency			string
	Status				CreditCardChargeStatus
}

type CreditCardChargeStatus int

const (
	Success		CreditCardChargeStatus = iota
	CardError
	ValidationError
	NetworkError
	OtherError
)

func (s CreditCardChargeStatus) String() string {
	switch s {
	case Success:
		return "Success"
	case CardError:
		return "CardError"
	case ValidationError:
		return "ValidationError"
	case NetworkError:
		return "NetworkError"
	case OtherError:
		return "OtherError"
	}

	return "OtherError"
}

// Service is the interface that provides the payment methods
type Service interface {
	// Charge creates a new credit card charge.
	Charge(charge CreditCardCharge) (CreditCardChargeStatus, error)
}

type service struct {

}

func (s *service) Charge(charge CreditCardCharge) (status CreditCardChargeStatus, err error) {
	// create random source
	randSource := rand.NewSource(time.Now().UnixNano())
	randNumber := rand.New(randSource)

	// in 80% of the cases we return a success status, in the other 20% an error
	success := randNumber.Float64() < 0.8

	// add a fake delay of 1 second to make it more realistic
	duration := time.Second
	time.Sleep(duration)

	if success {
		return Success, nil
	} else {
		return CardError, ErrCard
	}
}

func NewService() Service {
	return &service{

	}
}