package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/generate"
)

func TestDecodeCreateRequest(t *testing.T) {
	var (
		ctx     = context.TODO()
		name    = generate.RandomString(12)
		payload = bytes.NewBufferString(fmt.Sprintf(`{"name": "%s"}`, name))
		r       = httptest.NewRequest("POST", "/", payload)
	)

	raw, err := decodeCreateRequest(ctx, r)
	if err != nil {
		t.Fatal(err)
	}

	want := createRequest{
		name: name,
	}

	if have := raw.(createRequest); !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestDecodeCreateRequestInvalid(t *testing.T) {
	cases := []string{
		``,   // Empty body.
		`{}`, // name missing
	}

	for _, c := range cases {
		var (
			ctx     = context.TODO()
			payload = bytes.NewBufferString(c)
			r       = httptest.NewRequest("POST", "/", payload)
		)

		_, err := decodeCreateRequest(ctx, r)
		if have, want := errors.Cause(err), errors.ErrInvalidPayload; have != want {
			t.Errorf("have %v, want %v", have, want)
		}
	}

}

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
