package config

import (
	"time"
)

// ServiceUser provides user specific configs.
type ServiceUser interface {
	Get(baseConfig string) (config, error)
}

type config struct {
	baseConfig string
	createdAt  time.Time
}

type serviceUser struct{}

// NewServiceUser provides user specific configs.
func NewServiceUser() ServiceUser {
	return &serviceUser{}
}

func (s *serviceUser) Get(baseConfig string) (config, error) {
	// lookup base config
	// lookup current config
	// apply rules to base config
	// compare configs
	// return is newer
	return config{
		baseConfig: baseConfig,
		createdAt:  time.Now(),
	}, nil
}
