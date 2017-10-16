package dory

import (
	"github.com/pkg/errors"
)

var (
	ErrSignatureMissing   = errors.New("signature missing")
	ErrSignatureMissmatch = errors.New("signature missmatch")
	ErrUserIDMissing      = errors.New("userID missing")
)
