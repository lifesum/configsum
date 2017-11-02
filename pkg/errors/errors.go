package errors

import (
	"github.com/pkg/errors"
)

// Entity errors.
var (
	ErrExists   = errors.New("entity exists")
	ErrNotFound = errors.New("entity not found")
)

// Auth errors.
var (
	ErrClientNotFound     = errors.New("client not found")
	ErrSecretMissing      = errors.New("secret missing")
	ErrSignatureMissing   = errors.New("signature missing")
	ErrSignatureMissmatch = errors.New("signature missmatch")
	ErrUserIDMissing      = errors.New("userID missing")
)

// Transport errors.
var (
	ErrInvalidPayload = errors.New("payload invalid")
	ErrVarMissing     = errors.New("variable missing")
)

// Cause is a wraper over github.com/pkg/errors.Cause.
func Cause(err error) error {
	return errors.Cause(err)
}

// Wrap is a wrapper over github.com/pkg/errors.Wrap.
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Wrapf is a wrapper over github.com/pkg/errors.Wrapf.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}
