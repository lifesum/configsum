package config

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/oklog/ulid"
	"golang.org/x/text/language"

	"github.com/lifesum/configsum/pkg/auth"
	"github.com/lifesum/configsum/pkg/client"
	"github.com/lifesum/configsum/pkg/generate"
	"github.com/lifesum/configsum/pkg/rule"
)

func TestDecodeBaseUpdateRequest(t *testing.T) {
	var (
		seed    = rand.New(rand.NewSource(time.Now().UnixNano()))
		id, _   = ulid.New(ulid.Timestamp(time.Now()), seed)
		ctx     = context.WithValue(context.Background(), varID, id.String())
		payload = bytes.NewBufferString(`{"parameters": {"feature_decode_toggled": true}}`)
		target  = fmt.Sprintf("/%s", id)
		r       = httptest.NewRequest("PUT", target, payload)
	)

	raw, err := decodeBaseUpdateRequest(ctx, r)
	if err != nil {
		t.Fatal(err)
	}

	want := baseUpdateRequest{
		id: id.String(),
		parameters: rule.Parameters{
			"feature_decode_toggled": true,
		},
	}

	if have := raw.(baseUpdateRequest); !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave %#v\nwant %#v", have, want)
	}
}

func TestDecodeUserRenderRequest(t *testing.T) {
	var (
		baseConfig = generate.RandomString(6)
		ctx        = context.WithValue(context.Background(), varBaseConfig, baseConfig)
		locale     = language.MustParse("en_GB")
		payload    = bytes.NewBufferString(`{"device": {"location": {"locale": "en_GB"}}}`)
		target     = fmt.Sprintf("/%s", baseConfig)
		r          = httptest.NewRequest("PUT", target, payload)
	)

	raw, err := decodeUserRenderRequest(ctx, r)
	if err != nil {
		t.Fatal(err)
	}

	want := userRenderRequest{
		baseConfig: baseConfig,
		context: userRenderContext{
			Device: device{
				Location: location{
					locale: locale,
				},
			},
		},
	}

	if have := raw.(userRenderRequest); !reflect.DeepEqual(have, want) {
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

func TestUserRender(t *testing.T) {
	var (
		baseID     = generate.RandomString(16)
		baseName   = "some-base-config-4473"
		baseRepo   = preparePGBaseRepo(t)
		clientID   = generate.RandomString(12)
		userID     = generate.RandomString(12)
		us         = language.MustParseRegion("US")
		en         = language.MustParseBase("en")
		paramKey   = generate.RandomString(6)
		payload    = bytes.NewBufferString(`{"app": { "version": "8.6.7" }, "device": {"location": {"locale": "en_US", "timezoneOffset": 3600}, "os": { "platform": "iOS", "version": "8.0"}}, "user": { "subscription": 1 } }`)
		target     = fmt.Sprintf("/%s", baseName)
		parameters = rule.Parameters{
			paramKey: false,
		}
		req       = httptest.NewRequest("PUT", target, payload)
		rec       = httptest.NewRecorder()
		ruleRepo  = prepareRuleRepo(t)
		userRepo  = preparePGUserRepo(t)
		seed      = rand.New(rand.NewSource(time.Now().UnixNano()))
		svc       = NewUserService(baseRepo, userRepo, ruleRepo, generate.RandPercentage(seed))
		ruleID, _ = ulid.New(ulid.Timestamp(time.Now()), seed)
		router    = MakeHandler(svc, injectAuth(clientID, userID))
	)

	_, err := baseRepo.Create(baseID, clientID, baseName, parameters)
	if err != nil {
		t.Fatal(err)
	}

	enUS, err := language.Compose(en, us)
	if err != nil {
		t.Fatal(err)
	}

	override, err := rule.New(
		ruleID.String(),
		baseID,
		generate.RandomString(12),
		generate.RandomString(12),
		rule.KindOverride,
		true,
		rule.Criteria{
			rule.Criterion{
				Comparator: rule.ComparatorEQ,
				Key:        rule.DeviceLocationLocale,
				Value:      enUS,
			},
			rule.Criterion{
				Comparator: rule.ComparatorGT,
				Key:        rule.UserSubscription,
				Value:      0,
			},
		},
		[]rule.Bucket{
			{
				Name: "someBucketName",
				Parameters: rule.Parameters{
					paramKey: true,
				},
			},
		},
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	_, err = ruleRepo.Create(override)

	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(rec, req)

	t.Log(string(rec.Body.Bytes()))

	if have, want := rec.Code, http.StatusOK; have != want {
		t.Fatalf("have %v, want %v", have, want)
	}

	uc, err := userRepo.GetLatest(baseID, userID)
	if err != nil {
		t.Fatal(err)
	}

	want := rule.Parameters{
		paramKey: true,
	}

	if have := uc.rendered; !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func injectAuth(clientID, userID string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			ctx = context.WithValue(ctx, client.ContextKeyClientID, clientID)
			ctx = context.WithValue(ctx, auth.ContextKeyUserID, userID)

			return next(ctx, request)
		}
	}
}
