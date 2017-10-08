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

// Repo for Client interactions.
type Repo interface {
	lifecycle

	Lookup(id string) (Client, error)
	Store(id, name string) (Client, error)
}

// TokenRepo for Token interactions.
type TokenRepo interface {
	lifecycle

	Lookup(secret string) (Token, error)
	Store(clientID, secret string) (Token, error)
}

// Token is the relation between a Client secret and id.
type Token struct {
	clientID  string
	deleted   bool
	secret    string
	createdAt time.Time
}

type lifecycle interface {
	setup() error
	teardown() error
}
