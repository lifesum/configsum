package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/pg"
)

const (
	pgCreateSchema = `CREATE SCHEMA IF NOT EXISTS config`

	pgBaseCreateTable = `
		CREATE TABLE IF NOT EXISTS config.bases(
			client_id TEXT NOT NULL,
			deleted BOOL DEFAULT FALSE,
			id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			parameters JSONB NOT NULL,
			created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME ZONE 'utc'),
			updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME ZONE 'utc')
		)`
	pgBaseDropTable = `
		DROP TABLE IF EXISTS config.bases CASCADE`
	pgBaseCreate = `
		/* pgBaseCreate */
		INSERT INTO
			config.bases(client_id, id, name, parameters)
			VALUES(:clientId, :id, :name, :parameters)`
	pgBaseGetByID = `
		/* pgBaseGetByID */
		SELECT
			client_id, deleted, id, name, parameters, created_at, updated_at
		FROM
			config.bases
		WHERE
			id = :id
		ORDER BY
			created_at DESC
		LIMIT
			1`
	pgBaseGetByName = `
		/* pgBaseGetByName */
		SELECT
			client_id, deleted, id, name, parameters, created_at, updated_at
		FROM
			config.bases
		WHERE
			client_id = :clientId
			AND name = :name
		ORDER BY
			created_at DESC
		LIMIT
			1`
	pgBaseList = `
		/* pgBaseList */
		SELECT
			client_id, deleted, id, name, parameters, created_at, updated_at
		FROM
			config.bases
		WHERE
			deleted = :deleted
		ORDER BY
			created_at DESC`
	pgBaseUpdate = `
		/* pgBaseList */
		UPDATE
			config.bases
		SET
			deleted = :deleted,
			name = :name,
			parameters = :parameters,
			updated_at = :updatedAt
		WHERE
			id = :id
	`

	pgUserCreateTable = `
		CREATE TABLE IF NOT EXISTS config.users(
			id TEXT NOT NULL PRIMARY KEY,
			user_id TEXT NOT NULL,
			base_id TEXT NOT NULL,
			rendered JSONB NOT NULL,
			rule_decisions JSONB NOT NULL,
			created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME ZONE 'utc')
		)`
	pgUserDropTable = `DROP TABLE IF EXISTS config.users CASCADE`

	pgUserInsert = `
		/* pgUserInsert*/
		INSERT INTO
			config.users(base_id, id, rendered, rule_decisions, user_id) VALUES(
			:baseId,
			:id,
			:rendered,
			:ruleDecisions,
			:userId)`
	pgUserGetLatest = `
		/* pgUserGetLatest */
		SELECT
			id, user_id, base_id, rendered, rule_decisions, created_at
		FROM
			config.users
		WHERE
			base_id = :baseId
			AND user_id = :userId
		ORDER BY
			created_at DESC
		LIMIT
			1`

	pgUserIndexGetLatest = `
		CREATE INDEX
			users_get_latest
		ON
			config.users(base_id, user_id, created_at DESC)`
)

type pgBaseRepo struct {
	db *sqlx.DB
}

// NewPostgresBaseRepo returns a Postgres backed BaseRepo implementation.
func NewPostgresBaseRepo(db *sqlx.DB) BaseRepo {
	return &pgBaseRepo{db: db}
}

func (r *pgBaseRepo) Create(
	id, clientID, name string,
	parameters rendered,
) (BaseConfig, error) {
	rawParameters, err := json.Marshal(parameters)
	if err != nil {
		return BaseConfig{}, errors.Wrap(err, "marshal parameters")
	}

	_, err = r.db.NamedExec(pgBaseCreate, map[string]interface{}{
		"clientId":   clientID,
		"id":         id,
		"name":       name,
		"parameters": rawParameters,
	})
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrDuplicateKey:
			return BaseConfig{}, errors.Wrap(errors.ErrExists, "base config")
		case pg.ErrRelationNotFound:
			if err := r.setup(); err != nil {
				return BaseConfig{}, err
			}

			return r.Create(id, clientID, name, parameters)
		default:
			return BaseConfig{}, fmt.Errorf("named exec: %s", err)
		}
	}

	return BaseConfig{
		ClientID:   clientID,
		ID:         id,
		Name:       name,
		Parameters: parameters,
		CreatedAt:  time.Now(),
	}, nil
}

