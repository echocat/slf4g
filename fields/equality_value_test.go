package fields

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_DefaultValueEquality(t *testing.T) {
	funcA := func() {}
	funcB := func() {}

	cases := map[string]struct {
		left     interface{}
		right    interface{}
		expected bool
	}{
		"bothString":       {"hello", "hello", true},
		"bothNil":          {nil, nil, true},
		"stringAndNil":     {"hello", nil, false},
		"nilAndString":     {nil, "hello", false},
		"bothFunc":         {funcA, funcA, true},
		"funcAndNil":       {funcA, nil, false},
		"nilAndFunc":       {nil, funcA, false},
		"differentFunc":    {funcA, funcB, false},
		"bothLazy":         {LazyFormat("foo"), LazyFormat("foo"), true},
		"bothLazyMismatch": {LazyFormat("foo1"), LazyFormat("foo2"), false},
		"leftLazy":         {LazyFormat("foo"), "foo", true},
		"rightLazy":        {"foo", LazyFormat("foo"), true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, actualErr := DefaultValueEquality.AreValuesEqual(name, c.left, c.right)

			assert.ToBeNil(t, actualErr)
			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}

func Test_ValueEqualityFunc_AreValuesEqual(t *testing.T) {
	givenLeft := &struct{ foo string }{"aFoo"}
	givenRight := &struct{ bar string }{"aBar"}
	givenErr := errors.New("expected")

	instance := ValueEqualityFunc(func(name string, left, right interface{}) (bool, error) {
		assert.ToBeEqual(t, "foo", name)
		assert.ToBeSame(t, givenLeft, left)
		assert.ToBeSame(t, givenRight, right)
		return true, givenErr
	})

	actual, actualErr := instance.AreValuesEqual("foo", givenLeft, givenRight)

	assert.ToBeSame(t, givenErr, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_NewValueEqualityFacade(t *testing.T) {
	givenFunc := func() ValueEquality { panic("should never be called") }

	actual := NewValueEqualityFacade(givenFunc)

	assert.ToBeSame(t, valueEqualityFacade(givenFunc), actual)
}

func Test_valueEqualityFacade_AreValuesEqual(t *testing.T) {
	givenLeft := &struct{ foo string }{"aFoo"}
	givenRight := &struct{ bar string }{"aBar"}
	givenErr := errors.New("expected")
	givenEquality := ValueEqualityFunc(func(name string, left, right interface{}) (bool, error) {
		assert.ToBeEqual(t, "foo", name)
		assert.ToBeSame(t, givenLeft, left)
		assert.ToBeSame(t, givenRight, right)
		return true, givenErr
	})

	instance := valueEqualityFacade(func() ValueEquality { return givenEquality })

	actual, actualErr := instance.AreValuesEqual("foo", givenLeft, givenRight)

	assert.ToBeSame(t, givenErr, actualErr)
	assert.ToBeEqual(t, true, actual)
}
