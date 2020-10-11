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
