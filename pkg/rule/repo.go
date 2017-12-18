package rule

import (
	"time"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/generate"
)

// Supported kinds of rules.
const (
	KindOverride Kind = iota + 1
	KindExperiment
	KindRollout
)

// Kind defines the type of rule.
type Kind uint8

// Parameters is the set of keys and their new values that an applied rule sets.
type Parameters map[string]interface{}

// Bucket is a distinct set of parameters that can be used to control
// segregation by percentage split. Rules which are not of kind experiment will
// only have one bucket.
type Bucket struct {
	Name       string
	Parameters Parameters
	Percentage int
}

// Context carries information for rule decisions to match criteria.
type Context struct {
	User ContextUser
}

// ContextUser bundles user information for rule criteria to match.
type ContextUser struct {
	Age          uint8
	ID           string
	Registered   string
	Subscription int
}

// Decisions reflects a matrix of rules applied to a config and if present the
// results of dice rolls for percenatage based decisions.
type Decisions map[string][]int

// List is a collection of Rule.
type List []Rule

// Repo provides access to rules.
type Repo interface {
	lifecycle

	Create(input Rule) (Rule, error)
	GetByID(string) (Rule, error)
	UpdateWith(input Rule) (Rule, error)
	ListAll() ([]Rule, error)
	ListActive(configID string, now time.Time) ([]Rule, error)
}

// RepoMiddleware is a chainable behaviour modifier for Repo.
type RepoMiddleware func(Repo) Repo

type lifecycle interface {
	setup() error
	teardown() error
}

// Rule facilitates the overide of base configs with consumer provided parameters.
type Rule struct {
	active      bool
	activatedAt time.Time
	buckets     []Bucket
	configID    string
	createdAt   time.Time
	criteria    *Criteria
	description string
	deleted     bool
	endTime     time.Time
	ID          string
	kind        Kind
	name        string
	rollout     uint8
	startTime   time.Time
	updatedAt   time.Time
	RandFunc    generate.RandPercentageFunc
}

// New returns a valid rule.
func New(
	id, configID, name, description string,
	kind Kind,
	active bool,
	criteria *Criteria,
	buckets []Bucket,
	rollout *uint8,
	randFunc generate.RandPercentageFunc,
) (Rule, error) {
	r := Rule{
		active:      active,
		buckets:     buckets,
		configID:    configID,
		createdAt:   time.Now().UTC(),
		criteria:    criteria,
		description: description,
		ID:          id,
		kind:        kind,
		name:        name,
		RandFunc:    randFunc,
	}

	if rollout != nil {
		r.rollout = *rollout
	}

	err := r.validate()
	if err != nil {
		return Rule{}, err
	}

	return r, nil
}

func (r Rule) validate() error {
	if len(r.buckets) == 0 {
		return errors.Wrap(errors.ErrInvalidRule, "missing buckets")
	}

	if r.configID == "" {
		return errors.Wrap(errors.ErrInvalidRule, "missing configID")
	}

	if r.createdAt.IsZero() {
		return errors.Wrap(errors.ErrInvalidRule, "missing createdAt")
	}

	if r.ID == "" {
		return errors.Wrap(errors.ErrInvalidRule, "missing id")
	}

	if r.kind == 0 {
		return errors.Wrap(errors.ErrInvalidRule, "missing kind")
	}

	if r.name == "" {
		return errors.Wrap(errors.ErrInvalidRule, "missing metadate.name")
	}

	if r.rollout > 100 {
		return errors.Wrap(errors.ErrInvalidRule, "rollout percentage too high")
	}

	if len(r.buckets) > 1 {
		totalPercentage := 0
		for _, bucket := range r.buckets {
			totalPercentage = totalPercentage + bucket.Percentage
		}
		if totalPercentage != 100 {
			return errors.Wrap(errors.ErrInvalidRule, "bucket percentage not evenly distributed")
		}
	}

	return nil
}

// Run given an input params and context will try to match based on the rules
// Criteria and if matched overrides the input params with its own.
func (r Rule) Run(input Parameters, ctx Context, decisions []int, randInt generate.RandPercentageFunc) (Parameters, []int, error) {
	if r.criteria != nil && r.criteria.User != nil {
		if r.criteria.User.Age != nil {
			return nil, nil, errors.New("matching user age not implemented")
		}

		if r.criteria.User.ID != nil {
			ok, err := r.criteria.User.ID.match(ctx.User.ID)
			if err != nil {
				return nil, nil, errors.Wrap(err, "user id match")
			}

			if !ok {
				return nil, nil, errors.Wrap(errors.ErrRuleNoMatch, "user id")
			}
		}

		if r.criteria.User.Subscription != nil {
			ok, err := r.criteria.User.Subscription.match(ctx.User.Subscription)
			if err != nil {
				return nil, nil, errors.Wrap(err, "subscription match")
			}

			if !ok {
				return nil, nil, errors.Wrap(errors.ErrRuleNoMatch, "subscription")
			}
		}
	}

	var (
		params = Parameters{}
		d      = []int{}
	)

	diceRollout := randInt()
	if len(decisions) != 0 {
		diceRollout = decisions[0]
	}

	switch r.kind {
	case KindOverride:
		params = r.buckets[0].Parameters
	case KindExperiment:
		return Parameters{}, nil, errors.New("experiment based rules not implemented")
	case KindRollout:
		if len(decisions) != 0 {
			d = decisions
		} else {
			d = append(d, diceRollout)
		}

		if diceRollout <= int(r.rollout) {
			params = r.buckets[0].Parameters
		} else {
			return nil, d, errors.Wrap(errors.ErrRuleNotInRollout, "rollout percentage")
		}
	}

	for name, value := range params {
		input[name] = value
	}

	return input, d, nil
}
