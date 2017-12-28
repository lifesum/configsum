package client

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/pg"
)

const (
	pgDefaultSchema = "client"

	pgCreateSchema = `CREATE SCHEMA IF NOT EXISTS %s`

	pgClientCreateTable = `
		CREATE TABLE IF NOT EXISTS %s.clients(
			id TEXT NOT NULL PRIMARY KEY,
			deleted BOOL DEFAULT FALSE,
			name TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME ZONE 'utc')
		)`
	pgClientDropTable = `DROP TABLE IF EXISTS %s.clients CASCADE`
	pgClientList      = `
		/* pgClientList */
		SELECT
			created_at, deleted, id, name
		FROM
			%s.clients
		WHERE
			deleted = :deleted
		ORDER BY
			created_at DESC`
	pgClientLookup = `
		/* pgClientLookup */
		SELECT
			deleted, id, name, created_at
		FROM
			%s.clients
		WHERE
			deleted = :deleted
			AND id = :id
		LIMIT
			1`
	pgClientStore = `
		/* pgClientStore */
		INSERT INTO
			%s.clients(
				deleted,
				id,
				name
			)
			VALUES(
				:deleted,
				:id,
				:name
			)`

	pgTokenCreateTable = `
		CREATE TABLE IF NOT EXISTS %s.tokens(
			secret TEXT NOT NULL PRIMARY KEY,
			deleted BOOL DEFAULT FALSE,
			client_id TEXT NOT NULL,
			created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME ZONE 'utc')
		)`
	pgTokenDropTable = `DROP TABLE IF EXISTS %s.tokens CASCADE`
	pgTokenGetLatest = `
		/* pgTokenGetLatest */
		SELECT
			client_id, deleted, secret, created_at
		FROM
			%s.tokens
		WHERE
			client_id = :id
			AND deleted = :deleted
		ORDER BY
			created_at DESC
		LIMIT
			1`
	pgTokenLookup = `
		/* pgTokenLookup */
		SELECT
			client_id, deleted, secret, created_at
		FROM
			%s.tokens
		WHERE
			deleted = :deleted
			AND secret = :secret
		LIMIT
			1`
	pgTokenStore = `
		/* pgClientStore */
		INSERT INTO
			%s.tokens(
				client_id,
				secret
			)
			VALUES(
				:clientId,
				:secret
			)`
)

// PGRepoOption sets an optiomal parameter for the repo.
type PGRepoOption func(*PGRepo)

// PGRepoSchema sets the namespacing of the Postgres tables to a non-default
// schema.
func PGRepoSchema(schema string) PGRepoOption {
	return func(r *PGRepo) { r.schema = schema }
}

// PGRepo is Postgres backed Repo implementation.
type PGRepo struct {
	db     *sqlx.DB
	schema string
}

// NewPostgresRepo returns a Postgres backed Repo implementation.
func NewPostgresRepo(db *sqlx.DB, options ...PGRepoOption) *PGRepo {
	r := &PGRepo{
		db:     db,
		schema: pgDefaultSchema,
	}

	for _, option := range options {
		option(r)
	}

	return r
}

// List returns all clients.
func (r *PGRepo) List() (List, error) {
	rows, err := r.db.NamedQuery(
		r.prefixSchema(pgClientList),
		map[string]interface{}{
			"deleted": false,
		},
	)
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.Setup(); err != nil {
				return nil, err
			}

			return r.List()
		default:
			return nil, errors.Wrap(err, "List query")
		}
	}

	cs := List{}

	for rows.Next() {
		c := Client{}

		err := rows.Scan(
			&c.createdAt,
			&c.deleted,
			&c.id,
			&c.name,
		)
		if err != nil {
			return nil, errors.Wrap(err, "List scan")
		}

		cs = append(cs, c)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "List rows")
	}

	return cs, nil
}

// Lookup returns the client stored for the given id.
func (r *PGRepo) Lookup(id string) (Client, error) {
	query, args, err := r.db.BindNamed(
		r.prefixSchema(pgClientLookup),
		map[string]interface{}{
			"deleted": false,
			"id":      id,
		},
	)
	if err != nil {
		return Client{}, errors.Wrap(err, "bind named")
	}

	raw := struct {
		Deleted   bool      `db:"deleted"`
		ID        string    `db:"id"`
		Name      string    `db:"name"`
		CreatedAt time.Time `db:"created_at"`
	}{}

	err = r.db.Get(&raw, query, args...)
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.Setup(); err != nil {
				return Client{}, err
			}

			return r.Lookup(id)
		case sql.ErrNoRows:
			return Client{}, errors.Wrap(errors.ErrNotFound, "client lookup")
		default:
			return Client{}, errors.Wrap(err, "client lookup")
		}
	}

	return Client{
		deleted:   raw.Deleted,
		id:        raw.ID,
		name:      raw.Name,
		createdAt: raw.CreatedAt.UTC(),
	}, nil
}

// Store persists a new client with the given id and name.
func (r *PGRepo) Store(id, name string) (Client, error) {
	_, err := r.db.NamedExec(
		r.prefixSchema(pgClientStore),
		map[string]interface{}{
			"deleted": false,
			"id":      id,
			"name":    name,
		},
	)
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrDuplicateKey:
			return Client{}, errors.Wrap(errors.ErrExists, "client")
		case pg.ErrRelationNotFound:
			if err := r.Setup(); err != nil {
				return Client{}, err
			}

			return r.Store(id, name)
		default:
			return Client{}, errors.Wrap(err, "named exec")
		}
	}

	return Client{
		id:        id,
		name:      name,
		createdAt: time.Now().UTC(),
	}, nil
}

