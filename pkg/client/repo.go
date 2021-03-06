package client

import "time"

// Client represents distinct consumers like mobile apps, SPAs or other web
// servers.
type Client struct {
	deleted   bool
	id        string
	name      string
	createdAt time.Time
}

// List is a collection of Clients.
type List []Client

// Repo for Client interactions.
type Repo interface {
	lifecycle

	List() (List, error)
	Lookup(id string) (Client, error)
	Store(id, name string) (Client, error)
}

// RepoMiddleware is a chainable behaviour modifier for Repo.
type RepoMiddleware func(Repo) Repo

// TokenRepo for Token interactions.
type TokenRepo interface {
	lifecycle

	GetLatest(clientID string) (Token, error)
	Lookup(secret string) (Token, error)
	Store(clientID, secret string) (Token, error)
}

// TokenRepoMiddleware is a chainable behaviour modifier for TokenRepo.
type TokenRepoMiddleware func(next TokenRepo) TokenRepo

// Token is the relation between a Client secret and id.
type Token struct {
	clientID  string
	deleted   bool
	secret    string
	createdAt time.Time
}

type lifecycle interface {
	Setup() error
	Teardown() error
}
