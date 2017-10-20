package simple

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/lifesum/configsum/pkg/auth"
	"github.com/lifesum/configsum/pkg/errors"
)

type contextKey string

const (
	contextKeyUserID contextKey = "simpleUserID"
)

// AuthMiddleware returns a pluggable endpoint.Middleware which rejects the
// request if:
// * userID is missing from the context
func AuthMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			userID, ok := ctx.Value(contextKeyUserID).(string)
			if !ok {
				return nil, errors.Wrap(errors.ErrUserIDMissing, "request context")
			}

			ctx = context.WithValue(ctx, auth.ContextKeyUserID, userID)

			return next(ctx, request)
		}
	}
}
