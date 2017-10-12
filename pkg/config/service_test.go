package config

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestServiceUserRender(t *testing.T) {
	var (
		appID      = randString(characterSet)
		baseID     = randString(characterSet)
		baseName   = randString(characterSet)
		baseRender = rendered{
			randString(characterSet): false,
		}
		baseRepo, _ = NewInmemBaseRepo(inmemBaseState{
			appID: map[string]BaseConfig{
				baseName: BaseConfig{
					appID:    appID,
					id:       baseID,
					name:     baseName,
					rendered: baseRender,
				},
			},
		})
		userID      = randString(characterSet)
		userRepo, _ = NewInmemUserRepo()
		svc         = NewServiceUser(baseRepo, userRepo)
	)

	uc, err := svc.Render(appID, baseName, userID)
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

	rc, err := svc.Render(appID, baseName, userID)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := uc, rc; !reflect.DeepEqual(have, want) {
		t.Errorf("have %#v,want %#v", have, want)
	}
}

func TestServiceUserRenderConfigMissingBaseConfig(t *testing.T) {
	var (
		appID       = randString(characterSet)
		baseName    = randString(characterSet)
		baseRepo, _ = NewInmemBaseRepo(nil)
		userID      = randString(characterSet)
		userRepo, _ = NewInmemUserRepo()
		svc         = NewServiceUser(baseRepo, userRepo)
	)

	_, err := svc.Render(appID, baseName, userID)
	if have, want := errors.Cause(err), ErrNotFound; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}
