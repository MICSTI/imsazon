package hello

import (
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"net/http"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
)

// MakeHandler returns a handler for the hello service
func MakeHandler(hs Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	sayHelloHandler := kithttp.NewServer(
		makeSayHelloEndpoint(hs),
		decodeSayHelloRequest,
		encodeResponse,
		opts...,
	)

	sayHelloToHandler := kithttp.NewServer(
		makeSayHelloToEndpooint(hs),
		decodeSayHelloToRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/hello", sayHelloHandler).Methods("GET")
	r.Handle("/hello/{id}", sayHelloToHandler).Methods("GET")

	return r
}

var errBadRoute = errors.New("Invalid route")

func decodeSayHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	// there is nothing to decode here since there are no parameters
	return sayHelloRequest{}, nil
}

func decodeSayHelloToRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, errBadRoute
	}
	return sayHelloToRequest{Name: name}, nil
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