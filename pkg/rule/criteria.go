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

type comparator int8

type criteria struct {
	User *criteriaUser
}

type matcher interface {
	match(input interface{}) (bool, error)
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

type matcherListString struct {
	Value []string
}

func (m matcherListString) match(input interface{}) (bool, error) {
	t, ok := input.(string)
	if !ok {
		return false, errors.Wrapf(errors.ErrInvalidTypeToMatch, "missmatch %s != string", reflect.TypeOf(input).Kind())
	}

	for _, v := range m.Value {
		if t == v {
			return true, nil
		}
	}

	return false, nil
}

type criteriaUser struct {
	Age *matcher
	ID  *matcherListString
}
