package order

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"encoding/json"
	"net/http"
	"context"
	"github.com/gorilla/mux"
	"github.com/MICSTI/imsazon/models/user"
	"github.com/MICSTI/imsazon/models/product"
	orderModel "github.com/MICSTI/imsazon/models/order"
	"errors"
)

var ErrBadRoute = errors.New("Bad route")

// MakeHandler returns a handler for the order service.
func MakeHandler(os Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r := mux.NewRouter()

	r.Handle("/order/create", nil).Methods("POST")
	r.Handle("/order/{orderId}/update", nil).Methods("POST")
	r.Handle("/order/{orderId}", nil).Methods("GET")
	r.Handle("/order/all", nil).Methods("GET")
	r.Handle("/order/user/{userId}", nil).Methods("GET")

	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		UserId			user.UserId					`json:"userId"`
		Items			[]*product.SimpleProduct	`json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return createRequest{
		Order:		orderModel.New("", body.UserId, body.Items),
	}, nil
}

func decodeUpdateStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	id, ok := vars["orderId"]

	if !ok {
		return nil, ErrBadRoute
	}

	var body struct {
		Status			orderModel.OrderStatus		`json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return updateStatusRequest{
		Id:			orderModel.OrderId(id),
		NewStatus:	body.Status,
	}, nil
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

// encode errors from business logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}