package config

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/rule"
	confhttp "github.com/lifesum/configsum/pkg/transport/http"
)

// Headers.
const (
	headerContentType = "Content-Type"
	headerBaseID      = "X-Configsum-Base-Id"
	headerBaseName    = "X-Configsum-Base-Name"
	headerClientID    = "X-Configsum-Client-Id"
	headerID          = "X-Configsum-Id"
	headerCreatedAt   = "X-Configsum-Created"
)

// URL fragments.
const (
	varBaseConfig muxVar = "baseConfig"
	varID         muxVar = "id"
)

type muxVar string

// MakeBaseHandler returns an http.Handler for the base config service.
func MakeBaseHandler(
	svc BaseService,
	opts ...kithttp.ServerOption,
) http.Handler {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Methods("GET").Path(`/`).Name("configBaseList").Handler(
		kithttp.NewServer(
			baseListEndpoint(svc),
			decodeBaseListRequest,
			kithttp.EncodeJSONResponse,
			opts...,
		),
	)

	r.Methods("POST").Path(`/`).Name("configBaseCreate").Handler(
		kithttp.NewServer(
			baseCreateEndpoint(svc),
			confhttp.DecodeJSONSchema(decodeBaseCreateRequest, schemaBaseCreateRequest),
			kithttp.EncodeJSONResponse,
			opts...,
		),
	)

	r.Methods("GET").Path(`/{id:[a-zA-Z0-9]+}`).Name("configBaseGet").Handler(
		kithttp.NewServer(
			baseGetEndpoint(svc),
			decodeBaseGetRequest,
			kithttp.EncodeJSONResponse,
			append(
				opts,
				kithttp.ServerBefore(extractMuxVars(varID)),
			)...,
		),
	)

	r.Methods("PUT").Path(`/{id:[a-zA-Z0-9]+}`).Name("configBaseUpdate").Handler(
		kithttp.NewServer(
			baseUpdateEndpoint(svc),
			confhttp.DecodeJSONSchema(decodeBaseUpdateRequest, schemaBaseUpdateRequest),
			kithttp.EncodeJSONResponse,
			append(
				opts,
				kithttp.ServerBefore(extractMuxVars(varID)),
			)...,
		),
	)

	return r
}

// MakeHandler returns an http.Handler for the user config service.
func MakeHandler(
	svc UserService,
	auth endpoint.Middleware,
	opts ...kithttp.ServerOption,
) http.Handler {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Methods("PUT").Path(`/{baseConfig:[a-z0-9\-]+}`).Name("configUserRender").Handler(
		kithttp.NewServer(
			auth(userRenderEndpoint(svc)),
			confhttp.DecodeJSONSchema(decodeUserRenderRequest, schemaUserRenderRequest),
			encodeUserRenderResponse,
			append(
				opts,
				kithttp.ServerBefore(extractMuxVars(varBaseConfig)),
			)...,
		),
	)

	return r
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

func decodeBaseCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	v := struct {
		ClientID string `json:"client_id"`
		Name     string `json:"name"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return nil, errors.Wrapf(errors.ErrInvalidPayload, "%s", err)
	}

	return baseCreateRequest{clientID: v.ClientID, name: v.Name}, nil
}

func decodeBaseGetRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, ok := ctx.Value(varID).(string)
	if !ok {
		return nil, errors.Wrap(errors.ErrVarMissing, "id missing")
	}

	return baseGetRequest{id: id}, nil
}

func decodeBaseListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return baseListRequest{}, nil
}

func decodeBaseUpdateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, ok := ctx.Value(varID).(string)
	if !ok {
		return nil, errors.Wrap(errors.ErrVarMissing, "id missing")
	}

	v := struct {
		Parameters rule.Parameters `json:"parameters"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return nil, errors.Wrapf(errors.ErrInvalidPayload, "%s", err)
	}

	return baseUpdateRequest{
		id:         id,
		parameters: v.Parameters,
	}, nil
}

func decodeUserRenderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	baseConfig, ok := ctx.Value(varBaseConfig).(string)
	if !ok {
		return nil, errors.Wrap(errors.ErrVarMissing, "baseConfig missing")
	}

	c := userRenderContext{}

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		return nil, errors.Wrapf(errors.ErrInvalidPayload, "%s", err)
	}

	return userRenderRequest{
		baseConfig: baseConfig,
		context:    c,
	}, nil
}

func encodeUserRenderResponse(
	_ context.Context,
	w http.ResponseWriter,
	response interface{},
) error {
	r := response.(userRenderResponse)

	w.Header().Set(headerContentType, "application/json; charset=utf-8")
	w.Header().Set(headerBaseID, r.baseID)
	w.Header().Set(headerBaseName, r.baseName)
	w.Header().Set(headerClientID, r.clientID)
	w.Header().Set(headerID, r.id)
	w.Header().Set(headerCreatedAt, r.createdAt.Format(time.RFC3339Nano))

	return json.NewEncoder(w).Encode(r.rendered)
}
