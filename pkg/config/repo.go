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

// ruleDecisions reflects a matrix of rules applied to a config and if present
// the results of dice rolls for percenatage based decisions.
type ruleDecisions map[string][]int

// UserRepo provides access to user configs.
type UserRepo interface {
	lifecycle

	Append(
		id, baseID, userID string,
		decisiosn ruleDecisions,
		render rendered,
	) (UserConfig, error)
	GetLatest(baseID, userID string) (UserConfig, error)
}

// UserRepoMiddleware is chainable behaviour modifier for UserRepo.
type UserRepoMiddleware func(UserRepo) UserRepo

// UserConfig is a users rendered config.
type UserConfig struct {
	baseID        string
	id            string
	rendered      rendered
	ruleDecisions ruleDecisions
	userID        string
	createdAt     time.Time
}

type lifecycle interface {
	Setup() error
	Teardown() error
}
