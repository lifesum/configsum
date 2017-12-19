package config

import (
	"github.com/lifesum/configsum/pkg/errors"

	"testing"
)

func TestLocationInvalidLocale(t *testing.T) {
	var (
		input = []byte(`{"locale": "foobarz"}`)
		l     = location{}
	)

	err := l.UnmarshalJSON(input)
	if have, want := errors.Cause(err), errors.ErrInvalidPayload; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestUnmarshalUserContext(t *testing.T) {
	var (
		input = []byte(`{
			"age": 27,
			"registered": "2017-12-04T23:11:38Z",
			"subscription": 2
			}`)
		u = userInfo{}
	)

	err := u.UnmarshalJSON(input)
	if err != nil {
		t.Fatal(err)
	}
}
