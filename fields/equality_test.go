package fields

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_IsEqual_isEqual(t *testing.T) {
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := With("a", 1).With("b", 2).With("c", nil)

	actual, actualErr := AreEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_IsEqual_isNotEqual(t *testing.T) {
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := With("a", 1).With("b", 666).With("c", nil)

	actual, actualErr := AreEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_IsEqual_withNilDefaultEquality(t *testing.T) {
	v := DefaultEquality
	defer func() { DefaultEquality = v }()
	DefaultEquality = nil

	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := mapped{"a": 1, "b": 2, "c": nil}

	actual, actualErr := AreEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_EqualityImpl_AreFieldsEqual_withNilInstance(t *testing.T) {
	var instance *EqualityImpl
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := mapped{"a": 1, "b": 666, "c": nil}

	//goland:noinspection GoNilness
	actual, actualErr := instance.AreFieldsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_EqualityImpl_AreFieldsEqual_withNilValueEquality(t *testing.T) {
	instance := &EqualityImpl{ValueEquality: nil}
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := mapped{"a": 1, "b": 666, "c": nil}

	actual, actualErr := instance.AreFieldsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_EqualityImpl_AreFieldsEqual_moreInRight(t *testing.T) {
	instance := &EqualityImpl{ValueEquality: nil}
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := With("a", 1).With("b", 2).With("c", nil).With("d", 4)

	actual, actualErr := instance.AreFieldsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_EqualityImpl_AreFieldsEqual_lessInRight(t *testing.T) {
	instance := &EqualityImpl{ValueEquality: nil}
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := With("a", 1).With("b", 2)

	actual, actualErr := instance.AreFieldsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_EqualityImpl_oneValueDifferent(t *testing.T) {
	instance := &EqualityImpl{ValueEquality: DefaultValueEquality}
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := With("a", 1).With("b", 2).With("c", 3)

	actual, actualErr := instance.AreFieldsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_EqualityImpl_AreFieldsEqual_oneKeyDifferent(t *testing.T) {
	instance := &EqualityImpl{ValueEquality: DefaultValueEquality}
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := With("a", 1).With("b", 2).With("x", nil)

	actual, actualErr := instance.AreFieldsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_EqualityImpl_AreFieldsEqual_leftEmptyRightNil(t *testing.T) {
	instance := &EqualityImpl{ValueEquality: DefaultValueEquality}
	givenLeft := mapped{}
	givenRight := Fields(nil)

	actual, actualErr := instance.AreFieldsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_EqualityImpl_AreFieldsEqual_leftNilRightEmpty(t *testing.T) {
	instance := &EqualityImpl{ValueEquality: DefaultValueEquality}
	givenLeft := Fields(nil)
	givenRight := mapped{}

	actual, actualErr := instance.AreFieldsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_EqualityImpl_AreFieldsEqual_bothNil(t *testing.T) {
	instance := &EqualityImpl{ValueEquality: DefaultValueEquality}
	givenLeft := Fields(nil)
	givenRight := Fields(nil)

	actual, actualErr := instance.AreFieldsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_EqualityImpl_AreFieldsEqual_functionWithErr(t *testing.T) {
	givenLeft := mapped{"a": 1, "b": 2, "c": nil}
	givenRight := With("a", 1).With("b", 2).With("c", nil)
	givenErr := errors.New("expected")
	instance := &EqualityImpl{ValueEquality: ValueEqualityFunc(func(_ string, _, _ interface{}) (bool, error) {
		return false, givenErr
	})}

	actual, actualErr := instance.AreFieldsEqual(givenLeft, givenRight)

	assert.ToBeEqual(t, givenErr, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_EqualityFunc_AreFieldsEqual(t *testing.T) {
	givenLeft := With("a", 1)
	givenRight := With("a", 2)
	givenErr := errors.New("expected")

	instance := EqualityFunc(func(left, right Fields) (bool, error) {
		assert.ToBeSame(t, givenLeft, left)
		assert.ToBeSame(t, givenRight, right)
		return true, givenErr
	})

	actual, actualErr := instance.AreFieldsEqual(givenLeft, givenRight)

	assert.ToBeSame(t, givenErr, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_NewEqualityFacade(t *testing.T) {
	givenFunc := func() Equality { panic("should never be called") }

	actual := NewEqualityFacade(givenFunc)

	assert.ToBeSame(t, equalityFacade(givenFunc), actual)
}

func Test_equalityFacade_AreFieldsEqual(t *testing.T) {
	givenLeft := With("a", 1)
	givenRight := With("a", 2)
	givenErr := errors.New("expected")
	givenEquality := EqualityFunc(func(left, right Fields) (bool, error) {
		assert.ToBeSame(t, givenLeft, left)
		assert.ToBeSame(t, givenRight, right)
		return true, givenErr
	})

	instance := equalityFacade(func() Equality { return givenEquality })

	actual, actualErr := instance.AreFieldsEqual(givenLeft, givenRight)

	assert.ToBeSame(t, givenErr, actualErr)
	assert.ToBeEqual(t, true, actual)
}
