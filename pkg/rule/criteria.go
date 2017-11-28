package rule

import (
	"reflect"

	"github.com/lifesum/configsum/pkg/errors"
)

const (
	comparatorEQ comparator = iota
	comparatorNQ
	comparatorGT
	comparatorLT
	comparatorIN
)

type matcher interface {
	match(input interface{}) (bool, error)
}

type comparator int8

// Criteria defines if a rule will match on the given context data.
type Criteria struct {
	User *CriteriaUser
}

// CriteriaUser holds all relevant matchers concerning a user.
type CriteriaUser struct {
	Age *matcher
	ID  *MatcherStringList
}

type matcherBool struct {
	comparator comparator
	value      bool
}

type matcherInt struct {
	comparator comparator
	value      int
}

type matcherString struct {
	comparator comparator
	value      string
}

type matcherListInt struct {
	value []int
}

type MatcherStringList []string

func (m MatcherStringList) match(input interface{}) (bool, error) {
	t, ok := input.(string)
	if !ok {
		return false, errors.Wrapf(errors.ErrInvalidTypeToMatch, "missmatch %s != string", reflect.TypeOf(input).Kind())
	}

	for _, v := range m {
		if t == v {
			return true, nil
		}
	}

	return false, nil
}
