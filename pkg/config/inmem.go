package config

import "time"

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

func (r *inmemUserRepo) Get(baseName, id string) (*UserConfig, error) {
	return &UserConfig{
		baseConfig: baseName,
		userID:     id,
		createdAt:  time.Now(),
	}, nil
}
