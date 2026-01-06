package testlog

import (
	"testing"

	"github.com/echocat/slf4g/level"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/internal/test/assert"
)

func TestNewLogger(t *testing.T) {
	instance := NewLogger(t, Level(666))

	actualCoreLogger := log.UnwrapCoreLogger(instance)
	assert.ToBeOfType(t, &coreLogger{}, actualCoreLogger)

	assert.ToBeEqual(t, RootLoggerName, actualCoreLogger.GetName())
	assert.ToBeEqual(t, level.Level(666), instance.GetProvider().(*Provider).GetLevel())
}

func TestNewNamedLogger(t *testing.T) {
	instance := NewNamedLogger(t, "foo", Level(666))

	actualCoreLogger := log.UnwrapCoreLogger(instance)
	assert.ToBeOfType(t, &coreLogger{}, actualCoreLogger)

	assert.ToBeEqual(t, "foo", actualCoreLogger.GetName())
	assert.ToBeEqual(t, level.Level(666), instance.GetProvider().(*Provider).GetLevel())
}
