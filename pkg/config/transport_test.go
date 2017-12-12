package config

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/oklog/ulid"
	"golang.org/x/text/language"

	"github.com/lifesum/configsum/pkg/generate"
	"github.com/lifesum/configsum/pkg/rule"
)

func TestDecodeBaseUpdateRequest(t *testing.T) {
	var (
		seed    = rand.New(rand.NewSource(time.Now().UnixNano()))
		id, _   = ulid.New(ulid.Timestamp(time.Now()), seed)
		ctx     = context.WithValue(context.Background(), varID, id.String())
		payload = bytes.NewBufferString(`{"parameters": {"feature_decode_toggled": true}}`)
		target  = fmt.Sprintf("/%s", id)
		r       = httptest.NewRequest("PUT", target, payload)
	)

	raw, err := decodeBaseUpdateRequest(ctx, r)
	if err != nil {
		t.Fatal(err)
	}

	want := baseUpdateRequest{
		id: id.String(),
		parameters: rule.Parameters{
			"feature_decode_toggled": true,
		},
	}

	if have := raw.(baseUpdateRequest); !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave %#v\nwant %#v", have, want)
	}
}

func TestDecodeUserRenderRequest(t *testing.T) {
	var (
		baseConfig = generate.RandomString(6)
		ctx        = context.WithValue(context.Background(), varBaseConfig, baseConfig)
		locale     = language.MustParse("en_GB")
		payload    = bytes.NewBufferString(`{"device": {"location": {"locale": "en_GB"}}}`)
		target     = fmt.Sprintf("/%s", baseConfig)
		r          = httptest.NewRequest("PUT", target, payload)
	)

	raw, err := decodeUserRenderRequest(ctx, r)
	if err != nil {
		t.Fatal(err)
	}

	want := userRenderRequest{
		baseConfig: baseConfig,
		context: userRenderContext{
			Device: device{
				Location: location{
					locale: locale,
				},
			},
		},
	}

	if have := raw.(userRenderRequest); !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestExtractMuxVars(t *testing.T) {
	var (
		key = muxVar("testKey")
		val = generate.RandomString(12)
		req = httptest.NewRequest("GET", fmt.Sprintf("/root/%s", val), nil)
		r   = mux.NewRouter()
	)

	r.Methods("GET").Path(`/root/{testKey}`).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := extractMuxVars(key)(context.Background(), r)

		if have, want := ctx.Value(key), val; have != want {
			t.Errorf("have %v, want %v", have, want)
		}
	})

	r.ServeHTTP(httptest.NewRecorder(), req)
}
