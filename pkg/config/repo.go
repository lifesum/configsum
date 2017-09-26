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

// UserRepo provides access to user configs.
type UserRepo interface {
	Get(baseName, id string) (*UserConfig, error)
}

// UserConfig is a users rendered config.
type UserConfig struct {
	baseConfig string
	userID     string
	createdAt  time.Time
}
