package dory

import (
	"context"
	"net/http"
	"testing"

	"github.com/lifesum/configsum/pkg/generate"
)

func TestHTTPToContext(t *testing.T) {
	type expect struct {
		signature interface{}
		userID    interface{}
	}

	var (
		secret = generate.RandomString(32)
		userID = generate.RandomString(24)
	)

	signature, err := hashSignature(secret, userID)
	if err != nil {
		t.Fatal(err)
	}

	ts := map[*http.Request]expect{
		&http.Request{}: {nil, nil},
		&http.Request{
			Header: http.Header{
				headerSignature: []string{signature},
			},
		}: {nil, nil},
		&http.Request{
			Header: http.Header{
				headerSignature: []string{signature},
				headerUserID:    []string{userID},
			},
		}: {signature, userID},
	}

	for r, e := range ts {
		ctx := HTTPToContext(context.TODO(), r)

		if have, want := ctx.Value(contextKeySignature), e.signature; have != want {
			t.Errorf("have %v, want %v", have, want)
		}

		if have, want := ctx.Value(contextKeyUserID), e.userID; have != want {
			t.Errorf("have %v, want %v", have, want)
		}
	}
}
