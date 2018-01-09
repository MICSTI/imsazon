package cart

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"encoding/json"
	"context"
	"net/http"
	"github.com/MICSTI/imsazon/models/user"
	"github.com/MICSTI/imsazon/models/product"
	"github.com/gorilla/mux"
)

// MakeHandler returns a handler for the cart service
func MakeHandler(cs Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getCartHandler := kithttp.NewServer(
		makeGetCartEndpoint(cs),
		decodeGetCartRequest,
		encodeResponse,
		opts...,
	)

	putItemHandler := kithttp.NewServer(
		makePutItemEndpoint(cs),
		decodePutItemRequest,
		encodeResponse,
		opts...,
	)

	removeItemHandler := kithttp.NewServer(
		makeRemoveItemEndpoint(cs),
		decodeRemoveItemRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/cart/get", getCartHandler).Methods("POST")
	r.Handle("/cart/put", putItemHandler).Methods("POST")
	r.Handle("/cart/remove", removeItemHandler).Methods("POST")

	return r
}

func decodeGetCartRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		UserId			user.UserId			`json:"userId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return getCartRequest{
		UserId:			body.UserId,
	}, nil
}

func decodePutItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		UserId			user.UserId			`json:"userId"`
		ProductId		product.ProductId	`json:"productId"`
		Quantity		int					`json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return putItemRequest{
		UserId:			body.UserId,
		ProductId:		body.ProductId,
		Quantity:		body.Quantity,
	}, nil
}

func decodeRemoveItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		UserId			user.UserId			`json:"userId"`
		ProductId		product.ProductId	`json:"productId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return removeItemRequest{
		UserId:			body.UserId,
		ProductId:		body.ProductId,
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