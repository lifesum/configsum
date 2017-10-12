package client

import (
	"context"
	"testing"

	"github.com/pkg/errors"

	"github.com/lifesum/configsum/pkg/generate"
)

func TestAuthMiddleware(t *testing.T) {
	var (
		repo      = prepareInmemRepo(t)
		tokenRepo = prepareInmemTokenRepo(t)
		svc       = NewService(repo, tokenRepo)

		clientID   = generate.RandomString(24)
		clientName = generate.RandomString(12)
		secret     = generate.RandomString(32)
	)

	_, err := repo.Store(clientID, clientName)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tokenRepo.Store(clientID, secret)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()
	ctx = context.WithValue(ctx, contextKeySecret, secret)

	_, err = AuthMiddleware(svc)(nopEndpoint)(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthMiddlewareSecretMissing(t *testing.T) {
	var (
		repo      Repo
		tokenRepo TokenRepo

		svc = NewService(repo, tokenRepo)
	)

	ctx := context.TODO()

	_, err := AuthMiddleware(svc)(nopEndpoint)(ctx, nil)
	if have, want := errors.Cause(err), ErrSecretMissing; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func nopEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	return true, nil
}
