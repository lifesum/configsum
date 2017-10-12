package client

import (
	"context"
	"net/http"
)

const headerToken = "X-Configsum-Token"

// HTTPToContext moves the Client secret token from request header to context.
func HTTPToContext(
	ctx context.Context,
	r *http.Request,
) context.Context {
	secret := r.Header.Get(headerToken)
	if len(secret) != secretLen {
		return ctx
	}

	return context.WithValue(ctx, contextKeySecret, secret)
}
