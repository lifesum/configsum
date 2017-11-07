package config

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/text/language"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/generate"
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

func TestDecodeRenderRequest(t *testing.T) {
	var (
		baseConfig = generate.RandomString(6)
		ctx        = context.WithValue(context.Background(), varBaseConfig, baseConfig)
		locale     = language.MustParse("en_GB")
		payload    = bytes.NewBufferString(`{"device": {"location": {"locale": "en_GB"}}}`)
		target     = fmt.Sprintf("/%s", baseConfig)
		r          = httptest.NewRequest("PUT", target, payload)
	)

	raw, err := decodeRenderRequest(ctx, r)
	if err != nil {
		t.Fatal(err)
	}

	want := renderRequest{
		baseConfig: baseConfig,
		context: renderContext{
			Device: device{
				Location: location{
					locale: locale,
				},
			},
		},
	}

	if have := raw.(renderRequest); !reflect.DeepEqual(have, want) {
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
