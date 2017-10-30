package config

import (
	"time"

	"github.com/lifesum/configsum/pkg/instrument"
)

const labelRepoUser = "user"

type instrumentUserRepo struct {
	opObserve instrument.ObserveRepoFunc
	next      UserRepo
	store     string
}

// NewUserRepoInstrumentMiddleware wraps the next UserRepo Prometheus
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
		r.opObserve(r.store, labelRepoUser, "Append", begin, err)
	}(time.Now())

	return r.next.Append(id, baseID, userID, decisions, render)
}

func (r *instrumentUserRepo) GetLatest(
	baseID, userID string,
) (c UserConfig, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRepoUser, "GetLatest", begin, err)
	}(time.Now())

	return r.next.GetLatest(baseID, userID)
}

func (r *instrumentUserRepo) Setup() (err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRepoUser, "Setup", begin, err)
	}(time.Now())

	return r.next.Setup()
}

func (r *instrumentUserRepo) Teardown() (err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRepoUser, "Teardown", begin, err)
	}(time.Now())

	return r.next.Teardown()
}
