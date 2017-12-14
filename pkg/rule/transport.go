package rule

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/lifesum/configsum/pkg/errors"
)

// URL fragments.
const (
	varID muxVar = "id"
)

type muxVar string

func MakeHandler(svc Service, opts ...kithttp.ServerOption) http.Handler {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Methods("GET").Path(`/{id:[a-zA-Z0-9]+}`).Name("ruleGet").Handler(
		kithttp.NewServer(
			getEndpoint(svc),
			decodeGetRequest,
			kithttp.EncodeJSONResponse,
			append(
				opts,
				kithttp.ServerBefore(extractMuxVars(varID)),
			)...,
		),
	)

	r.Methods("PUT").Path(`/{id:[a-zA-Z0-9]+}/activate`).Name("ruleActivate").Handler(
		kithttp.NewServer(
			activateEndpoint(svc),
			decodeActivateRequest,
			kithttp.EncodeJSONResponse,
			append(
				opts,
				kithttp.ServerBefore(extractMuxVars(varID)),
			)...,
		),
	)

	return r
}

func decodeActivateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, ok := ctx.Value(varID).(string)
	if !ok {
		return nil, errors.Wrap(errors.ErrVarMissing, "id")
	}

	return activateRequest{id: id}, nil
}

func decodeGetRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, ok := ctx.Value(varID).(string)
	if !ok {
		return nil, errors.Wrap(errors.ErrVarMissing, "id")
	}

	return getRequest{id: id}, nil
}

func extractMuxVars(keys ...muxVar) kithttp.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		for _, k := range keys {
			if v, ok := mux.Vars(r)[string(k)]; ok {
				ctx = context.WithValue(ctx, k, v)
			}
		}

		return ctx
	}
}
