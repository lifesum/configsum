package config

import (
	"fmt"
	"time"
)

type inmemBaseRepo struct{}

// NewInmemBaseRepo returns an in-memory backed BaseRepo implementation.
func NewInmemBaseRepo() (BaseRepo, error) {
	return &inmemBaseRepo{}, nil
}

func (s *inmemBaseRepo) Get(name string) (*BaseConfig, error) {
	return &BaseConfig{
		name: name,
	}, nil
}

type inmemUserRepo struct{}

// NewInmemUserRepo returns an in-memory backed UserRepo implementation.
func NewInmemUserRepo() (UserRepo, error) {
	return &inmemUserRepo{}, nil
}

func (r *inmemUserRepo) Append(
	id, baseID, userID string,
	decisions ruleDecisions,
	render rendered,
) (UserConfig, error) {
	return UserConfig{}, fmt.Errorf("inmemUserRepo.Put() not implemented")
}

func (r *inmemUserRepo) GetLatest(baseID, id string) (UserConfig, error) {
	return UserConfig{
		baseID:    baseID,
		userID:    id,
		createdAt: time.Now(),
	}, nil
}

func (r *inmemUserRepo) Setup() error {
	return nil
}

func (r *inmemUserRepo) Teardown() error {
	return nil
}
