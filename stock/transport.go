package stock

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"encoding/json"
	"context"
	"net/http"
	"github.com/MICSTI/imsazon/product"
	"github.com/gorilla/mux"
)

// MakeHandler returns a handler for the stock service
func MakeHandler(sts Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getItemsHandler := kithttp.NewServer(
		makeGetItemsEndpoint(sts),
		decodeGetItemRequest,
		encodeResponse,
		opts...,
	)

	addHandler := kithttp.NewServer(
		makeAddEndpoint(sts),
		decodeAddRequest,
		encodeResponse,
		opts...,
	)

	withdrawHandler := kithttp.NewServer(
		makeWithdrawEndpoint(sts),
		decodeWithdrawRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/stock/items", getItemsHandler).Methods("GET")
	r.Handle("/stock/add", addHandler).Methods("POST")
	r.Handle("/stock/withdraw", withdrawHandler).Methods("POST")

	return r
}

func decodeGetItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
	// there are no parameters to the request, so we don't need to decode anything
	return getItemsRequest{}, nil
}

func decodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		ProductToAdd		product.Product		`json:"product"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return addRequest{
		Product:		body.ProductToAdd,
	}, nil
}

func decodeWithdrawRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		ProductToWithdraw	product.Product		`json:"product"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return withdrawRequest{
		Product:		body.ProductToWithdraw,
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