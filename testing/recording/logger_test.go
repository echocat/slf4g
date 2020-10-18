package recording

import (
	"testing"
	"time"

	"github.com/echocat/slf4g/level"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_NewLogger(t *testing.T) {
	actual := NewLogger()

	assert.ToBeNotNil(t, actual)
	assert.ToBeNotNil(t, actual.CoreLogger)
	assert.ToBeNotNil(t, actual.Logger)

	wrapped := log.UnwrapCoreLogger(actual.Logger)
	assert.ToBeSame(t, actual.CoreLogger, wrapped)
}

func Test_Logger_GetName_empty(t *testing.T) {
	instance := NewLogger()

	actual := instance.GetName()

	assert.ToBeEqual(t, RootLoggerName, actual)
}

func Test_Logger_GetName_configured(t *testing.T) {
	instance := NewLogger()
	instance.Name = "foo"

	actual := instance.GetName()

	assert.ToBeEqual(t, "foo", actual)
}

func Test_Logger_GetProvider_empty(t *testing.T) {
	instance := NewLogger()

	actual := instance.GetProvider()

	assert.ToBeSame(t, log.GetProvider(), actual)
}

func Test_Logger_GetProvider_configured(t *testing.T) {
	instance := NewLogger()
	instance.Provider = NewProvider()

	actual := instance.GetProvider()

	assert.ToBeSame(t, instance.Provider, actual)
}

func Test_Logger_IsLevelEnabled(t *testing.T) {
	instance := NewLogger()

	assert.ToBeEqual(t, false, instance.IsLevelEnabled(level.Debug))
	assert.ToBeEqual(t, true, instance.IsLevelEnabled(level.Info))
	assert.ToBeEqual(t, true, instance.IsLevelEnabled(level.Warn))

	instance.SetLevel(level.Warn)

	assert.ToBeEqual(t, false, instance.IsLevelEnabled(level.Debug))
	assert.ToBeEqual(t, false, instance.IsLevelEnabled(level.Info))
	assert.ToBeEqual(t, true, instance.IsLevelEnabled(level.Warn))
}

func Test_Logger_LogEvent(t *testing.T) {
	instance := NewLogger()
	instance.Provider = NewProvider()

	givenEvent := instance.NewEvent(level.Warn, nil).
		With("timestamp", time.Now()).
		With("logger", "foo")
	expected := givenEvent

	assert.ToBeEqual(t, 0, instance.Len())

	instance.Log(givenEvent, 666)

	assert.ToBeEqual(t, 1, instance.Len())
	assert.ToBeEqual(t, true, instance.MustContains(expected))
}
