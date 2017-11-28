package config

import (
	"fmt"
	"sort"
	"time"

	"github.com/lifesum/configsum/pkg/errors"

	"github.com/lifesum/configsum/pkg/rule"
)

// InmemBaseState is a container to pass initial state for the inmem repo.
type InmemBaseState map[string]map[string]BaseConfig

type inmemBaseRepo struct {
	configs InmemBaseState
}

// NewInmemBaseRepo returns an in-memory backed BaseRepo implementation.
func NewInmemBaseRepo(initial InmemBaseState) BaseRepo {
	if initial == nil {
		initial = InmemBaseState{}
	}

	return &inmemBaseRepo{
		configs: initial,
	}
}

func (s *inmemBaseRepo) Create(
	id, clientID, name string,
	parameters rule.Parameters,
) (BaseConfig, error) {
	_, ok := s.configs[clientID]
	if !ok {
		s.configs[clientID] = map[string]BaseConfig{}
	}

	_, ok = s.configs[clientID][name]
	if ok {
		return BaseConfig{}, errors.Wrapf(errors.ErrExists, "duplicate name '%s'", name)
	}

	c := BaseConfig{
		ClientID:   clientID,
		ID:         id,
		Name:       name,
		Parameters: parameters,
		CreatedAt:  time.Now(),
	}

	s.configs[clientID][name] = c

	return c, nil
}

func (s *inmemBaseRepo) GetByID(id string) (BaseConfig, error) {
	for _, cs := range s.configs {
		for _, c := range cs {
			if c.ID == id {
				return c, nil
			}
		}
	}

	return BaseConfig{}, errors.Wrapf(errors.ErrNotFound, "base config '%s'", id)
}

func (s *inmemBaseRepo) GetByName(clientID, name string) (BaseConfig, error) {
	client, ok := s.configs[clientID]
	if !ok {
		return BaseConfig{}, errors.Wrapf(errors.ErrNotFound, "client id '%s'", clientID)
	}

	bc, ok := client[name]
	if !ok {
		return BaseConfig{}, errors.Wrapf(errors.ErrNotFound, "base config name '%s'", name)
	}

	return bc, nil
}

func (s *inmemBaseRepo) List() (BaseList, error) {
	cs := BaseList{}

	for _, clientConfigs := range s.configs {
		for _, c := range clientConfigs {
			cs = append(cs, c)
		}
	}

	sort.Sort(cs)

	return cs, nil
}

func (s *inmemBaseRepo) Update(c BaseConfig) (BaseConfig, error) {
	client, ok := s.configs[c.ClientID]
	if !ok {
		return BaseConfig{}, errors.Wrapf(errors.ErrNotFound, "id '%s'", c.ID)
	}

	_, ok = client[c.Name]
	if !ok {
		return BaseConfig{}, errors.Wrapf(errors.ErrNotFound, "id '%s'", c.ID)
	}

	c = BaseConfig{
		ClientID:   c.ClientID,
		Deleted:    c.Deleted,
		ID:         c.ID,
		Name:       c.Name,
		Parameters: c.Parameters,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  time.Now(),
	}

	s.configs[c.ClientID][c.Name] = c

	return c, nil
}

func (s *inmemBaseRepo) setup() error {
	return nil
}

func (s *inmemBaseRepo) teardown() error {
	return nil
}

type inmemUserState map[string]map[string][]UserConfig

type inmemUserRepo struct {
	configs inmemUserState
	ids     map[string]struct{}
}

// NewInmemUserRepo returns an in-memory backed UserRepo implementation.
func NewInmemUserRepo() UserRepo {
	return &inmemUserRepo{
		configs: inmemUserState{},
		ids:     map[string]struct{}{},
	}
}

func (r *inmemUserRepo) Append(
	id, baseID, userID string,
	decisions rule.Decisions,
	render rule.Parameters,
) (UserConfig, error) {
	if _, ok := r.ids[id]; ok {
		return UserConfig{}, errors.Wrap(errors.ErrExists, "id")
	}

	if _, ok := r.configs[baseID]; !ok {
		r.configs[baseID] = map[string][]UserConfig{}
	}

	ul, ok := r.configs[baseID][userID]
	if !ok {
		ul = []UserConfig{}
	}

	c := UserConfig{
		baseID:        baseID,
		id:            id,
		rendered:      render,
		ruleDecisions: decisions,
		userID:        userID,
		createdAt:     time.Now(),
	}

	ul = append(ul, c)
	r.configs[baseID][userID] = ul
	r.ids[id] = struct{}{}

	return c, nil
}

func (r *inmemUserRepo) GetLatest(baseID, id string) (UserConfig, error) {
	_, ok := r.configs[baseID]
	if !ok {
		return UserConfig{}, errors.Wrap(errors.ErrNotFound, fmt.Sprintf("base id '%s'", baseID))
	}

	cs, ok := r.configs[baseID][id]
	if !ok {
		return UserConfig{}, errors.Wrap(errors.ErrNotFound, fmt.Sprintf("user id '%s'", id))
	}

	if len(cs) == 0 {
		return UserConfig{}, errors.Wrap(errors.ErrNotFound, fmt.Sprintf("no config '%s'", id))
	}

	return cs[len(cs)-1], nil
}

func (r *inmemUserRepo) setup() error {
	return nil
}

func (r *inmemUserRepo) teardown() error {
	return nil
}