func (r *pgBaseRepo) GetByID(id string) (BaseConfig, error) {
	query, args, err := r.db.BindNamed(pgBaseGetByID, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return BaseConfig{}, errors.Wrap(err, "named query")
	}

	raw := struct {
		ClientID   string    `db:"client_id"`
		Deleted    bool      `db:"deleted"`
		ID         string    `db:"id"`
		Name       string    `db:"name"`
		Parameters []byte    `db:"parameters"`
		CreatedAt  time.Time `db:"created_at"`
		UpdatedAt  time.Time `db:"updated_at"`
	}{}

	err = r.db.Get(&raw, query, args...)
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.setup(); err != nil {
				return BaseConfig{}, err
			}

			return r.GetByID(id)

		case sql.ErrNoRows:
			return BaseConfig{}, errors.Wrap(errors.ErrNotFound, "get base config by id")

		default:
			return BaseConfig{}, errors.Wrap(err, "get base config by id")
		}
	}

	params := rendered{}

	if err := json.Unmarshal(raw.Parameters, &params); err != nil {
		return BaseConfig{}, errors.Wrap(err, "unmarshal parameters")
	}

	return BaseConfig{
		ClientID:   raw.ClientID,
		ID:         raw.ID,
		Name:       raw.Name,
		Parameters: params,
		CreatedAt:  raw.CreatedAt,
	}, nil
}

func (r *pgBaseRepo) GetByName(clientID, name string) (BaseConfig, error) {
	query, args, err := r.db.BindNamed(pgBaseGetByName, map[string]interface{}{
		"clientId": clientID,
		"name":     name,
	})
	if err != nil {
		return BaseConfig{}, errors.Wrap(err, "named query")
	}

	raw := struct {
		ClientID   string    `db:"client_id"`
		Deleted    bool      `db:"deleted"`
		ID         string    `db:"id"`
		Name       string    `db:"name"`
		Parameters []byte    `db:"parameters"`
		CreatedAt  time.Time `db:"created_at"`
		UpdatedAt  time.Time `db:"updated_at"`
	}{}

	err = r.db.Get(&raw, query, args...)
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.setup(); err != nil {
				return BaseConfig{}, err
			}

			return r.GetByName(clientID, name)

		case sql.ErrNoRows:
			return BaseConfig{}, errors.Wrap(errors.ErrNotFound, "get base config by id")

		default:
			return BaseConfig{}, errors.Wrap(err, "get base config by id")
		}
	}

	params := rendered{}

	if err := json.Unmarshal(raw.Parameters, &params); err != nil {
		return BaseConfig{}, errors.Wrap(err, "unmarshal parameters")
	}

	return BaseConfig{
		ClientID:   raw.ClientID,
		ID:         raw.ID,
		Name:       raw.Name,
		Parameters: params,
		CreatedAt:  raw.CreatedAt,
	}, nil
}

func (r *pgBaseRepo) List() (BaseList, error) {
	rows, err := r.db.NamedQuery(pgBaseList, map[string]interface{}{
		"deleted": false,
	})
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.setup(); err != nil {
				return nil, err
			}

			return r.List()
		default:
			return nil, errors.Wrap(err, "baseRepo List")
		}
	}

	cs := BaseList{}

	for rows.Next() {
		var (
			c         = BaseConfig{}
			rawParams = []byte{}
		)

		// client_id, deleted, id, name, parameters, created_at, updated_at
		err := rows.Scan(
			&c.ClientID,
			&c.Deleted,
			&c.ID,
			&c.Name,
			&rawParams,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "List scan")
		}

		if err := json.Unmarshal(rawParams, &c.Parameters); err != nil {
			return nil, errors.Wrap(err, "unmarshal parameters")
		}

		cs = append(cs, c)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "List rows")
	}

	return cs, nil
}

func (r *pgBaseRepo) Update(c BaseConfig) (BaseConfig, error) {
	rawParameters, err := json.Marshal(c.Parameters)
	if err != nil {
		return BaseConfig{}, errors.Wrap(err, "marshal parameters")
	}

	updatedAt := time.Now()

	res, err := r.db.NamedExec(pgBaseUpdate, map[string]interface{}{
		"id":         c.ID,
		"deleted":    c.Deleted,
		"name":       c.Name,
		"parameters": rawParameters,
		"updatedAt":  updatedAt,
	})
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.setup(); err != nil {
				return BaseConfig{}, err
			}

			return r.Update(c)
		default:
			return BaseConfig{}, fmt.Errorf("named exec: %s", err)
		}
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return BaseConfig{}, err
	}

	if rows == 0 {
		return BaseConfig{}, errors.Wrapf(errors.ErrNotFound, "id '%s'", c.ID)
	}

	return BaseConfig{
		ClientID:   c.ClientID,
		Deleted:    c.Deleted,
		ID:         c.ID,
		Name:       c.Name,
		Parameters: c.Parameters,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  updatedAt,
	}, nil
}

