package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

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
)

// DecodeJSONSchema validates the request payload against the given schema and
// returns an invalid payload error in case the validation fails.
func DecodeJSONSchema(
	next kithttp.DecodeRequestFunc,
	schema *gojsonschema.Schema,
) kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		// Attempts to work with gojsonschema.NewReaderLoader turned out to
		// lead to inexplicable io errors. While less elegant it works for
		// now. Be aware of the re-assignment of the request body so
		// subsequent decode functions can function properly.
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
	case errors.ErrInvalidPayload, errors.ErrParametersInvalid:
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
