package client

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/oklog/ulid"

	"github.com/lifesum/configsum/pkg/generate"
)

func TestServiceCreate(t *testing.T) {
	var (
		name      = generate.RandomString(12)
		repo      = prepareInmemRepo(t)
		tokenRepo = prepareInmemTokenRepo(t)
		svc       = NewService(repo, tokenRepo)
	)

	clientSVC, err := svc.Create(name)
	if err != nil {
		t.Fatal(err)
	}

	clientRepo, err := repo.Lookup(clientSVC.id)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := clientSVC, clientRepo; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	token, err := tokenRepo.GetLatest(clientSVC.id)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := token.clientID, clientSVC.id; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestListWithToken(t *testing.T) {
	var (
		numClients = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(24)
		repo       = prepareInmemRepo(t)
		tokenRepo  = prepareInmemTokenRepo(t)
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
