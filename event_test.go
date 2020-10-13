package log

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/fields"

	"github.com/echocat/slf4g/internal/test/assert"

	"github.com/echocat/slf4g/level"
)

func Test_NewEvent_withoutFields(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLevel := level.Error
	givenCallDepth := 66

	actual := NewEvent(givenProvider, givenLevel, givenCallDepth)

	assert.ToBeOfType(t, &eventImpl{}, actual)
	assert.ToBeSame(t, givenProvider, actual.(*eventImpl).provider)
	assert.ToBeEqual(t, givenLevel, actual.GetLevel())
	assert.ToBeEqual(t, givenCallDepth, actual.GetCallDepth())
	assert.ToBeEqual(t, fields.Empty(), actual.(*eventImpl).fields)
	assert.ToBeNil(t, actual.GetContext())
}

func Test_NewEvent_withOneFields(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLevel := level.Error
	givenCallDepth := 66

	givenFields1 := fields.With("a", "1")

	actual := NewEvent(givenProvider, givenLevel, givenCallDepth, givenFields1)

	assert.ToBeOfType(t, &eventImpl{}, actual)
	assert.ToBeSame(t, givenProvider, actual.(*eventImpl).provider)
	assert.ToBeEqual(t, givenLevel, actual.GetLevel())
	assert.ToBeEqual(t, givenCallDepth, actual.GetCallDepth())
	assert.ToBeSame(t, givenFields1, actual.(*eventImpl).fields)
	assert.ToBeNil(t, actual.GetContext())
}

func Test_NewEvent_with3Fields(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLevel := level.Error
	givenCallDepth := 66

	givenFields1 := fields.With("a", 1)
	givenFields2 := fields.With("a", 2).With("b", 2)
	givenFields3 := fields.With("a", 3).With("c", 3)

	actual := NewEvent(givenProvider, givenLevel, givenCallDepth, givenFields1, givenFields2, givenFields3)

	assert.ToBeOfType(t, &eventImpl{}, actual)
	assert.ToBeSame(t, givenProvider, actual.(*eventImpl).provider)
	assert.ToBeEqual(t, givenLevel, actual.GetLevel())
	assert.ToBeEqual(t, givenCallDepth, actual.GetCallDepth())
	assert.ToBeEqualUsing(t, fields.With("a", 3).With("b", 2).With("c", 3), actual.(*eventImpl).fields, fields.AreEqual)
	assert.ToBeNil(t, actual.GetContext())
}

func Test_NewEvent_withErrorInFieldsPanics(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLevel := level.Error
	givenCallDepth := 66
	givenError := errors.New("expected")

	givenFields1 := fields.With("a", 1)
	givenFields2 := &fieldsThatErrors{fields.With("a", 2), givenError}

	assert.Execution(t, func() {
		NewEvent(givenProvider, givenLevel, givenCallDepth, givenFields1, givenFields2)
	}).WillPanicWith("expected")
}

type fieldsThatErrors struct {
	fields.Fields
	err error
}

func (instance *fieldsThatErrors) ForEach(func(key string, value interface{}) error) error {
	return instance.err
}

type entries []entry

func (instance *entries) add(key string, value interface{}) {
	*instance = append(*instance, entry{key, value})
}

func (instance *entries) consumer() func(key string, value interface{}) error {
	return func(key string, value interface{}) error {
		instance.add(key, value)
		return nil
	}
}

type entry struct {
	key   string
	value interface{}
}
