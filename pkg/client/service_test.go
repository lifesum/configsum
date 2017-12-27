package client

import (
	"math/rand"
	"testing"
	"time"

	"github.com/oklog/ulid"

	"github.com/lifesum/configsum/pkg/generate"
)

func TestServiceCreate(t *testing.T) {
	t.Parallel()

	var (
		name      = generate.RandomString(12)
		repo      = preparePGRepo(t)
		tokenRepo = preparePGTokenRepo(t)
		svc       = NewService(repo, tokenRepo)
	)

	clientSVC, secret, err := svc.Create(name)
	if err != nil {
		t.Fatal(err)
	}

	clientRepo, err := repo.Lookup(clientSVC.id)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := clientSVC.deleted, clientRepo.deleted; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := clientSVC.id, clientRepo.id; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := clientSVC.name, clientRepo.name; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	token, err := tokenRepo.GetLatest(clientSVC.id)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := token.clientID, clientSVC.id; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := token.secret, secret; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestListWithToken(t *testing.T) {
	t.Parallel()

	var (
		numClients = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(24)
		repo       = preparePGRepo(t)
		tokenRepo  = preparePGTokenRepo(t)
		svc        = NewService(repo, tokenRepo)
	)

	for i := 0; i < numClients; i++ {
		testCreateClientWithToken(repo, tokenRepo, t)
	}

	ct, err := svc.ListWithToken()
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(ct), numClients; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestServiceLookupBySecret(t *testing.T) {
	t.Parallel()

	var (
		repo      = preparePGRepo(t)
		tokenRepo = preparePGTokenRepo(t)
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

func testCreateClientWithToken(repo Repo, tokenRepo TokenRepo, t *testing.T) {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

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
}
