package rule

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/oklog/ulid"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/generate"
)

var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

type prepareFunc func(t *testing.T) Repo

func TestRule(t *testing.T) {
	var (
		userID = generate.RandomString(24)
		ctx    = context{
			user: contextUser{
				age: uint8(rand.Intn(99)),
				id:  userID,
			},
		}
		ids = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
			userID,
			generate.RandomString(24),
			generate.RandomString(24),
		}
		input = parameters{
			"feature_x": false,
			"feature_y": false,
		}
		overrideKind = kindOverride
		r            = Rule{
			criteria: &criteria{
				User: &criteriaUser{
					ID: &matcherListString{
						Value: ids,
					},
				},
			},
			buckets: []bucket{
				bucket{
					Parameters: parameters{
						"feature_x": true,
					},
				},
			},
			kind: overrideKind,
		}
	)

	have, err := r.run(input, ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, bucket := range r.buckets {
		want := bucket

		if reflect.DeepEqual(have, want) {
			t.Errorf("have %v, want %v", have, want)
		}
	}
}

func testRepoGet(t *testing.T, p prepareFunc) {
	var (
		repo      = p(t)
		configID  = generate.RandomString(24)
		name      = generate.RandomString(32)
		endTime   = time.Now().Add(1000)
		startTime = time.Now().Add(100)
		ids       = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
		}
		buckets = []bucket{
			bucket{
				Name: generate.RandomString(24),
				Parameters: parameters{
					"feature_x": true,
				},
				Percentage: 100,
			},
		}
		criteria = criteria{
			User: &criteriaUser{
				ID: &matcherListString{
					Value: ids,
				},
			},
		}
	)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	rule := generateRule(
		false,
		id.String(),
		configID,
		name,
		false,
		kindOverride,
		startTime,
		endTime,
		buckets,
		&criteria,
	)

	_, err = repo.Create(rule)
	if err != nil {
		t.Fatal(err)
	}

	r, err := repo.GetByName(rule.configID, rule.name)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := r.configID, rule.configID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := r.id, rule.id; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := r.criteria, rule.criteria; !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave %#v, \nwant %#v", have, want)
	}

	if have, want := r.buckets, rule.buckets; !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave %#v, \nwant %#v", have, want)
	}

	if have, want := r.name, rule.name; !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave %#v, \nwant %#v", have, want)
	}

	if have, want := r.activatedAt.IsZero(), true; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := r.deleted, false; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testRepoGetNotFound(t *testing.T, p prepareFunc) {
	var (
		configID = generate.RandomString(24)
		name     = generate.RandomString(24)
		repo     = p(t)
	)

	_, err := repo.GetByName(configID, name)
	if have, want := errors.Cause(err), errors.ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testRepoListDeleted(t *testing.T, p prepareFunc) {
	var (
		repo        = p(t)
		configID    = generate.RandomString(24)
		nameRuleOne = generate.RandomString(32)
		nameRuleTwo = generate.RandomString(32)
		endTime     = time.Now().Add(1000)
		startTime   = time.Now().Add(100)
		ids         = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
		}
		buckets = []bucket{
			bucket{
				Name: generate.RandomString(24),
				Parameters: parameters{
					"feature_x": true,
				},
				Percentage: 100,
			},
		}
		crit = criteria{
			User: &criteriaUser{
				ID: &matcherListString{
					Value: ids,
				},
			},
		}
		updateIds = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
		}
		updateCriteria = criteria{
			User: &criteriaUser{
				ID: &matcherListString{
					Value: updateIds,
				},
			},
		}
		updateBuckets = []bucket{
			bucket{
				Name: generate.RandomString(24),
				Parameters: parameters{
					"feature_x": true,
				},
				Percentage: 40,
			},
			bucket{
				Name: generate.RandomString(24),
				Parameters: parameters{
					"feature_y": false,
					"feature_z": true,
				},
				Percentage: 60,
			},
		}
	)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	ruleOne := generateRule(
		false,
		id.String(),
		configID,
		nameRuleOne,
		false,
		kindOverride,
		startTime,
		endTime,
		buckets,
		&crit,
	)

	_, err = repo.Create(ruleOne)
	if err != nil {
		t.Fatal(err)
	}

	id, err = ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	ruleTwo := generateRule(
		false,
		id.String(),
		configID,
		nameRuleTwo,
		false,
		kindOverride,
		startTime,
		endTime,
		buckets,
		&crit,
	)

	_, err = repo.Create(ruleTwo)
	if err != nil {
		t.Fatal(err)
	}

	updatedRule := generateRule(
		true,
		ruleOne.id,
		ruleOne.configID,
		ruleOne.name,
		true,
		kindExperiment,
		time.Now().Add(200),
		time.Now().Add(2000),
		updateBuckets,
		&updateCriteria,
	)

	_, err = repo.UpdateWith(updatedRule)
	if err != nil {
		t.Fatal(err)
	}

	rl, err := repo.ListAll(configID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(rl), 1; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testRepoCreateDuplicate(t *testing.T, p prepareFunc) {
	var (
		repo      = p(t)
		configID  = generate.RandomString(24)
		name      = generate.RandomString(32)
		endTime   = time.Now().Add(1000)
		startTime = time.Now().Add(100)
		ids       = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
		}
		buckets = []bucket{
			bucket{
				Name: generate.RandomString(24),
				Parameters: parameters{
					"feature_x": true,
				},
				Percentage: 100,
			},
		}
		criteria = criteria{
			User: &criteriaUser{
				ID: &matcherListString{
					Value: ids,
				},
			},
		}
	)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	rule := generateRule(
		false,
		id.String(),
		configID,
		name,
		false,
		kindOverride,
		startTime,
		endTime,
		buckets,
		&criteria,
	)

	_, err = repo.Create(rule)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Create(rule)
	if have, want := errors.Cause(err), errors.ErrExists; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testRepoUpdateWith(t *testing.T, p prepareFunc) {
	var (
		repo      = p(t)
		configID  = generate.RandomString(24)
		name      = generate.RandomString(32)
		endTime   = time.Now().Add(1000)
		startTime = time.Now().Add(100)
		ids       = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
		}
		buckets = []bucket{
			bucket{
				Name: generate.RandomString(24),
				Parameters: parameters{
					"feature_x": true,
				},
				Percentage: 100,
			},
		}
		crit = criteria{
			User: &criteriaUser{
				ID: &matcherListString{
					Value: ids,
				},
			},
		}
		updateIds = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
		}
		updateCriteria = criteria{
			User: &criteriaUser{
				ID: &matcherListString{
					Value: updateIds,
				},
			},
		}
		updateBuckets = []bucket{
			bucket{
				Name: generate.RandomString(24),
				Parameters: parameters{
					"feature_x": true,
				},
				Percentage: 40,
			},
			bucket{
				Name: generate.RandomString(24),
				Parameters: parameters{
					"feature_y": false,
					"feature_z": true,
				},
				Percentage: 60,
			},
		}
	)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	rule := generateRule(
		false,
		id.String(),
		configID,
		name,
		false,
		kindOverride,
		startTime,
		endTime,
		buckets,
		&crit,
	)

	_, err = repo.Create(rule)
	if err != nil {
		t.Fatal(err)
	}

	updatedRule := generateRule(
		true,
		id.String(),
		rule.configID,
		rule.name,
		rule.deleted,
		kindExperiment,
		time.Now().Add(200),
		time.Now().Add(2000),
		updateBuckets,
		&updateCriteria,
	)

	updatedRule.activatedAt = time.Now().AddDate(0, -1, 0)

	ur, err := repo.UpdateWith(updatedRule)
	if err != nil {
		t.Fatal(err)
	}

	rl, err := repo.GetByName(updatedRule.configID, updatedRule.name)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := ur.configID, rl.configID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := ur.id, rl.id; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := (*ur.criteria), updateCriteria; !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave %#v, \nwant %#v", have, want)
	}

	if have, want := ur.buckets, updateBuckets; !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave %#v, \nwant %#v", have, want)
	}

	if have, want := ur.name, rl.name; !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave %#v, \nwant %#v", have, want)
	}

	if have, want := ur.activatedAt.IsZero(), false; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testRepoListAllEmpty(t *testing.T, p prepareFunc) {
	var (
		repo     = p(t)
		configID = generate.RandomString(24)
	)

	rl, err := repo.ListAll(configID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(rl), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testRepoListAll(t *testing.T, p prepareFunc) {
	var (
		repo        = p(t)
		configID    = generate.RandomString(24)
		nameRuleOne = generate.RandomString(32)
		nameRuleTwo = generate.RandomString(32)
		endTime     = time.Now().Add(1000)
		startTime   = time.Now().Add(100)
		ids         = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
		}
		buckets = []bucket{
			bucket{
				Name: generate.RandomString(24),
				Parameters: parameters{
					"feature_x": true,
				},
				Percentage: 100,
			},
		}
		crit = criteria{
			User: &criteriaUser{
				ID: &matcherListString{
					Value: ids,
				},
			},
		}
	)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	ruleOne := generateRule(
		false,
		id.String(),
		configID,
		nameRuleOne,
		false,
		kindOverride,
		startTime,
		endTime,
		buckets,
		&crit,
	)

	_, err = repo.Create(ruleOne)
	if err != nil {
		t.Fatal(err)
	}

	id, err = ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	ruleTwo := generateRule(
		false,
		id.String(),
		configID,
		nameRuleTwo,
		false,
		kindOverride,
		startTime,
		endTime,
		buckets,
		&crit,
	)

	_, err = repo.Create(ruleTwo)
	if err != nil {
		t.Fatal(err)
	}

	rl, err := repo.ListAll(configID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(rl), 2; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testRepoListActiveEmpty(t *testing.T, p prepareFunc) {
	var (
		repo     = p(t)
		configID = generate.RandomString(24)
	)

	rl, err := repo.ListActive(configID, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(rl), 0; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testRepoListActive(t *testing.T, p prepareFunc) {
	var (
		repo          = p(t)
		configID      = generate.RandomString(24)
		nameRuleOne   = generate.RandomString(32)
		nameRuleTwo   = generate.RandomString(32)
		nameRuleThree = generate.RandomString(32)
		endTime       = time.Now().AddDate(0, 1, 0)
		startTime     = time.Now().AddDate(0, -1, 0)
		ids           = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
		}
		buckets = []bucket{
			bucket{
				Name: generate.RandomString(24),
				Parameters: parameters{
					"feature_x": true,
				},
				Percentage: 100,
			},
		}
		crit = criteria{
			User: &criteriaUser{
				ID: &matcherListString{
					Value: ids,
				},
			},
		}
	)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	ruleOne := generateRule(
		true,
		id.String(),
		configID,
		nameRuleOne,
		false,
		kindOverride,
		startTime,
		endTime,
		buckets,
		&crit,
	)

	_, err = repo.Create(ruleOne)
	if err != nil {
		t.Fatal(err)
	}

	id, err = ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	ruleTwo := generateRule(
		true,
		id.String(),
		configID,
		nameRuleTwo,
		false,
		kindOverride,
		startTime,
		endTime,
		buckets,
		&crit,
	)

	_, err = repo.Create(ruleTwo)
	if err != nil {
		t.Fatal(err)
	}

	id, err = ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	ruleThree := generateRule(
		false,
		id.String(),
		configID,
		nameRuleThree,
		false,
		kindOverride,
		startTime,
		endTime,
		buckets,
		&crit,
	)

	_, err = repo.Create(ruleThree)
	if err != nil {
		t.Fatal(err)
	}

	rl, err := repo.ListActive(configID, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(rl), 2; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func generateRule(
	active bool,
	id, configID, name string,
	deleted bool,
	kind kind,
	startTime, endTime time.Time,
	buckets []bucket,
	criteria *criteria,
) Rule {
	return Rule{
		active:      active,
		buckets:     buckets,
		configID:    configID,
		createdAt:   time.Now(),
		criteria:    criteria,
		deleted:     deleted,
		description: generate.RandomString(24),
		endTime:     endTime,
		id:          id,
		kind:        kind,
		name:        name,
		startTime:   startTime,
	}
}
