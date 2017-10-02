// +build integration

package config

import (
	"flag"
	"fmt"
	"math/rand"
	"os/user"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	// Blank import for Postgres capabilities.
	_ "github.com/lib/pq"

	"github.com/lifesum/configsum/pkg/pg"
)

const (
	characterSet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numCharacterSet = "0123456789"
)

var pgURI string

func TestPostgresUserRepoGet(t *testing.T) {
	var (
		baseConfig = randString(characterSet)
		userID     = randString(numCharacterSet)
		repo       = preparePGUserRepo(t)
	)

	_, err := repo.Get(baseConfig, userID)
	if err == nil {
		t.Fatalf("expected ErrNotFound")
	}

	if errors.Cause(err) != ErrNotFound {
		t.Fatalf("expected ErrNotFound")
	}
}

func randString(charset string) string {
	var (
		s = rand.New(rand.NewSource(time.Now().UnixNano()))
		b = make([]byte, len(charset))
	)

	for i := range b {
		b[i] = charset[s.Intn(len(charset))]
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
