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
		baseID = randString(characterSet)
		userID = randString(numCharacterSet)
		render = rendered{
			randString(numCharacterSet): rand.Intn(128),
		}
		repo = preparePGUserRepo(t)
	)

	_, err := repo.GetLatest(baseID, userID)
	if errors.Cause(err) != ErrNotFound {
		t.Fatalf("expected ErrNotFound")
	}

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Put(
		randString(characterSet),
		randString(characterSet),
		randString(numCharacterSet),
		render,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Put(randString(characterSet), baseID, userID, map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Put(id.String(), baseID, userID, render)
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

	if have, want := c.rendered, render; reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestPostgresUserRepoGetNotFound(t *testing.T) {
	t.Fail()
}

func TestPostgresUserRepoPutDuplicate(t *testing.T) {
	t.Fail()
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
