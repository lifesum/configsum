package config

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xeipuuv/gojsonschema"

	"github.com/lifesum/configsum/pkg/errors"
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
