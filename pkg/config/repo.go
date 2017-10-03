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

func (r rendered) SetBool(key string, value bool) {
	r[key] = value
}

func (r rendered) setNumber(key string, value float64) {
	r[key] = value
}

func (r rendered) setNumberList(key string, value []float64) {
	r[key] = value
}

func (r rendered) setString(key, value string) {
	r[key] = value
}

func (r rendered) setStringList(key string, value []string) {
	r[key] = value
}

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
