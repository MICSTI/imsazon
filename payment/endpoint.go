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
	Status				CreditCardChargeStatus		`json:"status"`
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
			Status:			status,
			Err:			err,
		}, nil
	}
}