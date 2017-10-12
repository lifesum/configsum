package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/lifesum/configsum/pkg/generate"
)

func TestHTTPToContext(t *testing.T) {
	secret, err := generate.SecureToken(secretByteLen)
	if err != nil {
		t.Fatal(err)
	}

	ts := map[*http.Request]interface{}{
		&http.Request{}: nil,
		&http.Request{
			Header: http.Header{
				headerToken: []string{"invalidSecret"},
			},
		}: nil,
		&http.Request{
			Header: http.Header{
				headerToken: []string{secret},
			},
		}: secret,
	}

	for r, want := range ts {
		ctx := HTTPToContext(context.TODO(), r)

		if have := ctx.Value(contextKeySecret); have != want {
			t.Errorf("have %v, want %v", have, want)
		}
	}

}
