package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/lifesum/configsum/pkg/pg"
)

const (
	pgCreateSchema = `CREATE SCHEMA IF NOT EXISTS config`
	pgCreateTable  = `
		CREATE TABLE IF NOT EXISTS config.users(
			id TEXT NOT NULL PRIMARY KEY,
			user_id TEXT NOT NULL,
			base_id TEXT NOT NULL,
			rendered JSONB NOT NULL,
			created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME ZONE 'utc'),
			activated_at TIMESTAMP WITHOUT TIME ZONE
		)`
	pgDropTable = `DROP TABLE IF EXISTS config.users CASCADE`

	pgInsertUser = `
		INSERT INTO
			config.users(base_id, id, rendered, user_id) VALUES(
			:base_id,
			:id,
			:rendered,
			:user_id)`
	pgSelectUsers = `
		SELECT
			id, user_id, base_id, rendered, created_at
		FROM
			config.users
		LIMIT
			:limit`
)

type pgUserRepo struct {
	db *sqlx.DB
}

// NewPostgresUserRepo returns a Postgres backed UserRepo implementation.
func NewPostgresUserRepo(db *sqlx.DB) (UserRepo, error) {
	return &pgUserRepo{
		db: db,
	}, nil
}

func (r *pgUserRepo) Get(baseID, userID string) (UserConfig, error) {
	query, args, err := r.db.BindNamed(pgSelectUsers, map[string]interface{}{
		"limit": 1,
	})
	if err != nil {
		return UserConfig{}, fmt.Errorf("named query: %s", err)
	}

	raw := struct {
		BaseID    string    `db:"base_id"`
		ID        string    `db:"id"`
		Rendered  []byte    `db:"rendered"`
		UserID    string    `db:"user_id"`
		CreatedAt time.Time `db:"created_at"`
	}{}

	err = r.db.Get(&raw, query, args...)
	if err != nil {
		if pg.IsRelationNotFound(pg.Wrap(err)) {
			if err := r.Setup(); err != nil {
				return UserConfig{}, err
			}

			return r.Get(baseID, userID)
		}

		if err == sql.ErrNoRows {
			return UserConfig{}, errors.Wrap(ErrNotFound, "get user config")
		}

		return UserConfig{}, fmt.Errorf("get: %s", err)
	}

	render := rendered{}

	fmt.Println(string(raw.Rendered))
	err = json.Unmarshal(raw.Rendered, &render)
	if err != nil {
		return UserConfig{}, fmt.Errorf("rendered unmarshal: %s", err)
	}

	return UserConfig{
		baseID:    raw.BaseID,
		id:        raw.ID,
		rendered:  render,
		userID:    raw.UserID,
		createdAt: raw.CreatedAt,
	}, nil
}

func (r *pgUserRepo) Put(
	id, baseID, userID string,
	render rendered,
) (UserConfig, error) {
	raw, err := json.Marshal(render)
	if err != nil {
		return UserConfig{}, fmt.Errorf("marashl rendered: %s", err)
	}

	_, err = r.db.NamedExec(pgInsertUser, map[string]interface{}{
		"base_id":  baseID,
		"id":       id,
		"rendered": raw,
		"user_id":  userID,
	})
	if err != nil {
		return UserConfig{}, fmt.Errorf("named exec: %s", err)
	}

	return UserConfig{
		baseID:    baseID,
		id:        id,
		userID:    userID,
		rendered:  render,
		createdAt: time.Now(),
	}, nil
}

func (r *pgUserRepo) Setup() error {
	for _, q := range []string{
		pgCreateSchema,
		pgCreateTable,
	} {
		_, err := r.db.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *pgUserRepo) Teardown() error {
	for _, q := range []string{
		pgDropTable,
	} {
		_, err := r.db.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}
