package client

import (
	"math/rand"
	"testing"
	"time"

	"github.com/oklog/ulid"

	"github.com/lifesum/configsum/pkg/generate"
)

func TestServiceLookupBySecret(t *testing.T) {
	var (
		repo      = prepareInmemRepo(t)
		tokenRepo = prepareInmemTokenRepo(t)
		seed      = rand.New(rand.NewSource(time.Now().UnixNano()))
		svc       = NewService(repo, tokenRepo)
	)

	clientID, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	secret, err := generate.SecureToken(secretByteLen)
	if err != nil {
		t.Fatal(err)
	}

	c, err := repo.Store(clientID.String(), generate.RandomString(12))
	if err != nil {
		t.Fatal(err)
	}

	_, err = tokenRepo.Store(c.id, secret)
	if err != nil {
		t.Fatal(err)
	}

	c, err = svc.LookupBySecret(secret)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := c.id, clientID; err != nil {
		t.Errorf("have %v, want %v", have, want)
	}
}
