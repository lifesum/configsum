package rule

import (
	"reflect"

	"github.com/lifesum/configsum/pkg/errors"
)

const (
	comparatorGT comparator = iota
)

type comparator int8

// Criteria defines if a rule will match on the given context data.
type Criteria struct {
	User   *CriteriaUser
	Locale *MatcherString
}

// CriteriaUser holds all relevant matchers concerning a user.
type CriteriaUser struct {
	Age          *MatcherInt
	ID           *MatcherStringList
	Subscription *MatcherInt
}

// MatcherBool defines methods for matching a rule on bool type
type MatcherBool struct {
	value bool
}

func (m MatcherBool) match(input interface{}) (bool, error) {
	t, ok := input.(bool)
	if !ok {
		return false, errors.Wrapf(errors.ErrInvalidTypeToMatch, "missmatch %s != bool", reflect.TypeOf(input).Kind())
	}

	return t == m.value, nil
}

// MatcherInt defines methods for matching a rule on int type
type MatcherInt struct {
	Comparator comparator
	Value      int
}

func (m MatcherInt) match(input interface{}) (bool, error) {
	t, ok := input.(int)
	if !ok {
		return false, errors.Wrapf(errors.ErrInvalidTypeToMatch, "missmatch %s != int", reflect.TypeOf(input).Kind())
	}

	switch m.Comparator {
	case comparatorGT:
		return t > m.Value, nil
	default:
		return false, nil
	}
}

// MatcherString defines methods for matching a rule on string type
type MatcherString string

func (m MatcherString) match(input interface{}) (bool, error) {
	t, ok := input.(string)
	if !ok {
		return false, errors.Wrapf(errors.ErrInvalidTypeToMatch, "missmatch %s != string", reflect.TypeOf(input).Kind())
	}

	if t == string(m) {
		return true, nil
	}

	return false, nil
}

// MatcherStringList defines methods for matching a rule on string list type
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
