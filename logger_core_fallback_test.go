package log

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_fallbackCoreLogger_Log(t *testing.T) {
	givenError := errors.New("expected")
	instance, buf := newFallbackCoreLogger("foo")

	t2 := time.Now()
	t1 := t2.Add(-2 * time.Minute)

	// WARNING! Do not move these lines, because the test relies on it.
	// I know this could be better... ;-)
	instance.Log(NewEvent(instance.GetProvider(), level.Trace, 0).
		With("a", 11).
		With("b", 12).
		With("message", "hello").
		With("timestamp", t1))
	instance.Log(NewEvent(instance.GetProvider(), level.Info, 0).
		With("a", 11).
		With("b", 12).
		With("timestamp", t1).
		With("logger", fallbackRootLoggerName))
	instance.Log(NewEvent(instance.GetProvider(), level.Error, 0).
		With("a", 21).
		With("c", 23).
		With("message", "  hello    ").
		WithError(givenError).
		With("timestamp", t2))

	assert.ToBeEqual(t, fmt.Sprintf(strings.TrimLeft(`
I%s %d logger_core_fallback_test.go:30] a=11 b=12
E%s %d logger_core_fallback_test.go:35]   hello a=21 c=23 error="expected" logger="foo"
`, "\n"),
		t1.Format(simpleTimeLayout), pid,
		t2.Format(simpleTimeLayout), pid,
	), buf.String())
}

func Test_fallbackCoreLogger_Log_withoutTimestamp(t *testing.T) {
	instance, buf := newFallbackCoreLogger("foo")

	instance.Log(NewEvent(instance.GetProvider(), level.Info, 0))

	assert.ToBeMatching(t, `^I\d{2}\d{2} \d{2}:\d{2}:\d{2}\.\d{6} \d+ logger_core_fallback_test.go:\d+] logger="foo"`, buf.String())
}

func Test_fallbackCoreLogger_Log_withLazyValue(t *testing.T) {
	instance, buf := newFallbackCoreLogger("foo")

	instance.Log(NewEvent(instance.GetProvider(), level.Info, 0).
		With("foo", lazyMock(666)))

	assert.ToBeMatching(t, `^I.+logger_core_fallback_test.go:\d+] foo=666 logger="foo"`, buf.String())
}

func Test_fallbackCoreLogger_Log_brokenCallDepth(t *testing.T) {
	instance, buf := newFallbackCoreLogger("foo")

	instance.Log(NewEvent(instance.GetProvider(), level.Info, 10000))

	assert.ToBeMatching(t, `^I.+ \d+ \?\?\?:1] logger="foo"`, buf.String())
}

func Test_fallbackCoreLogger_Log_withErrorWhileMarshalling(t *testing.T) {
	instance, buf := newFallbackCoreLogger("foo")

	instance.Log(NewEvent(instance.GetProvider(), level.Info, 0).
		With("foo", failingJsonMarshalling("expected")))

	assert.ToBeMatching(t, `^ERR!! Cannot format event.+: expected`, buf.String())
}

func Test_fallbackCoreLogger_Log_levels(t *testing.T) {
	cases := []struct {
		expectedC string
		level     level.Level
	}{
		{"T", level.Trace},
		{"D", level.Debug},
		{"I", level.Info},
		{"W", level.Warn},
		{"E", level.Error},
		{"F", level.Fatal},
		{"?", level.Level(6666)},
	}
	for _, c := range cases {
		t.Run(c.expectedC, func(t *testing.T) {
			instance, buf := newFallbackCoreLogger("foo")
			instance.level = 1

			instance.Log(NewEvent(instance.GetProvider(), c.level, 0))

			assert.ToBeMatching(t, `^`+c.expectedC+`\d{2}\d{2} \d{2}:\d{2}:\d{2}\.\d{6} \d+ logger_core_fallback_test.go:\d+] logger="foo"`, buf.String())
		})
	}
}

func Test_fallbackCoreLogger_GetName(t *testing.T) {
	instance := &fallbackCoreLogger{name: "foo"}

	actual := instance.GetName()

	assert.ToBeEqual(t, "foo", actual)
}

func Test_fallbackCoreLogger_GetProvider(t *testing.T) {
	givenProvider := &fallbackProvider{}
	instance := &fallbackCoreLogger{fallbackProvider: givenProvider}

	actual := instance.GetProvider()

	assert.ToBeSame(t, givenProvider, actual)
}

func Test_fallbackCoreLogger_GetLevel(t *testing.T) {
	instance, _ := newFallbackCoreLogger("foo")

	assert.ToBeEqual(t, level.Info, instance.GetLevel())
	instance.fallbackProvider.level = level.Warn
	assert.ToBeEqual(t, level.Warn, instance.GetLevel())

	for _, l := range instance.GetAllLevels() {
		instance.level = l
		assert.ToBeEqual(t, l, instance.GetLevel())
	}

	instance.level = 0
	assert.ToBeEqual(t, level.Warn, instance.GetLevel())
}

func Test_fallbackCoreLogger_SetLevel(t *testing.T) {
	instance, _ := newFallbackCoreLogger("foo")

	assert.ToBeEqual(t, level.Level(0), instance.level)

	for _, l := range instance.GetAllLevels() {
		instance.SetLevel(l)
		assert.ToBeEqual(t, l, instance.level)
	}

	instance.SetLevel(0)
	assert.ToBeEqual(t, level.Level(0), instance.level)
}

func Test_IsFallbackLogger(t *testing.T) {
	givenMockedProvider := newMockProvider("foo")
	givenMockedProvider.rootProvider = func() Logger {
		return newMockLogger("root")
	}

	previous := SetProvider(nil)
	defer SetProvider(previous)

	assert.ToBeEqual(t, true, IsFallbackLogger(GetRootLogger()))
	assert.ToBeEqual(t, []Provider{}, GetAllProviders())

	SetProvider(givenMockedProvider)

	assert.ToBeEqual(t, false, IsFallbackLogger(GetRootLogger()))

	SetProvider(nil)

	assert.ToBeEqual(t, true, IsFallbackLogger(GetRootLogger()))
}

func newFallbackCoreLogger(name string) (*fallbackCoreLogger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	provider := &fallbackProvider{
		out: buf,
	}
	provider.cache = NewLoggerCache(provider.rootFactory, provider.factory)
	return &fallbackCoreLogger{
		fallbackProvider: provider,
		name:             name,
	}, buf
}

type failingJsonMarshalling string

func (instance failingJsonMarshalling) MarshalJSON() ([]byte, error) {
	return nil, stringError(instance)
}