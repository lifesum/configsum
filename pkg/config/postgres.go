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
	pgUserCreateSchema = `CREATE SCHEMA IF NOT EXISTS config`
	pgUserCreateTable  = `
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

type pgUserRepo struct {
	db *sqlx.DB
}

// NewPostgresUserRepo returns a Postgres backed UserRepo implementation.
func NewPostgresUserRepo(db *sqlx.DB) (UserRepo, error) {
	return &pgUserRepo{
		db: db,
	}, nil
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
			if err := r.Setup(); err != nil {
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
			if err := r.Setup(); err != nil {
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

func (r *pgUserRepo) Setup() error {
	for _, q := range []string{
		pgUserCreateSchema,
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

func (r *pgUserRepo) Teardown() error {
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
