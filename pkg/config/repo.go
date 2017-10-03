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

	Append(
		id, baseID, userID string,
		ruleIDs []string,
		render rendered,
	) (UserConfig, error)
	GetLatest(baseID, userID string) (UserConfig, error)
}

// UserConfig is a users rendered config.
type UserConfig struct {
	baseID    string
	id        string
	rendered  rendered
	ruleIDs   []string
	userID    string
	createdAt time.Time
}

type lifecycle interface {
	Setup() error
	Teardown() error
}
