package rule

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/pg"
)

const (
	pgRuleCreateSchema = `CREATE SCHEMA IF NOT EXISTS rule`
	pgRuleCreateTable  = `
		CREATE TABLE IF NOT EXISTS rule.rules(
			id TEXT NOT NULL PRIMARY KEY,
			active BOOLEAN NOT NULL DEFAULT FALSE,
			buckets JSONB NOT NULL,
			config_id TEXT NOT NULL,	
			criteria JSONB NOT NULL,
			description TEXT NOT NULL,
			deleted BOOLEAN NOT NULL DEFAULT FALSE,
			kind INT8 NOT NULL,
			name TEXT NOT NULL,
			activated_at TIMESTAMP WITHOUT TIME ZONE,
			created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() AT TIME ZONE 'utc'),
			end_time TIMESTAMP WITHOUT TIME ZONE,
			start_time TIMESTAMP WITHOUT TIME ZONE,
			updated_at TIMESTAMP WITHOUT TIME ZONE
		)`
	pgRuleDropTable = `DROP TABLE IF EXISTS rule.rules CASCADE`

	pgRuleInsert = `
		INSERT INTO
		rule.rules(
			id, 
			active, 
			buckets, 
			config_id, 
			created_at, 
			criteria, 
			description,
			end_time,
			kind,
			name,
			start_time, 
			updated_at) 
			VALUES(
				:id,
				:active,
				:buckets,
				:configId,
				:createdAt,
				:criteria,
				:description,
				:endTime,
				:kind,
				:name,
				:startTime,
				:updatedAt
		)`

	pgRuleGetByName = `
		SELECT
			id, 
			active, 
			activated_at, 
			buckets, 
			config_id, 
			created_at, 
			criteria, 
			description, 
			deleted,
			end_time, 
			kind, 
			name, 
			start_time, 
			updated_at
		FROM
			rule.rules
		WHERE
			config_id = :configId	
			AND name = :name
			AND deleted = false
		ORDER BY
			created_at DESC
		LIMIT
			1`

	pgRuleUpdate = `
		UPDATE rule.rules
		SET
			active = :active,
			activated_at = :activatedAt,
			buckets = :buckets,
			criteria = :criteria,
			description = :description,
			deleted = :deleted,
			end_time = :endTime,
			kind = :kind,
			name = :name,
			start_time = :startTime,
			updated_at = :updatedAt
		WHERE
			config_id = :configId
			AND name = :name
	`

	pgRuleListAll = `
		SELECT 		
			id, 
			active, 
			activated_at, 
			buckets, 
			config_id, 
			created_at, 
			criteria, 
			description, 
			end_time, 
			kind, 
			name, 
			start_time, 
			updated_at
		FROM
			rule.rules
		WHERE
			deleted = false
	`

	pgRuleListActive = `
		SELECT	
			id, 
			active, 
			activated_at, 
			buckets, 
			config_id, 
			created_at, 
			criteria, 
			description, 
			end_time, 
			kind, 
			name, 
			start_time, 
			updated_at
		FROM
			rule.rules
		WHERE
			active = true
			AND start_time <= $1
			AND end_time > $1
			AND deleted = false
	`
)

type pgRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo returns a Postgres backed Repo implementation.
func NewPostgresRepo(db *sqlx.DB) Repo {
	return &pgRepo{
		db: db,
	}
}

func (r *pgRepo) Create(input Rule) (Rule, error) {
	rawBuckets, err := json.Marshal(input.buckets)
	if err != nil {
		return Rule{}, errors.Wrap(err, "marshal buckets")
	}

	rawCriteria, err := json.Marshal(input.criteria)
	if err != nil {
		return Rule{}, errors.Wrap(err, "marshal criteria")
	}

	_, err = r.db.NamedExec(pgRuleInsert, map[string]interface{}{
		"id":          input.id,
		"active":      input.active,
		"buckets":     rawBuckets,
		"configId":    input.configID,
		"createdAt":   input.createdAt,
		"criteria":    rawCriteria,
		"description": input.description,
		"endTime":     input.endTime,
		"kind":        input.kind,
		"name":        input.name,
		"startTime":   input.startTime,
		"updatedAt":   time.Now().UTC(),
	})
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrDuplicateKey:
			return Rule{}, errors.Wrap(errors.ErrExists, "rule")
		case pg.ErrRelationNotFound:
			if err := r.setup(); err != nil {
				return Rule{}, err
			}

			return r.Create(input)
		default:
			return Rule{}, fmt.Errorf("named exec: %s", err)
		}
	}

	return input, nil
}

func (r *pgRepo) GetByName(configID, name string) (Rule, error) {
	query, args, err := r.db.BindNamed(pgRuleGetByName, map[string]interface{}{
		"configId": configID,
		"name":     name,
	})
	if err != nil {
		return Rule{}, fmt.Errorf("named query: %s", err)
	}

	raw := struct {
		ID          string      `db:"id"`
		Active      bool        `db:"active"`
		ActivatedAt pq.NullTime `db:"activated_at"`
		Buckets     []byte      `db:"buckets"`
		ConfigID    string      `db:"config_id"`
		CreatedAt   time.Time   `db:"created_at"`
		Criteria    []byte      `db:"criteria"`
		Description string      `db:"description"`
		Deleted     bool        `db:"deleted"`
		EndTime     time.Time   `db:"end_time"`
		Kind        kind        `db:"kind"`
		Name        string      `db:"name"`
		StartTime   time.Time   `db:"start_time"`
		UpdatedAt   time.Time   `db:"updated_at"`
	}{}

	err = r.db.Get(&raw, query, args...)
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.setup(); err != nil {
				return Rule{}, err
			}

			return r.GetByName(configID, name)
		case sql.ErrNoRows:
			return Rule{}, errors.Wrap(errors.ErrNotFound, "get rule")

		default:
			return Rule{}, fmt.Errorf("get: %s", err)
		}
	}

	buckets := []bucket{}

	if err := json.Unmarshal(raw.Buckets, &buckets); err != nil {
		return Rule{}, errors.Wrap(err, "unmarshal buckets")
	}

	criteria := criteria{}

	if err := json.Unmarshal(raw.Criteria, &criteria); err != nil {
		return Rule{}, errors.Wrap(err, "unmarshal criteria")
	}

	var activatedAt time.Time
	if raw.ActivatedAt.Valid {
		activatedAt = (raw.ActivatedAt).Time
	}

	return Rule{
		id:          raw.ID,
		active:      raw.Active,
		activatedAt: activatedAt,
		buckets:     buckets,
		configID:    raw.ConfigID,
		createdAt:   raw.CreatedAt,
		criteria:    &criteria,
		description: raw.Description,
		deleted:     raw.Deleted,
		endTime:     raw.EndTime,
		kind:        raw.Kind,
		name:        raw.Name,
		startTime:   raw.StartTime,
		updatedAt:   raw.UpdatedAt,
	}, nil
}

