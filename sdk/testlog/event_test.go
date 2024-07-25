package testlog

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_event_ForEach(t *testing.T) {
	instance := &event{fields: fields.
		With("a", 1).
		With("b", 2).
		With("c", 3)}

	actual := map[string]interface{}{}
	actualErr := instance.ForEach(func(key string, value interface{}) error {
		actual[key] = value
		return nil
	})

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, map[string]interface{}{"c": 3, "b": 2, "a": 1}, actual)
}

func Test_event_Get(t *testing.T) {
	instance := &event{fields: fields.
		With("a", 1).
		With("b", 2)}

	actual1, actualExists1 := instance.Get("a")
	assert.ToBeEqual(t, true, actualExists1)
	assert.ToBeEqual(t, 1, actual1)
	actual2, actualExists2 := instance.Get("b")
	assert.ToBeEqual(t, true, actualExists2)
	assert.ToBeEqual(t, 2, actual2)
	actual3, actualExists3 := instance.Get("c")
	assert.ToBeEqual(t, false, actualExists3)
	assert.ToBeEqual(t, nil, actual3)
}

func Test_event_Len(t *testing.T) {
	instance := &event{fields: fields.
		With("a", 1).
		With("b", 2)}

	actual := instance.Len()

	assert.ToBeEqual(t, 2, actual)
}

func Test_event_GetLevel(t *testing.T) {
	instance := &event{level: level.Error}

	actual := instance.GetLevel()

	assert.ToBeEqual(t, level.Error, actual)
}

func Test_event_With(t *testing.T) {
	expected := fields.With("a", 1).With("b", 2)
	instance := &event{fields: fields.With("a", 1)}

	actual := instance.With("b", 2)

	assert.ToBeEqualUsing(t, expected, actual.(*event).fields, fields.AreEqual)
}

func Test_event_Withf(t *testing.T) {
	expected := fields.With("a", 1).Withf("b", "%d", 2)
	instance := &event{fields: fields.With("a", 1)}

	actual := instance.Withf("b", "%d", 2)

	assert.ToBeEqualUsing(t, expected, actual.(*event).fields, fields.AreEqual)
}

func Test_event_WithError(t *testing.T) {
	givenError := errors.New("expected")
	expected := fields.With("a", 1).With("error", givenError)
	instance := &event{
		fields:   fields.With("a", 1),
		provider: &Provider{},
	}

	actual := instance.WithError(givenError)

	assert.ToBeEqualUsing(t, expected, actual.(*event).fields, fields.AreEqual)
}

func Test_event_WithAll(t *testing.T) {
	givenMap := map[string]interface{}{
		"b": 2,
		"c": 3,
	}
	expected := fields.With("a", 1).WithAll(givenMap)
	instance := &event{fields: fields.With("a", 1)}

	actual := instance.WithAll(givenMap)

	assert.ToBeEqualUsing(t, expected, actual.(*event).fields, fields.AreEqual)
}

func Test_event_Without(t *testing.T) {
	expected := fields.With("b", 2).With("d", 4)
	instance := &event{fields: fields.
		With("a", 1).
		With("b", 2).
		With("c", 3).
		With("d", 4)}

	actual := instance.Without("a", "c")

	assert.ToBeEqualUsing(t, expected, actual.(*event).fields, fields.AreEqual)
}
