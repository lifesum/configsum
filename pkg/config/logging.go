package config

import (
	"time"

	"github.com/go-kit/kit/log"
)

// Log fields.
const (
	logBaseID        = "baseId"
	logClientID      = "clientId"
	logDuration      = "duration"
	logElements      = "elements"
	logErr           = "err"
	logID            = "id"
	logName          = "name"
	logOp            = "op"
	logParameters    = "parameters"
	logPkg           = "pkg"
	logRendered      = "rendered"
	logRepo          = "repo"
	logRuleDecisions = "ruleDecisions"
	logStore         = "store"
	logUserID        = "userId"
)

type logBaseRepo struct {
	logger log.Logger
	next   BaseRepo
}

// NewBaseRepoLogMiddleware wraps the next BaseRepo with logging capabilities.
func NewBaseRepoLogMiddleware(logger log.Logger, store string) BaseRepoMiddleware {
	return func(next BaseRepo) BaseRepo {
		return &logBaseRepo{
			logger: log.With(
				logger,
				logPkg, "config",
				logRepo, "base",
				logStore, store,
			),
			next: next,
		}
	}
}

func (r *logBaseRepo) Create(
	id, clientID, name string,
	parameters rendered,
) (c BaseConfig, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logClientID, clientID,
			logDuration, time.Since(begin).Nanoseconds(),
			logID, id,
			logName, name,
			logOp, "Create",
			logParameters, parameters,
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.Create(id, clientID, name, parameters)
}

func (r *logBaseRepo) GetByID(id string) (c BaseConfig, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logDuration, time.Since(begin).Nanoseconds(),
			logID, id,
			logOp, "GetByID",
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.GetByID(id)
}

func (r *logBaseRepo) GetByName(clientID, name string) (c BaseConfig, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logClientID, clientID,
			logDuration, time.Since(begin).Nanoseconds(),
			logName, name,
			logOp, "GetByName",
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.GetByName(clientID, name)
}

func (r *logBaseRepo) List() (l BaseList, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logDuration, time.Since(begin).Nanoseconds(),
			logElements, len(l),
			logOp, "List",
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.List()
}

func (r *logBaseRepo) Update(input BaseConfig) (bc BaseConfig, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logDuration, time.Since(begin).Nanoseconds(),
			logOp, "Update",
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.Update(input)
}

func (r *logBaseRepo) setup() (err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logDuration, time.Since(begin).Nanoseconds(),
			logOp, "setup",
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.setup()
}

func (r *logBaseRepo) teardown() (err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logDuration, time.Since(begin).Nanoseconds(),
			logOp, "setup",
		}

		if err != nil {
			ps = append(ps, logErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.teardown()
}

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

func (r *logUserRepo) setup() (err error) {
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

func (r *logUserRepo) teardown() (err error) {
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
