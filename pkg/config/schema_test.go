package config

import (
	"testing"

	"github.com/xeipuuv/gojsonschema"
)

var (
	schema            = gojsonschema.NewStringLoader(requestCapabilities)
	testFailingInputs = []string{
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
	testSuccessInputs = []string{
		`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en_GB", "timezoneOffset": 7201}, "os": {"platform": "WatchOS", "version": "9.4"}}, "user": {"age": 23}}`,                   // Working case
		`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en_GB", "timezoneOffset": 7201}, "os": {"platform": "WatchOS", "version": "9.4"}}, "user": {}}`,                            // User age optional
		`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en", "timezoneOffset": 7201}, "os": {"platform": "WatchOS", "version": "9.4"}}, "user": {"age": 27}}`,                      // Only region provided
		`{"app": {"version": "6.4.1"}, "device": {"location": {"locale": "en_GB", "timezoneOffset": 7201}, "metadata": null, "os": {"platform": "WatchOS", "version": "9.4"}}, "user": {"age": 23}}`, // Metadata value is null
	}
)

func TestPayloadValidationFailure(t *testing.T) {
	want := "invalid JSON error"

	for _, input := range testFailingInputs {
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

func TestPayloadValidationSuccess(t *testing.T) {
	want := "valid JSON input"

	for _, input := range testSuccessInputs {
		json := gojsonschema.NewStringLoader(input)

		result, err := gojsonschema.Validate(schema, json)
		if err != nil {
			t.Fatal(err)
		}

		if !result.Valid() {
			t.Errorf("have %v, want %v", input, want)
		}
	}
}
