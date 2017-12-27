// +build integration

package client

import (
	"flag"
	"fmt"
	"os/user"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/lifesum/configsum/pkg/pg"
)

var pgURI string

func TestPGRepoList(t *testing.T) {
	testRepoList(t, preparePGRepo)
}

func TestPGRepoListEmpty(t *testing.T) {
	testRepoListEmpty(t, preparePGRepo)
}

func TestPGRepoLookup(t *testing.T) {
	testRepoLookup(t, preparePGRepo)
}

func TestPGRepoLookupNotFound(t *testing.T) {
	testRepoLookupNotFound(t, preparePGRepo)
}

func TestPGTokenRepoGetLatest(t *testing.T) {
	testTokenRepoGetLatest(t, preparePGTokenRepo)
}

func TestPGTokenRepoLookup(t *testing.T) {
	testTokenRepoLookup(t, preparePGTokenRepo)
}

func TestPGTokenRepoLookupNotFound(t *testing.T) {
	testTokenRepoLookupNotFound(t, preparePGTokenRepo)
}

func preparePGRepo(t *testing.T) Repo {
	db, err := sqlx.Connect("postgres", pgURI)
	if err != nil {
		t.Fatal(err)
	}

	r := NewPostgresRepo(db)

	if err := r.Teardown(); err != nil {
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
