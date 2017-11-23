package config

import (
	"time"

	"github.com/lifesum/configsum/pkg/instrument"
)

const (
	labelBaseRepo = "base"
	labelUserRepo = "user"
)

type instrumentBaseRepo struct {
	next      BaseRepo
	opObserve instrument.ObserveRepoFunc
	store     string
}

// NewBaseRepoInstrumentMiddleware wraps the next BaseRepo and add Prometheus
// instrumentation capabilities.
func NewBaseRepoInstrumentMiddleware(
	opObserve instrument.ObserveRepoFunc,
	store string,
) BaseRepoMiddleware {
	return func(next BaseRepo) BaseRepo {
		return &instrumentBaseRepo{
			next:      next,
			opObserve: opObserve,
			store:     store,
		}
	}
}

func (r *instrumentBaseRepo) Create(
	id, clientID, name string,
	parameters rendered,
) (c BaseConfig, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelBaseRepo, "Create", begin, err)
	}(time.Now())

	return r.next.Create(id, clientID, name, parameters)
}

func (r *instrumentBaseRepo) GetByID(id string) (c BaseConfig, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelBaseRepo, "GetByID", begin, err)
	}(time.Now())

	return r.next.GetByID(id)
}

func (r *instrumentBaseRepo) GetByName(clientID, name string) (c BaseConfig, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelBaseRepo, "GetByName", begin, err)
	}(time.Now())

	return r.next.GetByName(clientID, name)
}

func (r *instrumentBaseRepo) List() (l BaseList, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelBaseRepo, "List", begin, err)
	}(time.Now())

	return r.next.List()
}

func (r *instrumentBaseRepo) Update(input BaseConfig) (bc BaseConfig, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelBaseRepo, "Update", begin, err)
	}(time.Now())

	return r.next.Update(input)
}

func (r *instrumentBaseRepo) setup() (err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelBaseRepo, "setup", begin, err)
	}(time.Now())

	return r.next.setup()
}

func (r *instrumentBaseRepo) teardown() (err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelBaseRepo, "teardown", begin, err)
	}(time.Now())

	return r.next.teardown()
}

type instrumentUserRepo struct {
	next      UserRepo
	opObserve instrument.ObserveRepoFunc
	store     string
}

// NewUserRepoInstrumentMiddleware wraps the next BaseRepo and add Prometheus
// instrumentation capabilities.
func NewUserRepoInstrumentMiddleware(
	opObserve instrument.ObserveRepoFunc,
	store string,
) UserRepoMiddleware {
	return func(next UserRepo) UserRepo {
		return &instrumentUserRepo{
			next:      next,
			opObserve: opObserve,
			store:     store,
		}
	}
}

func (r *instrumentUserRepo) Append(
	id, baseID, userID string,
	decisions ruleDecisions,
	render rendered,
) (c UserConfig, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelUserRepo, "Append", begin, err)
	}(time.Now())

	return r.next.Append(id, baseID, userID, decisions, render)
}

func (r *instrumentUserRepo) GetLatest(
	baseID, userID string,
) (c UserConfig, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelUserRepo, "GetLatest", begin, err)
	}(time.Now())

	return r.next.GetLatest(baseID, userID)
}

func (r *instrumentUserRepo) setup() (err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelUserRepo, "Setup", begin, err)
	}(time.Now())

	return r.next.setup()
}

func (r *instrumentUserRepo) teardown() (err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelUserRepo, "Teardown", begin, err)
	}(time.Now())

	return r.next.teardown()
}
