package rule

import (
	"time"

	"github.com/lifesum/configsum/pkg/instrument"
)

const (
	labelRuleRepo = "rule"
)

type instrumentRuleRepo struct {
	next      Repo
	opObserve instrument.ObserveRepoFunc
	store     string
}

// NewRuleRepoInstrumentMiddleware wraps the next Repo and adds Prometheus
// instrumentation capabilities.
func NewRuleRepoInstrumentMiddleware(
	opObserve instrument.ObserveRepoFunc,
	store string,
) RepoMiddleware {
	return func(next Repo) Repo {
		return &instrumentRuleRepo{
			next:      next,
			opObserve: opObserve,
			store:     store,
		}
	}
}

func (r *instrumentRuleRepo) Create(input Rule) (rl Rule, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRuleRepo, "Create", begin, err)
	}(time.Now())

	return r.next.Create(input)
}

func (r *instrumentRuleRepo) GetByID(id string) (rl Rule, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRuleRepo, "GetByID", begin, err)
	}(time.Now())

	return r.next.GetByID(id)
}

func (r *instrumentRuleRepo) UpdateWith(input Rule) (rl Rule, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRuleRepo, "UpdateWith", begin, err)
	}(time.Now())

	return r.next.UpdateWith(input)
}

func (r *instrumentRuleRepo) ListAll() (rs []Rule, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRuleRepo, "ListAll", begin, err)
	}(time.Now())

	return r.next.ListAll()
}

func (r *instrumentRuleRepo) ListActive(
	configID string,
	now time.Time,
) (rls []Rule, err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRuleRepo, "ListActive", begin, err)
	}(time.Now())

	return r.next.ListActive(configID, now)
}

func (r *instrumentRuleRepo) Setup() (err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRuleRepo, "Setup", begin, err)
	}(time.Now())

	return r.next.Setup()
}

func (r *instrumentRuleRepo) Teardown() (err error) {
	defer func(begin time.Time) {
		r.opObserve(r.store, labelRuleRepo, "Teardown", begin, err)
	}(time.Now())

	return r.next.Teardown()
}
