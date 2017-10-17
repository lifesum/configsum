package dory

import (
	"context"
	"net/http"
)

const (
	headerSignature = "X-Dory-Signature"
	headerUserID    = "X-Dory-Userid"
)

// HTTPToContext moves the Dory signature and userID from the request headers
// to the context.
func HTTPToContext(ctx context.Context, r *http.Request) context.Context {
	var (
		signature = r.Header.Get(headerSignature)
		userID    = r.Header.Get(headerUserID)
	)

	if signature == "" || userID == "" {
		return ctx
	}

	ctx = context.WithValue(ctx, contextKeySignature, signature)

	return context.WithValue(ctx, contextKeyUserID, userID)
}
