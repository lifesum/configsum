package rule

// Service for Rule interactions.
type Service interface {
	Activate(id string) error
	Deactivate(id string) error
	GetByID(id string) (Rule, error)
	UpdateRollout(id string, rollout uint8) error
}

type service struct {
	repo Repo
}

// NewService for Rule interactions.
func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Activate(id string) error {
	r, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if r.active {
		return nil
	}

	r.active = true

	_, err = s.repo.UpdateWith(r)

	return err
}

func (s *service) Deactivate(id string) error {
	r, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if !r.active {
		return nil
	}

	r.active = false

	_, err = s.repo.UpdateWith(r)

	return err
}

func (s *service) GetByID(id string) (Rule, error) {
	return s.repo.GetByID(id)
}

func (s *service) UpdateRollout(id string, rollout uint8) error {
	r, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if r.rollout == rollout {
		return nil
	}

	r.rollout = rollout

	_, err = s.repo.UpdateWith(r)

	return err
}
