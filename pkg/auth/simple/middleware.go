package simple

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"

	"github.com/lifesum/configsum/pkg/auth"
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
				return nil, errors.Wrap(ErrUserIDMissing, "request context")
			}

			ctx = context.WithValue(ctx, auth.ContextKeyUserID, userID)

			return next(ctx, request)
		}
	}
}
