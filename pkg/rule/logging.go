package rule

import (
	"time"

	"github.com/go-kit/kit/log"
)

// Log fields.
const (
	logBuckets   = "buckets"
	logConfigID  = "configID"
	logCreatedAt = "createdAt"
	logDuration  = "duration"
	logElements  = "elements"
	logEndTime   = "endTime"
	logErr       = "err"
	logID        = "id"
	logKind      = "kind"
	logName      = "name"
	logOp        = "op"
	logPkg       = "pkg"
	logRepo      = "repo"
	logStartTime = "startTime"
	logStore     = "store"
	logUpdatedAt = "updatedAt"
)

type logRuleRepo struct {
	logger log.Logger
	next   Repo
}

// NewRuleRepoLogMiddleware wraps the next Repo with logging capabilities.
func NewRuleRepoLogMiddleware(logger log.Logger, store string) RepoMiddleware {
	return func(next Repo) Repo {
		return &logRuleRepo{
			logger: log.With(
				logger,
				logPkg, "rule",
				logRepo, "repo",
				logStore, store,
			),
			next: next,
		}
	}
}

func (r *logRuleRepo) Create(input Rule) (rl Rule, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logBuckets, input.buckets,
			logConfigID, input.configID,
			logCreatedAt, input.createdAt,
			logDuration, time.Since(begin).Nanoseconds(),
			logEndTime, input.endTime,
			logID, input.ID,
			logKind, input.kind,
			logName, input.name,
			logOp, "Create",
			logStartTime, input.startTime,
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.Create(input)
}

func (r *logRuleRepo) GetByName(configID, name string) (rl Rule, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logDuration, time.Since(begin).Nanoseconds(),
			logConfigID, configID,
			logName, name,
			logOp, "GetByName",
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.GetByName(configID, name)
}

func (r *logRuleRepo) UpdateWith(input Rule) (rl Rule, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logDuration, time.Since(begin).Nanoseconds(),
			logName, input.name,
			logOp, "UpdatedWith",
			logUpdatedAt, input.updatedAt,
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.UpdateWith(input)
}

func (r *logRuleRepo) ListAll() (rls []Rule, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logDuration, time.Since(begin).Nanoseconds(),
			logElements, len(rls),
			logOp, "ListAll",
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.ListAll()
}

func (r *logRuleRepo) ListActive(configID string, now time.Time) (rls []Rule, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logDuration, time.Since(begin).Nanoseconds(),
			logElements, len(rls),
			logOp, "ListAll",
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.ListActive(configID, now)
}

func (r *logRuleRepo) setup() (err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logDuration, time.Since(begin).Nanoseconds(),
			logOp, "Setup",
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.setup()
}

func (r *logRuleRepo) teardown() (err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logDuration, time.Since(begin).Nanoseconds(),
			logOp, "Teardown",
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.teardown()
}
