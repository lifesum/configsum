package client

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/generate"
)

// Service provides Clients.
type Service interface {
	Create(name string) (Client, error)
	ListWithToken() (clientTokens, error)
	LookupBySecret(secret string) (Client, error)
}

type service struct {
	repo      Repo
	seed      *rand.Rand
	tokenRepo TokenRepo
}

// NewService provides Clients.
func NewService(repo Repo, tokenRepo TokenRepo) Service {
	return &service{
		repo:      repo,
		seed:      rand.New(rand.NewSource(time.Now().UnixNano())),
		tokenRepo: tokenRepo,
	}
}

func (s *service) Create(name string) (Client, error) {
	clientID, err := ulid.New(ulid.Timestamp(time.Now()), s.seed)
	if err != nil {
		return Client{}, err
	}

	c, err := s.repo.Store(clientID.String(), name)
	if err != nil {
		return Client{}, err
	}

	secret, err := generate.SecureToken(secretByteLen)
	if err != nil {
		return Client{}, err
	}

	_, err = s.tokenRepo.Store(clientID.String(), secret)
	if err != nil {
		return Client{}, err
	}

	return c, nil
}

func (s *service) ListWithToken() (clientTokens, error) {
	cs, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	ct := clientTokens{}

	for _, c := range cs {
		t, err := s.tokenRepo.GetLatest(c.id)
		if err != nil {
			continue
		}

		ct[c] = t
	}

	return ct, nil
}

func (s *service) LookupBySecret(secret string) (Client, error) {
	t, err := s.tokenRepo.Lookup(secret)
	if err != nil {
		return Client{}, errors.Wrap(err, "service lookup")
	}

	return s.repo.Lookup(t.clientID)
}

// clientTokens is a mapping of Client to single Token, usually the latest one
// for the specific Client.
type clientTokens map[Client]Token
