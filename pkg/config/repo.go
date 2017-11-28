package config

import (
	"time"

	"github.com/lifesum/configsum/pkg/rule"
)

// BaseRepo provides access to base configs.
type BaseRepo interface {
	lifecycle

	Create(id, clientID, name string, parameters rule.Parameters) (BaseConfig, error)
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
	Parameters rule.Parameters
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

// UserRepo provides access to user configs.
type UserRepo interface {
	lifecycle

	Append(
		id, baseID, userID string,
		decisiosn rule.Decisions,
		render rule.Parameters,
	) (UserConfig, error)
	GetLatest(baseID, userID string) (UserConfig, error)
}

// UserRepoMiddleware is chainable behaviour modifier for UserRepo.
type UserRepoMiddleware func(UserRepo) UserRepo

// UserConfig is a users rendered config.
type UserConfig struct {
	baseID        string
	id            string
	rendered      rule.Parameters
	ruleDecisions rule.Decisions
	userID        string
	createdAt     time.Time
}

type lifecycle interface {
	setup() error
	teardown() error
}
