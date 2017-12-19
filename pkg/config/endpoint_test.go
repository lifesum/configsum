package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/lifesum/configsum/pkg/errors"
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

	timestamp, err := time.Parse(time.RFC3339, "2017-12-04T23:11:38Z")
	if err != nil {
		t.Fatal(err)
	}

	err = u.UnmarshalJSON(input)
	if err != nil {
		t.Fatal(err)
	}

	have := u
	want := userInfo{
		Age:          27,
		Registered:   timestamp,
		Subscription: 2,
	}

	if !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}
