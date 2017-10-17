package simple

import (
	"context"
	"net/http"
)

const headerUserID = "X-Configsum-Userid"

// HTTPToContext moves the userID from the X-Configsum-Userid header into the
// context of the request.
func HTTPToContext(ctx context.Context, r *http.Request) context.Context {
	userID := r.Header.Get(headerUserID)

	if userID == "" {
		return ctx
	}

	return context.WithValue(ctx, contextKeyUserID, userID)
}
