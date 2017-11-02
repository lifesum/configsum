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
