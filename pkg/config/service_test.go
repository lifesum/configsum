package config

import (
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/generate"
	"github.com/lifesum/configsum/pkg/rule"
)

func randIntGenerateTest() int {
	return 61
}

func TestBaseServiceUpdate(t *testing.T) {
	t.Parallel()

	var (
		clientID   = generate.RandomString(12)
		baseID     = generate.RandomString(16)
		baseName   = generate.RandomString(6)
		baseParams = rule.Parameters{
			generate.RandomString(6): true,
		}
		baseRepo = preparePGBaseRepo(t)
		svc      = NewBaseService(baseRepo, nil)
	)

	_, err := baseRepo.Create(baseID, clientID, baseName, nil)
	if err != nil {
		t.Fatal(err)
	}

	updated, err := svc.Update(baseID, baseParams)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := updated.Parameters, baseParams; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestUserServiceRender(t *testing.T) {
	t.Parallel()

	var (
		clientID   = generate.RandomString(24)
		baseID     = generate.RandomString(24)
		baseName   = generate.RandomString(24)
		featureKey = generate.RandomString(24)
		baseParams = rule.Parameters{
			featureKey: false,
		}
		baseRepo = preparePGBaseRepo(t)
		userID   = generate.RandomString(24)
		userRepo = preparePGUserRepo(t)
		ruleID   = generate.RandomString(24)
		ruleRepo = prepareRuleRepo(t)
		svc      = NewUserService(baseRepo, userRepo, ruleRepo, randIntGenerateTest)
		matchIDs = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
			userID,
			generate.RandomString(24),
			generate.RandomString(24),
		}
	)

	_, err := baseRepo.Create(baseID, clientID, baseName, baseParams)
	if err != nil {
		t.Fatal(err)
	}

	r, err := rule.New(
		ruleID,
		baseID,
		"override",
		"",
		rule.KindOverride,
		true,
		rule.Criteria{
			rule.Criterion{
				Comparator: rule.ComparatorIN,
				Key:        rule.UserID,
				Value:      matchIDs,
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

	if have, want := uc.rendered, c.rendered; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	rc, err := svc.Render(clientID, baseName, userID, userRenderContext{})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := uc.rendered, rc.rendered; !reflect.DeepEqual(have, want) {
		t.Errorf("have %#v, want %#v", have, want)
	}
}

func TestUserServiceNotInRollout(t *testing.T) {
	t.Parallel()

	var (
		clientID   = generate.RandomString(24)
		baseID     = generate.RandomString(24)
		baseName   = generate.RandomString(24)
		baseParams = rule.Parameters{
			"feature_one": false,
			"feature_two": false,
		}
		baseRepo     = preparePGBaseRepo(t)
		userID       = generate.RandomString(24)
		userRepo     = preparePGUserRepo(t)
		rpOne        = uint8(25) // rule not in rollout
		rpTwo        = uint8(70) // rule in rollout
		ruleOneID    = generate.RandomString(24)
		ruleTwoID    = generate.RandomString(24)
		ruleRepo     = prepareRuleRepo(t)
		svc          = NewUserService(baseRepo, userRepo, ruleRepo, randIntGenerateTest)
		ruleOneParam = rule.Parameters{
			"feature_one": true,
		}
		ruleTwoParam = rule.Parameters{
			"feature_two": true,
		}
		expected = rule.Parameters{
			"feature_one": false,
			"feature_two": true,
		}
		matchIDs = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
			userID,
			generate.RandomString(24),
			generate.RandomString(24),
		}
	)

	_, err := baseRepo.Create(baseID, clientID, baseName, baseParams)
	if err != nil {
		t.Fatal(err)
	}

	ruleOne, err := rule.New(
		ruleOneID,
		baseID,
		"ruleOneRollout",
		"",
		rule.KindRollout,
		true,
		rule.Criteria{
			rule.Criterion{
				Comparator: rule.ComparatorIN,
				Key:        rule.UserID,
				Value:      matchIDs,
			},
		},
		[]rule.Bucket{
			{
				Name:       "defualt",
				Parameters: ruleOneParam,
			},
		},
		&rpOne,
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
		rule.Criteria{
			rule.Criterion{
				Comparator: rule.ComparatorIN,
				Key:        rule.UserID,
				Value:      matchIDs,
			},
		},
		[]rule.Bucket{
			{
				Name:       "defualt",
				Parameters: ruleTwoParam,
			},
		},
		&rpTwo,
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

	if have, want := uc1.rendered, baseParams; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	_, err = ruleRepo.Create(ruleTwo)
	if err != nil {
		t.Fatal(err)
	}

	uc2, err := svc.Render(clientID, baseName, userID, userRenderContext{})
	if err != nil {
		t.Fatal(err)
	}

	if have, want := uc2.ruleDecisions[ruleTwoID], uc1.ruleDecisions[ruleTwoID]; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := uc2.rendered, expected; !reflect.DeepEqual(have, want) {
		t.Errorf("have %#v, want %#v", have, want)
	}
}

func TestUserServiceRenderFailingRule(t *testing.T) {
	t.Parallel()

	var (
		clientID   = generate.RandomString(24)
		baseID     = generate.RandomString(24)
		baseName   = generate.RandomString(24)
		baseParams = rule.Parameters{
			generate.RandomString(24): false,
		}
		baseRepo = preparePGBaseRepo(t)
		matchIDs = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
		}
		ruleID   = generate.RandomString(24)
		ruleRepo = prepareRuleRepo(t)
		userID   = generate.RandomString(24)
		userRepo = preparePGUserRepo(t)
		svc      = NewUserService(baseRepo, userRepo, ruleRepo, randIntGenerateTest)
	)

	_, err := baseRepo.Create(baseID, clientID, baseName, baseParams)
	if err != nil {
		t.Fatal(err)
	}

	r, err := rule.New(
		ruleID,
		baseID,
		"broken rule",
		"",
		rule.KindOverride,
		true,
		rule.Criteria{
			rule.Criterion{
				Comparator: rule.ComparatorIN,
				Key:        rule.UserID,
				Value:      matchIDs,
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
	t.Parallel()

	var (
		clientID = generate.RandomString(24)
		baseName = generate.RandomString(24)
		baseRepo = preparePGBaseRepo(t)
		userID   = generate.RandomString(24)
		userRepo = preparePGUserRepo(t)
		ruleRepo = prepareRuleRepo(t)
		svc      = NewUserService(baseRepo, userRepo, ruleRepo, randIntGenerateTest)
	)

	_, err := svc.Render(clientID, baseName, userID, userRenderContext{})
	if have, want := errors.Cause(err), errors.ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestValidateParamDelta(t *testing.T) {
	t.Parallel()

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

func prepareRuleRepo(t *testing.T) rule.Repo {
	db, err := sqlx.Connect("postgres", pgURI)
	if err != nil {
		t.Fatal(err)
	}

	r := rule.NewPostgresRepo(db, rule.PGRepoSchema(t.Name()))

	return r
}
