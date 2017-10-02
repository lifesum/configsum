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
	// Blank import for Postgres capabilities.
	_ "github.com/lib/pq"

	"github.com/lifesum/configsum/pkg/pg"
)

const (
	characterSet = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numCharacterSet = "0123456789"
)

var pgURI string
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randString(charset string) string {
	b := make([]byte, len(charset))

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func TestPostgresUserRepoGet(t *testing.T) {
	repo := preparePGUserRepo(t)

	_, err := repo.Get(randString(characterSet), randString(numCharacterSet))
	if err != nil {
		t.Fatal(err)
	}
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
