package rule

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/oklog/ulid"
	"golang.org/x/text/language"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/generate"
)

var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

type prepareFunc func(t *testing.T) Repo

func randIntGenerateTest() int {
	return 61
}

func TestRuleMatch(t *testing.T) {
	t.Parallel()

	var (
		userID    = generate.RandomString(24)
		ruleID, _ = ulid.New(ulid.Timestamp(time.Now()), seed)
		baseID    = generate.RandomString(16)
		gb        = language.MustParseRegion("GB")
		en        = language.MustParseBase("en")
		rp        = uint8(77) // rule in rollout
		ctx       = Context{
			User: ContextUser{
				Age:          uint8(rand.Intn(99)),
				ID:           userID,
				Subscription: 2,
			},
			Locale: ContextLocale{
				Locale: language.Make("en_GB"),
			},
		}
		ids = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			userID,
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
		}
		input = Parameters{
			"feature_x": false,
			"feature_y": false,
		}
	)

	enGB, err := language.Compose(en, gb)
	if err != nil {
		t.Fatal(err)
	}

	rule, err := New(
		ruleID.String(),
		baseID,
		generate.RandomString(12),
		generate.RandomString(12),
		KindRollout,
		true,
		Criteria{
			Criterion{
				Comparator: ComparatorGT,
				Key:        UserSubscription,
				Value:      1,
			},
			Criterion{
				Comparator: ComparatorIN,
				Key:        UserID,
				Value:      ids,
			},
			Criterion{
				Comparator: ComparatorEQ,
				Key:        DeviceLocationLocale,
				Value:      enGB,
			},
		},
		[]Bucket{
			Bucket{
				Parameters: Parameters{
					"feature_x": true,
				},
			},
		},
		&rp,
	)

	if err != nil {
		t.Fatal(err)
	}

	have, _, err := rule.Run(input, ctx, nil, randIntGenerateTest)
	if err != nil {
		t.Fatal(err)
	}

	want := Parameters{
		"feature_x": true,
		"feature_y": false,
	}

	if !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestRuleNoMatch(t *testing.T) {
	t.Parallel()

	var (
		userID    = generate.RandomString(24)
		ruleID, _ = ulid.New(ulid.Timestamp(time.Now()), seed)
		baseID    = generate.RandomString(16)
		rp        = uint8(63) // rule in rollout
		ctx       = Context{
			User: ContextUser{
				Age:          uint8(rand.Intn(99)),
				ID:           userID,
				Subscription: 0,
			},
		}
		input = Parameters{
			"feature_x": false,
			"feature_y": false,
		}
		ids = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
		}
	)

	rule, err := New(
		ruleID.String(),
		baseID,
		generate.RandomString(12),
		generate.RandomString(12),
		KindRollout,
		true,
		Criteria{
			Criterion{
				Comparator: ComparatorIN,
				Key:        UserID,
				Value:      ids,
			},
			Criterion{
				Comparator: ComparatorGT,
				Key:        UserSubscription,
				Value:      0,
			},
		},
		[]Bucket{
			Bucket{
				Parameters: Parameters{
					"feature_x": true,
				},
			},
		},
		&rp,
	)

	if err != nil {
		t.Fatal(err)
	}

	_, _, err = rule.Run(input, ctx, nil, randIntGenerateTest)
	if have, want := errors.Cause(err), errors.ErrCriterionNotMatch; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestRuleDecisions(t *testing.T) {
	t.Parallel()

	var (
		userID    = generate.RandomString(24)
		ruleID, _ = ulid.New(ulid.Timestamp(time.Now()), seed)
		baseID    = generate.RandomString(16)
		rp        = uint8(69) // rule in rollout
		ctx       = Context{
			User: ContextUser{
				Age: uint8(rand.Intn(99)),
				ID:  userID,
			},
		}
		input = Parameters{
			"feature_x": false,
			"feature_y": false,
		}
		ids = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
			userID,
			generate.RandomString(24),
			generate.RandomString(24),
		}
	)

	r, err := New(
		ruleID.String(),
		baseID,
		generate.RandomString(12),
		generate.RandomString(12),
		KindRollout,
		true,
		Criteria{
			Criterion{
				Comparator: ComparatorIN,
				Key:        UserID,
				Value:      ids,
			},
		},
		[]Bucket{
			Bucket{
				Parameters: Parameters{
					"feature_x": true,
				},
			},
		},
		&rp,
	)

	if err != nil {
		t.Fatal(err)
	}

	_, d, err := r.Run(input, ctx, nil, randIntGenerateTest)
	if err != nil {
		t.Fatal(err)
	}

	have := d[0]

	if want := randIntGenerateTest(); have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestRuleNoDecisions(t *testing.T) {
	t.Parallel()

	var (
		userID    = generate.RandomString(24)
		ruleID, _ = ulid.New(ulid.Timestamp(time.Now()), seed)
		baseID    = generate.RandomString(16)
		rp        = uint8(41) // rule not in rollout
		ctx       = Context{
			User: ContextUser{
				Age: uint8(rand.Intn(99)),
				ID:  userID,
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
		input = Parameters{
			"feature_x": false,
			"feature_y": false,
		}
	)

	r, err := New(
		ruleID.String(),
		baseID,
		generate.RandomString(12),
		generate.RandomString(12),
		KindRollout,
		true,
		Criteria{
			Criterion{
				Comparator: ComparatorIN,
				Key:        UserID,
				Value:      ids,
			},
		},
		[]Bucket{
			Bucket{
				Parameters: Parameters{
					"feature_x": true,
				},
			},
		},
		&rp,
	)

	if err != nil {
		t.Fatal(err)
	}

	_, _, err = r.Run(input, ctx, nil, randIntGenerateTest)
	if have, want := errors.Cause(err), errors.ErrRuleNotInRollout; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestRuleRollout(t *testing.T) {
	t.Parallel()

	var (
		userID    = generate.RandomString(24)
		ruleID, _ = ulid.New(ulid.Timestamp(time.Now()), seed)
		baseID    = generate.RandomString(16)
		rp        = uint8(83) // rule in rollout
		ctx       = Context{
			User: ContextUser{
				Age: uint8(rand.Intn(99)),
				ID:  userID,
			},
		}
		input = Parameters{
			"feature_x": false,
			"feature_y": false,
		}
		ids = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
			userID,
			generate.RandomString(24),
			generate.RandomString(24),
		}
	)

	r, err := New(
		ruleID.String(),
		baseID,
		generate.RandomString(12),
		generate.RandomString(12),
		KindRollout,
		true,
		Criteria{
			Criterion{
				Comparator: ComparatorIN,
				Key:        UserID,
				Value:      ids,
			},
		},
		[]Bucket{
			Bucket{
				Parameters: Parameters{
					"feature_x": true,
				},
			},
		},
		&rp,
	)

	if err != nil {
		t.Fatal(err)
	}

	have, _, err := r.Run(input, ctx, nil, randIntGenerateTest)
	if err != nil {
		t.Fatal(err)
	}

	want := Parameters{
		"feature_x": true,
		"feature_y": false,
	}

	if !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestRuleRolloutOutsideBucket(t *testing.T) {
	t.Parallel()

	var (
		userID    = generate.RandomString(24)
		ruleID, _ = ulid.New(ulid.Timestamp(time.Now()), seed)
		baseID    = generate.RandomString(16)
		rp        = uint8(37) // rule not in rollout
		ctx       = Context{
			User: ContextUser{
				Age: uint8(rand.Intn(99)),
				ID:  userID,
			},
		}
		input = Parameters{
			"feature_x": false,
			"feature_y": false,
		}
		ids = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
			userID,
			generate.RandomString(24),
			generate.RandomString(24),
		}
	)

	r, err := New(
		ruleID.String(),
		baseID,
		generate.RandomString(12),
		generate.RandomString(12),
		KindRollout,
		true,
		Criteria{
			Criterion{
				Comparator: ComparatorIN,
				Key:        UserID,
				Value:      ids,
			},
		},
		[]Bucket{
			Bucket{
				Parameters: Parameters{
					"feature_x": true,
				},
			},
		},
		&rp,
	)

	if err != nil {
		t.Fatal(err)
	}

	_, _, err = r.Run(input, ctx, nil, randIntGenerateTest)
	if have, want := errors.Cause(err), errors.ErrRuleNotInRollout; have != want {
		t.Errorf("have %v, want %v", have, want)
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
		buckets = []Bucket{
			Bucket{
				Name: generate.RandomString(24),
				Parameters: Parameters{
					"feature_x": true,
				},
				Percentage: 100,
			},
		}
		criteria = Criteria{
			Criterion{
				Comparator: ComparatorIN,
				Key:        UserID,
				Value:      ids,
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
		KindOverride,
		startTime,
		endTime,
		buckets,
		criteria,
	)

	_, err = repo.Create(rule)
	if err != nil {
		t.Fatal(err)
	}

	r, err := repo.GetByID(rule.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := r.configID, rule.configID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := r.ID, rule.ID; have != want {
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

func testRepoGetByIDNotFound(t *testing.T, p prepareFunc) {
	_, err := p(t).GetByID(generate.RandomString(12))
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
		updateIds = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
		}
		buckets = []Bucket{
			Bucket{
				Name: generate.RandomString(24),
				Parameters: Parameters{
					"feature_x": true,
				},
				Percentage: 100,
			},
		}
		criteria = Criteria{
			Criterion{
				Comparator: ComparatorIN,
				Key:        UserID,
				Value:      ids,
			},
		}
		updatedCriteria = Criteria{
			Criterion{
				Comparator: ComparatorIN,
				Key:        UserID,
				Value:      updateIds,
			},
		}
		updateBuckets = []Bucket{
			Bucket{
				Name: generate.RandomString(24),
				Parameters: Parameters{
					"feature_x": true,
				},
				Percentage: 40,
			},
			Bucket{
				Name: generate.RandomString(24),
				Parameters: Parameters{
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
		KindOverride,
		startTime,
		endTime,
		buckets,
		criteria,
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
		KindOverride,
		startTime,
		endTime,
		buckets,
		criteria,
	)

	_, err = repo.Create(ruleTwo)
	if err != nil {
		t.Fatal(err)
	}

	updatedRule := generateRule(
		true,
		ruleOne.ID,
		ruleOne.configID,
		ruleOne.name,
		true,
		KindExperiment,
		time.Now().Add(200),
		time.Now().Add(2000),
		updateBuckets,
		updatedCriteria,
	)

	_, err = repo.UpdateWith(updatedRule)
	if err != nil {
		t.Fatal(err)
	}

	rl, err := repo.ListAll()
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(rl), 1; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testRepoCreateDuplicate(t *testing.T, p prepareFunc) {
	var (
		repo     = p(t)
		configID = generate.RandomString(24)
		name     = generate.RandomString(32)
		buckets  = []Bucket{
			Bucket{
				Name: generate.RandomString(24),
				Parameters: Parameters{
					"feature_x": true,
				},
			},
		}
	)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	rule, err := New(
		id.String(),
		configID,
		name, "",
		KindOverride,
		false,
		nil,
		buckets,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

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
		updateIds = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
		}
		buckets = []Bucket{
			Bucket{
				Name: generate.RandomString(24),
				Parameters: Parameters{
					"feature_x": true,
				},
				Percentage: 100,
			},
		}
		criteria = Criteria{
			Criterion{
				Comparator: ComparatorIN,
				Key:        UserID,
				Value:      ids,
			},
		}
		updateCriteria = Criteria{
			Criterion{
				Comparator: ComparatorIN,
				Key:        UserID,
				Value:      updateIds,
			},
		}
		updateBuckets = []Bucket{
			Bucket{
				Name: generate.RandomString(24),
				Parameters: Parameters{
					"feature_x": true,
				},
				Percentage: 40,
			},
			Bucket{
				Name: generate.RandomString(24),
				Parameters: Parameters{
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
		KindOverride,
		startTime,
		endTime,
		buckets,
		criteria,
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
		KindExperiment,
		time.Now().Add(200),
		time.Now().Add(2000),
		updateBuckets,
		updateCriteria,
	)

	updatedRule.activatedAt = time.Now().AddDate(0, -1, 0)

	ur, err := repo.UpdateWith(updatedRule)
	if err != nil {
		t.Fatal(err)
	}

	rl, err := repo.GetByID(updatedRule.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := ur.configID, rl.configID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := ur.ID, rl.ID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := ur.criteria, updateCriteria; !reflect.DeepEqual(have, want) {
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
	repo := p(t)

	rl, err := repo.ListAll()
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
		buckets = []Bucket{
			Bucket{
				Name: generate.RandomString(24),
				Parameters: Parameters{
					"feature_x": true,
				},
				Percentage: 100,
			},
		}
		criteria = Criteria{
			Criterion{
				Comparator: ComparatorIN,
				Key:        UserID,
				Value:      ids,
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
		KindOverride,
		startTime,
		endTime,
		buckets,
		criteria,
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
		KindOverride,
		startTime,
		endTime,
		buckets,
		criteria,
	)

	_, err = repo.Create(ruleTwo)
	if err != nil {
		t.Fatal(err)
	}

	rl, err := repo.ListAll()
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
		buckets = []Bucket{
			Bucket{
				Name: generate.RandomString(24),
				Parameters: Parameters{
					"feature_x": true,
				},
				Percentage: 100,
			},
		}
		criteria = Criteria{
			Criterion{
				Comparator: ComparatorIN,
				Key:        UserID,
				Value:      ids,
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
		KindOverride,
		startTime,
		endTime,
		buckets,
		criteria,
	)

	_, err = repo.Create(ruleOne)
	if err != nil {
		t.Fatal(err)
	}

	id, err = ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	var zeroTime time.Time

	ruleTwo := generateRule(
		true,
		id.String(),
		configID,
		nameRuleTwo,
		false,
		KindOverride,
		zeroTime,
		zeroTime,
		buckets,
		criteria,
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
		KindOverride,
		startTime,
		endTime,
		buckets,
		criteria,
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

func testRepoNoCriteria(t *testing.T, p prepareFunc) {
	var (
		repo      = p(t)
		configID  = generate.RandomString(24)
		nameRule  = generate.RandomString(32)
		endTime   = time.Now().AddDate(0, 1, 0)
		startTime = time.Now().AddDate(0, -1, 0)
		buckets   = []Bucket{
			Bucket{
				Name: generate.RandomString(24),
				Parameters: Parameters{
					"feature_x": true,
				},
				Percentage: 100,
			},
		}
		criteria = Criteria{}
	)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	r := generateRule(
		true,
		id.String(),
		configID,
		nameRule,
		false,
		KindOverride,
		startTime,
		endTime,
		buckets,
		criteria,
	)

	_, err = repo.Create(r)
	if err != nil {
		t.Fatal(err)
	}

	rl, err := repo.ListActive(configID, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(rl), 1; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	for _, rule := range rl {
		if have, want := rule.ID, r.ID; have != want {
			t.Errorf("have %v, want %v", have, want)
		}
	}
}

func testRepoCreateRollout(t *testing.T, p prepareFunc) {
	var (
		repo     = p(t)
		configID = generate.RandomString(24)
		name     = generate.RandomString(32)
		buckets  = []Bucket{
			Bucket{
				Name: generate.RandomString(24),
				Parameters: Parameters{
					"feature_x": true,
				},
			},
		}
	)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	rule, err := New(
		id.String(),
		configID,
		name,
		"",
		KindRollout,
		false,
		nil,
		buckets,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	rule.rollout = 57

	created, err := repo.Create(rule)
	if err != nil {
		t.Fatal(err)
	}

	retrieved, err := repo.GetByID(rule.ID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := retrieved.kind, created.kind; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := retrieved.rollout, created.rollout; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func generateRule(
	active bool,
	id, configID, name string,
	deleted bool,
	kind Kind,
	startTime, endTime time.Time,
	buckets []Bucket,
	criteria Criteria,
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
		ID:          id,
		kind:        kind,
		name:        name,
		startTime:   startTime,
	}
}
