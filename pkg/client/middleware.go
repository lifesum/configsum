package client

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/lifesum/configsum/pkg/errors"
)

type contextKey string

// Context keys.
const (
	ContextKeyClientID contextKey = "clientID"
	contextKeySecret   contextKey = "clientSecret"
)

const (
	secretLen     = 44
	secretByteLen = 32
)

// AuthMiddleware returns a client token based authentication.
func AuthMiddleware(svc Service) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			secret, ok := ctx.Value(contextKeySecret).(string)
			if !ok {
				return nil, errors.Wrap(errors.ErrSecretMissing, "request context")
			}

			c, err := svc.LookupBySecret(secret)
			if err != nil {
				return nil, errors.Wrap(errors.ErrClientNotFound, err.Error())
			}

			ctx = context.WithValue(ctx, ContextKeyClientID, c.id)

			return next(ctx, request)
		}
	}
}
