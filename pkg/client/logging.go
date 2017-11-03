package client

import (
	"time"

	"github.com/go-kit/kit/log"
)

// Log fields.
const (
	logFieldClientID = "client_id"
	logFieldDuration = "duration"
	logFieldElements = "elements"
	logFieldErr      = "err"
	logFieldID       = "id"
	logFieldOp       = "op"
	logFieldPkg      = "pkg"
	logFieldRepo     = "repo"
	logFieldSecret   = "secret"
	logFieldStore    = "store"
)

type logRepo struct {
	logger log.Logger
	next   Repo
}

// NewRepoLogMiddleware wraps the next Repo with logging capabilities.
func NewRepoLogMiddleware(logger log.Logger, store string) RepoMiddleware {
	return func(next Repo) Repo {
		return &logRepo{
			logger: log.With(
				logger,
				logFieldPkg, "client",
				logFieldRepo, labelRepo,
			),
			next: next,
		}
	}
}

func (r *logRepo) List() (cs List, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logFieldDuration, time.Since(begin).Nanoseconds(),
			logFieldElements, len(cs),
			logFieldOp, "List",
		}

		if err != nil {
			ps = append(ps, logFieldErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.List()
}

func (r *logRepo) Lookup(id string) (client Client, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logFieldDuration, time.Since(begin).Nanoseconds(),
			logFieldID, id,
			logFieldOp, "Lookup",
		}

		if err != nil {
			ps = append(ps, logFieldErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.Lookup(id)
}

func (r *logRepo) Store(id, name string) (client Client, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logFieldDuration, time.Since(begin).Nanoseconds(),
			logFieldID, id,
			logFieldOp, "Store",
		}

		if err != nil {
			ps = append(ps, logFieldErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.Store(id, name)
}

func (r *logRepo) setup() (err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logFieldDuration, time.Since(begin).Nanoseconds(),
			logFieldOp, "setup",
		}

		if err != nil {
			ps = append(ps, logFieldErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.setup()
}

func (r *logRepo) teardown() (err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logFieldDuration, time.Since(begin).Nanoseconds(),
			logFieldOp, "teardown",
		}

		if err != nil {
			ps = append(ps, logFieldErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.teardown()
}

type logTokenRepo struct {
	logger log.Logger
	next   TokenRepo
}

// NewTokenRepoLogMiddleware wraps the next TokenRepo with logging capabilities.
func NewTokenRepoLogMiddleware(logger log.Logger, store string) TokenRepoMiddleware {
	return func(next TokenRepo) TokenRepo {
		return &logTokenRepo{
			logger: log.With(
				logger,
				logFieldPkg, "client",
				logFieldRepo, labelRepoToken,
			),
			next: next,
		}
	}
}

func (r *logTokenRepo) GetLatest(clientID string) (token Token, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logFieldClientID, clientID,
			logFieldDuration, time.Since(begin).Nanoseconds(),
			logFieldOp, "GetLatest",
		}

		if err != nil {
			ps = append(ps, logFieldErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.GetLatest(clientID)
}

func (r *logTokenRepo) Lookup(secret string) (token Token, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logFieldDuration, time.Since(begin).Nanoseconds(),
			logFieldOp, "Lookup",
			logFieldSecret, secret,
		}

		if err != nil {
			ps = append(ps, logFieldErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.Lookup(secret)
}

func (r *logTokenRepo) Store(clientID, secret string) (token Token, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logFieldDuration, time.Since(begin).Nanoseconds(),
			logFieldClientID, clientID,
			logFieldOp, "Store",
			logFieldSecret, secret,
		}

		if err != nil {
			ps = append(ps, logFieldErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.Store(clientID, secret)
}

func (r *logTokenRepo) setup() (err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logFieldDuration, time.Since(begin).Nanoseconds(),
			logFieldOp, "setup",
		}

		if err != nil {
			ps = append(ps, logFieldErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.setup()
}

func (r *logTokenRepo) teardown() (err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			logFieldDuration, time.Since(begin).Nanoseconds(),
			logFieldOp, "teardown",
		}

		if err != nil {
			ps = append(ps, logFieldErr, err)
		}

		_ = r.logger.Log(ps...)
	}(time.Now())

	return r.next.teardown()
}
