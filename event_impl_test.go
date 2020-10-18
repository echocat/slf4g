package log

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_eventImpl_ForEach(t *testing.T) {
	instance := &eventImpl{fields: fields.
		With("a", 1).
		With("b", 2).
		With("c", 3)}

	actual := entries{}
	actualErr := instance.ForEach(actual.consumer())

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, entries{{"c", 3}, {"b", 2}, {"a", 1}}, actual)
}

func Test_eventImpl_Get(t *testing.T) {
	instance := &eventImpl{fields: fields.
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

func Test_eventImpl_Len(t *testing.T) {
	instance := &eventImpl{fields: fields.
		With("a", 1).
		With("b", 2)}

	actual := instance.Len()

	assert.ToBeEqual(t, 2, actual)
}

func Test_eventImpl_GetLevel(t *testing.T) {
	instance := &eventImpl{level: level.Error}

	actual := instance.GetLevel()

	assert.ToBeEqual(t, level.Error, actual)
}

func Test_eventImpl_GetContext(t *testing.T) {
	givenContext := &struct{ foo string }{"bar"}
	instance := &eventImpl{context: givenContext}

	actual := instance.GetContext()

	assert.ToBeEqual(t, givenContext, actual)
}

func Test_eventImpl_With(t *testing.T) {
	expected := fields.With("a", 1).With("b", 2)
	instance := &eventImpl{fields: fields.With("a", 1)}

	actual := instance.With("b", 2)

	assert.ToBeEqualUsing(t, expected, actual.(*eventImpl).fields, fields.AreEqual)
}

func Test_eventImpl_Withf(t *testing.T) {
	expected := fields.With("a", 1).Withf("b", "%d", 2)
	instance := &eventImpl{fields: fields.With("a", 1)}

	actual := instance.Withf("b", "%d", 2)

	assert.ToBeEqualUsing(t, expected, actual.(*eventImpl).fields, fields.AreEqual)
}

func Test_eventImpl_WithError(t *testing.T) {
	givenError := errors.New("expected")
	expected := fields.With("a", 1).With("anErrorKey", givenError)
	instance := &eventImpl{
		fields:   fields.With("a", 1),
		provider: &mockProvider{fieldKeysSpec: &mockFieldKeysSpec{error: "anErrorKey"}},
	}

	actual := instance.WithError(givenError)

	assert.ToBeEqualUsing(t, expected, actual.(*eventImpl).fields, fields.AreEqual)
}

func Test_eventImpl_WithAll(t *testing.T) {
	givenMap := map[string]interface{}{
		"b": 2,
		"c": 3,
	}
	expected := fields.With("a", 1).WithAll(givenMap)
	instance := &eventImpl{fields: fields.With("a", 1)}

	actual := instance.WithAll(givenMap)

	assert.ToBeEqualUsing(t, expected, actual.(*eventImpl).fields, fields.AreEqual)
}

func Test_eventImpl_Without(t *testing.T) {
	expected := fields.With("b", 2).With("d", 4)
	instance := &eventImpl{fields: fields.
		With("a", 1).
		With("b", 2).
		With("c", 3).
		With("d", 4)}

	actual := instance.Without("a", "c")

	assert.ToBeEqualUsing(t, expected, actual.(*eventImpl).fields, fields.AreEqual)
}

func Test_eventImpl_WithContext(t *testing.T) {
	givenContextBefore := &struct{ foo string }{"bar"}
	givenContextAfter := &struct{ foo string }{"other"}
	instance := &eventImpl{context: givenContextBefore}

	actual := instance.WithContext(givenContextAfter)

	assert.ToBeSame(t, givenContextAfter, actual.(*eventImpl).context)
}
