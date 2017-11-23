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

func TestPostgresBaseRepoCreateDuplicate(t *testing.T) {
	testBaseRepoCreateDuplicate(t, preparePGBaseRepo)
}

func TestPostgresBaseRepoGetByID(t *testing.T) {
	testBaseRepoGetByID(t, preparePGBaseRepo)
}

func TestPostgresBaseRepoGetByIDNotFound(t *testing.T) {
	testBaseRepoGetByIDNotFound(t, preparePGBaseRepo)
}

func TestPostgresBaseRepoGetByName(t *testing.T) {
	testBaseRepoGetByName(t, preparePGBaseRepo)
}

func TestPostgresBaseRepoGetByNameNotFound(t *testing.T) {
	testBaseRepoGetByNameNotFound(t, preparePGBaseRepo)
}

func TestPostgresBaseRepoUpdate(t *testing.T) {
	testBaseRepoUpdate(t, preparePGBaseRepo)
}

func TestPostgresBaseRepoList(t *testing.T) {
	testBaseRepoList(t, preparePGBaseRepo)
}

func TestPostgresUserRepoGetLatest(t *testing.T) {
	testUserRepoGetLatest(t, preparePGUserRepo)
}

func TestPostgresUserRepoGetLatestNotFound(t *testing.T) {
	testUserRepoGetLatestNotFound(t, preparePGUserRepo)
}

func TestPostgresUserRepoAppendDuplicate(t *testing.T) {
	testUserRepoAppendDuplicate(t, preparePGUserRepo)
}

func preparePGBaseRepo(t *testing.T) BaseRepo {
	db, err := sqlx.Connect("postgres", pgURI)
	if err != nil {
		t.Fatal(err)
	}

	r := NewPostgresBaseRepo(db)

	if err := r.teardown(); err != nil {
		t.Fatal(err)
	}

	return r
}

func preparePGUserRepo(t *testing.T) UserRepo {
	db, err := sqlx.Connect("postgres", pgURI)
	if err != nil {
		t.Fatal(err)
	}

	r := NewPostgresUserRepo(db)

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
