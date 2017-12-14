package rule

// Service for Rule interactions.
type Service interface {
	Activate(id string) error
	Deactivate(id string) error
	GetByID(id string) (Rule, error)
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
