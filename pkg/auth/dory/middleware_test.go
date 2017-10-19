package dory

import (
	"context"
	"testing"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/generate"
)

func TestAuthMiddleware(t *testing.T) {
	var (
		ctx    = context.TODO()
		secret = generate.RandomString(32)
		userID = generate.RandomString(24)
	)

	signature, err := hashSignature(secret, userID)
	if err != nil {
		t.Fatal(err)
	}

	ctx = context.WithValue(ctx, contextKeySignature, signature)
	ctx = context.WithValue(ctx, contextKeyUserID, userID)

	_, err = AuthMiddleware(secret)(nopEndpoint)(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthMiddlewareSignatureMissing(t *testing.T) {
	var (
		ctx    = context.TODO()
		secret = generate.RandomString(32)
	)

	_, err := AuthMiddleware(secret)(nopEndpoint)(ctx, nil)
	if have, want := errors.Cause(err), errors.ErrSignatureMissing; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestAuthMiddlewareSignatureMissmatch(t *testing.T) {
	var (
		ctx    = context.TODO()
		secret = generate.RandomString(32)
		userID = generate.RandomString(24)
	)

	ctx = context.WithValue(ctx, contextKeySignature, generate.RandomString(12))
	ctx = context.WithValue(ctx, contextKeyUserID, userID)

	_, err := AuthMiddleware(secret)(nopEndpoint)(ctx, nil)
	if have, want := errors.Cause(err), errors.ErrSignatureMissmatch; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestAuthMiddlewareUserIDMissing(t *testing.T) {
	var (
		ctx    = context.TODO()
		secret = generate.RandomString(32)
		userID = generate.RandomString(24)
	)

	signature, err := hashSignature(secret, userID)
	if err != nil {
		t.Fatal(err)
	}

	ctx = context.WithValue(ctx, contextKeySignature, signature)

	_, err = AuthMiddleware(secret)(nopEndpoint)(ctx, nil)
	if have, want := errors.Cause(err), errors.ErrUserIDMissing; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func nopEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	return true, nil
}
