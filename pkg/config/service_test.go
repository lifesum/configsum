package config

import (
	"reflect"
	"testing"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/generate"
	"github.com/lifesum/configsum/pkg/rule"
)

func TestBaseServiceUpdate(t *testing.T) {
	var (
		clientID = generate.RandomString(12)
		baseID   = generate.RandomString(16)
		baseName = generate.RandomString(6)
		baseRepo = NewInmemBaseRepo(InmemBaseState{
			clientID: map[string]BaseConfig{
				baseName: BaseConfig{
					ClientID:   clientID,
					ID:         baseID,
					Name:       baseName,
					Parameters: nil,
				},
			},
		})
		svc = NewBaseService(baseRepo, nil)
	)

	_, err := svc.Update(baseID, rendered{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserServiceRender(t *testing.T) {
	var (
		clientID   = generate.RandomString(24)
		baseID     = generate.RandomString(24)
		baseName   = generate.RandomString(24)
		baseRender = rendered{
			generate.RandomString(24): false,
		}
		baseRepo = NewInmemBaseRepo(InmemBaseState{
			clientID: map[string]BaseConfig{
				baseName: BaseConfig{
					ClientID:   clientID,
					ID:         baseID,
					Name:       baseName,
					Parameters: baseRender,
				},
			},
		})
		userID   = generate.RandomString(24)
		userRepo = NewInmemUserRepo()
		ruleRepo = rule.NewInmemRuleRepo()
		svc      = NewUserService(baseRepo, userRepo, ruleRepo)
	)

	uc, err := svc.Render(clientID, baseName, userID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := uc.rendered, baseRender; !reflect.DeepEqual(have, want) {
		t.Errorf("have %#v,want %#v", have, want)
	}

	c, err := userRepo.GetLatest(baseID, userID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := uc, c; !reflect.DeepEqual(have, want) {
		t.Errorf("have %#v,want %#v", have, want)
	}

	rc, err := svc.Render(clientID, baseName, userID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := uc, rc; !reflect.DeepEqual(have, want) {
		t.Errorf("have %#v,want %#v", have, want)
	}
}

func TestUserServiceRenderConfigMissingBaseConfig(t *testing.T) {
	var (
		clientID = generate.RandomString(24)
		baseName = generate.RandomString(24)
		baseRepo = NewInmemBaseRepo(nil)
		userID   = generate.RandomString(24)
		userRepo = NewInmemUserRepo()
		ruleRepo = rule.NewInmemRuleRepo()
		svc      = NewUserService(baseRepo, userRepo, ruleRepo)
	)

	_, err := svc.Render(clientID, baseName, userID)
	if have, want := errors.Cause(err), errors.ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestValidateParamDelta(t *testing.T) {
	var (
		key   = generate.RandomString(6)
		cases = []struct {
			base rendered
			new  rendered
		}{
			{
				base: rendered{
					key: false,
				},
			}, // New missing.
			{
				base: rendered{
					key: false,
				},
				new: rendered{
					key: 12,
				},
			}, // Invalid change of types.
		}
	)

	for _, c := range cases {
		err := validateParamDelta(c.base, c.new)
		if have, want := errors.Cause(err), errors.ErrParametersInvalid; have != want {
			t.Errorf("have %v, want %v", have, want)
		}
	}
}
