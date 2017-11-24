package config

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/oklog/ulid"

	"github.com/lifesum/configsum/pkg/errors"
	"github.com/lifesum/configsum/pkg/generate"
)

type prepareBaseRepoFunc func(t *testing.T) BaseRepo

type prepareUserRepoFunc func(t *testing.T) UserRepo

func testBaseRepoCreateDuplicate(t *testing.T, p prepareBaseRepoFunc) {
	var (
		clientID   = generate.RandomString(24)
		name       = generate.RandomString(12)
		parameters = rendered{
			"feature_awesome-sauce_toggle": true,
		}
		repo = p(t)
		seed = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Create(id.String(), clientID, name, parameters)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Create(id.String(), clientID, name, parameters)
	if have, want := errors.Cause(err), errors.ErrExists; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testBaseRepoGetByID(t *testing.T, p prepareBaseRepoFunc) {
	var (
		clientID   = generate.RandomString(24)
		name       = generate.RandomString(12)
		parameters = rendered{
			"feature_awesome-sauce_toggle": true,
		}
		repo = p(t)
		seed = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Create(id.String(), clientID, name, parameters)
	if err != nil {
		t.Fatal(err)
	}

	c, err := repo.GetByID(id.String())
	if err != nil {
		t.Fatal(err)
	}

	if have, want := c.ClientID, clientID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := c.ID, id.String(); have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := c.Name, name; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := c.Parameters, parameters; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testBaseRepoGetByIDNotFound(t *testing.T, p prepareBaseRepoFunc) {
	repo := p(t)

	_, err := repo.GetByID(generate.RandomString(24))
	if have, want := errors.Cause(err), errors.ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testBaseRepoGetByName(t *testing.T, p prepareBaseRepoFunc) {
	var (
		clientID   = generate.RandomString(24)
		name       = generate.RandomString(12)
		parameters = rendered{
			"feature_awesome-sauce_toggle": true,
		}
		repo = p(t)
		seed = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Create(id.String(), clientID, name, parameters)
	if err != nil {
		t.Fatal(err)
	}

	c, err := repo.GetByName(clientID, name)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := c.ClientID, clientID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := c.ID, id.String(); have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := c.Name, name; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := c.Parameters, parameters; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testBaseRepoGetByNameNotFound(t *testing.T, p prepareBaseRepoFunc) {
	var (
		clientID = generate.RandomString(24)
		name     = generate.RandomString(12)
		repo     = p(t)
	)

	_, err := repo.GetByName(clientID, name)
	if have, want := errors.Cause(err), errors.ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testBaseRepoList(t *testing.T, p prepareBaseRepoFunc) {
	var (
		repo       = p(t)
		seed       = rand.New(rand.NewSource(time.Now().UnixNano()))
		numConfigs = seed.Intn(24)

		es = BaseList{}
	)

	for i := 0; i < numConfigs; i++ {
		var (
			clientID   = generate.RandomString(24)
			name       = generate.RandomString(12)
			parameters = rendered{
				generate.RandomString(8):  false,
				generate.RandomString(12): float64(seed.Intn(64)),
				generate.RandomString(16): generate.RandomString(24),
			}
		)

		id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
		if err != nil {
			t.Fatal(err)
		}

		c, err := repo.Create(id.String(), clientID, name, parameters)
		if err != nil {
			t.Fatal(err)
		}

		es = append(es, c)
	}

	cs, err := repo.List()
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(cs), numConfigs; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	sort.Sort(es)

	for i, expect := range es {
		c := cs[i]

		if have, want := c.ClientID, expect.ClientID; have != want {
			t.Errorf("have %v, want %v", have, want)
		}

		if have, want := c.ID, expect.ID; have != want {
			t.Errorf("have %v, want %v", have, want)
		}

		if have, want := c.Name, expect.Name; have != want {
			t.Errorf("have %v, want %v", have, want)
		}

		if have, want := c.Parameters, expect.Parameters; !reflect.DeepEqual(have, want) {
			t.Errorf("\nhave %#v\nwant %#v", have, want)
		}
	}
}

func testBaseRepoUpdate(t *testing.T, p prepareBaseRepoFunc) {
	var (
		clientID   = generate.RandomString(24)
		name       = generate.RandomString(12)
		parameters = rendered{
			"feature_awesome-sauce_toggle": true,
		}
		repo = p(t)
		seed = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Update(BaseConfig{ID: id.String()})
	if have, want := errors.Cause(err), errors.ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	c, err := repo.Create(id.String(), clientID, name, parameters)
	if err != nil {
		t.Fatal(err)
	}

	newParams := rendered{
		"feature_awesome-sauce_toggle": true,
		"feature_awesome-sauce_desc":   generate.RandomString(24),
	}

	_, err = repo.Update(BaseConfig{
		ClientID:   clientID,
		Deleted:    false,
		ID:         id.String(),
		Name:       name,
		Parameters: newParams,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	})
	if err != nil {
		t.Fatal(err)
	}

	c, err = repo.GetByID(id.String())
	if err != nil {
		t.Fatal(err)
	}

	if have, want := c.Parameters, newParams; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testUserRepoGetLatest(t *testing.T, p prepareUserRepoFunc) {
	var (
		baseID    = generate.RandomString(24)
		userID    = generate.RandomString(24)
		render    = rendered{}
		repo      = p(t)
		seed      = rand.New(rand.NewSource(time.Now().UnixNano()))
		decisions = ruleDecisions{
			generate.RandomString(24): []int{seed.Intn(100), seed.Intn(100)},
			generate.RandomString(24): []int{seed.Intn(100)},
			generate.RandomString(24): []int{},
		}
	)

	render.setNumber(generate.RandomString(24), seed.Float64())

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Append(
		generate.RandomString(24),
		generate.RandomString(24),
		generate.RandomString(24),
		decisions,
		render,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Append(
		generate.RandomString(24),
		baseID,
		userID,
		nil,
		rendered{},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Append(id.String(), baseID, userID, decisions, render)
	if err != nil {
		t.Fatal(err)
	}

	c, err := repo.GetLatest(baseID, userID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := c.baseID, baseID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := c.id, id.String(); have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := c.userID, userID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := c.rendered, render; !reflect.DeepEqual(c.rendered, render) {
		t.Errorf("\nhave %#v,\nwant %#v", have, want)
	}

	if have, want := c.ruleDecisions, decisions; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testUserRepoGetLatestNotFound(t *testing.T, p prepareUserRepoFunc) {
	var (
		baseID = generate.RandomString(24)
		userID = generate.RandomString(24)
		repo   = p(t)
	)

	_, err := repo.GetLatest(baseID, userID)
	if have, want := errors.Cause(err), errors.ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func testUserRepoAppendDuplicate(t *testing.T, p prepareUserRepoFunc) {
	var (
		baseID    = generate.RandomString(24)
		userID    = generate.RandomString(24)
		render    = rendered{}
		repo      = p(t)
		seed      = rand.New(rand.NewSource(time.Now().UnixNano()))
		decisions = ruleDecisions{
			generate.RandomString(32): []int{seed.Int(), seed.Int()},
			generate.RandomString(32): []int{seed.Int()},
			generate.RandomString(32): []int{},
		}
	)

	render.SetBool(generate.RandomString(24), false)

	id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Append(id.String(), baseID, userID, decisions, render)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Append(id.String(), baseID, userID, decisions, render)
	if have, want := errors.Cause(err), errors.ErrExists; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}
