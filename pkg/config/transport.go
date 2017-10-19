package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/lifesum/configsum/pkg/auth/dory"
	"github.com/lifesum/configsum/pkg/auth/simple"
	"github.com/lifesum/configsum/pkg/client"
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

// MakeHandler returns an http.Handler for the config service.
func MakeHandler(
	logger log.Logger,
	svc ServiceUser,
	auth endpoint.Middleware,
	opts ...kithttp.ServerOption,
) http.Handler {
	r := mux.NewRouter()
	r.StrictSlash(true)

	opts = append(
		opts,
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerErrorLogger(log.With(logger, "route", "configUser")),
		kithttp.ServerFinalizer(serverFinalizer(log.With(logger, "route", "configUser"))),
	)

	r.Methods("PUT").Path(`/{baseConfig:[a-z0-9\-]+}`).Name("configUser").Handler(
		kithttp.NewServer(
			auth(userEndpoint(svc)),
			decodeUserRequest,
			encodeUserResponse,
			opts...,
		),
	)

	return r
}

func decodeUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	c, ok := mux.Vars(r)["baseConfig"]
	if !ok {
		return nil, fmt.Errorf("Baseconfig missing")
	}

	return userRequest{baseConfig: c}, nil
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
	case ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case client.ErrNotFound, client.ErrSecretMissing:
		w.WriteHeader(http.StatusUnauthorized)
	case dory.ErrSignatureMissing, dory.ErrSignatureMissmatch, dory.ErrUserIDMissing:
		w.WriteHeader(http.StatusUnauthorized)
	case simple.ErrUserIDMissing:
		w.WriteHeader(http.StatusUnauthorized)
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

func serverFinalizer(logger log.Logger) kithttp.ServerFinalizerFunc {
	return func(ctx context.Context, code int, r *http.Request) {
		_ = logger.Log(
			"request", map[string]interface{}{
				"authorization":    ctx.Value(kithttp.ContextKeyRequestAuthorization),
				"header":           r.Header,
				"host":             ctx.Value(kithttp.ContextKeyRequestHost),
				"method":           ctx.Value(kithttp.ContextKeyRequestMethod),
				"path":             ctx.Value(kithttp.ContextKeyRequestPath),
				"proto":            ctx.Value(kithttp.ContextKeyRequestProto),
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
		)
	}
}
