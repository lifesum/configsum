// +build integration

package config

import (
	"flag"
	"fmt"
	"os/user"
	"testing"

	"github.com/jmoiron/sqlx"
	// Blank import for Postgres capabilities.
	_ "github.com/lib/pq"

	"github.com/lifesum/configsum/pkg/pg"
)

var pgURI string

func TestPostgresUserRepoGetLatest(t *testing.T) {
	testUserRepoGetLatest(t, preparePGUserRepo)
}

func TestPostgresUserRepoGetLatestNotFound(t *testing.T) {
	testUserRepoGetLatestNotFound(t, preparePGUserRepo)
}

func TestPostgresUserRepoAppendDuplicate(t *testing.T) {
	testUserRepoAppendDuplicate(t, preparePGUserRepo)
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
