package config

import "fmt"

// ServiceUser provides user specific configs.
type ServiceUser interface {
	Get() error
}

type serviceUser struct{}

// NewServiceUser provides user specific configs.
func NewServiceUser() ServiceUser {
	return &serviceUser{}
}

func (s *serviceUser) Get() error {
	return fmt.Errorf("serviceUser.Get() not implemented")
}
