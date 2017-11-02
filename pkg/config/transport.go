package config

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/instrument"
)

type contextKey string

const (
	ctxKeyTimeBegin contextKey = "begin"
	ctxKeyRoute     contextKey = "route"
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
)

type muxVar string

// MakeHandler returns an http.Handler for the config service.
func MakeHandler(
	logger log.Logger,
	reqObserve instrument.ObserveRequestFunc,
	svc ServiceUser,
	auth endpoint.Middleware,
	opts ...kithttp.ServerOption,
) http.Handler {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.Methods("PUT").Path(`/{baseConfig:[a-z0-9\-]+}`).Name("configUser").Handler(
		kithttp.NewServer(
			auth(userEndpoint(svc)),
			decodeJSONSchema(decodeUserRequest, decodeClientPayloadSchema),
			encodeUserResponse,
			append(
				opts,
				kithttp.ServerBefore(kithttp.PopulateRequestContext),
				kithttp.ServerBefore(populateRequestContext),
				kithttp.ServerBefore(extractMuxVars(varBaseConfig)),
				kithttp.ServerErrorEncoder(encodeError),
				kithttp.ServerErrorLogger(log.With(logger, "route", "configUser")),
				kithttp.ServerFinalizer(
					serverFinalizer(
						log.With(logger, "route", "configUser"),
						reqObserve,
					),
				),
			)...,
		),
	)

	return r
}

func populateRequestContext(ctx context.Context, r *http.Request) context.Context {
	route := "unknown"

	if current := mux.CurrentRoute(r); current != nil {
		route = current.GetName()
	}

	ctx = context.WithValue(ctx, ctxKeyTimeBegin, time.Now())
	ctx = context.WithValue(ctx, ctxKeyRoute, route)

	return ctx
}

func decodeJSONSchema(
	next kithttp.DecodeRequestFunc,
	schema *gojsonschema.Schema,
) kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		raw, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		if len(raw) == 0 {
			return nil, errors.Wrap(errors.ErrInvalidPayload, "empty body")
		}
		res, err := schema.Validate(gojsonschema.NewBytesLoader(raw))
		if nil != err {
			return nil, errors.Wrap(errors.ErrInvalidPayload, err.Error())
		}

		if !res.Valid() {
			err := errors.ErrInvalidPayload

			for _, e := range res.Errors() {
				err = errors.Wrap(err, e.String())
			}

			return nil, err
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(raw))

		return next(ctx, r)
	}
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

func decodeUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	baseConfig, ok := ctx.Value(varBaseConfig).(string)
	if !ok {
		return nil, errors.Wrap(errors.ErrVarMissing, "baseConfig missing")
	}

	c := userContext{}

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		return nil, err
	}

	return userRequest{
		baseConfig: baseConfig,
		context:    c,
	}, nil
}

func encodeUserResponse(
	_ context.Context,
	w http.ResponseWriter,
	response interface{},
) error {
	r := response.(userResponse)

	w.Header().Set(headerContentType, "application/json; charset=utf-8")
	w.Header().Set(headerBaseID, r.baseID)
	w.Header().Set(headerBaseName, r.baseName)
	w.Header().Set(headerClientID, r.clientID)
	w.Header().Set(headerID, r.id)
	w.Header().Set(headerCreatedAt, r.createdAt.Format(time.RFC3339Nano))

	return json.NewEncoder(w).Encode(r.rendered)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch errors.Cause(err) {
	case errors.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case errors.ErrClientNotFound, errors.ErrSecretMissing:
		w.WriteHeader(http.StatusUnauthorized)
	case errors.ErrSignatureMissing, errors.ErrSignatureMissmatch, errors.ErrUserIDMissing:
		w.WriteHeader(http.StatusUnauthorized)
	case errors.ErrInvalidPayload:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set(headerContentType, "application/json; charset=utf-8")

	_ = json.NewEncoder(w).Encode(struct {
		Reason string `json:"reason"`
	}{
		Reason: err.Error(),
	})
}

func serverFinalizer(
	logger log.Logger,
	reqObserve instrument.ObserveRequestFunc,
) kithttp.ServerFinalizerFunc {
	return func(ctx context.Context, code int, r *http.Request) {
		var (
			begin  = ctx.Value(ctxKeyTimeBegin).(time.Time)
			method = ctx.Value(kithttp.ContextKeyRequestMethod).(string)
			host   = ctx.Value(kithttp.ContextKeyRequestHost).(string)
			proto  = ctx.Value(kithttp.ContextKeyRequestProto).(string)
			route  = ctx.Value(ctxKeyRoute).(string)
		)

		_ = logger.Log(
			"duration", time.Since(begin).Nanoseconds(),
			"request", map[string]interface{}{
				"authorization":    ctx.Value(kithttp.ContextKeyRequestAuthorization),
				"header":           r.Header,
				"host":             host,
				"method":           method,
				"path":             ctx.Value(kithttp.ContextKeyRequestPath),
				"proto":            proto,
				"referer":          ctx.Value(kithttp.ContextKeyRequestReferer),
				"remoteAddr":       ctx.Value(kithttp.ContextKeyRequestRemoteAddr),
				"requestId":        ctx.Value(kithttp.ContextKeyRequestXRequestID),
				"requestUri":       ctx.Value(kithttp.ContextKeyRequestURI),
				"transferEncoding": r.TransferEncoding,
			},
			"response", map[string]interface{}{
				"header":     ctx.Value(kithttp.ContextKeyResponseHeaders).(http.Header),
				"size":       ctx.Value(kithttp.ContextKeyResponseSize).(int64),
				"statusCode": code,
			},
			"route", route,
		)

		reqObserve(code, host, method, proto, route, begin)
	}
}
