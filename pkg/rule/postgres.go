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
			rollout INT8 NOT NULL,
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
			rollout,
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
				:rollout,
				:startTime,
				:updatedAt
		)`

	pgRuleGetByID = `
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
			rollout,
			start_time,
			updated_at
		FROM
			rule.rules
		WHERE
			id = :id
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
			rollout = :rollout,
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
			rollout,
			start_time,
			updated_at
		FROM
			rule.rules
		WHERE
			deleted = false`

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
			rollout,
			start_time,
			updated_at
		FROM
			rule.rules
		WHERE
			active = true
			AND config_id = :configId
			AND deleted = false
			AND (
				end_time IS NULL
				OR end_time >= :now
			)
			AND (
				start_time IS NULL
				OR start_time <= :now
			)`
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

	input.createdAt = input.createdAt.UTC()
	input.updatedAt = time.Now().UTC()

	args := map[string]interface{}{
		"id":          input.ID,
		"active":      input.active,
		"buckets":     rawBuckets,
		"configId":    input.configID,
		"createdAt":   input.createdAt,
		"criteria":    rawCriteria,
		"description": input.description,
		"endTime":     input.endTime,
		"kind":        input.kind,
		"name":        input.name,
		"rollout":     input.rollout,
		"startTime":   input.startTime,
		"updatedAt":   input.updatedAt,
	}

	if input.endTime.IsZero() {
		args["endTime"] = nil
	}

	if input.startTime.IsZero() {
		args["startTime"] = nil
	}

	_, err = r.db.NamedExec(pgRuleInsert, args)
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

func (r *pgRepo) GetByID(id string) (Rule, error) {
	query, args, err := r.db.BindNamed(pgRuleGetByID, map[string]interface{}{
		"id": id,
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
		EndTime     pq.NullTime `db:"end_time"`
		Kind        Kind        `db:"kind"`
		Name        string      `db:"name"`
		Rollout     uint8       `db:"rollout"`
		StartTime   pq.NullTime `db:"start_time"`
		UpdatedAt   time.Time   `db:"updated_at"`
	}{}

	err = r.db.Get(&raw, query, args...)
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.setup(); err != nil {
				return Rule{}, err
			}

			return r.GetByID(id)
		case sql.ErrNoRows:
			return Rule{}, errors.Wrap(errors.ErrNotFound, "get rule")

		default:
			return Rule{}, fmt.Errorf("get: %s", err)
		}
	}

	buckets := []Bucket{}

	if err := json.Unmarshal(raw.Buckets, &buckets); err != nil {
		return Rule{}, errors.Wrap(err, "unmarshal buckets")
	}

	// TODO(xla): If the the value in the column is NULL criteria will be non
	// nil if we don't have this extra check in place.
	var criteria *Criteria

	if len(raw.Criteria) > 0 && string(raw.Criteria) != "null" {
		criteria = &Criteria{}

		if err := json.Unmarshal(raw.Criteria, criteria); err != nil {
			return Rule{}, errors.Wrap(err, "unmarshal criteria")
		}
	}

	var activatedAt time.Time
	if raw.ActivatedAt.Valid {
		activatedAt = (raw.ActivatedAt).Time
	}

	var endTime time.Time
	if raw.EndTime.Valid {
		endTime = (raw.EndTime).Time
	}

	var startTime time.Time
	if raw.StartTime.Valid {
		startTime = (raw.StartTime).Time
	}

	return Rule{
		ID:          raw.ID,
		active:      raw.Active,
		activatedAt: activatedAt,
		buckets:     buckets,
		configID:    raw.ConfigID,
		createdAt:   raw.CreatedAt.UTC(),
		criteria:    criteria,
		description: raw.Description,
		deleted:     raw.Deleted,
		endTime:     endTime,
		kind:        raw.Kind,
		name:        raw.Name,
		rollout:     raw.Rollout,
		startTime:   startTime,
		updatedAt:   raw.UpdatedAt.UTC(),
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
		"id":          input.ID,
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
		"rollout":     input.rollout,
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

func (r *pgRepo) ListAll() ([]Rule, error) {
	rows, err := r.db.Queryx(pgRuleListAll)
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.setup(); err != nil {
				return []Rule{}, err
			}

			return r.ListAll()
		case sql.ErrNoRows:
			return []Rule{}, errors.Wrap(errors.ErrNotFound, "list all rules")

		default:
			return []Rule{}, fmt.Errorf("list all rules: %s", err)
		}
	}

	return buildList(rows)
}

func (r *pgRepo) ListActive(configID string, now time.Time) ([]Rule, error) {
	query, args, err := r.db.BindNamed(pgRuleListActive, map[string]interface{}{
		"configId": configID,
		"now":      now,
	})
	if err != nil {
		return nil, fmt.Errorf("named query: %s", err)
	}

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		switch errors.Cause(pg.Wrap(err)) {
		case pg.ErrRelationNotFound:
			if err := r.setup(); err != nil {
				return []Rule{}, err
			}

			return r.ListActive(configID, now)
		case sql.ErrNoRows:
			return []Rule{}, errors.Wrap(errors.ErrNotFound, "list all active rules")

		default:
			return []Rule{}, fmt.Errorf("list all active rules: %s", err)
		}
	}

	return buildList(rows)
}

func buildList(rows *sqlx.Rows) ([]Rule, error) {
	defer rows.Close()

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
			EndTime     pq.NullTime `db:"end_time"`
			Kind        Kind        `db:"kind"`
			Name        string      `db:"name"`
			Rollout     uint8       `db:"rollout"`
			StartTime   pq.NullTime `db:"start_time"`
			UpdatedAt   time.Time   `db:"updated_at"`
		}{}

		err := rows.StructScan(&raw)
		if err != nil {
			return []Rule{}, fmt.Errorf("scan rule: %s", err)
		}

		buckets := []Bucket{}

		if err := json.Unmarshal(raw.Buckets, &buckets); err != nil {
			return []Rule{}, errors.Wrap(err, "unmarshal buckets in rule scan")
		}

		criteria := Criteria{}

		if err := json.Unmarshal(raw.Criteria, &criteria); err != nil {
			return []Rule{}, errors.Wrap(err, "unmarshal criteria in rule scan")
		}

		var activatedAt time.Time
		if raw.ActivatedAt.Valid {
			activatedAt = (raw.ActivatedAt).Time
		}

		var endTime time.Time
		if raw.EndTime.Valid {
			endTime = (raw.EndTime).Time
		}

		var startTime time.Time
		if raw.StartTime.Valid {
			startTime = (raw.StartTime).Time
		}

		rules = append(rules, Rule{
			ID:          raw.ID,
			active:      raw.Active,
			activatedAt: activatedAt,
			buckets:     buckets,
			configID:    raw.ConfigID,
			createdAt:   raw.CreatedAt,
			criteria:    &criteria,
			description: raw.Description,
			endTime:     endTime,
			kind:        raw.Kind,
			name:        raw.Name,
			rollout:     raw.Rollout,
			startTime:   startTime,
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
