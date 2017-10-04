package config

import (
	"github.com/pkg/errors"
)

// Common errors.
var (
	ErrExists   = errors.New("entity exists")
	ErrNotFound = errors.New("entity not found")
)
