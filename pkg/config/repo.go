package config

import (
	"time"
)

// BaseRepo provides access to base configs.
type BaseRepo interface {
	lifecycle

	Create(id, clientID, name string, parameters rendered) (BaseConfig, error)
	GetByID(id string) (BaseConfig, error)
	GetByName(clientID, name string) (BaseConfig, error)
	List() (BaseList, error)
	Update(BaseConfig) (BaseConfig, error)
}

// BaseRepoMiddleware is chainable behaviour modifier for BaseRepo.
type BaseRepoMiddleware func(BaseRepo) BaseRepo

// BaseConfig is the entire space of available parameters.
type BaseConfig struct {
	ClientID   string
	Deleted    bool
	ID         string
	Name       string
	Parameters rendered
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// BaseList is a collection of BaseConfig.
type BaseList []BaseConfig

func (l BaseList) Len() int {
	return len(l)
}

func (l BaseList) Less(i, j int) bool {
	return l[i].CreatedAt.After(l[j].CreatedAt)
}

func (l BaseList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
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
	setup() error
	teardown() error
}
