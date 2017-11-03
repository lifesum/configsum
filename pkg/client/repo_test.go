package client

import (
	"math/rand"
	"testing"
	"time"

	"github.com/oklog/ulid"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/generate"
)

type prepareRepoFunc func(t *testing.T) Repo

func testRepoList(t *testing.T, p prepareRepoFunc) {
	var (
		repo = p(t)
		seed = rand.New(rand.NewSource(time.Now().UnixNano()))
		want = rand.Intn(12)
	)

	for i := 0; i < want; i++ {
		name := generate.RandomString(24)

		id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
		if err != nil {
			t.Fatal(err)
		}

		_, err = repo.Store(id.String(), name)
		if err != nil {
			t.Fatal(err)
		}
	}

	cs, err := repo.List()
	if err != nil {
		t.Fatal(err)
	}

	if have := len(cs); have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testRepoListEmpty(t *testing.T, p prepareRepoFunc) {
	repo := p(t)

	cs, err := repo.List()
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(cs), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testRepoLookup(t *testing.T, p prepareRepoFunc) {
	var (
		name = generate.RandomString(24)
		repo = p(t)
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

func testRepoLookupNotFound(t *testing.T, p prepareRepoFunc) {
	var (
		id   = generate.RandomString(24)
		repo = p(t)
	)

	_, err := repo.Lookup(id)
	if have, want := errors.Cause(err), errors.ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

type prepareTokenRepoFunc func(t *testing.T) TokenRepo

func testTokenRepoGetLatest(t *testing.T, p prepareTokenRepoFunc) {
	var (
		repo = p(t)
		seed = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	secret, err := generate.SecureToken(secretByteLen)
	if err != nil {
		t.Fatal(err)
	}

	clientID, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Store(clientID.String(), secret)
	if err != nil {
		t.Fatal(err)
	}

	token, err := repo.GetLatest(clientID.String())
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

func testTokenRepoLookup(t *testing.T, p prepareTokenRepoFunc) {
	var (
		repo = p(t)
		seed = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	secret, err := generate.SecureToken(secretByteLen)
	if err != nil {
		t.Fatal(err)
	}

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

func testTokenRepoLookupNotFound(t *testing.T, p prepareTokenRepoFunc) {
	repo := p(t)

	secret, err := generate.SecureToken(secretByteLen)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Lookup(secret)
	if have, want := errors.Cause(err), errors.ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}
