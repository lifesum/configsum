package client

import (
	"time"

	"github.com/lifesum/configsum/pkg/instrument"
)

const (
	labelRepo      = "client"
	labelRepoToken = "token"
)

type instrumentRepo struct {
	opObserve instrument.ObserveRepoFunc
	next      Repo
	store     string
}

// NewRepoInstrumentMiddleware wraps the next Repo with Prometheus
// instrumenation capabilities.
func NewRepoInstrumentMiddleware(
	opObserve instrument.ObserveRepoFunc,
	store string,
) RepoMiddleware {
	return func(next Repo) Repo {
		return &instrumentRepo{
			next:      next,
			opObserve: opObserve,
			store:     store,
		}
	}
}

func (r *instrumentRepo) List() (cs List, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRepo, "List", begin, err)
	}(time.Now())

	return r.next.List()
}

func (r *instrumentRepo) Lookup(id string) (client Client, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRepo, "Lookup", begin, err)
	}(time.Now())

	return r.next.Lookup(id)
}

func (r *instrumentRepo) Store(id, name string) (client Client, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRepo, "Store", begin, err)
	}(time.Now())

	return r.next.Store(id, name)
}

func (r *instrumentRepo) Setup() (err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRepo, "Setup", begin, err)
	}(time.Now())

	return r.next.Setup()
}

func (r *instrumentRepo) Teardown() (err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRepo, "Teardown", begin, err)
	}(time.Now())

	return r.next.Teardown()
}

type instrumentTokenRepo struct {
	opObserve instrument.ObserveRepoFunc
	next      TokenRepo
	store     string
}

// NewTokenRepoInstrumentMiddleware wraps the next TokenRepo with Prometheus
// instrumenation capabilities.
func NewTokenRepoInstrumentMiddleware(
	opObserve instrument.ObserveRepoFunc,
	store string,
) TokenRepoMiddleware {
	return func(next TokenRepo) TokenRepo {
		return &instrumentTokenRepo{
			next:      next,
			opObserve: opObserve,
			store:     store,
		}
	}
}

func (r *instrumentTokenRepo) GetLatest(clientID string) (token Token, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRepoToken, "GetLatest", begin, err)
	}(time.Now())

	return r.next.GetLatest(clientID)
}

func (r *instrumentTokenRepo) Lookup(secret string) (token Token, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRepoToken, "Lookup", begin, err)
	}(time.Now())

	return r.next.Lookup(secret)
}

func (r *instrumentTokenRepo) Store(clientID, secret string) (token Token, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRepoToken, "Store", begin, err)
	}(time.Now())

	return r.next.Store(clientID, secret)
}

func (r *instrumentTokenRepo) Setup() (err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRepoToken, "Setup", begin, err)
	}(time.Now())

	return r.next.Setup()
}

func (r *instrumentTokenRepo) Teardown() (err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRepoToken, "Teardown", begin, err)
	}(time.Now())

	return r.next.Teardown()
}