func (r *pgRepo) UpdateWith(input Rule) (Rule, error) {
	rawBuckets, err := json.Marshal(input.buckets)
	if err != nil {
		return Rule{}, errors.Wrap(err, "marshal buckets")
	}

	rawCriteria, err := json.Marshal(input.criteria)
	if err != nil {
		return Rule{}, errors.Wrap(err, "marshal criteria")
	}

	_, err = r.db.NamedExec(pgRuleUpdate, map[string]interface{}{
		"id":          input.id,
		"active":      input.active,
		"activatedAt": input.activatedAt,
		"configId":    input.configID,
		"buckets":     rawBuckets,
		"createdAt":   input.createdAt,
		"criteria":    rawCriteria,
		"description": input.description,
		"deleted":     input.deleted,
		"endTime":     input.endTime,
		"kind":        input.kind,
		"name":        input.name,
		"startTime":   input.startTime,
		"updatedAt":   time.Now().UTC(),
	})
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.setup(); err != nil {
				return Rule{}, err
			}

			return r.UpdateWith(input)
		case sql.ErrNoRows:
			return Rule{}, errors.Wrap(errors.ErrNotFound, "update rule")

		default:
			return Rule{}, fmt.Errorf("update named exec: %s", err)
		}
	}

	return input, nil
}

func (r *pgRepo) ListAll(configID string) ([]Rule, error) {
	rows, err := r.db.Queryx(pgRuleListAll)
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.setup(); err != nil {
				return []Rule{}, err
			}

			return r.ListAll(configID)
		case sql.ErrNoRows:
			return []Rule{}, errors.Wrap(errors.ErrNotFound, "list all rules")

		default:
			return []Rule{}, fmt.Errorf("list all rules: %s", err)
		}
	}

	return buildList(rows)
}

func (r *pgRepo) ListActive(configID string, now time.Time) ([]Rule, error) {
	rows, err := r.db.Queryx(pgRuleListActive, now)
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.setup(); err != nil {
				return []Rule{}, err
			}

			return r.ListAll(configID)
		case sql.ErrNoRows:
			return []Rule{}, errors.Wrap(errors.ErrNotFound, "list all active rules")

		default:
			return []Rule{}, fmt.Errorf("list all active rules: %s", err)
		}
	}

	return buildList(rows)
}

func buildList(rows *sqlx.Rows) ([]Rule, error) {
	rules := []Rule{}

	for rows.Next() {
		raw := struct {
			ID          string      `db:"id"`
			Active      bool        `db:"active"`
			ActivatedAt pq.NullTime `db:"activated_at"`
			Buckets     []byte      `db:"buckets"`
			ConfigID    string      `db:"config_id"`
			CreatedAt   time.Time   `db:"created_at"`
			Criteria    []byte      `db:"criteria"`
			Description string      `db:"description"`
			EndTime     time.Time   `db:"end_time"`
			Kind        kind        `db:"kind"`
			Name        string      `db:"name"`
			StartTime   time.Time   `db:"start_time"`
			UpdatedAt   time.Time   `db:"updated_at"`
		}{}

		err := rows.StructScan(&raw)
		if err != nil {
			return []Rule{}, fmt.Errorf("scan rule: %s", err)
		}

		buckets := []bucket{}

		if err := json.Unmarshal(raw.Buckets, &buckets); err != nil {
			return []Rule{}, errors.Wrap(err, "unmarshal buckets in rule scan")
		}

		criteria := criteria{}

		if err := json.Unmarshal(raw.Criteria, &criteria); err != nil {
			return []Rule{}, errors.Wrap(err, "unmarshal criteria in rule scan")
		}

		var activatedAt time.Time
		if raw.ActivatedAt.Valid {
			activatedAt = (raw.ActivatedAt).Time
		}

		rules = append(rules, Rule{
			id:          raw.ID,
			active:      raw.Active,
			activatedAt: activatedAt,
			buckets:     buckets,
			configID:    raw.ConfigID,
			createdAt:   raw.CreatedAt,
			criteria:    &criteria,
			description: raw.Description,
			endTime:     raw.EndTime,
			kind:        raw.Kind,
			name:        raw.Name,
			startTime:   raw.StartTime,
			updatedAt:   raw.UpdatedAt,
		})
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return rules, nil

}

func (r *pgRepo) setup() error {
	for _, q := range []string{
		pgRuleCreateSchema,
		pgRuleCreateTable,
	} {
		_, err := r.db.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *pgRepo) teardown() error {
	for _, q := range []string{
		pgRuleDropTable,
	} {
		_, err := r.db.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}
