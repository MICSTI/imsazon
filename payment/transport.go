package payment

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
	"encoding/json"
	"context"
	"errors"
	"github.com/gorilla/mux"
)

// MakeHandler returns a handler for the payment service
func MakeHandler(ps Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	chargeHandler := kithttp.NewServer(
		makeChargeEndpoint(ps),
		decodeChargeRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/payment/charge", chargeHandler).Methods("POST")

	return r
}

// encode the JSON response
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(erroer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type erroer interface {
	error() error
}

var errBadRoute = errors.New("Bad route")

func decodeChargeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Id				string		`json:"transactionId"`
		CardNumber		string		`json:"creditCard"`
		Amount			float32		`json:"amount"`
		Currency		string		`json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return chargeRequest{
		Id:					body.Id,
		CardNumber:			body.CardNumber,
		Amount:				body.Amount,
		Currency:			body.Currency,
	}, nil
}

// encode errors from business logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
		break
	case ErrCard:
		w.WriteHeader(http.StatusBadRequest)
		break
	case ErrValidation:
		w.WriteHeader(http.StatusBadRequest)
		break
	case ErrNetwork:
		w.WriteHeader(http.StatusInternalServerError)
		break
	case ErrOther:
		w.WriteHeader(http.StatusBadRequest)
		break
	default:
		w.WriteHeader(http.StatusInternalServerError)
		break
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}