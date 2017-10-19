package simple

import (
	"context"
	"testing"

	"github.com/go-kit/kit/endpoint"

	"github.com/lifesum/configsum/pkg/auth"
	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/generate"
)

func TestAuthMiddleware(t *testing.T) {
	var (
		userID = generate.RandomString(24)
		ctx    = context.WithValue(context.TODO(), contextKeyUserID, userID)
	)

	_, err := AuthMiddleware()(nopEndpoint(t, userID))(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthMiddlewareUserIDMissing(t *testing.T) {
	_, err := AuthMiddleware()(nopEndpoint(t, ""))(context.TODO(), nil)
	if have, want := errors.Cause(err), errors.ErrUserIDMissing; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func nopEndpoint(t *testing.T, want string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		have, ok := ctx.Value(auth.ContextKeyUserID).(string)
		if !ok {
			t.Fatalf("userID missing")
		}

		if have != want {
			t.Errorf("have %v, want %v", have, want)
		}

		return true, nil
	}
}
