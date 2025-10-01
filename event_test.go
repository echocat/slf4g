package log

import (
	"testing"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
)

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

func TestNewEvent(t *testing.T) {
	someValues := map[string]interface{}{
		"foo": 1,
		"bar": 2,
	}
	aLevel := level.Fatal
	anEventFactory := eventFactory(func(l level.Level, values map[string]interface{}) Event {
		return &fallbackEvent{nil, fields.WithAll(values), l}
	})

	actual := NewEvent(anEventFactory, aLevel, someValues)
	assert.ToBeNotNil(t, actual)
	assert.ToBeEqual(t, aLevel, actual.GetLevel())
	actualValues, actualValuesErr := fields.AsMap(actual)
	assert.ToBeNoError(t, actualValuesErr)
	assert.ToBeEqual(t, someValues, actualValues)
}
func TestNewEventWithFields_direct(t *testing.T) {
	t.Run("hasFunction_NewEventWithFields", func(t *testing.T) {
		someFields := fields.WithAll(map[string]interface{}{
			"foo": 1,
			"bar": 2,
		})
		aLevel := level.Fatal
		anEventFactory := eventFactoryWithFields(func(l level.Level, fds fields.ForEachEnabled) Event {
			target, err := fields.AsFields(fds)
			assert.ToBeNoError(t, err)
			return &fallbackEvent{nil, target, l}
		})

		actual := NewEventWithFields(anEventFactory, aLevel, someFields)
		assert.ToBeNotNil(t, actual)
		assert.ToBeOfType(t, (*fallbackEvent)(nil), actual)
		cActual := actual.(*fallbackEvent)
		assert.ToBeEqual(t, aLevel, cActual.GetLevel())
		assert.ToBeEqual(t, someFields, cActual.fields)
	})

	t.Run("fallback_to_NewEvent", func(t *testing.T) {
		someFields := fields.WithAll(map[string]interface{}{
			"foo": 1,
			"bar": 2,
		})
		aLevel := level.Fatal
		anEventFactory := eventFactory(func(l level.Level, values map[string]interface{}) Event {
			return &fallbackEvent{nil, fields.WithAll(values), l}
		})

		actual := NewEventWithFields(anEventFactory, aLevel, someFields)
		assert.ToBeNotNil(t, actual)
		assert.ToBeOfType(t, (*fallbackEvent)(nil), actual)
		cActual := actual.(*fallbackEvent)
		assert.ToBeEqual(t, aLevel, cActual.GetLevel())
		assert.ToBeEqual(t, someFields, cActual.fields)
	})
}

type eventFactory func(level.Level, map[string]interface{}) Event

func (instance eventFactory) NewEvent(l level.Level, values map[string]interface{}) Event {
	return instance(l, values)
}

type eventFactoryWithFields func(level.Level, fields.ForEachEnabled) Event

func (instance eventFactoryWithFields) NewEvent(level.Level, map[string]interface{}) Event {
	panic("should not be called")
}

func (instance eventFactoryWithFields) NewEventWithFields(l level.Level, fds fields.ForEachEnabled) Event {
	return instance(l, fds)
}
