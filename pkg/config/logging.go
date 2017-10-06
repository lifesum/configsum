package config

import (
	"time"

	"github.com/go-kit/kit/log"
)

// Log fields.
const (
	logBaseID        = "baseId"
	logDuration      = "duration"
	logErr           = "err"
	logID            = "id"
	logOp            = "op"
	logPkg           = "pkg"
	logRendered      = "rendered"
	logRepo          = "repo"
	logRuleDecisions = "ruleDecisions"
	logStore         = "store"
	logUserID        = "userId"
)

type logUserRepo struct {
	logger log.Logger
	next   UserRepo
}

// NewUserRepoLogMiddleware wraps the next UserRepo with logging capabilities.
func NewUserRepoLogMiddleware(logger log.Logger, store string) UserRepoMiddleware {
	return func(next UserRepo) UserRepo {
		return &logUserRepo{
			logger: log.With(
				logger,
				logPkg, "config",
				logRepo, "user",
				logStore, store,
			),
			next: next,
		}
	}
}

func (r *logUserRepo) Append(
	id, baseID, userID string,
	decisions ruleDecisions,
	render rendered,
) (c UserConfig, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logBaseID, baseID,
			logDuration, time.Since(begin).Nanoseconds(),
			logID, id,
			logOp, "Append",
			logRendered, render,
			logRuleDecisions, decisions,
			logUserID, userID,
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.Append(id, baseID, userID, decisions, render)
}

func (r *logUserRepo) GetLatest(baseID, userID string) (c UserConfig, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logBaseID, baseID,
			logDuration, time.Since(begin).Nanoseconds(),
			logOp, "GetLatest",
			logUserID, userID,
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.GetLatest(baseID, userID)
}

func (r *logUserRepo) Setup() (err error) {
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

	return r.next.Setup()
}

func (r *logUserRepo) Teardown() (err error) {
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

	return r.next.Teardown()
}
