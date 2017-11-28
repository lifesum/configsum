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
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/text/language"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/generate"
	"github.com/lifesum/configsum/pkg/rule"
)

const (
	testSchema = `
	{
		  "$schema": "http://json-schema.org/draft-06/schema#",
		  "title": "Client payload",
		  "description": "Common set of information to determine device capabilities and user provided info.",
		  "type": "object",
		  "properties": {
		  	"app": {
		  		"type": "object",
		  		"properties": {
		  			"version": {
		  				"description": "The version of the client application.",
		  				"type": "string"
		  			}
		  		},
		  		"required": ["version"]
		  	}
		  },
		  "required": ["app"]
	}`
)

func TestDecodeJSONSchemaValidInput(t *testing.T) {
	var (
		next = func(ctx context.Context, r *http.Request) (interface{}, error) {
			return true, nil
		}
		payload = bytes.NewBufferString(`{"app": {"version": "6.4.1"}}`)
		req     = httptest.NewRequest("PUT", "/v1/config/foo", payload)
	)

	schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(testSchema))
	if err != nil {
		t.Fatal(err)
	}

	have, _ := decodeJSONSchema(next, schema)(context.Background(), req)
	if want := true; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestDecodeJSONSchemaEmptyBody(t *testing.T) {
	var (
		next = func(ctx context.Context, r *http.Request) (interface{}, error) {
			return nil, nil
		}
		req = httptest.NewRequest("PUT", "/v1/config/foo", nil)
	)

	schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(testSchema))
	if err != nil {
		t.Fatal(err)
	}

	_, have := decodeJSONSchema(next, schema)(context.Background(), req)
	if want := errors.ErrInvalidPayload; errors.Cause(have) != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestDecodeJSONSchemaMissingField(t *testing.T) {
	var (
		next = func(ctx context.Context, r *http.Request) (interface{}, error) {
			return nil, nil
		}
		payload = bytes.NewBufferString(`{"platform": "WatchOS"}`)
		req     = httptest.NewRequest("PUT", "/v1/config/foo", payload)
	)

	schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(testSchema))
	if err != nil {
		t.Fatal(err)
	}

	_, have := decodeJSONSchema(next, schema)(context.Background(), req)
	if want := errors.ErrInvalidPayload; errors.Cause(have) != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

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
	)

	r := mux.NewRouter()

	r.Methods("GET").Path(`/root/{testKey}`).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := extractMuxVars(key)(context.Background(), r)

		if have, want := ctx.Value(key), val; have != want {
			t.Errorf("have %v, want %v", have, want)
		}
	})

	req := httptest.NewRequest("GET", fmt.Sprintf("/root/%s", val), nil)

	r.ServeHTTP(httptest.NewRecorder(), req)
}
