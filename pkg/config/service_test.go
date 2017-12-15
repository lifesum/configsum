package config

import (
	"reflect"
	"testing"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/generate"
	"github.com/lifesum/configsum/pkg/rule"
)

func randIntGenerateTest() int {
	return 61
}

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

	_, err := svc.Update(baseID, rule.Parameters{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserServiceRender(t *testing.T) {
	var (
		clientID   = generate.RandomString(24)
		baseID     = generate.RandomString(24)
		baseName   = generate.RandomString(24)
		featureKey = generate.RandomString(24)
		baseRender = rule.Parameters{
			featureKey: false,
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
		ruleID   = generate.RandomString(24)
		ruleRepo = rule.NewInmemRuleRepo()
		svc      = NewUserService(baseRepo, userRepo, ruleRepo)
		matchIDs = rule.MatcherStringList{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
			userID,
			generate.RandomString(24),
			generate.RandomString(24),
		}
	)

	r, err := rule.New(
		ruleID,
		baseID,
		"override",
		"",
		rule.KindOverride,
		true,
		&rule.Criteria{
			User: &rule.CriteriaUser{
				ID: &matchIDs,
			},
		},
		[]rule.Bucket{
			{
				Name: "defualt",
				Parameters: rule.Parameters{
					featureKey: true,
				},
			},
		},
		nil,
		randIntGenerateTest,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ruleRepo.Create(r)
	if err != nil {
		t.Fatal(err)
	}

	uc, err := svc.Render(clientID, baseName, userID, userRenderContext{})
	if err != nil {
		t.Fatal(err)
	}

	want := rule.Parameters{
		featureKey: true,
	}

	if have := uc.rendered; !reflect.DeepEqual(have, want) {
		t.Errorf("have %#v,want %#v", have, want)
	}

	c, err := userRepo.GetLatest(baseID, userID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := uc, c; !reflect.DeepEqual(have, want) {
		t.Errorf("have %#v,want %#v", have, want)
	}

	rc, err := svc.Render(clientID, baseName, userID, userRenderContext{})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := uc, rc; !reflect.DeepEqual(have, want) {
		t.Errorf("have %#v, want %#v", have, want)
	}
}

func TestUserServiceNotInRollout(t *testing.T) {
	var (
		clientID   = generate.RandomString(24)
		baseID     = generate.RandomString(24)
		baseName   = generate.RandomString(24)
		featureKey = generate.RandomString(24)
		baseRender = rule.Parameters{
			featureKey: false,
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
		userID    = generate.RandomString(24)
		userRepo  = NewInmemUserRepo()
		rpOne     = uint8(73) // rule not in rollout
		rpTwo     = uint8(49) // rule in rollout
		ruleOneID = generate.RandomString(24)
		ruleTwoID = generate.RandomString(24)
		ruleRepo  = rule.NewInmemRuleRepo()
		svc       = NewUserService(baseRepo, userRepo, ruleRepo)
		matchIDs  = rule.MatcherStringList{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
			userID,
			generate.RandomString(24),
			generate.RandomString(24),
		}
		ruleOneParam = rule.Parameters{
			featureKey: true,
		}
		ruleTwoParam = rule.Parameters{
			featureKey: false,
		}
	)

	ruleOne, err := rule.New(
		ruleOneID,
		baseID,
		"ruleOneRollout",
		"",
		rule.KindRollout,
		true,
		&rule.Criteria{
			User: &rule.CriteriaUser{
				ID: &matchIDs,
			},
		},
		[]rule.Bucket{
			{
				Name:       "defualt",
				Parameters: ruleOneParam,
			},
		},
		&rpOne,
		randIntGenerateTest,
	)

	if err != nil {
		t.Fatal(err)
	}

	ruleTwo, err := rule.New(
		ruleTwoID,
		baseID,
		"ruleTwoRollout",
		"",
		rule.KindRollout,
		true,
		&rule.Criteria{
			User: &rule.CriteriaUser{
				ID: &matchIDs,
			},
		},
		[]rule.Bucket{
			{
				Name:       "defualt",
				Parameters: ruleTwoParam,
			},
		},
		&rpTwo,
		randIntGenerateTest,
	)

	if err != nil {
		t.Fatal(err)
	}

	_, err = ruleRepo.Create(ruleOne)
	if err != nil {
		t.Fatal(err)
	}

	uc1, err := svc.Render(clientID, baseName, userID, userRenderContext{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = ruleRepo.Create(ruleTwo)
	if err != nil {
		t.Fatal(err)
	}

	uc2, err := svc.Render(clientID, baseName, userID, userRenderContext{})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := uc2.ruleDecisions[ruleOneID], uc1.ruleDecisions[ruleOneID]; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := uc1.rendered, ruleOneParam; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := uc2.rendered, ruleTwoParam; reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestUserServiceRenderFailingRule(t *testing.T) {
	var (
		clientID   = generate.RandomString(24)
		baseID     = generate.RandomString(24)
		baseName   = generate.RandomString(24)
		baseRender = rule.Parameters{
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
		matchIDs = rule.MatcherStringList{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
		}
		ruleID   = generate.RandomString(24)
		ruleRepo = rule.NewInmemRuleRepo()
		userID   = generate.RandomString(24)
		userRepo = NewInmemUserRepo()
		svc      = NewUserService(baseRepo, userRepo, ruleRepo)
	)

	r, err := rule.New(
		ruleID,
		baseID,
		"broken rule",
		"",
		rule.KindOverride,
		true,
		&rule.Criteria{
			User: &rule.CriteriaUser{
				ID: &matchIDs,
			},
		},
		[]rule.Bucket{
			{
				Name: "default",
				Parameters: rule.Parameters{
					"feature_focus_toggle": true,
				},
			},
		},
		nil,
		randIntGenerateTest,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ruleRepo.Create(r)
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Render(clientID, baseName, userID, userRenderContext{})
	if err != nil {
		t.Fatal(err)
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

	_, err := svc.Render(clientID, baseName, userID, userRenderContext{})
	if have, want := errors.Cause(err), errors.ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestValidateParamDelta(t *testing.T) {
	var (
		key   = generate.RandomString(6)
		cases = []struct {
			base rule.Parameters
			new  rule.Parameters
		}{
			{
				base: rule.Parameters{
					key: false,
				},
			}, // New missing.
			{
				base: rule.Parameters{
					key: false,
				},
				new: rule.Parameters{
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
