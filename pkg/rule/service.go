package rule

// Service for Rule interactions.
type Service interface {
	GetByID(string) (Rule, error)
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

func (s *service) GetByID(id string) (Rule, error) {
	return s.repo.GetByID(id)
}
