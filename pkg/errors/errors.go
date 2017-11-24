package errors

import (
	"github.com/pkg/errors"
)

// Entity errors.
var (
	ErrID       = errors.New("id creation")
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

// Config errors.
var (
	ErrParametersInvalid = errors.New("parameters invalid")
)

// Transport errors.
var (
	ErrInvalidPayload = errors.New("payload invalid")
	ErrVarMissing     = errors.New("variable missing")
)

// Rule errors.
var (
	ErrInvalidRule        = errors.New("invalid rule")
	ErrInvalidTypeToMatch = errors.New("invalid input type")
	ErrNoRuleForID        = errors.New("no rules for this ID")
	ErrNoRuleWithName     = errors.New("no rule with name")
)

// Cause is a wraper over github.com/pkg/errors.Cause.
func Cause(err error) error {
	return errors.Cause(err)
}

// Errorf is a wraper over github.com/pkg/errors.Errorf.
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// New is a wraper over github.com/pkg/errors.New.
func New(message string) error {
	return errors.New(message)
}

// Wrap is a wrapper over github.com/pkg/errors.Wrap.
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Wrapf is a wrapper over github.com/pkg/errors.Wrapf.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}
