package payment

import (
	"github.com/go-kit/kit/endpoint"
	"context"
)

type chargeRequest struct {
	Id					string
	CardNumber			string
	Amount				float32
	Currency			string
}

type chargeResponse struct {
	Id					string						`json:"transactionId"`
	CardNumber			string						`json:"creditCard"`
	Amount				float32						`json:"amount"`
	Currency			string						`json:"currency"`
	Status				string						`json:"status"`
	Err					error						`json:"error,omitempty"`
}

func (r chargeResponse) error() error { return r.Err }

func makeChargeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(chargeRequest)

		creditCardCharge := CreditCardCharge{
			Id:				req.Id,
			CardNumber:		req.CardNumber,
			Amount:			req.Amount,
			Currency:		req.Currency,
		}

		status, err := s.Charge(creditCardCharge)

		return chargeResponse{
			Id:				req.Id,
			CardNumber:		req.CardNumber,
			Amount:			req.Amount,
			Currency:		req.Currency,
			Status:			status.String(),
			Err:			err,
		}, nil
	}
}