package config

import "time"

// BaseRepo provides access to base configs.
type BaseRepo interface {
	Get(name string) (*BaseConfig, error)
}

// BaseConfig is the entire space of available parameters.
type BaseConfig struct {
	name string
}

type rendered map[string]interface{}

// UserRepo provides access to user configs.
type UserRepo interface {
	lifecycle

	Get(baseID, userID string) (UserConfig, error)
	Put(id, baseID, userID string, render rendered) (UserConfig, error)
}

// UserConfig is a users rendered config.
type UserConfig struct {
	baseID      string
	id          string
	rendered    rendered
	ruleIDs     []string
	userID      string
	createdAt   time.Time
	activatedAt time.Time
}

type lifecycle interface {
	Setup() error
	Teardown() error
}
