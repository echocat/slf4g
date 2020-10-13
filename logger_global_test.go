package log

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/echocat/slf4g/fields"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_GetLogger(t *testing.T) {
	givenLogger := newMockLogger("foo")
	givenProvider := newMockProvider("bar")
	defer setProvider(givenProvider)()

	givenProvider.provider = func(name string) Logger {
		assert.ToBeEqual(t, "foo", name)
		return givenLogger
	}

	actual := GetLogger("foo")

	assert.ToBeSame(t, givenLogger, actual)
}

func Test_GetRootLogger(t *testing.T) {
	givenLogger := newMockLogger("foo")
	defer setRootLogger(givenLogger)()

	actual := GetRootLogger()

	assert.ToBeSame(t, givenLogger, actual)
}

func Test_IsLevelEnabled(t *testing.T) {
	givenLogger := newMockLogger("foo")
	defer setRootLogger(givenLogger)()

	for _, l := range level.GetProvider().GetLevels() {
		t.Run(fmt.Sprintf("level-%d", l), func(t *testing.T) {
			givenLogger.setLevel(l)

			assert.ToBeEqual(t, false, IsLevelEnabled(l-1))
			assert.ToBeEqual(t, true, IsLevelEnabled(l))
			assert.ToBeEqual(t, true, IsLevelEnabled(l+1))
		})
	}
}

func Test_Log(t *testing.T) {
	cases := []struct {
		logFunc func(args ...interface{})
		level   level.Level
	}{
		{Trace, level.Trace},
		{Debug, level.Debug},
		{Info, level.Info},
		{Warn, level.Warn},
		{Error, level.Error},
		{Fatal, level.Fatal},
	}

	for _, c := range cases {
		t.Run(levelToName(c.level), func(t *testing.T) {
			givenLogger := newMockLogger("foo")
			givenLogger.initLoggedEvents()
			givenLogger.setLevel(1)
			defer setRootLogger(givenLogger)()
			messageKey := givenLogger.getFieldKeysSpec().GetMessage()

			c.logFunc()
			c.logFunc(1)
			c.logFunc(1, 2, 3)

			givenLogger.setLevel(c.level + 1)
			c.logFunc("should not appear because level disabled")

			assert.ToBeEqual(t, 3, len(givenLogger.loggedEvents()))
			assert.ToBeEqualUsing(t,
				NewEvent(givenLogger.GetProvider(), c.level, 3),
				givenLogger.loggedEvent(0),
				AreEventsEqual,
			)
			assert.ToBeEqualUsing(t,
				NewEvent(givenLogger.GetProvider(), c.level, 3).
					With(messageKey, 1),
				givenLogger.loggedEvent(1),
				AreEventsEqual,
			)
			assert.ToBeEqualUsing(t,
				NewEvent(givenLogger.GetProvider(), c.level, 3).
					With(messageKey, []interface{}{1, 2, 3}),
				givenLogger.loggedEvent(2),
				AreEventsEqual,
			)
		})
	}
}

func Test_Logf(t *testing.T) {
	cases := []struct {
		logFunc func(fmt string, args ...interface{})
		level   level.Level
	}{
		{Tracef, level.Trace},
		{Debugf, level.Debug},
		{Infof, level.Info},
		{Warnf, level.Warn},
		{Errorf, level.Error},
		{Fatalf, level.Fatal},
	}

	for _, c := range cases {
		t.Run(levelToName(c.level), func(t *testing.T) {
			givenLogger := newMockLogger("foo")
			givenLogger.initLoggedEvents()
			givenLogger.setLevel(1)
			defer setRootLogger(givenLogger)()
			messageKey := givenLogger.getFieldKeysSpec().GetMessage()

			c.logFunc("hello")
			c.logFunc("hello %d", 1)

			givenLogger.setLevel(c.level + 1)
			c.logFunc("should not appear because level disabled")

			assert.ToBeEqual(t, 2, len(givenLogger.loggedEvents()))
			assert.ToBeEqualUsing(t,
				NewEvent(givenLogger.GetProvider(), c.level, 3).
					With(messageKey, fields.LazyFormat("hello")),
				givenLogger.loggedEvent(0),
				AreEventsEqual,
			)
			assert.ToBeEqualUsing(t,
				NewEvent(givenLogger.GetProvider(), c.level, 3).
					With(messageKey, fields.LazyFormat("hello %d", 1)),
				givenLogger.loggedEvent(1),
				AreEventsEqual,
			)
		})
	}
}