func (r *pgBaseRepo) setup() error {
	for _, q := range []string{
		pgCreateSchema,
		pgBaseCreateTable,
	} {
		_, err := r.db.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *pgBaseRepo) teardown() error {
	for _, q := range []string{
		pgBaseDropTable,
	} {
		_, err := r.db.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}

type pgUserRepo struct {
	db *sqlx.DB
}

// NewPostgresUserRepo returns a Postgres backed UserRepo implementation.
func NewPostgresUserRepo(db *sqlx.DB) UserRepo {
	return &pgUserRepo{db: db}
}

func (r *pgUserRepo) Append(
	id, baseID, userID string,
	decisions ruleDecisions,
	render rendered,
) (UserConfig, error) {
	rawRendered, err := json.Marshal(render)
	if err != nil {
		return UserConfig{}, errors.Wrap(err, "marshal rendered")
	}

	rawDecisions, err := json.Marshal(decisions)
	if err != nil {
		return UserConfig{}, errors.Wrap(err, "marshal decisions")
	}

	_, err = r.db.NamedExec(pgUserInsert, map[string]interface{}{
		"baseId":        baseID,
		"id":            id,
		"rendered":      rawRendered,
		"ruleDecisions": rawDecisions,
		"userId":        userID,
	})
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrDuplicateKey:
			return UserConfig{}, errors.Wrap(errors.ErrExists, "user config")
		case pg.ErrRelationNotFound:
			if err := r.setup(); err != nil {
				return UserConfig{}, err
			}

			return r.Append(id, baseID, userID, decisions, render)
		default:
			return UserConfig{}, fmt.Errorf("named exec: %s", err)
		}
	}

	return UserConfig{
		baseID:    baseID,
		id:        id,
		userID:    userID,
		rendered:  render,
		createdAt: time.Now(),
	}, nil
}

func (r *pgUserRepo) GetLatest(baseID, userID string) (UserConfig, error) {
	query, args, err := r.db.BindNamed(pgUserGetLatest, map[string]interface{}{
		"baseId": baseID,
		"userId": userID,
	})
	if err != nil {
		return UserConfig{}, fmt.Errorf("named query: %s", err)
	}

	raw := struct {
		BaseID        string    `db:"base_id"`
		ID            string    `db:"id"`
		Rendered      []byte    `db:"rendered"`
		RuleDecisions []byte    `db:"rule_decisions"`
		UserID        string    `db:"user_id"`
		CreatedAt     time.Time `db:"created_at"`
	}{}

	err = r.db.Get(&raw, query, args...)
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.setup(); err != nil {
				return UserConfig{}, err
			}

			return r.GetLatest(baseID, userID)
		case sql.ErrNoRows:
			return UserConfig{}, errors.Wrap(errors.ErrNotFound, "get user config")

		default:
			return UserConfig{}, fmt.Errorf("get: %s", err)
		}
	}

	render := rendered{}

	if err := json.Unmarshal(raw.Rendered, &render); err != nil {
		return UserConfig{}, errors.Wrap(err, "unmarshal rendered")
	}

	decisions := ruleDecisions{}

	if err := json.Unmarshal(raw.RuleDecisions, &decisions); err != nil {
		return UserConfig{}, errors.Wrap(err, "unmarshal decisons")
	}

	return UserConfig{
		baseID:        raw.BaseID,
		id:            raw.ID,
		rendered:      render,
		ruleDecisions: decisions,
		userID:        raw.UserID,
		createdAt:     raw.CreatedAt,
	}, nil
}

func (r *pgUserRepo) setup() error {
	for _, q := range []string{
		pgCreateSchema,
		pgUserCreateTable,
		pgUserIndexGetLatest,
	} {
		_, err := r.db.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *pgUserRepo) teardown() error {
	for _, q := range []string{
		pgUserDropTable,
	} {
		_, err := r.db.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}
