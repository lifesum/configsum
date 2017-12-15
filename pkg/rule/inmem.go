package rule

import (
	"time"

	"github.com/lifesum/configsum/pkg/errors"
)

type inmemRuleState map[string]map[string]Rule

type inmemRepo struct {
	rules inmemRuleState
	ids   map[string]struct{}
}

// NewInmemRuleRepo returns an in-memory backed Repo implementation.
func NewInmemRuleRepo() Repo {
	return &inmemRepo{
		rules: inmemRuleState{},
		ids:   map[string]struct{}{},
	}
}

func (r *inmemRepo) GetByID(id string) (Rule, error) {
	for _, cs := range r.rules {
		for _, rule := range cs {
			if rule.ID == id {
				return rule, nil
			}
		}
	}

	return Rule{}, errors.Wrapf(errors.ErrNotFound, "rule id '%s'", id)
}

func (r *inmemRepo) Create(input Rule) (Rule, error) {
	if _, ok := r.ids[input.ID]; ok {
		return Rule{}, errors.Wrap(errors.ErrExists, "id")
	}

	if _, ok := r.rules[input.configID]; !ok {
		r.rules[input.configID] = map[string]Rule{}
	}

	if _, ok := r.rules[input.configID][input.name]; ok {
		return Rule{}, errors.Wrap(errors.ErrExists, "name")
	}

	r.rules[input.configID][input.name] = input

	return r.rules[input.configID][input.name], nil
}

func (r *inmemRepo) UpdateWith(input Rule) (Rule, error) {
	if _, ok := r.rules[input.configID]; !ok {
		return Rule{}, errors.Wrapf(errors.ErrNoRuleForID, ": %s", input.configID)
	}

	if _, ok := r.rules[input.configID][input.name]; !ok {
		return Rule{}, errors.Wrapf(errors.ErrNoRuleWithName, ": %s", input.name)
	}

	r.rules[input.configID][input.name] = input

	return r.rules[input.configID][input.name], nil
}

func (r *inmemRepo) ListAll() ([]Rule, error) {
	rules := []Rule{}

	for _, c := range r.rules {
		for _, rule := range c {
			if !rule.deleted {
				rules = append(rules, rule)
			}
		}
	}

	return rules, nil
}

func (r *inmemRepo) ListActive(configID string, now time.Time) ([]Rule, error) {
	rn, ok := r.rules[configID]
	if !ok {
		return []Rule{}, nil
	}

	rules := []Rule{}

	for _, rule := range rn {
		if !rule.deleted && rule.active {
			if !rule.startTime.IsZero() && rule.startTime.After(now) {
				continue
			}

			if !rule.endTime.IsZero() && rule.endTime.Before(now) {
				continue
			}

			rules = append(rules, rule)
		}
	}

	return rules, nil
}

func (r *inmemRepo) setup() error {
	return nil
}

func (r *inmemRepo) teardown() error {
	return nil
}
