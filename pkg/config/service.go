package config

import (
	"math/rand"
	"reflect"
	"time"

	"github.com/oklog/ulid"

	"github.com/lifesum/configsum/pkg/errors"
)

// ServiceUser provides user specific configs.
type ServiceUser interface {
	Render(clientID, baseName, userID string) (UserConfig, error)
}

type serviceUser struct {
	baseRepo BaseRepo
	userRepo UserRepo
	seed     *rand.Rand
}

// NewServiceUser provides user specific configs.
func NewServiceUser(baseRepo BaseRepo, userRepo UserRepo) ServiceUser {
	return &serviceUser{
		baseRepo: baseRepo,
		userRepo: userRepo,
		seed:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *serviceUser) Render(clientID, baseName, userID string) (UserConfig, error) {
	bc, err := s.baseRepo.Get(clientID, baseName)
	if err != nil {
		return UserConfig{}, errors.Wrap(err, "baseRepo.Get")
	}

	uc, err := s.userRepo.GetLatest(bc.ID, userID)
	if err != nil {
		switch errors.Cause(err) {
		case errors.ErrNotFound:
			uc = UserConfig{}
		default:
			return UserConfig{}, errors.Wrap(err, "userRepo.GetLatest")
		}
	}

	// TODO(nabilm): Create temp config with rules applied

	if reflect.DeepEqual(bc.Rendered, uc.rendered) {
		return uc, nil
	}

	id, err := ulid.New(ulid.Timestamp(time.Now()), s.seed)
	if err != nil {
		return UserConfig{}, errors.Wrap(err, "create ulid")
	}

	return s.userRepo.Append(id.String(), bc.ID, userID, nil, bc.Rendered)
}
