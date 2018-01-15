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

func TestPostgresRepoGetByIDNotFound(t *testing.T) {
	t.Parallel()

	testRepoGetByIDNotFound(t, preparePGRepo)
}

func TestPostgresRepoCreateDuplicate(t *testing.T) {
	t.Parallel()

	testRepoCreateDuplicate(t, preparePGRepo)
}

func TestPostgresRepoGet(t *testing.T) {
	t.Parallel()

	testRepoGet(t, preparePGRepo)
}

func TestPostgresRepoUpdateWith(t *testing.T) {
	t.Parallel()

	testRepoUpdateWith(t, preparePGRepo)
}

func TestPostgresRepoListAll(t *testing.T) {
	t.Parallel()

	testRepoListAll(t, preparePGRepo)
}

func TestPostgresRepoListAllEmpty(t *testing.T) {
	t.Parallel()

	testRepoListAllEmpty(t, preparePGRepo)
}

func TestPostgresRepoListDeleted(t *testing.T) {
	t.Parallel()

	testRepoListDeleted(t, preparePGRepo)
}

func TestPostgresRepoListActive(t *testing.T) {
	t.Parallel()

	testRepoListActive(t, preparePGRepo)
}

func TestPostgresRepoListActiveEmpty(t *testing.T) {
	t.Parallel()

	testRepoListActiveEmpty(t, preparePGRepo)
}

func TestPostgresRepoCreateRollout(t *testing.T) {
	t.Parallel()

	testRepoCreateRollout(t, preparePGRepo)
}

func preparePGRepo(t *testing.T) Repo {
	db, err := sqlx.Connect("postgres", pgURI)
	if err != nil {
		t.Fatal(err)
	}

	r := NewPostgresRepo(db, PGRepoSchema(t.Name()))

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
