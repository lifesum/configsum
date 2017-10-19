package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/oklog/ulid"

	"github.com/lifesum/configsum/pkg/errors"
)

func testUserRepoGetLatest(t *testing.T, p prepareFunc) {
	var (
		baseID    = randString(characterSet)
		userID    = randString(numCharacterSet)
		render    = rendered{}
		repo      = p(t)
		decisions = ruleDecisions{
			randString(numCharacterSet): []int{seed.Intn(100), seed.Intn(100)},
			randString(numCharacterSet): []int{seed.Intn(100)},
			randString(numCharacterSet): []int{},
		}
	)

	render.setNumber(randString(characterSet), seed.Float64())

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Append(
		randString(characterSet),
		randString(characterSet),
		randString(numCharacterSet),
		decisions,
		render,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Append(
		randString(characterSet),
		baseID,
		userID,
		nil,
		rendered{},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Append(id.String(), baseID, userID, decisions, render)
	if err != nil {
		t.Fatal(err)
	}

	c, err := repo.GetLatest(baseID, userID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := c.baseID, baseID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := c.id, id.String(); have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := c.userID, userID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := c.rendered, render; !reflect.DeepEqual(c.rendered, render) {
		t.Errorf("\nhave %#v,\nwant %#v", have, want)
	}

	if have, want := c.ruleDecisions, decisions; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testUserRepoGetLatestNotFound(t *testing.T, p prepareFunc) {
	var (
		baseID = randString(characterSet)
		userID = randString(numCharacterSet)
		repo   = p(t)
	)

	_, err := repo.GetLatest(baseID, userID)
	if have, want := errors.Cause(err), errors.ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testUserRepoAppendDuplicate(t *testing.T, p prepareFunc) {
	var (
		baseID    = randString(characterSet)
		userID    = randString(numCharacterSet)
		render    = rendered{}
		repo      = p(t)
		decisions = ruleDecisions{
			randString(numCharacterSet): []int{seed.Int(), seed.Int()},
			randString(numCharacterSet): []int{seed.Int()},
			randString(numCharacterSet): []int{},
		}
	)

	render.SetBool(randString(numCharacterSet), false)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Append(id.String(), baseID, userID, decisions, render)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Append(id.String(), baseID, userID, decisions, render)
	if have, want := errors.Cause(err), errors.ErrExists; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}
