package rule

import (
	"testing"

	"github.com/lifesum/configsum/pkg/generate"
)

func TestMatcherListString(t *testing.T) {
	var (
		input    = generate.RandomString(24)
		goodVals = []string{
			input,
			generate.RandomString(24),
			generate.RandomString(24),
		}
		badVals = []string{
			generate.RandomString(24),
			generate.RandomString(24),
			generate.RandomString(24),
		}
	)

	m := MatcherStringList(goodVals)

	ok, err := m.match(input)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Errorf("expect input to match")
	}

	m = MatcherStringList(badVals)

	ok, err = m.match(input)
	if err != nil {
		t.Fatal(err)
	}

	if ok {
		t.Errorf("expect input not to match")
	}
}
