package rule

import (
	"encoding/json"
	"reflect"
	"testing"

	"golang.org/x/text/language"
)

func TestCriterionDeviceLocationLocaleMarshal(t *testing.T) {
	var (
		tag  = language.MustParse("en-US")
		want = Criterion{
			Comparator: ComparatorEQ,
			Key:        DeviceLocationLocale,
			Value:      tag,
			Path:       "",
		}
	)

	raw, err := json.Marshal(&want)
	if err != nil {
		t.Fatal(err)
	}

	var have Criterion

	err = json.Unmarshal(raw, &have)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestCriterionUserSubscriptionMarshal(t *testing.T) {
	var (
		want = Criterion{
			Comparator: ComparatorEQ,
			Key:        UserSubscription,
			Value:      0,
			Path:       "",
		}
	)

	raw, err := json.Marshal(&want)
	if err != nil {
		t.Fatal(err)
	}

	var have Criterion

	err = json.Unmarshal(raw, &have)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(have, want) {
		t.Log(reflect.TypeOf(want.Value))
		t.Log(reflect.TypeOf(have.Value))
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestValidDateCriteriaMarshal(t *testing.T) {
	var (
		want = Criterion{
			Comparator: ComparatorEQ,
			Key:        ValidDate,
			Value:      0,
			Path:       "",
		}
	)

	raw, err := json.Marshal(&want)
	if err != nil {
		t.Fatal(err)
	}

	var have Criterion

	err = json.Unmarshal(raw, &have)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(have, want) {
		t.Log(reflect.TypeOf(want.Value))
		t.Log(reflect.TypeOf(have.Value))
		t.Errorf("have %v, want %v", have, want)
	}
}
