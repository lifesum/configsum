package rule

import (
	"fmt"
	"time"

	"github.com/lifesum/configsum/pkg/errors"
)

type inmemRuleState map[string]map[string]rule

type inmemRepo struct {
	rules inmemRuleState
	ids   map[string]struct{}
}

// NewInmemRepo returns an in-memory backed Repo implementation.
func NewInmemRepo() (Repo, error) {
	return &inmemRepo{
		rules: inmemRuleState{},
		ids:   map[string]struct{}{},
	}, nil
}

func (r *inmemRepo) GetByName(configID, name string) (rule, error) {
	rules, ok := r.rules[configID]
	if !ok {
		return rule{}, errors.Wrap(errors.ErrNotFound, fmt.Sprintf("config id '%s'", configID))
	}

	rl, ok := rules[name]
	if !ok {
		return rule{}, errors.Wrap(errors.ErrNotFound, fmt.Sprintf("rule name '%s'", name))
	}

	return rl, nil
}

func (r *inmemRepo) Create(input rule) (rule, error) {
	if _, ok := r.ids[input.id]; ok {
		return rule{}, errors.Wrap(errors.ErrExists, "id")
	}

	if _, ok := r.rules[input.configID]; !ok {
		r.rules[input.configID] = map[string]rule{}
	}

	if _, ok := r.rules[input.configID][input.name]; ok {
		return rule{}, errors.Wrap(errors.ErrExists, "name")
	}

	r.rules[input.configID][input.name] = input

	return r.rules[input.configID][input.name], nil
}

func (r *inmemRepo) UpdateWith(input rule) (rule, error) {
	if _, ok := r.rules[input.configID]; !ok {
		return rule{}, errors.Wrapf(errors.ErrNoRuleForID, ": %s", input.configID)
	}

	if _, ok := r.rules[input.configID][input.name]; !ok {
		return rule{}, errors.Wrapf(errors.ErrNoRuleWithName, ": %s", input.name)
	}

	r.rules[input.configID][input.name] = input

	return r.rules[input.configID][input.name], nil
}

func (r *inmemRepo) ListAll(configID string) ([]rule, error) {
	rn, ok := r.rules[configID]
	if !ok {
		return []rule{}, nil
	}

	rules := []rule{}

	for _, rule := range rn {
		if !rule.deleted {
			rules = append(rules, rule)
		}
	}

	return rules, nil
}

func (r *inmemRepo) ListActive(configID string, now time.Time) ([]rule, error) {
	rn, ok := r.rules[configID]
	if !ok {
		return []rule{}, nil
	}

	rules := []rule{}

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
