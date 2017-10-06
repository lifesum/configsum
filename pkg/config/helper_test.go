package config

import (
	"math/rand"
	"testing"
	"time"
)

const (
	characterSet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numCharacterSet = "0123456789"
)

var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

type prepareFunc func(t *testing.T) UserRepo

func randString(charset string) string {
	b := make([]byte, len(charset))

	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}

	return string(b)
}
