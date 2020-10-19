package fields

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_mustAsMap_withMapped(t *testing.T) {
	given := mapped{"foo": 1, "bar": 2}

	actual := mustAsMap(given)

	assert.ToBeEqual(t, given, actual)
}

func Test_mustAsMap_withPMapped(t *testing.T) {
	given := mapped{"foo": 1, "bar": 2}

	actual := mustAsMap(&given)

	assert.ToBeEqual(t, given, actual)
}

func Test_mustAsMap_copies(t *testing.T) {
	given := &lineage{With("a", 1), With("b", 2)}

	actual := mustAsMap(given)

	assert.ToBeEqual(t, mapped{"a": 1, "b": 2}, actual)
}

func Test_mustAsMap_withNil(t *testing.T) {
	actual := mustAsMap(nil)

	assert.ToBeEqual(t, mapped{}, actual)
}

func Test_asMap_withError(t *testing.T) {
	givenErr := errors.New("expected")
	given := ForEachFunc(func(func(key string, value interface{}) error) error {
		return givenErr
	})

	actual, actualErr := asMap(given)

	assert.ToBeEqual(t, givenErr, actualErr)
	assert.ToBeEqual(t, mapped(nil), actual)
}

func Test_mustAsMap_withError(t *testing.T) {
	givenErr := errors.New("expected")
	given := ForEachFunc(func(func(key string, value interface{}) error) error {
		return givenErr
	})

	assert.Execution(t, func() {
		mustAsMap(given)
	}).WillPanicWith("expected")
}

func Test_AsFields_withNil(t *testing.T) {
	actual, actualErr := AsFields(nil)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, Empty(), actual)
}

func Test_AsFields_withFields(t *testing.T) {
	givenFields := With("foo", "bar")
	actual, actualErr := AsFields(givenFields)

	assert.ToBeNil(t, actualErr)
	assert.ToBeSame(t, givenFields, actual)
}

func Test_AsFields_withForEachEnabled(t *testing.T) {
	givenForEachEnabled := aMap{"foo": 1, "bar": 2}
	actual, actualErr := AsFields(givenForEachEnabled)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": 2}, actual)
}

type aMap map[string]interface{}

func (instance aMap) ForEach(consumer func(string, interface{}) error) error {
	for k, v := range instance {
		if err := consumer(k, v); err != nil {
			return err
		}
	}
	return nil
}
