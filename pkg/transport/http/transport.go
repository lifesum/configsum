package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

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
)

// ErrorEncoder translates domain specific errors to HTTP status codes.
func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	switch errors.Cause(err) {
	case errors.ErrExists:
		w.WriteHeader(http.StatusConflict)
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

// PopulateRequestContext extracts common information about a request and stores
// it in the context.
func PopulateRequestContext(ctx context.Context, r *http.Request) context.Context {
	route := "unknown"

	if current := mux.CurrentRoute(r); current != nil {
		route = current.GetName()
	}

	ctx = context.WithValue(ctx, ctxKeyTimeBegin, time.Now())
	ctx = context.WithValue(ctx, ctxKeyRoute, route)

	return ctx
}

// ServerFinalizer instruments handler calls to expose Prometheus metrics and
/// log request/response information.
func ServerFinalizer(
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
