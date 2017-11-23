package config

import (
	"testing"

	"github.com/xeipuuv/gojsonschema"
)

func TestSchemaBaseCreateInvalid(t *testing.T) {
	cases := []string{
		`{}`, // Empty.
		`{"client_id": "clientID"}`,  // Name missing,
		`{"name": "baseConfigName"}`, // ClientID missing.
	}

	for _, c := range cases {
		r, err := schemaBaseCreateRequest.Validate(gojsonschema.NewStringLoader(c))
		if err != nil {
			t.Fatal(err)
		}

		if r.Valid() {
			t.Errorf("invalid: %s", c)
		}
	}
}

func TestSchemaBaseUpdateInvalid(t *testing.T) {
	cases := []string{
		`{}`,                                                      // Missing parameters.
		`{"parameters": {}}`,                                      // Empty parameters.
		`{"parameters": {"feature_array_toggled": {}} }`,          // Array as parameter type.
		`{"parameters": {"feature_object_toggled": []} }`,         // Object as parameter type.
		`{"parameters": {"feature_inv4l1d$char_toggled": true} }`, // Invalid character in parameter key.
	}

	for _, input := range cases {
		res, err := schemaBaseUpdateRequest.Validate(gojsonschema.NewStringLoader(input))
		if err != nil {
			t.Fatal(err)
		}

		if res.Valid() {
			t.Errorf("invalid: %s", input)
		}
	}
}

func TestSchemaUserRenderInvalid(t *testing.T) {
	var (
		want  = "invalid JSON error"
		cases = []string{
			`{}`,                                                                                                                                // App missing
			`{"app": {}}`,                                                                                                                       // Empty App object
			`{"app": {"version": "6.4.1"}}`,                                                                                                     // Device missing
			`{"app": {"version": "6.4.1"}, "device": {}, "os": {"platform": "WatchOS", "version": "9.4"}}`,                                      // Empty Device object
			`{"app": {"version": "6.4.1"}, "device": {"location": {}, "os": {"platform": "WatchOS", "version": "9.4"}}}`,                        // Empty Location object
			`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en_GB"}, "os": {"platform": "WatchOS", "version": "9.4"}}}`,       // TimezoneOffset missing
			`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en_GB", "timezoneOffset": 7200}}}`,                                // OS missing
			`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en_GB", "timezoneOffset": 7200}, "os": {}}}`,                      // Empty Os object
			`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en_GB", "timezoneOffset": 7200}, "os": {"platform": "WatchOS"}}}`, // Version missing
		}
	)

	for _, input := range cases {
		res, err := schemaUserRenderRequest.Validate(gojsonschema.NewStringLoader(input))
		if err != nil {
			t.Fatal(err)
		}

		if res.Valid() {
			t.Errorf("have %v, want %v", input, want)
		}
	}
}

func TestSchemaUserRenderValid(t *testing.T) {
	var (
		want  = "valid JSON input"
		cases = []string{
			`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en_GB", "timezoneOffset": 7201}, "os": {"platform": "WatchOS", "version": "9.4"}}, "user": {"age": 23}}`,                   // Working case
			`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en_GB", "timezoneOffset": 7201}, "os": {"platform": "WatchOS", "version": "9.4"}}, "user": {}}`,                            // User age optional
			`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en", "timezoneOffset": 7201}, "os": {"platform": "WatchOS", "version": "9.4"}}, "user": {"age": 27}}`,                      // Only region provided
			`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en_GB", "timezoneOffset": 7201}, "metadata": null, "os": {"platform": "WatchOS", "version": "9.4"}}, "user": {"age": 23}}`, // Metadata value is null
		}
	)

	for _, input := range cases {
		res, err := schemaUserRenderRequest.Validate(gojsonschema.NewStringLoader(input))
		if err != nil {
			t.Fatal(err)
		}

		if !res.Valid() {
			t.Errorf("have %v, want %v", input, want)
		}
	}
}