func Test_IsEnabled(t *testing.T) {
	cases := []struct {
		logFunc func() bool
		level   level.Level
	}{
		{IsTraceEnabled, level.Trace},
		{IsDebugEnabled, level.Debug},
		{IsInfoEnabled, level.Info},
		{IsWarnEnabled, level.Warn},
		{IsErrorEnabled, level.Error},
		{IsFatalEnabled, level.Fatal},
	}

	for _, c := range cases {
		t.Run(levelToName(c.level), func(t *testing.T) {
			givenLogger := newMockLogger("foo")
			defer setRootLogger(givenLogger)()

			givenLogger.setLevel(1)
			assert.ToBeEqual(t, true, c.logFunc())
			givenLogger.setLevel(c.level)
			assert.ToBeEqual(t, true, c.logFunc())
			givenLogger.setLevel(10000)
			assert.ToBeEqual(t, false, c.logFunc())
		})
	}
}

func Test_With(t *testing.T) {
	givenLogger := newMockLogger("foo")
	defer setRootLogger(givenLogger)()

	actual := With("a", 1).With("b", 2)

	assert.ToBeOfType(t, &loggerImpl{}, actual)
	assert.ToBeEqual(t, fields.
		With("a", 1).
		With("b", 2),
		actual.(*loggerImpl).fields)
}

func Test_Withf(t *testing.T) {
	givenLogger := newMockLogger("foo")
	defer setRootLogger(givenLogger)()

	actual := Withf("a", "%d", 1).With("b", 2)

	assert.ToBeOfType(t, &loggerImpl{}, actual)
	assert.ToBeEqual(t, fields.
		Withf("a", "%d", 1).
		With("b", 2),
		actual.(*loggerImpl).fields)
}

func Test_WithError(t *testing.T) {
	givenError := errors.New("expected")
	givenLogger := newMockLogger("foo")
	defer setRootLogger(givenLogger)()
	errorKey := givenLogger.getFieldKeysSpec().GetError()

	actual := WithError(givenError).With("b", 2)

	assert.ToBeOfType(t, &loggerImpl{}, actual)
	assert.ToBeEqual(t, fields.
		With(errorKey, givenError).
		With("b", 2),
		actual.(*loggerImpl).fields)
}

func Test_WithAll(t *testing.T) {
	givenLogger := newMockLogger("foo")
	defer setRootLogger(givenLogger)()

	actual := WithAll(map[string]interface{}{"a": 1, "b": 2}).With("c", 3)

	assert.ToBeOfType(t, &loggerImpl{}, actual)
	assert.ToBeEqualUsing(t, fields.
		With("a", 1).
		With("b", 2).
		With("c", 3),
		actual.(*loggerImpl).fields, fields.AreEqual)
}

var providerVLock = new(sync.Mutex)

func setProvider(to Provider) func() {
	providerVLock.Lock()
	oldProviderV := providerV

	oldProvider := (*Provider)(atomic.LoadPointer(&providerPointer))

	SetProvider(to)
	providerV = to
	return func() {
		defer providerVLock.Unlock()
		providerV = oldProviderV
		if oldProvider != nil && *oldProvider != nil {
			SetProvider(*oldProvider)
		} else {
			SetProvider(nil)
		}
	}
}

func setRootLogger(to *mockLogger) func() {
	provider := to.getProvider()
	provider.rootProvider = func() Logger {
		return to
	}
	return setProvider(provider)
}

func levelToName(in level.Level) string {
	switch in {
	case level.Trace:
		return "Trace"
	case level.Debug:
		return "Debug"
	case level.Info:
		return "Info"
	case level.Warn:
		return "Warn"
	case level.Error:
		return "Error"
	case level.Fatal:
		return "Fatal"
	default:
		panic(fmt.Sprintf("unknown level: %d", in))
	}
}
