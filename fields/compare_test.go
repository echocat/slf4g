package fields

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_Equal_isEqual(t *testing.T) {
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := With("a", 1).With("b", 2).With("c", nil)

	actual, actualErr := IsEqualCustom(DefaultEntryEqualityFunction, givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_Equal_withNilFunctionAndDefaultToo(t *testing.T) {
	v := DefaultEntryEqualityFunction
	defer func() {
		DefaultEntryEqualityFunction = v
	}()
	DefaultEntryEqualityFunction = nil

	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := mapped{"a": 1, "b": 2, "c": nil}

	actual, actualErr := IsEqualCustom(nil, givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_Equal_withNilFunction(t *testing.T) {
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := With("a", 1).With("b", 2).With("c", nil)

	actual, actualErr := IsEqualCustom(nil, givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_Equal_moreInRight(t *testing.T) {
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := With("a", 1).With("b", 2).With("c", nil).With("d", 4)

	actual, actualErr := IsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_Equal_lessInRight(t *testing.T) {
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := With("a", 1).With("b", 2)

	actual, actualErr := IsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_Equal_oneValueDifferent(t *testing.T) {
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := With("a", 1).With("b", 2).With("c", 3)

	actual, actualErr := IsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_Equal_oneKeyDifferent(t *testing.T) {
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := With("a", 1).With("b", 2).With("x", nil)

	actual, actualErr := IsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_Equal_leftEmptyRightNil(t *testing.T) {
	givenLeft := mapped{}
	givenRight := Fields(nil)

	actual, actualErr := IsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_Equal_leftNilRightEmpty(t *testing.T) {
	givenLeft := Fields(nil)
	givenRight := mapped{}

	actual, actualErr := IsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_Equal_bothNil(t *testing.T) {
	givenLeft := Fields(nil)
	givenRight := Fields(nil)

	actual, actualErr := IsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_Equal_functionWithErr(t *testing.T) {
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := With("a", 1).With("b", 2).With("c", nil)
	givenErr := errors.New("expected")

	actual, actualErr := IsEqualCustom(func(string, interface{}, interface{}) (bool, error) {
		return false, givenErr
	}, givenLeft, givenRight)

	assert.ToBeEqual(t, givenErr, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_DefaultEntryEqualityFunction(t *testing.T) {
	funcA := func() {}
	funcB := func() {}

	cases := map[string]struct {
		left     interface{}
		right    interface{}
		expected bool
	}{
		"bothString":    {"hello", "hello", true},
		"bothNil":       {nil, nil, true},
		"stringAndNil":  {"hello", nil, false},
		"nilAndString":  {nil, "hello", false},
		"bothFunc":      {funcA, funcA, true},
		"funcAndNil":    {funcA, nil, false},
		"nilAndFunc":    {nil, funcA, false},
		"differentFunc": {funcA, funcB, false},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, actualErr := DefaultEntryEqualityFunction(name, c.left, c.right)

			assert.ToBeNil(t, actualErr)
			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}
