package config

import (
	"math/rand"
	"time"
)

const (
	characterSet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numCharacterSet = "0123456789"
)

func randString(charset string) string {
	var (
		b    = make([]byte, len(charset))
		seed = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}

	return string(b)
}
