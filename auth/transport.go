package auth

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
	"encoding/json"
	"context"
	"errors"
	"github.com/gorilla/mux"
)

// MakeHandler returns the handler for the auth service
func MakeHandler(as Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	loginHandler := kithttp.NewServer(
		makeLoginEndpoint(as),
		decodeLoginRequest,
		encodeResponse,
		opts...,
	)

	checkHandler := kithttp.NewServer(
		makeCheckEndpoint(as),
		decodeCheckRequest,
		encodeResponse,
		opts...
	)

	r := mux.NewRouter()

	r.Handle("/auth/login", loginHandler).Methods("POST")
	r.Handle("/auth/check", checkHandler).Methods("POST")

	return r
}

var errBadRoute = errors.New("Bad route")

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Username	string	`json:"username"`
		Password	string	`json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return loginRequest{
		Username: 	body.Username,
		Password:	body.Password,
	}, nil
}

func decodeCheckRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Token		string	`json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return checkRequest{
		Token:		body.Token,
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
