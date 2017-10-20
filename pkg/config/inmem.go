package config

import (
	"fmt"
	"time"

	"github.com/lifesum/configsum/pkg/errors"
)

// InmemBaseState is a container to pass initial state for the inmem repo.
type InmemBaseState map[string]map[string]BaseConfig

type inmemBaseRepo struct {
	configs InmemBaseState
}

// NewInmemBaseRepo returns an in-memory backed BaseRepo implementation.
func NewInmemBaseRepo(initial InmemBaseState) (BaseRepo, error) {
	if initial == nil {
		initial = InmemBaseState{}
	}

	return &inmemBaseRepo{
		configs: initial,
	}, nil
}

func (s *inmemBaseRepo) Get(clientID, name string) (BaseConfig, error) {
	client, ok := s.configs[clientID]
	if !ok {
		return BaseConfig{}, errors.Wrap(errors.ErrNotFound, fmt.Sprintf("client id '%s'", clientID))
	}

	bc, ok := client[name]
	if !ok {
		return BaseConfig{}, errors.Wrap(errors.ErrNotFound, fmt.Sprintf("base config name '%s'", name))
	}

	return bc, nil
}

type inmemUserState map[string]map[string][]UserConfig

type inmemUserRepo struct {
	configs inmemUserState
	ids     map[string]struct{}
}

// NewInmemUserRepo returns an in-memory backed UserRepo implementation.
func NewInmemUserRepo() (UserRepo, error) {
	return &inmemUserRepo{
		configs: inmemUserState{},
		ids:     map[string]struct{}{},
	}, nil
}

func (r *inmemUserRepo) Append(
	id, baseID, userID string,
	decisions ruleDecisions,
	render rendered,
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

func (r *inmemUserRepo) Setup() error {
	return nil
}

func (r *inmemUserRepo) Teardown() error {
	return nil
}
