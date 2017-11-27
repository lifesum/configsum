package rule

import (
	"fmt"
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

func (r *inmemRepo) GetByName(configID, name string) (Rule, error) {
	rules, ok := r.rules[configID]
	if !ok {
		return Rule{}, errors.Wrap(errors.ErrNotFound, fmt.Sprintf("config id '%s'", configID))
	}

	rl, ok := rules[name]
	if !ok {
		return Rule{}, errors.Wrap(errors.ErrNotFound, fmt.Sprintf("rule name '%s'", name))
	}

	return rl, nil
}

func (r *inmemRepo) Create(input Rule) (Rule, error) {
	if _, ok := r.ids[input.id]; ok {
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

func (r *inmemRepo) ListAll(configID string) ([]Rule, error) {
	rn, ok := r.rules[configID]
	if !ok {
		return []Rule{}, nil
	}

	rules := []Rule{}

	for _, rule := range rn {
		if !rule.deleted {
			rules = append(rules, rule)
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
		if !rule.deleted &&
			rule.active &&
			rule.startTime.Before(now) &&
			rule.endTime.After(now) {
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
