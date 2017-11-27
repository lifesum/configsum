package rule

import (
	"time"

	"github.com/lifesum/configsum/pkg/errors"
)

const (
	kindOverride kind = iota + 1
	kindExperiment
	kindRollout
)

type kind uint8

type parameters map[string]interface{}

type bucket struct {
	Name       string
	Parameters parameters
	Percentage int
}

type context struct {
	user contextUser
}
type contextUser struct {
	age uint8
	id  string
}

// Repo provides access to rules.
type Repo interface {
	lifecycle

	Create(input Rule) (Rule, error)
	GetByName(configID, name string) (Rule, error)
	UpdateWith(input Rule) (Rule, error)
	ListAll(configID string) ([]Rule, error)
	ListActive(configID string, now time.Time) ([]Rule, error)
}

// RepoMiddleware is a chainable behaviour modifier for Repo.
type RepoMiddleware func(Repo) Repo

// Rule facilitates the overide of base configs with consumer provided parameters.
type Rule struct {
	active      bool
	activatedAt time.Time
	buckets     []bucket
	configID    string
	createdAt   time.Time
	criteria    *criteria
	description string
	deleted     bool
	endTime     time.Time
	id          string
	kind        kind
	name        string
	startTime   time.Time
	updatedAt   time.Time
}

type lifecycle interface {
	setup() error
	teardown() error
}

func (r Rule) validate() (bool, error) {
	if r.buckets == nil {
		return false, errors.Wrap(errors.ErrInvalidRule, "missing buckets attribute")
	}

	if r.configID == "" {
		return false, errors.Wrap(errors.ErrInvalidRule, "missing configID attribute")
	}

	if r.createdAt.IsZero() {
		return false, errors.Wrap(errors.ErrInvalidRule, "missing createdAt attribute")
	}

	if r.id == "" {
		return false, errors.Wrap(errors.ErrInvalidRule, "missing id attribute")
	}

	if r.kind == 0 {
		return false, errors.Wrap(errors.ErrInvalidRule, "missing kind attribute")
	}

	if r.name == "" {
		return false, errors.Wrap(errors.ErrInvalidRule, "missing metadate.name attribute")
	}

	totalPercentage := 0
	for _, bucket := range r.buckets {
		totalPercentage = totalPercentage + bucket.Percentage
	}
	if totalPercentage != 100 {
		return false, errors.Wrap(errors.ErrInvalidRule, "percentage not evenly distributed")
	}

	return true, nil
}

func (r Rule) run(input parameters, ctx context) (parameters, error) {
	if r.criteria.User != nil {
		if r.criteria.User.Age != nil {
			return nil, errors.New("matching user age not implemented")
		}

		if r.criteria.User.ID != nil {
			ok, err := r.criteria.User.ID.match(ctx.user.id)
			if err != nil {
				return nil, errors.Wrap(err, "user id match")
			}

			if !ok {
				return nil, errors.New("rule didn't match")
			}
		}
	}

	params := parameters{}

	switch r.kind {
	case kindOverride:
		params = r.buckets[0].Parameters
	case kindExperiment:
		return parameters{}, errors.New("experiment based rules not implemented")
	case kindRollout:
		return parameters{}, errors.New("rollout based rules not implemented")
	}

	for name, value := range params {
		input[name] = value
	}

	return params, nil
}
