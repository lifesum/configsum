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
	t.Parallel()

	testBaseRepoCreateDuplicate(t, preparePGBaseRepo)
}

func TestPostgresBaseRepoGetByID(t *testing.T) {
	t.Parallel()

	testBaseRepoGetByID(t, preparePGBaseRepo)
}

func TestPostgresBaseRepoGetByIDNotFound(t *testing.T) {
	t.Parallel()

	testBaseRepoGetByIDNotFound(t, preparePGBaseRepo)
}

func TestPostgresBaseRepoGetByName(t *testing.T) {
	t.Parallel()

	testBaseRepoGetByName(t, preparePGBaseRepo)
}

func TestPostgresBaseRepoGetByNameNotFound(t *testing.T) {
	t.Parallel()

	testBaseRepoGetByNameNotFound(t, preparePGBaseRepo)
}

func TestPostgresBaseRepoUpdate(t *testing.T) {
	t.Parallel()

	testBaseRepoUpdate(t, preparePGBaseRepo)
}

func TestPostgresBaseRepoList(t *testing.T) {
	t.Parallel()

	testBaseRepoList(t, preparePGBaseRepo)
}

func TestPostgresUserRepoGetLatest(t *testing.T) {
	t.Parallel()

	testUserRepoGetLatest(t, preparePGUserRepo)
}

func TestPostgresUserRepoGetLatestNotFound(t *testing.T) {
	t.Parallel()

	testUserRepoGetLatestNotFound(t, preparePGUserRepo)
}

func TestPostgresUserRepoAppendDuplicate(t *testing.T) {
	t.Parallel()

	testUserRepoAppendDuplicate(t, preparePGUserRepo)
}

func preparePGBaseRepo(t *testing.T) BaseRepo {
	db, err := sqlx.Connect("postgres", pgURI)
	if err != nil {
		t.Fatal(err)
	}

	r := NewPostgresBaseRepo(db, PGBaseRepoSchema(t.Name()))

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

	r := NewPostgresUserRepo(db, PGUserRepoSchema(t.Name()))

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
