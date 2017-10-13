package client

import (
	"github.com/pkg/errors"
)

// Service provides Clients.
type Service interface {
	LookupBySecret(secret string) (Client, error)
}

type service struct {
	repo      Repo
	tokenRepo TokenRepo
}

// NewService provides Clients.
func NewService(repo Repo, tokenRepo TokenRepo) Service {
	return &service{
		repo:      repo,
		tokenRepo: tokenRepo,
	}
}

func (s *service) LookupBySecret(secret string) (Client, error) {
	t, err := s.tokenRepo.Lookup(secret)
	if err != nil {
		return Client{}, errors.Wrap(err, "service lookup")
	}

	return s.repo.Lookup(t.clientID)
}
