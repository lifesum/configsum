package config

import (
	"database/sql"
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
			rule_ids TEXT[],
			rendered JSONB NOT NULL,
			created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME ZONE 'utc'),
			activated_at TIMESTAMP WITHOUT TIME ZONE
		)`
	pgDropTable = `DROP TABLE IF EXISTS config.users CASCADE`

	pgSelectUsers = `
		SELECT
			id, user_id, base_id, rule_ids, rendered, created_at, activated_at
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

func (r *pgUserRepo) Get(baseName, id string) (UserConfig, error) {
	query, args, err := r.db.BindNamed(pgSelectUsers, map[string]interface{}{
		"limit": 1,
	})
	if err != nil {
		return UserConfig{}, fmt.Errorf("named query: %s", err)
	}

	raw := struct {
		BaseID      string                 `db:"base_id"`
		ID          string                 `db:"id"`
		Rendered    map[string]interface{} `db:"rendered"`
		RuleIDs     []string               `db:"rule_ids"`
		UserID      string                 `db:"user_id"`
		CreatedAt   time.Time              `db:"created_at"`
		ActivatedAt time.Time              `db:"activated_at"`
	}{}

	err = r.db.Get(&raw, query, args...)
	if err != nil {
		if pg.IsRelationNotFound(pg.Wrap(err)) {
			if err := r.Setup(); err != nil {
				return UserConfig{}, err
			}

			return r.Get(baseName, id)
		}

		if err == sql.ErrNoRows {
			return UserConfig{}, errors.Wrap(ErrNotFound, "get user config")
		}

		return UserConfig{}, fmt.Errorf("get: %s", err)
	}

	return UserConfig{
		baseID: raw.BaseID,
		id:     raw.ID,
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
