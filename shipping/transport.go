package shipping

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"encoding/json"
	"context"
	"net/http"
	"errors"
	"github.com/gorilla/mux"
	orderModel "github.com/MICSTI/imsazon/models/order"
)

// MakeHandler returns a handler for the shipping service.
func MakeHandler(shs Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	shipHandler := kithttp.NewServer(
		makeShipEndpoint(shs),
		decodeShipRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/ship/{orderId}", shipHandler).Methods("POST")

	return r
}

var ErrBadRoute = errors.New("Bad route")

func decodeShipRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["orderId"]
	if !ok {
		return nil, ErrBadRoute
	}
	return shipRequest{OrderId: orderModel.OrderId(id)}, nil
}

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

// encode errors from business logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	case ErrBadRoute:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}