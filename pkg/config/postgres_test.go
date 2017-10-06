// +build integration

package config

import (
	"flag"
	"fmt"
	"math/rand"
	"os/user"
	"reflect"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid"
	"github.com/pkg/errors"
	// Blank import for Postgres capabilities.
	_ "github.com/lib/pq"

	"github.com/lifesum/configsum/pkg/pg"
)

const (
	characterSet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numCharacterSet = "0123456789"
)

var (
	seed = rand.New(rand.NewSource(time.Now().UnixNano()))

	pgURI string
)

func TestPostgresUserRepoGetLatest(t *testing.T) {
	var (
		baseID    = randString(characterSet)
		userID    = randString(numCharacterSet)
		render    = rendered{}
		repo      = preparePGUserRepo(t)
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

func TestPostgresUserRepoGetLatestNotFound(t *testing.T) {
	var (
		baseID = randString(characterSet)
		userID = randString(numCharacterSet)
		repo   = preparePGUserRepo(t)
	)

	_, err := repo.GetLatest(baseID, userID)
	if have, want := errors.Cause(err), ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestPostgresUserRepoAppendDuplicate(t *testing.T) {
	var (
		baseID    = randString(characterSet)
		userID    = randString(numCharacterSet)
		render    = rendered{}
		repo      = preparePGUserRepo(t)
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
	if have, want := errors.Cause(err), ErrExists; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func randString(charset string) string {
	b := make([]byte, len(charset))

	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}

	return string(b)
}

func preparePGUserRepo(t *testing.T) UserRepo {
	db, err := sqlx.Connect("postgres", pgURI)
	if err != nil {
		t.Fatal(err)
	}

	r, err := NewPostgresUserRepo(db)
	if err != nil {
		t.Fatal(err)
	}

	if err := r.Teardown(); err != nil {
		t.Fatal(err)
	}

	return r
}

func init() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	uri := flag.String("postgres.uri", fmt.Sprintf(pg.DefaultTestURI, u.Username), "Postgres connection URL")

	flag.Parse()

	pgURI = *uri
}
