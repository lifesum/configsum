package config

// ServiceUser provides user specific configs.
type ServiceUser interface {
	Get(baseConfig, userID string) (UserConfig, error)
}

type serviceUser struct {
	baseRepo BaseRepo
	userRepo UserRepo
}

// NewServiceUser provides user specific configs.
func NewServiceUser(baseRepo BaseRepo, userRepo UserRepo) ServiceUser {
	return &serviceUser{
		baseRepo: baseRepo,
		userRepo: userRepo,
	}
}

func (s *serviceUser) Get(baseConfig, userID string) (UserConfig, error) {
	// lookup base config
	bc, err := s.baseRepo.Get(baseConfig)
	if err != nil {
		return UserConfig{}, err
	}

	// lookup current config by userID
	return s.userRepo.Get(bc.name, userID)

	// compare configs
	// store config
	// return rendered config
}