// Setup prepares the PGRepo for operation.
func (r *PGRepo) Setup() error {
	for _, q := range []string{
		r.prefixSchema(pgCreateSchema),
		r.prefixSchema(pgClientCreateTable),
	} {
		_, err := r.db.Exec(q)
		if err != nil {
			return errors.Wrap(err, "pgRepo.setup()")
		}
	}

	return nil
}

// Teardown deconstructs all dependencies of the repo.
func (r *PGRepo) Teardown() error {
	for _, q := range []string{
		r.prefixSchema(pgClientDropTable),
	} {
		_, err := r.db.Exec(q)
		if err != nil {
			return errors.Wrap(err, "pgRepo.teardowm()")
		}
	}

	return nil
}

func (r *PGRepo) prefixSchema(query string) string {
	return fmt.Sprintf(query, r.schema)
}

// PGTokenRepoOption sets an optiomal parameter for the token repo.
type PGTokenRepoOption func(*PGTokenRepo)

// PGTokenRepoSchema sets the namespacing of the Postgres tables to a
// non-default schema.
func PGTokenRepoSchema(schema string) PGTokenRepoOption {
	return func(r *PGTokenRepo) { r.schema = schema }
}

// PGTokenRepo is a Postgres backed TokenRepo implementation.
type PGTokenRepo struct {
	db     *sqlx.DB
	schema string
}

// NewPostgresTokenRepo returns a Postgres backed TokenRepo implementation.
func NewPostgresTokenRepo(db *sqlx.DB, options ...PGTokenRepoOption) *PGTokenRepo {
	r := &PGTokenRepo{
		db:     db,
		schema: pgDefaultSchema,
	}

	for _, option := range options {
		option(r)
	}

	return r
}

// GetLatest returns the newest token for the given client id.
func (r *PGTokenRepo) GetLatest(clientID string) (Token, error) {
	query, args, err := r.db.BindNamed(
		r.prefixSchema(pgTokenGetLatest),
		map[string]interface{}{
			"id":      clientID,
			"deleted": false,
		},
	)
	if err != nil {
		return Token{}, errors.Wrap(err, "bind named")
	}

	raw := struct {
		ClientID  string    `db:"client_id"`
		Deleted   bool      `db:"deleted"`
		Secret    string    `db:"secret"`
		CreatedAt time.Time `db:"created_at"`
	}{}

	err = r.db.Get(&raw, query, args...)
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.Setup(); err != nil {
				return Token{}, err
			}

			return r.GetLatest(clientID)
		case sql.ErrNoRows:
			return Token{}, errors.Wrap(errors.ErrNotFound, "token lookup")
		default:
			return Token{}, errors.Wrap(err, "token lookup")
		}
	}

	return Token{
		clientID:  raw.ClientID,
		deleted:   raw.Deleted,
		secret:    raw.Secret,
		createdAt: raw.CreatedAt,
	}, nil
}

// Lookup given a secret returns the associated token.
func (r *PGTokenRepo) Lookup(secret string) (Token, error) {
	query, args, err := r.db.BindNamed(
		r.prefixSchema(pgTokenLookup),
		map[string]interface{}{
			"deleted": false,
			"secret":  secret,
		},
	)
	if err != nil {
		return Token{}, errors.Wrap(err, "bind named")
	}

	raw := struct {
		ClientID  string    `db:"client_id"`
		Deleted   bool      `db:"deleted"`
		Secret    string    `db:"secret"`
		CreatedAt time.Time `db:"created_at"`
	}{}

	err = r.db.Get(&raw, query, args...)
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.Setup(); err != nil {
				return Token{}, err
			}

			return r.Lookup(secret)
		case sql.ErrNoRows:
			return Token{}, errors.Wrap(errors.ErrNotFound, "token lookup")
		default:
			return Token{}, errors.Wrap(err, "token lookup")
		}
	}

	return Token{
		clientID:  raw.ClientID,
		deleted:   raw.Deleted,
		secret:    raw.Secret,
		createdAt: raw.CreatedAt,
	}, nil
}

// Store persists a new token with the given client id and secret.
func (r *PGTokenRepo) Store(clientID, secret string) (Token, error) {
	_, err := r.db.NamedExec(
		r.prefixSchema(pgTokenStore),
		map[string]interface{}{
			"clientId": clientID,
			"secret":   secret,
		},
	)
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.Setup(); err != nil {
				return Token{}, err
			}

			return r.Store(clientID, secret)
		default:
			return Token{}, errors.Wrap(err, "named exec")
		}
	}

	return Token{
		clientID:  clientID,
		deleted:   false,
		secret:    secret,
		createdAt: time.Now(),
	}, nil
}

// Setup prepares all dependencies for the Postgres repo.
func (r *PGTokenRepo) Setup() error {
	for _, q := range []string{
		r.prefixSchema(pgCreateSchema),
		r.prefixSchema(pgTokenCreateTable),
	} {
		_, err := r.db.Exec(q)
		if err != nil {
			return errors.Wrap(err, "pgTokenRepo.setup()")
		}
	}

	return nil
}

// Teardown removes all dependencies of the Postgres repo.
func (r *PGTokenRepo) Teardown() error {
	for _, q := range []string{
		r.prefixSchema(pgTokenDropTable),
	} {
		_, err := r.db.Exec(q)
		if err != nil {
			return errors.Wrap(err, "pgTokenRepo.teardown()")
		}
	}

	return nil
}

func (r *PGTokenRepo) prefixSchema(query string) string {
	return fmt.Sprintf(query, r.schema)
}
