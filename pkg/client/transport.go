package client

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/lifesum/configsum/pkg/errors"
)

const headerToken = "X-Configsum-Token"

// MakeHandler returns an http.Handler for Service.
func MakeHandler(svc Service, opts ...kithttp.ServerOption) http.Handler {
	r := mux.NewRouter()
	r.StrictSlash(false)

	r.Methods("GET").Path(`/`).Name("clientList").Handler(
		kithttp.NewServer(
			listEndpoint(svc),
			decodeListRequest,
			kithttp.EncodeJSONResponse,
			opts...,
		),
	)

	r.Methods("POST").Path("/").Name("clientCreate").Handler(
		kithttp.NewServer(
			createEndpoint(svc),
			decodeCreateRequest,
			kithttp.EncodeJSONResponse,
			opts...,
		),
	)

	return r
}

// HTTPToContext moves the Client secret token from request header to context.
func HTTPToContext(
	ctx context.Context,
	r *http.Request,
) context.Context {
	secret := r.Header.Get(headerToken)
	if len(secret) != secretLen {
		return ctx
	}

	return context.WithValue(ctx, contextKeySecret, secret)
}
func decodeCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	v := struct {
		Name string `json:"name"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, errors.Wrap(errors.ErrInvalidPayload, err.Error())
	}

	if v.Name == "" {
		return nil, errors.Wrap(errors.ErrInvalidPayload, "missing name")
	}

	return createRequest{
		name: v.Name,
	}, nil
}

func decodeListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}
