package simple

import (
	"context"
	"net/http"
	"testing"

	"github.com/lifesum/configsum/pkg/generate"
)

func TestHTTPToContext(t *testing.T) {
	var (
		userID = generate.RandomString(24)
		ts     = map[*http.Request]interface{}{
			&http.Request{}: nil,
			&http.Request{
				Header: http.Header{
					headerUserID: []string{userID},
				},
			}: userID,
		}
	)

	for r, want := range ts {
		ctx := HTTPToContext(context.TODO(), r)

		if have := ctx.Value(contextKeyUserID); have != want {
			t.Errorf("have %v, want %v", have, want)
		}
	}
}
