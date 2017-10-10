package config

import (
	"math/rand"
	"reflect"
	"time"

	"github.com/oklog/ulid"
	"github.com/pkg/errors"
)

// ServiceUser provides user specific configs.
type ServiceUser interface {
	Render(appID, baseName, userID string) (UserConfig, error)
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

func (s *serviceUser) Render(appID, baseName, userID string) (UserConfig, error) {
	bc, err := s.baseRepo.Get(appID, baseName)
	if err != nil {
		return UserConfig{}, errors.Wrap(err, "baseRepo.Get")
	}

	uc, err := s.userRepo.GetLatest(bc.id, userID)
	if err != nil {
		switch errors.Cause(err) {
		case ErrNotFound:
			uc = UserConfig{}
		default:
			return UserConfig{}, errors.Wrap(err, "userRepo.GetLatest")
		}
	}

	// TODO(nabilm): Create temp config with rules applied

	if reflect.DeepEqual(bc.rendered, uc.rendered) {
		return uc, nil
	}

	id, err := ulid.New(ulid.Timestamp(time.Now()), s.seed)
	if err != nil {
		return UserConfig{}, errors.Wrap(err, "create ulid")
	}

	return s.userRepo.Append(id.String(), bc.id, userID, nil, bc.rendered)
}
