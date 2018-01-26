package rule

import (
	"encoding/json"

	"golang.org/x/text/language"

	"github.com/lifesum/configsum/pkg/errors"
)

// Comparators.
const (
	ComparatorGT Comparator = iota
	ComparatorEQ
	ComparatorNQ
	ComparatorIN
)

// Comparator defines the type of comparison for a Criterion.
type Comparator int8

func (c Comparator) String() string {
	switch c {
	case ComparatorEQ:
		return "ComparatorEQ"
	case ComparatorNQ:
		return "ComparatorNQ"
	case ComparatorGT:
		return "ComparatorGT"
	case ComparatorIN:
		return "ComparatorIN"
	default:
		return "unknown comparator"
	}
}

// App context keys.
const (
	AppVersion CriterionKey = iota + 1
)

// Device context keys.
const (
	DeviceLocationLocale CriterionKey = iota + 101
	DeviceLocationOffset
	DeviceOSPlatform
	DeviceOSVersion
)

// Metadata context keys.
const (
	MetadataBool CriterionKey = iota + 201
	MetadataNumber
	MetadataString
)

// User context keys.
const (
	UserAge CriterionKey = iota + 301
	UserRegistered
	UserID
	UserSubscription
)

// CriterionKey is the set of possible input to match on.
type CriterionKey int

func (k CriterionKey) String() string {
	switch k {
	case AppVersion:
		return "AppVersion"
	case UserID:
		return "UserID"
	case UserSubscription:
		return "UserSubscription"
	default:
		return "unknown"
	}
}

// Criteria is a collection of Criterion.
type Criteria []Criterion

// Criterion is a single decision which can be evaluated to decide if a Rule
// should be applied.
type Criterion struct {
	Comparator Comparator
	Key        CriterionKey
	Value      interface{}
	Path       string
}

// MarshalJSON to satisfy json.Marshaler.
func (c Criterion) MarshalJSON() ([]byte, error) {
	var value interface{}

	switch c.Key {
	case DeviceLocationLocale:
		t, ok := c.Value.(language.Tag)
		if !ok {
			return nil, errors.Wrapf(errors.ErrInvalidTypeToMatch, "%s: value not a language.Tag", c.Key)
		}

		value = t.String()
	case UserSubscription:
		t, ok := c.Value.(int)
		if !ok {
			return nil, errors.Wrapf(errors.ErrInvalidTypeToMatch, "%s: value noat an int", c.Key)
		}

		value = t
	case UserID:
		t, ok := c.Value.([]string)
		if !ok {
			return nil, errors.Wrapf(errors.ErrInvalidTypeToMatch, "%s: value not a string slice", c.Key)
		}

		value = t
	default:
		return nil, errors.Errorf("marshaling for '%s' not supported", c.Key)
	}

	return json.Marshal(struct {
		Comparator Comparator   `json:"comparator"`
		Key        CriterionKey `json:"key"`
		Value      interface{}  `json:"value"`
		Path       string       `json:"path"`
	}{
		Comparator: c.Comparator,
		Key:        c.Key,
		Value:      value,
		Path:       c.Path,
	})
}

// UnmarshalJSON to satisfy json.Unmarshaler.
func (c *Criterion) UnmarshalJSON(raw []byte) error {
	v := struct {
		Comparator Comparator   `json:"comparator"`
		Key        CriterionKey `json:"key"`
		Value      interface{}  `json:"value"`
		Path       string       `json:"path"`
	}{}

	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}

	c.Comparator = v.Comparator
	c.Key = v.Key
	c.Path = v.Path

	switch c.Key {
	case DeviceLocationLocale:
		s, ok := v.Value.(string)
		if !ok {
			return errors.Wrapf(errors.ErrInvalidTypeToMatch, "%s: value not a string", c.Key)
		}

		t, err := language.Parse(s)
		if err != nil {
			// TODO(xla): Wrap in proper marshal error.
			return err
		}

		c.Value = t
	case UserSubscription:
		s, ok := v.Value.(float64)
		if !ok {
			return errors.Wrapf(errors.ErrInvalidTypeToMatch, "%s: value not an int", c.Key)
		}

		c.Value = int(s)
	case UserID:
		s, ok := constructSlice(v.Value)
		if !ok {
			return errors.Wrapf(errors.ErrInvalidTypeToMatch, "%s: value not a string slice", c.Key)
		}

		c.Value = s
	default:
		return errors.Errorf("unmarshaling for '%s' not supported", c.Key)
	}

	return nil
}

func (c Criterion) match(ctx Context) error {
	switch c.Key {
	case DeviceLocationLocale:
		expected, ok := c.Value.(language.Tag)
		if !ok {
			return errors.Wrapf(errors.ErrInvalidTypeToMatch, "%s: value not a language.Tag", c.Key)
		}

		return matchLocationLocale(c.Comparator, expected, ctx.Locale.Locale)
	case UserSubscription:
		expected, ok := c.Value.(int)
		if !ok {
			return errors.Wrapf(errors.ErrInvalidTypeToMatch, "%s: value not an int", c.Key)
		}
		return matchUserSubscription(c.Comparator, expected, ctx.User.Subscription)
	case UserID:
		expected, ok := constructSlice(c.Value)
		if !ok {
			return errors.Wrapf(errors.ErrInvalidTypeToMatch, "%s: value not a string slice", c.Key)
		}

		return matchUserID(c.Comparator, expected, ctx.User.ID)
	default:
		errors.Errorf("unsupported Key '%d'", c.Key)
	}

	return nil
}

func matchUserID(comparator Comparator, expected []string, input string) error {
	switch comparator {
	case ComparatorIN:
		for _, id := range expected {
			if id == input {
				return nil
			}
		}
	default:
		return errors.Errorf("comparator '%s' not supported", comparator)
	}

	return errors.Wrap(errors.ErrCriterionNotMatch, "id cannot be found in the id list")
}

func matchUserSubscription(comparator Comparator, expected, input int) error {
	switch comparator {
	case ComparatorGT:
		if input <= expected {
			return errors.Wrap(errors.ErrCriterionNotMatch, "input value smaller or equal than criterion value")
		}
	default:
		return errors.Errorf("comparator '%s' not supported", comparator)
	}

	return nil
}

func matchLocationLocale(comparator Comparator, expected, input language.Tag) error {
	var (
		er, _ = expected.Region()
		ir, _ = input.Region()
		ok    = er.Contains(ir)
	)

	switch comparator {
	case ComparatorEQ:
		if !ok {
			return errors.Wrap(errors.ErrCriterionNotMatch, "input region not part of expected region")
		}
	case ComparatorNQ:
		if ok {
			return errors.Wrap(errors.ErrCriterionNotMatch, "input region is part of expected region")
		}
	default:
		return errors.Errorf("comparator '%s' not supproted", comparator)
	}

	return nil
}

func constructSlice(input interface{}) ([]string, bool) {
	r := []string{}
	switch v := input.(type) {
	case []string:
		for _, s := range v {
			r = append(r, s)
		}

		return r, true
	case []interface{}:
		for _, s := range v {
			r = append(r, s.(string))
		}

		return r, true
	default:
		return nil, false
	}
}
