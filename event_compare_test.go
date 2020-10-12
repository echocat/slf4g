package log

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_IsEventEqual_isEqual(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)

	actual, actualErr := IsEventEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_IsEventEqual_withNilFunctionAndDefaultToo(t *testing.T) {
	v := fields.DefaultEntryEqualityFunction
	defer func() {
		fields.DefaultEntryEqualityFunction = v
	}()
	fields.DefaultEntryEqualityFunction = nil

	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)

	actual, actualErr := IsEventEqualCustom(nil, givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_IsEventEqual_withNilFunction(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)

	actual, actualErr := IsEventEqualCustom(nil, givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_IsEventEqual_differentLevel(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Info, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)

	actual, actualErr := IsEventEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_IsEventEqual_differentCallDepth(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 4).
		With("a", 1).
		With("b", 2).
		With("c", 3)

	actual, actualErr := IsEventEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_IsEventEqual_moreEntries(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3).
		With("d", 4)

	actual, actualErr := IsEventEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_IsEventEqual_lessEntries(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2)

	actual, actualErr := IsEventEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_IsEventEqual_oneDifferentValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 666)

	actual, actualErr := IsEventEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_IsEventEqual_oneDifferentKeys(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("xyz", 3)

	actual, actualErr := IsEventEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_IsEventEqual_leftEmptyRightNil(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	var givenRight Event

	actual, actualErr := IsEventEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_IsEventEqual_leftNilRightEmpty(t *testing.T) {
	givenProvider := newMockProvider("test")
	var givenLeft Event
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("xyz", 3)

	actual, actualErr := IsEventEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_IsEventEqual_bothNil(t *testing.T) {
	var givenLeft Event
	var givenRight Event

	actual, actualErr := IsEventEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_IsEventEqual_functionWithErr(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenErr := errors.New("expected")

	actual, actualErr := IsEventEqualCustom(func(string, interface{}, interface{}) (bool, error) {
		return false, givenErr
	}, givenLeft, givenRight)

	assert.ToBeEqual(t, givenErr, actualErr)
	assert.ToBeEqual(t, false, actual)
}
