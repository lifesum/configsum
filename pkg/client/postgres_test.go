// +build integration

package client

import (
	"flag"
	"fmt"
	"math/rand"
	"os/user"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid"
	"github.com/pkg/errors"

	"github.com/lifesum/configsum/pkg/generate"
	"github.com/lifesum/configsum/pkg/pg"
)

var pgURI string

func TestPGRepoLookup(t *testing.T) {
	var (
		name = generate.RandomString(24)
		repo = preparePGRepo(t)
		seed = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Store(id.String(), name)
	if err != nil {
		t.Fatal(err)
	}

	c, err := repo.Lookup(id.String())
	if err != nil {
		t.Fatal(err)
	}

	if have, want := c.id, id.String(); have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := c.name, name; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := c.deleted, false; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestPGRepoLookupNotFound(t *testing.T) {
	var (
		id   = generate.RandomString(24)
		repo = preparePGRepo(t)
	)

	_, err := repo.Lookup(id)
	if have, want := errors.Cause(err), ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestPGTokenRepoLookup(t *testing.T) {
	var (
		secret = generate.RandomString(32)
		repo   = preparePGTokenRepo(t)
		seed   = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	clientID, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Store(clientID.String(), secret)
	if err != nil {
		t.Fatal(err)
	}

	token, err := repo.Lookup(secret)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := token.clientID, clientID.String(); have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := token.deleted, false; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := token.secret, secret; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestPGTokenRepoLookupNotFound(t *testing.T) {
	var (
		secret = generate.RandomString(32)
		repo   = preparePGTokenRepo(t)
	)

	_, err := repo.Lookup(secret)
	if have, want := errors.Cause(err), ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func preparePGRepo(t *testing.T) Repo {
	db, err := sqlx.Connect("postgres", pgURI)
	if err != nil {
		t.Fatal(err)
	}

	r := NewPostgresRepo(db)

	if err := r.teardown(); err != nil {
		t.Fatal(err)
	}

	return r
}

func preparePGTokenRepo(t *testing.T) TokenRepo {
	db, err := sqlx.Connect("postgres", pgURI)
	if err != nil {
		t.Fatal(err)
	}

	r := NewPostgresTokenRepo(db)

	if err := r.teardown(); err != nil {
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
