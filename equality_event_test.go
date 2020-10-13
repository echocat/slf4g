package log

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/fields"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_AreEventsEqual_isEqual(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)

	actual, actualErr := AreEventsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_AreEventsEqual_isNotEqual(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 666)

	actual, actualErr := AreEventsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_AreEventsEqual_withNilDefaultEventEquality(t *testing.T) {
	v := DefaultEventEquality
	defer func() {
		DefaultEventEquality = v
	}()
	DefaultEventEquality = nil

	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)

	actual, actualErr := AreEventsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_EventEqualityImpl_AreEventsEqual_withNilInstance(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	var instance *EventEqualityImpl

	//goland:noinspection GoNilness
	actual, actualErr := instance.AreEventsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}
func Test_EventEqualityImpl_AreEventsEqual_withNilFunction(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 666)
	instance := EventEqualityImpl{
		CompareLevel:       true,
		CompareCallDepth:   true,
		CompareValuesUsing: nil,
	}

	actual, actualErr := instance.AreEventsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_EventEqualityImpl_AreEventsEqual_Level_respecting(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight1 := NewEvent(givenProvider, level.Info, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight2 := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	instance := EventEqualityImpl{
		CompareLevel:       true,
		CompareCallDepth:   true,
		CompareValuesUsing: nil,
	}

	actual1, actualErr1 := instance.AreEventsEqual(givenLeft, givenRight1)
	assert.ToBeNil(t, actualErr1)
	assert.ToBeEqual(t, false, actual1)

	actual2, actualErr2 := instance.AreEventsEqual(givenLeft, givenRight2)
	assert.ToBeNil(t, actualErr2)
	assert.ToBeEqual(t, true, actual2)
}

func Test_EventEqualityImpl_AreEventsEqual_Level_ignoring(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight1 := NewEvent(givenProvider, level.Info, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight2 := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	instance := EventEqualityImpl{
		CompareLevel:       false,
		CompareCallDepth:   true,
		CompareValuesUsing: nil,
	}

	actual1, actualErr1 := instance.AreEventsEqual(givenLeft, givenRight1)
	assert.ToBeNil(t, actualErr1)
	assert.ToBeEqual(t, true, actual1)

	actual2, actualErr2 := instance.AreEventsEqual(givenLeft, givenRight2)
	assert.ToBeNil(t, actualErr2)
	assert.ToBeEqual(t, true, actual2)
}

func Test_EventEqualityImpl_AreEventsEqual_callDepth_respecting(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight1 := NewEvent(givenProvider, level.Error, 4).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight2 := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	instance := EventEqualityImpl{
		CompareLevel:       true,
		CompareCallDepth:   true,
		CompareValuesUsing: nil,
	}

	actual1, actualErr1 := instance.AreEventsEqual(givenLeft, givenRight1)
	assert.ToBeNil(t, actualErr1)
	assert.ToBeEqual(t, false, actual1)

	actual2, actualErr2 := instance.AreEventsEqual(givenLeft, givenRight2)
	assert.ToBeNil(t, actualErr2)
	assert.ToBeEqual(t, true, actual2)
}

func Test_EventEqualityImpl_AreEventsEqual_callDepth_ignoring(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight1 := NewEvent(givenProvider, level.Error, 4).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight2 := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	instance := EventEqualityImpl{
		CompareLevel:       true,
		CompareCallDepth:   false,
		CompareValuesUsing: nil,
	}

	actual1, actualErr1 := instance.AreEventsEqual(givenLeft, givenRight1)
	assert.ToBeNil(t, actualErr1)
	assert.ToBeEqual(t, true, actual1)

	actual2, actualErr2 := instance.AreEventsEqual(givenLeft, givenRight2)
	assert.ToBeNil(t, actualErr2)
	assert.ToBeEqual(t, true, actual2)
}

func Test_EventEqualityImpl_AreEventsEqual_moreEntries(t *testing.T) {
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
	instance := EventEqualityImpl{
		CompareLevel:       true,
		CompareCallDepth:   true,
		CompareValuesUsing: fields.DefaultValueEquality,
	}

	actual, actualErr := instance.AreEventsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_EventEqualityImpl_AreEventsEqual_lessEntries(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2)
	instance := EventEqualityImpl{
		CompareLevel:       true,
		CompareCallDepth:   true,
		CompareValuesUsing: fields.DefaultValueEquality,
	}

	actual, actualErr := instance.AreEventsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_EventEqualityImpl_AreEventsEqual_oneDifferentValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 666)
	instance := EventEqualityImpl{
		CompareLevel:       true,
		CompareCallDepth:   true,
		CompareValuesUsing: fields.DefaultValueEquality,
	}

	actual, actualErr := instance.AreEventsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_EventEqualityImpl_AreEventsEqual_oneDifferentKeys(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("xyz", 3)
	instance := EventEqualityImpl{
		CompareLevel:       true,
		CompareCallDepth:   true,
		CompareValuesUsing: fields.DefaultValueEquality,
	}

	actual, actualErr := instance.AreEventsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_EventEqualityImpl_AreEventsEqual_leftEmptyRightNil(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenLeft := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("c", 3)
	var givenRight Event
	instance := EventEqualityImpl{
		CompareLevel:       true,
		CompareCallDepth:   true,
		CompareValuesUsing: fields.DefaultValueEquality,
	}

	actual, actualErr := instance.AreEventsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_EventEqualityImpl_AreEventsEqual_leftNilRightEmpty(t *testing.T) {
	givenProvider := newMockProvider("test")
	var givenLeft Event
	givenRight := NewEvent(givenProvider, level.Error, 3).
		With("a", 1).
		With("b", 2).
		With("xyz", 3)

	actual, actualErr := AreEventsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_EventEqualityImpl_AreEventsEqual_bothNil(t *testing.T) {
	var givenLeft Event
	var givenRight Event
	instance := EventEqualityImpl{
		CompareLevel:       true,
		CompareCallDepth:   true,
		CompareValuesUsing: fields.DefaultValueEquality,
	}

	actual, actualErr := instance.AreEventsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_EventEqualityImpl_AreEventsEqual_functionWithErr(t *testing.T) {
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
	instance := EventEqualityImpl{
		CompareLevel:     true,
		CompareCallDepth: true,
		CompareValuesUsing: fields.ValueEqualityFunc(func(key string, left, right interface{}) (bool, error) {
			return false, givenErr
		}),
	}

	actual, actualErr := instance.AreEventsEqual(givenLeft, givenRight)

	assert.ToBeEqual(t, givenErr, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_EventEqualityImpl_AreEventsEqual_WithIgnoringKeys(t *testing.T) {
	instance := &EventEqualityImpl{}

	actual := instance.WithIgnoringKeys("a", "b")

	assert.ToBeEqual(t, &ignoringKeysEventEquality{instance, []string{"a", "b"}}, actual)
}

func Test_EventEqualityFunc_AreEventsEqual(t *testing.T) {
	givenProvider := newMockProvider("foo")
	givenLeft := NewEvent(givenProvider, level.Info, 1)
	givenRight := NewEvent(givenProvider, level.Warn, 2)
	givenErr := errors.New("expected")

	instance := EventEqualityFunc(func(left, right Event) (bool, error) {
		assert.ToBeSame(t, givenLeft, left)
		assert.ToBeSame(t, givenRight, right)
		return true, givenErr
	})

	actual, actualErr := instance.AreEventsEqual(givenLeft, givenRight)

	assert.ToBeSame(t, givenErr, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_EventEqualityFunc_AreEventsEqual_WithIgnoringKeys(t *testing.T) {
	instance := EventEqualityFunc(func(left, right Event) (bool, error) {
		panic("should never be called")
	})

	actual := instance.WithIgnoringKeys("a", "b")

	assert.ToBeOfType(t, &ignoringKeysEventEquality{}, actual)
	assert.ToBeSame(t, instance, actual.(*ignoringKeysEventEquality).parent)
	assert.ToBeEqual(t, []string{"a", "b"}, actual.(*ignoringKeysEventEquality).keysToIgnore)
}

func Test_NewEventEqualityFacade(t *testing.T) {
	givenFunc := func() EventEquality { panic("should never be called") }

	actual := NewEventEqualityFacade(givenFunc)

	assert.ToBeSame(t, eventEqualityFacade(givenFunc), actual)
}

func Test_eventEqualityFacade_AreEventsEqual(t *testing.T) {
	givenProvider := newMockProvider("foo")
	givenLeft := NewEvent(givenProvider, level.Info, 1)
	givenRight := NewEvent(givenProvider, level.Warn, 2)
	givenErr := errors.New("expected")
	givenEquality := EventEqualityFunc(func(left, right Event) (bool, error) {
		assert.ToBeSame(t, givenLeft, left)
		assert.ToBeSame(t, givenRight, right)
		return true, givenErr
	})

	instance := eventEqualityFacade(func() EventEquality { return givenEquality })

	actual, actualErr := instance.AreEventsEqual(givenLeft, givenRight)

	assert.ToBeSame(t, givenErr, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_eventEqualityFacade_AreEventsEqual_WithIgnoringKeys(t *testing.T) {
	instance := eventEqualityFacade(func() EventEquality {
		panic("should never be called")
	})

	actual := instance.WithIgnoringKeys("a", "b")

	assert.ToBeOfType(t, &ignoringKeysEventEquality{}, actual)
	assert.ToBeSame(t, instance, actual.(*ignoringKeysEventEquality).parent)
	assert.ToBeEqual(t, []string{"a", "b"}, actual.(*ignoringKeysEventEquality).keysToIgnore)
}

func Test_ignoringKeysEventEquality_AreEventsEqual(t *testing.T) {
	givenProvider := newMockProvider("foo")
	givenLeft := NewEvent(givenProvider, level.Info, 1).
		With("foo", 1).
		With("bar", 2)
	givenRight := NewEvent(givenProvider, level.Warn, 2).
		With("foo", 1).
		With("bar", 2)
	givenError := errors.New("expected")

	givenEquality := EventEqualityFunc(func(left, right Event) (bool, error) {
		assert.ToBeEqualUsing(t, givenLeft.Without("foo"), left, AreEventsEqual)
		assert.ToBeEqualUsing(t, givenRight.Without("foo"), right, AreEventsEqual)
		return true, givenError
	})

	instance := &ignoringKeysEventEquality{givenEquality, []string{"foo", "xyz"}}

	actual, actualErr := instance.AreEventsEqual(givenLeft, givenRight)

	assert.ToBeSame(t, givenError, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_ignoringKeysEventEquality_AreEventsEqual_bothNil(t *testing.T) {
	var givenLeft Event
	var givenRight Event

	givenEquality := EventEqualityFunc(func(left, right Event) (bool, error) {
		panic("should never be called")
	})

	instance := &ignoringKeysEventEquality{givenEquality, []string{"foo", "xyz"}}

	actual, actualErr := instance.AreEventsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, true, actual)
}

func Test_ignoringKeysEventEquality_AreEventsEqual_leftNil(t *testing.T) {
	givenProvider := newMockProvider("foo")
	givenLeft := NewEvent(givenProvider, level.Info, 1).
		With("foo", 1).
		With("bar", 2)
	var givenRight Event

	givenEquality := EventEqualityFunc(func(left, right Event) (bool, error) {
		panic("should never be called")
	})

	instance := &ignoringKeysEventEquality{givenEquality, []string{"foo", "xyz"}}

	actual, actualErr := instance.AreEventsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_ignoringKeysEventEquality_AreEventsEqual_rightNil(t *testing.T) {
	givenProvider := newMockProvider("foo")
	var givenLeft Event
	givenRight := NewEvent(givenProvider, level.Warn, 2).
		With("foo", 1).
		With("bar", 2)

	givenEquality := EventEqualityFunc(func(left, right Event) (bool, error) {
		panic("should never be called")
	})

	instance := &ignoringKeysEventEquality{givenEquality, []string{"foo", "xyz"}}

	actual, actualErr := instance.AreEventsEqual(givenLeft, givenRight)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, false, actual)
}

func Test_ignoringKeysEventEquality_AreEventsEqual_WithIgnoringKeys(t *testing.T) {
	delegate := &EventEqualityImpl{}
	instance := &ignoringKeysEventEquality{delegate, []string{"a", "b"}}

	actual := instance.WithIgnoringKeys("c", "d")

	assert.ToBeOfType(t, &ignoringKeysEventEquality{}, actual)
	assert.ToBeSame(t, delegate, actual.(*ignoringKeysEventEquality).parent)
	assert.ToBeEqual(t, []string{"a", "b", "c", "d"}, actual.(*ignoringKeysEventEquality).keysToIgnore)
}
