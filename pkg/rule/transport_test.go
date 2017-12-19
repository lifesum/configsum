package rule

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/oklog/ulid"

	"github.com/lifesum/configsum/pkg/generate"
)

func randIntGenTest() int {
	return 61
}

func TestRuleActivate(t *testing.T) {
	var (
		configID = generate.RandomString(12)
		repo     = preparePGRepo(t)
		svc      = NewService(repo)
		id, _    = ulid.New(ulid.Timestamp(time.Now()), seed)
		target   = fmt.Sprintf("/%s/activate", id.String())
		req      = httptest.NewRequest("PUT", target, nil)
		rec      = httptest.NewRecorder()
		r        = MakeHandler(svc)
	)

	rule, err := New(
		id.String(),
		configID,
		generate.RandomString(12),
		generate.RandomString(42),
		KindOverride,
		false,
		nil,
		[]Bucket{
			{
				Name: "default",
				Parameters: Parameters{
					"feature_funky_toggle": true,
				},
			},
		},
		nil,
		randIntGenTest,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Create(rule)
	if err != nil {
		t.Fatal(err)
	}

	r.ServeHTTP(rec, req)

	if have, want := rec.Code, http.StatusNoContent; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	updated, err := repo.GetByID(id.String())
	if err != nil {
		t.Fatal(err)
	}

	if have, want := updated.active, true; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	// Check for idempotency.
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if have, want := rec.Code, http.StatusNoContent; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestRuleDeactivate(t *testing.T) {
	var (
		configID = generate.RandomString(12)
		repo     = preparePGRepo(t)
		svc      = NewService(repo)
		id, _    = ulid.New(ulid.Timestamp(time.Now()), seed)
		target   = fmt.Sprintf("/%s/deactivate", id.String())
		req      = httptest.NewRequest("PUT", target, nil)
		rec      = httptest.NewRecorder()
		r        = MakeHandler(svc)
	)

	rule, err := New(
		id.String(),
		configID,
		generate.RandomString(12),
		generate.RandomString(42),
		KindOverride,
		true,
		nil,
		[]Bucket{
			{
				Name: "default",
				Parameters: Parameters{
					"feature_funky_toggle": true,
				},
			},
		},
		nil,
		randIntGenTest,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Create(rule)
	if err != nil {
		t.Fatal(err)
	}

	r.ServeHTTP(rec, req)

	if have, want := rec.Code, http.StatusNoContent; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	updated, err := repo.GetByID(id.String())
	if err != nil {
		t.Fatal(err)
	}

	if have, want := updated.active, false; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	// Check for idempotency.
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if have, want := rec.Code, http.StatusNoContent; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestRuleGet(t *testing.T) {
	var (
		configID = generate.RandomString(12)
		repo     = preparePGRepo(t)
		svc      = NewService(repo)
		id, _    = ulid.New(ulid.Timestamp(time.Now()), seed)
		target   = fmt.Sprintf("/%s", id.String())
		req      = httptest.NewRequest("GET", target, nil)
		rec      = httptest.NewRecorder()
		r        = MakeHandler(svc)
		locale   = MatcherString("en_GB")
	)

	rule, err := New(
		id.String(),
		configID,
		"override_funky_staff",
		"Overrides funky feature for all staff memebers",
		KindOverride,
		true,
		&Criteria{
			Locale: &locale,
			User: &CriteriaUser{
				ID: &MatcherStringList{
					generate.RandomString(12),
					generate.RandomString(12),
					generate.RandomString(12),
				},
			},
		},
		[]Bucket{
			{
				Name: "default",
				Parameters: Parameters{
					"feature_funky_toggle": true,
				},
			},
		},
		nil,
		randIntGenTest,
	)
	if err != nil {
		t.Fatal(err)
	}

	created, err := repo.Create(rule)
	if err != nil {
		t.Fatal(err)
	}

	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	resp := responseRule{}

	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	if have, want := resp.rule.active, created.active; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := resp.rule.activatedAt, created.activatedAt; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := resp.rule.buckets, created.buckets; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := resp.rule.configID, created.configID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := resp.rule.criteria.Locale, created.criteria.Locale; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := resp.rule.criteria.User, created.criteria.User; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := resp.rule.description, created.description; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := resp.rule.deleted, created.deleted; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := resp.rule.endTime, created.endTime; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := resp.rule.ID, created.ID; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := resp.rule.kind, created.kind; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := resp.rule.name, created.name; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := resp.rule.rollout, created.rollout; have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	if have, want := resp.rule.startTime, created.startTime; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestRuleList(t *testing.T) {
	var (
		seed     = rand.New(rand.NewSource(time.Now().UnixNano()))
		numRules = seed.Intn(24) + seed.Intn(6)
		configID = generate.RandomString(12)
		repo     = preparePGRepo(t)
		svc      = NewService(repo)
		req      = httptest.NewRequest("GET", "/", nil)
		rec      = httptest.NewRecorder()
		r        = MakeHandler(svc)
	)

	r.ServeHTTP(rec, req)

	if have, want := rec.Code, http.StatusNoContent; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	for i := 0; i < numRules; i++ {
		id, err := ulid.New(ulid.Timestamp(time.Now()), seed)
		if err != nil {
			t.Fatal(err)
		}

		rule, err := New(
			id.String(),
			configID,
			generate.RandomString(12),
			generate.RandomString(42),
			KindOverride,
			true,
			nil,
			[]Bucket{
				{
					Name: "default",
					Parameters: Parameters{
						"feature_funky_toggle": true,
					},
				},
			},
			nil,
			randIntGenTest,
		)
		if err != nil {
			t.Fatal(err)
		}

		_, err = repo.Create(rule)
		if err != nil {
			t.Fatal(err)
		}
	}

	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if have, want := rec.Code, http.StatusOK; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	resp := responseList{}

	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	if have, want := len(resp.rules), numRules; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestRuleUpdateRollout(t *testing.T) {
	var (
		configID = generate.RandomString(12)
		repo     = preparePGRepo(t)
		svc      = NewService(repo)
		id, _    = ulid.New(ulid.Timestamp(time.Now()), seed)
		payload  = bytes.NewBufferString(`{"rollout": 13}`)
		target   = fmt.Sprintf("/%s/rollout", id.String())
		req      = httptest.NewRequest("PUT", target, payload)
		rec      = httptest.NewRecorder()
		r        = MakeHandler(svc)
	)

	rule, err := New(
		id.String(),
		configID,
		generate.RandomString(12),
		generate.RandomString(42),
		KindOverride,
		true,
		nil,
		[]Bucket{
			{
				Name: "default",
				Parameters: Parameters{
					"feature_funky_toggle": true,
				},
			},
		},
		nil,
		randIntGenTest,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Create(rule)
	if err != nil {
		t.Fatal(err)
	}

	r.ServeHTTP(rec, req)

	if have, want := rec.Code, http.StatusNoContent; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	updated, err := repo.GetByID(id.String())
	if err != nil {
		t.Fatal(err)
	}

	if have, want := updated.rollout, uint8(13); have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	// Check for idempotency.
	payload = bytes.NewBufferString(`{"rollout": 13}`)
	req = httptest.NewRequest("PUT", target, payload)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if have, want := rec.Code, http.StatusNoContent; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestDecodeGetRequest(t *testing.T) {
	var (
		seed   = rand.New(rand.NewSource(time.Now().UnixNano()))
		id, _  = ulid.New(ulid.Timestamp(time.Now()), seed)
		ctx    = context.WithValue(context.Background(), varID, id.String())
		target = fmt.Sprintf("/%s", id.String())
		r      = httptest.NewRequest("GET", target, nil)
	)

	raw, err := decodeGetRequest(ctx, r)
	if err != nil {
		t.Fatal(err)
	}

	want := getRequest{id: id.String()}

	if have := raw.(getRequest); !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestExtractMuxVars(t *testing.T) {
	var (
		key = muxVar("testKey")
		val = generate.RandomString(12)
		req = httptest.NewRequest("GET", fmt.Sprintf("/root/%s", val), nil)
		r   = mux.NewRouter()
	)

	r.Methods("GET").Path(`/root/{testKey}`).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := extractMuxVars(key)(context.Background(), r)

		if have, want := ctx.Value(key), val; have != want {
			t.Errorf("have %v, want %v", have, want)
		}
	})

	r.ServeHTTP(httptest.NewRecorder(), req)
}
