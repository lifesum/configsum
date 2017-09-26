package config

import (
	"time"
)

// ServiceUser provides user specific configs.
type ServiceUser interface {
	Get() (config, error)
}

type config struct {
	createdAt time.Time
}

type serviceUser struct{}

// NewServiceUser provides user specific configs.
func NewServiceUser() ServiceUser {
	return &serviceUser{}
}

func (s *serviceUser) Get() (config, error) {
	return config{createdAt: time.Now()}, nil
}
