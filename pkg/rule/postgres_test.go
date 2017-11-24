package rule

import (
	"flag"
	"fmt"
	"os/user"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/lifesum/configsum/pkg/pg"
)

var pgURI string

func TestPostgresRepoGetNotFound(t *testing.T) {
	testRepoGetNotFound(t, preparePGRepo)
}

func TestPostgresRepoCreateDuplicate(t *testing.T) {
	testRepoCreateDuplicate(t, preparePGRepo)
}

func TestPostgresRepoGet(t *testing.T) {
	testRepoGet(t, preparePGRepo)
}

func TestPostgresRepoUpdateWith(t *testing.T) {
	testRepoUpdateWith(t, preparePGRepo)
}

func TestPostgresRepoListAll(t *testing.T) {
	testRepoListAll(t, preparePGRepo)
}

func TestPostgresRepoListActive(t *testing.T) {
	testRepoListActive(t, preparePGRepo)
}

func TestPostgresRepoListActiveEmpty(t *testing.T) {
	testRepoListActiveEmpty(t, preparePGRepo)
}

func TestPostgresRepoListAllEmpty(t *testing.T) {
	testRepoListAllEmpty(t, preparePGRepo)
}

func TestPostgresRepoListDeleted(t *testing.T) {
	testRepoListDeleted(t, preparePGRepo)
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

func init() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	uri := flag.String("postgres.uri", fmt.Sprintf(pg.DefaultTestURI, u.Username), "Postgres connection URL")

	flag.Parse()

	pgURI = *uri
}
