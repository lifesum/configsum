package client

import (
	"time"

	"github.com/lifesum/configsum/pkg/errors"
)

type inmemRepo struct {
	clients map[string]Client
}

// NewInmemRepo returns a memory backed Repo implementation.
func NewInmemRepo() Repo {
	return &inmemRepo{
		clients: map[string]Client{},
	}
}

func (r *inmemRepo) Lookup(id string) (Client, error) {
	c, ok := r.clients[id]
	if !ok {
		return Client{}, errors.Wrap(errors.ErrNotFound, "client lookup")
	}

	return c, nil
}

func (r *inmemRepo) Store(id, name string) (Client, error) {
	c := Client{
		id:        id,
		name:      name,
		createdAt: time.Now(),
	}

	r.clients[id] = c

	return c, nil
}

func (r *inmemRepo) setup() error {
	return nil
}

func (r *inmemRepo) teardown() error {
	return nil
}

type inmemTokenRepo struct {
	tokens map[string]Token
}

// NewInmemTokenRepo returns a memory backed Repo implementation.
func NewInmemTokenRepo() TokenRepo {
	return &inmemTokenRepo{
		tokens: map[string]Token{},
	}
}

func (r *inmemTokenRepo) Lookup(secret string) (Token, error) {
	t, ok := r.tokens[secret]
	if !ok {
		return Token{}, errors.Wrap(errors.ErrNotFound, "token lookup")
	}

	return t, nil
}

func (r *inmemTokenRepo) Store(clientID, secret string) (Token, error) {
	t := Token{
		clientID:  clientID,
		secret:    secret,
		createdAt: time.Now(),
	}

	r.tokens[secret] = t

	return t, nil
}

func (r *inmemTokenRepo) setup() error {
	return nil
}

func (r *inmemTokenRepo) teardown() error {
	return nil
}
