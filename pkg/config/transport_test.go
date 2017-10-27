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
		  	},
		    "metadata": {
		      "type": "object",
		      "additionalProperties": {
		        "anyOf": [
		          {
		            "type": "string"
		          },
		          {
		            "type": "integer"
		          },
		          {
		            "type": "array",
		            "items": {
		              "type": "string"
		            }
		          },
		          {
		            "type": "array",
		            "items": {
		              "type": "integer"
		            }
		          }
		        ]
		      }
		    },
		    "device": {
		    	"type": "object",
		    	"properties": {
		    		"location": {
		    			"type": "object",
		    			"properties": {
		    				"locale": {
		    					"description": "The device's locale setting according to ISO 639-3 and/or BCP 47.",
		    					"type": "string",
		    					"pattern": "^[a-z]{2,3}_([A-Z]{2}|[0-9]{3})$"
		    				},
		    				"timezoneOffset": {
		    					"description": "time offset from GMT",
		    					"type": "integer"
		    				}
		    			},
		    			"required": ["locale", "timezoneOffset"]
		    		},
		    		"os": {
		    			"type": "object",
		    			"properties": {
		    				"platform": {
		    					"description": "The client platform that makes this request.",
		    					"enum": [ "Android", "iOS", "WatchOS" ]
		    				},
		    				"version": {
		    					"description": "Version of the os that runs on the client platform",
		    					"type": "string"
		    				}
		    			},
		    			"required": ["platform", "version"]
		    		}
		    	},
		    	"required": ["location", "os"]
		    },
		    "user": {
		    	"type": "object",
		    	"properties": {
		    		"age": {
		    			"description": "Age of the application's logged in user.",
		    			"type": "integer"
		    		}
		    	},
		    	"required": ["age"]
		    }
		  },
		  "required": ["app", "device"]
	}`

	failingSchema = `
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

var testInputs = []string{
	`{}`,                                                                                                                                                                   // App missing
	`{"app": {}}`,                                                                                                                                                          // Empty App object
	`{"app": {"version": "6.4.1"}}`,                                                                                                                                        // Device missing
	`{"app": {"version": "6.4.1"}, "device": {}, "os": {"platform": "WatchOS", "version": "9.4"}}`,                                                                         // Empty Device object
	`{"app": {"version": "6.4.1"}, "device": {"location": {}, "os": {"platform": "WatchOS", "version": "9.4"}}}`,                                                           // Empty Location object
	`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en_GB"}, "os": {"platform": "WatchOS", "version": "9.4"}}}`,                                          // TimezoneOffset missing
	`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en_GB", "timezoneOffset": 7200}}}`,                                                                   // OS missing
	`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en_GB", "timezoneOffset": 7200}, "os": {}}}`,                                                         // Empty Os object
	`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en_GB", "timezoneOffset": 7200}, "os": {"platform": "WatchOS"}}}`,                                    // Version missing
	`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "testLocale", "timezoneOffset": 7201}, "os": {"platform": "WatchOS", "version": "9.4"}}}`,             // Wrong locale format
	`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "testLocale", "timezoneOffset": 7201}, "os": {"platform": "WatchOS", "version": "9.4"}}, "user": {}}`, // Empty User object
}

func TestPayloadValidation(t *testing.T) {

	schema := gojsonschema.NewStringLoader(testSchema)

	var (
		want = "invalid JSON error"
	)

	for _, input := range testInputs {
		json := gojsonschema.NewStringLoader(input)

		result, err := gojsonschema.Validate(schema, json)
		if err != nil {
			t.Fatal(err)
		}

		if result.Valid() {
			t.Errorf("have %v, want %v", input, want)
		}
	}
}

func TestDecodeJSONSchemaEmptyBody(t *testing.T) {

	var (
		next = func(ctx context.Context, r *http.Request) (interface{}, error) {
			return nil, nil
		}
		req = httptest.NewRequest("PUT", "/v1/config/foo", nil)
	)

	schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(failingSchema))
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
