package rule

import (
	"context"
	"encoding/json"
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

// MakeHandler sets up an http.Handler with all public API endpoints mounted.
func MakeHandler(svc Service, opts ...kithttp.ServerOption) http.Handler {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Methods("GET").Path(`/`).Name("ruleList").Handler(
		kithttp.NewServer(
			listEndpoint(svc),
			decodeListRequest,
			kithttp.EncodeJSONResponse,
			opts...,
		),
	)

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

	r.Methods("PUT").Path(`/{id:[a-zA-Z0-9]+}/deactivate`).Name("ruleDeactivate").Handler(
		kithttp.NewServer(
			deactivateEndpoint(svc),
			decodeDeactivateRequest,
			kithttp.EncodeJSONResponse,
			append(
				opts,
				kithttp.ServerBefore(extractMuxVars(varID)),
			)...,
		),
	)

	r.Methods("PUT").Path(`/{id:[a-zA-Z0-9]+}/rollout`).Name("ruleUpdateRollout").Handler(
		kithttp.NewServer(
			updateRolloutEndpoint(svc),
			decodeUpdateRolloutRequest,
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

func decodeDeactivateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, ok := ctx.Value(varID).(string)
	if !ok {
		return nil, errors.Wrap(errors.ErrVarMissing, "id")
	}

	return deactivateRequest{id: id}, nil
}

func decodeGetRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, ok := ctx.Value(varID).(string)
	if !ok {
		return nil, errors.Wrap(errors.ErrVarMissing, "id")
	}

	return getRequest{id: id}, nil
}

func decodeListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return struct{}{}, nil
}

func decodeUpdateRolloutRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, ok := ctx.Value(varID).(string)
	if !ok {
		return nil, errors.Wrap(errors.ErrVarMissing, "id")
	}

	v := struct {
		Rollout uint8 `json:"rollout"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, errors.Wrapf(errors.ErrInvalidPayload, "%s", err)
	}

	return updateRolloutRequest{id: id, rollout: v.Rollout}, nil
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
