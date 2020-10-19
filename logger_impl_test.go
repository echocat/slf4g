package log

import (
	"errors"
	"fmt"
	"testing"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_loggerImpl_Unwrap(t *testing.T) {
	givenCoreLogger := newMockCoreLogger("foo")
	instance := newLoggerImpl(givenCoreLogger)

	actual := instance.Unwrap()

	assert.ToBeSame(t, givenCoreLogger, actual)
}

func Test_loggerImpl_GetName(t *testing.T) {
	givenCoreLogger := newMockCoreLogger("foo")
	instance := newLoggerImpl(givenCoreLogger)

	actual := instance.GetName()

	assert.ToBeEqual(t, "foo", actual)
}

func Test_loggerImpl_Log(t *testing.T) {

	for _, l := range level.GetProvider().GetLevels() {
		t.Run(fmt.Sprintf("level-%d", l), func(t *testing.T) {
			givenCoreLogger := newMockCoreLogger("foo")
			givenCoreLogger.initLoggedEvents()
			instance := newLoggerImpl(givenCoreLogger)

			givenEvent := givenCoreLogger.NewEvent(l, nil).
				With("a", 1)
			expectedEvent := givenCoreLogger.NewEvent(l, nil).
				With("a", 1)

			givenCoreLogger.level = 1
			instance.Log(givenEvent, 3)
			givenCoreLogger.level = l
			instance.Log(givenEvent, 3)
			givenCoreLogger.level = l + 1
			instance.Log(givenEvent, 3)

			assert.ToBeEqual(t, 3, len(*givenCoreLogger.loggedEvents))
			assert.ToBeEqual(t, expectedEvent, givenCoreLogger.loggedEvent(0))
			assert.ToBeEqual(t, expectedEvent, givenCoreLogger.loggedEvent(1))
			assert.ToBeEqual(t, expectedEvent, givenCoreLogger.loggedEvent(2))
		})
	}
}

func Test_loggerImpl_IsLevelEnabled(t *testing.T) {
	for _, l := range level.GetProvider().GetLevels() {
		t.Run(fmt.Sprintf("level-%d", l), func(t *testing.T) {
			givenCoreLogger := newMockCoreLogger("foo")
			instance := newLoggerImpl(givenCoreLogger)
			givenCoreLogger.level = l

			assert.ToBeEqual(t, false, instance.IsLevelEnabled(l-1))
			assert.ToBeEqual(t, true, instance.IsLevelEnabled(l))
			assert.ToBeEqual(t, true, instance.IsLevelEnabled(l+1))
		})
	}
}

func Test_loggerImpl_log(t *testing.T) {
	givenLogger := newMockLogger("foo")
	cases := []struct {
		logFunc func(args ...interface{})
		level   level.Level
	}{
		{givenLogger.Trace, level.Trace},
		{givenLogger.Debug, level.Debug},
		{givenLogger.Info, level.Info},
		{givenLogger.Warn, level.Warn},
		{givenLogger.Error, level.Error},
		{givenLogger.Fatal, level.Fatal},
	}

	for _, c := range cases {
		t.Run(levelToName(c.level), func(t *testing.T) {
			givenLogger.initLoggedEvents()
			givenLogger.setLevel(1)
			messageKey := givenLogger.getFieldKeysSpec().GetMessage()

			c.logFunc()
			c.logFunc(1)
			c.logFunc(1, 2, 3)

			givenLogger.setLevel(c.level + 1)
			c.logFunc("should not appear because level disabled")

			assert.ToBeEqual(t, 3, len(givenLogger.loggedEvents()))
			assert.ToBeEqualUsing(t,
				givenLogger.NewEvent(c.level, nil),
				givenLogger.loggedEvent(0),
				AreEventsEqual,
			)
			assert.ToBeEqualUsing(t,
				givenLogger.NewEvent(c.level, nil).
					With(messageKey, 1),
				givenLogger.loggedEvent(1),
				AreEventsEqual,
			)
			assert.ToBeEqualUsing(t,
				givenLogger.NewEvent(c.level, nil).
					With(messageKey, []interface{}{1, 2, 3}),
				givenLogger.loggedEvent(2),
				AreEventsEqual,
			)
		})
	}
}

func Test_loggerImpl_logf(t *testing.T) {
	givenLogger := newMockLogger("foo")
	cases := []struct {
		logFunc func(format string, args ...interface{})
		level   level.Level
	}{
		{givenLogger.Tracef, level.Trace},
		{givenLogger.Debugf, level.Debug},
		{givenLogger.Infof, level.Info},
		{givenLogger.Warnf, level.Warn},
		{givenLogger.Errorf, level.Error},
		{givenLogger.Fatalf, level.Fatal},
	}

	for _, c := range cases {
		t.Run(levelToName(c.level), func(t *testing.T) {
			givenLogger.initLoggedEvents()
			givenLogger.setLevel(1)
			messageKey := givenLogger.getFieldKeysSpec().GetMessage()

			c.logFunc("hello")
			c.logFunc("hello %d", 1)

			givenLogger.setLevel(c.level + 1)
			c.logFunc("should not appear because level disabled")

			assert.ToBeEqual(t, 2, len(givenLogger.loggedEvents()))
			assert.ToBeEqualUsing(t,
				givenLogger.NewEvent(c.level, nil).
					With(messageKey, fields.LazyFormat("hello")),
				givenLogger.loggedEvent(0),
				AreEventsEqual,
			)
			assert.ToBeEqualUsing(t,
				givenLogger.NewEvent(c.level, nil).
					With(messageKey, fields.LazyFormat("hello %d", 1)),
				givenLogger.loggedEvent(1),
				AreEventsEqual,
			)
		})
	}
}

func Test_loggerImpl_IsEnabled(t *testing.T) {
	givenLogger := newMockLogger("foo")
	cases := []struct {
		logFunc func() bool
		level   level.Level
	}{
		{givenLogger.IsTraceEnabled, level.Trace},
		{givenLogger.IsDebugEnabled, level.Debug},
		{givenLogger.IsInfoEnabled, level.Info},
		{givenLogger.IsWarnEnabled, level.Warn},
		{givenLogger.IsErrorEnabled, level.Error},
		{givenLogger.IsFatalEnabled, level.Fatal},
	}

	for _, c := range cases {
		t.Run(levelToName(c.level), func(t *testing.T) {
			givenLogger.setLevel(1)
			assert.ToBeEqual(t, true, c.logFunc())
			givenLogger.setLevel(c.level)
			assert.ToBeEqual(t, true, c.logFunc())
			givenLogger.setLevel(10000)
			assert.ToBeEqual(t, false, c.logFunc())
		})
	}
}

func Test_loggerImpl_With(t *testing.T) {
	givenLogger := newMockLogger("foo")

	actual := givenLogger.With("a", 1).With("b", 2)

	assert.ToBeOfType(t, &loggerImpl{}, actual)
	assert.ToBeEqual(t, fields.
		With("a", 1).
		With("b", 2),
		actual.(*loggerImpl).fields)
}

func Test_loggerImpl_Withf(t *testing.T) {
	givenLogger := newMockLogger("foo")

	actual := givenLogger.Withf("a", "%d", 1).With("b", 2)

	assert.ToBeOfType(t, &loggerImpl{}, actual)
	assert.ToBeEqual(t, fields.
		Withf("a", "%d", 1).
		With("b", 2),
		actual.(*loggerImpl).fields)
}

func Test_loggerImpl_WithError(t *testing.T) {
	givenError := errors.New("expected")
	givenLogger := newMockLogger("foo")
	errorKey := givenLogger.getFieldKeysSpec().GetError()

	actual := givenLogger.WithError(givenError).With("b", 2)

	assert.ToBeOfType(t, &loggerImpl{}, actual)
	assert.ToBeEqual(t, fields.
		With(errorKey, givenError).
		With("b", 2),
		actual.(*loggerImpl).fields)
}

func Test_loggerImpl_WithAll(t *testing.T) {
	givenLogger := newMockLogger("foo")

	actual := givenLogger.WithAll(map[string]interface{}{"a": 1, "b": 2}).With("c", 3)

	assert.ToBeOfType(t, &loggerImpl{}, actual)
	assert.ToBeEqualUsing(t, fields.
		With("a", 1).
		With("b", 2).
		With("c", 3),
		actual.(*loggerImpl).fields, fields.AreEqual)
}

func Test_loggerImpl_Without(t *testing.T) {
	givenLogger := newMockLogger("foo").
		With("a", 1).
		With("b", 2).
		With("c", 3).
		With("d", 4)

	actual := givenLogger.Without("b", "d")

	assert.ToBeOfType(t, &loggerImpl{}, actual)
	assert.ToBeEqualUsing(t, fields.
		With("a", 1).
		With("c", 3),
		actual.(*loggerImpl).fields, fields.AreEqual)
}

func Test_loggerImpl_Accepts(t *testing.T) {
	givenCoreLogger := newMockCoreLogger("foo")
	instance := newLoggerImpl(givenCoreLogger)
	var givenEvent1 Event
	givenEvent2 := givenCoreLogger.NewEvent(level.Fatal, nil)

	assert.ToBeEqual(t, givenCoreLogger.Accepts(givenEvent1), instance.Accepts(givenEvent1))
	assert.ToBeEqual(t, givenCoreLogger.Accepts(givenEvent2), instance.Accepts(givenEvent2))
}

func Test_loggerImpl_NewEvent(t *testing.T) {
	givenCoreLogger := newMockCoreLogger("foo")
	instance := newLoggerImpl(givenCoreLogger)
	givenValues := map[string]interface{}{"foo": "bar"}

	assert.ToBeEqual(t, givenCoreLogger.NewEvent(level.Fatal, givenValues), instance.NewEvent(level.Fatal, givenValues))
}

func Test_loggerImpl_NewEventWithFields_usingAsMap(t *testing.T) {
	givenCoreLogger := newMockCoreLogger("foo")
	instance := newLoggerImpl(givenCoreLogger)
	giveValues := map[string]interface{}{"foo": "bar"}
	givenFields := fields.WithAll(giveValues)

	assert.ToBeEqual(t, givenCoreLogger.NewEvent(level.Fatal, giveValues), instance.NewEventWithFields(level.Fatal, givenFields))
}
func Test_loggerImpl_NewEventWithFields_usingFields(t *testing.T) {
	givenCoreLogger := &mockCoreLoggerWithNewEventWithFields{newMockCoreLogger("foo")}
	instance := newLoggerImpl(givenCoreLogger.mockCoreLogger)
	givenFields := fields.WithAll(map[string]interface{}{"foo": "bar"})

	assert.ToBeEqual(t, givenCoreLogger.NewEventWithFields(level.Fatal, givenFields), instance.NewEventWithFields(level.Fatal, givenFields))
}

func Test_loggerImpl_NewEventWithFields_panicsOnErrors(t *testing.T) {
	givenCoreLogger := newMockCoreLogger("foo")
	instance := newLoggerImpl(givenCoreLogger)

	assert.Execution(t, func() {
		instance.NewEventWithFields(level.Fatal, fields.ForEachFunc(func(func(string, interface{}) error) error {
			return errors.New("expected")
		}))
	}).WillPanicWith("^cannot make .+: expected$")
}

func newLoggerImpl(in *mockCoreLogger) *loggerImpl {
	return &loggerImpl{
		coreProvider: func() CoreLogger {
			return in
		},
		fields: fields.Empty(),
	}
}

type mockCoreLoggerWithNewEventWithFields struct {
	*mockCoreLogger
}

func (instance *mockCoreLoggerWithNewEventWithFields) NewEventWithFields(l level.Level, f fields.ForEachEnabled) Event {
	asFields, err := fields.AsFields(f)
	if err != nil {
		panic(err)
	}
	return &fallbackEvent{
		provider: instance.provider,
		level:    l,
		fields:   asFields,
	}
}
