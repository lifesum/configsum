package pg

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// Default URIs.
const (
	DefaultDevURI  = "postgres://%s@127.0.0.1:5432/configsum_dev?sslmode=disable&connect_timeout=5"
	DefaultTestURI = "postgres://%s@127.0.0.1:5432/configsum_test?sslmode=disable&connect_timeout=5"
)

const (
	codeDuplicateKeyViolation = "23505"
	codeRelationshipNotFound  = "42P01"
)

// Errors.
var (
	ErrDuplicateKey     = errors.New("duplicate key")
	ErrRelationNotFound = errors.New("relation not found")
)

// IsRelationNotFound indicates if err is ErrRelationNotFound.
func IsRelationNotFound(err error) bool {
	return err == ErrRelationNotFound
}

// Wrap translates *pq.Error codes into actionable errors.
func Wrap(err error) error {
	if e, ok := err.(*pq.Error); ok {
		switch e.Code {
		case codeDuplicateKeyViolation:
			return ErrDuplicateKey
		case codeRelationshipNotFound:
			return ErrRelationNotFound
		}
	}

	return err
}
