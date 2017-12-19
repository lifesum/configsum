package config

import (
	"math/rand"
	"reflect"
	"time"

	"github.com/oklog/ulid"

	"github.com/lifesum/configsum/pkg/client"
	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/rule"
)

// BaseService provides base configs.
type BaseService interface {
	Create(clientID, name string) (BaseConfig, error)
	Get(id string) (BaseConfig, error)
	List() ([]BaseConfig, error)
	Update(id string, parameters rule.Parameters) (BaseConfig, error)
}

type baseService struct {
	baseRepo   BaseRepo
	clientRepo client.Repo
	seed       *rand.Rand
}

// NewBaseService provides base configs.
func NewBaseService(baseRepo BaseRepo, clientRepo client.Repo) BaseService {
	return &baseService{
		baseRepo:   baseRepo,
		clientRepo: clientRepo,
		seed:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *baseService) Create(clientID, name string) (BaseConfig, error) {
	id, err := ulid.New(ulid.Timestamp(time.Now()), s.seed)
	if err != nil {
		return BaseConfig{}, errors.Wrap(errors.ErrID, err.Error())
	}

	_, err = s.clientRepo.Lookup(clientID)
	if err != nil {
		return BaseConfig{}, err
	}

	return s.baseRepo.Create(id.String(), clientID, name, nil)
}

func (s *baseService) Get(id string) (BaseConfig, error) {
	return s.baseRepo.GetByID(id)
}

func (s *baseService) List() ([]BaseConfig, error) {
	cs, err := s.baseRepo.List()
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func (s *baseService) Update(id string, params rule.Parameters) (BaseConfig, error) {
	bc, err := s.baseRepo.GetByID(id)
	if err != nil {
		return BaseConfig{}, err
	}

	err = validateParamDelta(bc.Parameters, params)
	if err != nil {
		return BaseConfig{}, err
	}

	return s.baseRepo.Update(BaseConfig{
		ClientID:   bc.ClientID,
		Deleted:    bc.Deleted,
		ID:         bc.ID,
		Name:       bc.Name,
		Parameters: params,
		CreatedAt:  bc.CreatedAt,
		UpdatedAt:  bc.UpdatedAt,
	})
}

// UserService provides user specific configs.
type UserService interface {
	Render(clientID, baseName, userID string, ctx userRenderContext) (UserConfig, error)
}

type userService struct {
	baseRepo BaseRepo
	userRepo UserRepo
	ruleRepo rule.Repo
	seed     *rand.Rand
}

// NewUserService provides user specific configs.
func NewUserService(baseRepo BaseRepo, userRepo UserRepo, ruleRepo rule.Repo) UserService {
	return &userService{
		baseRepo: baseRepo,
		userRepo: userRepo,
		ruleRepo: ruleRepo,
		seed:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *userService) Render(
	clientID, baseName, userID string,
	ctx userRenderContext,
) (UserConfig, error) {
	bc, err := s.baseRepo.GetByName(clientID, baseName)
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

	rs, err := s.ruleRepo.ListActive(bc.ID, time.Now())
	if err != nil {
		return UserConfig{}, err
	}

	var (
		decisions = rule.Decisions{}
		params    = rule.Parameters(bc.Parameters)
	)

	for _, r := range rs {
		ctx := rule.Context{
			User: rule.ContextUser{
				ID:           userID,
				Age:          ctx.User.Age,
				Registered:   ctx.User.Registered,
				Subscription: ctx.User.Subscription,
			},
			Locale: rule.ContextLocale{
				Locale: ctx.Device.Location.locale,
			},
		}

		pm, d, err := r.Run(params, ctx, uc.ruleDecisions[r.ID], r.RandFunc)
		if err != nil {
			switch errors.Cause(err) {
			case errors.ErrRuleNoMatch:
				continue
			case errors.ErrRuleNotInRollout:
				decisions[r.ID] = d
				continue
			default:
				return UserConfig{}, err
			}
		}

		params = pm
	}

	if reflect.DeepEqual(params, uc.rendered) && reflect.DeepEqual(uc.ruleDecisions, decisions) {
		return uc, nil
	}

	id, err := ulid.New(ulid.Timestamp(time.Now()), s.seed)
	if err != nil {
		return UserConfig{}, errors.Wrap(err, "create ulid")
	}

	return s.userRepo.Append(id.String(), bc.ID, userID, decisions, params)
}

// validateParamDelta given a base and the new version of the parameters
// returns an error if:
// * a key from base is missing in the new version
// * the type of a key was changed
func validateParamDelta(base, new rule.Parameters) error {
	for key, val := range base {
		v, ok := new[key]
		if !ok {
			return errors.Wrapf(errors.ErrParametersInvalid, "key missing '%s'", key)
		}

		if reflect.TypeOf(val).Kind() != reflect.TypeOf(v).Kind() {
			return errors.Wrapf(
				errors.ErrParametersInvalid,
				"value for '%s' missmatch '%s' != '%s'",
				key,
				reflect.TypeOf(val).Kind(),
				reflect.TypeOf(v).Kind(),
			)
		}
	}

	return nil
}
