package mail

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
	"encoding/json"
	"context"
	"github.com/gorilla/mux"
)

// MakeHandler returns a handler for the mail service
func MakeHandler(ms Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	sendHandler := kithttp.NewServer(
		makeSendEndpoint(ms),
		decodeSendRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/mail/send", sendHandler).Methods("POST")

	return r
}

func decodeSendRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		From		string		`json:"from"`
		To			string		`json:"to"`
		Subject		string		`json:"subject"`
		Body		string		`json:"body"`
		ContentType	string		`json:"contentType"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return sendRequest{
		email:	Email{
			From:			body.From,
			To:				body.To,
			Subject:		body.Subject,
			Body:			body.Body,
			ContentType:	body.ContentType,
			},
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