package dory

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"

	"github.com/lifesum/configsum/pkg/auth"
)

type contextKey string

// Context keys.
const (
	contextKeySignature contextKey = "dorySignature"
	contextKeyUserID    contextKey = "doryUserID"
)

// AuthMiddleware returns a pluggable endpoint.Middleware which transparently
// inspects Dory specific Authentication information and rejects the request if:
// * signature or userID are missing
// * the signature does not match
func AuthMiddleware(secret string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			signature, ok := ctx.Value(contextKeySignature).(string)
			if !ok {
				return nil, errors.Wrap(ErrSignatureMissing, "request context")
			}

			userID, ok := ctx.Value(contextKeyUserID).(string)
			if !ok {
				return nil, errors.Wrap(ErrUserIDMissing, "request context")
			}

			s, err := hashSignature(secret, userID)
			if err != nil {
				return nil, errors.Wrap(err, "signature hash")
			}

			if s != signature {
				return nil, errors.Wrap(ErrSignatureMissmatch, "auth")
			}

			ctx = context.WithValue(ctx, auth.ContextKeyUserID, userID)

			return next(ctx, request)
		}
	}
}

func hashSignature(secret, userID string) (string, error) {
	h := sha256.New()

	_, err := h.Write([]byte(secret))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum([]byte(userID))), nil
}
