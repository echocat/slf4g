package log

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_fallbackCoreLogger_Log(t *testing.T) {
	instance, buf := newFallbackCoreLogger("foo")

	t2 := time.Now()
	t1 := t2.Add(-2 * time.Minute)

	// WARNING! Do not move these lines, because the test relies on it.
	// I know this could be better... ;-)
	instance.Log(NewEvent(instance.GetProvider(), level.Trace, 0).
		With("a", 11).
		With("b", 12).
		With(instance.GetProvider().GetFieldKeysSpec().GetTimestamp(), t1))
	instance.Log(NewEvent(instance.GetProvider(), level.Info, 0).
		With("a", 11).
		With("b", 12).
		With(instance.GetProvider().GetFieldKeysSpec().GetTimestamp(), t1))
	instance.Log(NewEvent(instance.GetProvider(), level.Error, 0).
		With("a", 21).
		With("c", 23).
		With(instance.GetProvider().GetFieldKeysSpec().GetTimestamp(), t2))

	assert.ToBeEqual(t, fmt.Sprintf(strings.TrimLeft(`
I%s %d logger_core_fallback_test.go:27] a=11 b=12 logger="foo"
E%s %d logger_core_fallback_test.go:31] a=21 c=23 logger="foo"
`, "\n"),
		t1.Format(simpleTimeLayout), pid,
		t2.Format(simpleTimeLayout), pid,
	), buf.String())
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
