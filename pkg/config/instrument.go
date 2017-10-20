package config

import (
	"time"

	"github.com/lifesum/configsum/pkg/instrument"
)

const labelRepoUser = "user"

type instrumentUserRepo struct {
	errCount  instrument.CountRepoFunc
	opCount   instrument.CountRepoFunc
	opObserve instrument.ObserveRepoFunc
	next      UserRepo
	store     string
}

// NewUserRepoInstrumentMiddleware wraps the next UserRepo Prometheus
// instrumentation capabilities.
func NewUserRepoInstrumentMiddleware(
	errCount instrument.CountRepoFunc,
	opCount instrument.CountRepoFunc,
	opObserve instrument.ObserveRepoFunc,
	store string,
) UserRepoMiddleware {
	return func(next UserRepo) UserRepo {
		return &instrumentUserRepo{
			errCount:  errCount,
			next:      next,
			opCount:   opCount,
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
		r.track(begin, err, "Append")
	}(time.Now())

	return r.next.Append(id, baseID, userID, decisions, render)
}

func (r *instrumentUserRepo) GetLatest(
	baseID, userID string,
) (c UserConfig, err error) {
	defer func(begin time.Time) {
		r.track(begin, err, "GetLatest")
	}(time.Now())

	return r.next.GetLatest(baseID, userID)
}

func (r *instrumentUserRepo) Setup() (err error) {
	defer func(begin time.Time) {
		r.track(begin, err, "Setup")
	}(time.Now())

	return r.next.Setup()
}

func (r *instrumentUserRepo) Teardown() (err error) {
	defer func(begin time.Time) {
		r.track(begin, err, "Teardown")
	}(time.Now())

	return r.next.Teardown()
}

func (r *instrumentUserRepo) track(begin time.Time, err error, op string) {
	if err != nil {
		r.errCount(r.store, labelRepoUser, op)

		return
	}

	r.opCount(r.store, labelRepoUser, op)
	r.opObserve(r.store, labelRepoUser, op, begin)
}
