package testlog

import (
	"testing"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
)

func TestHook(t *testing.T) {
	provider := Hook(t)

	log.Info("log.Info(..)")

	provider.GetRootLogger().Info("provider.GetRootLogger().Info(..)")
	provider.GetLogger("foo").Info("provider.GetLogger(foo).Info(..)")

	fooLogger := log.GetLogger("foo")
	actual, actualOk := level.Get(fooLogger)
	assert.ToBeEqual(t, level.Debug, actual)
	assert.ToBeEqual(t, true, actualOk)
	assert.ToBeEqual(t, true, level.Set(fooLogger, level.Level(666)))

	fooLogger2 := log.GetLogger("foo")
	actual2, actualOk2 := level.Get(fooLogger2)
	assert.ToBeEqual(t, level.Level(666), actual2)
	assert.ToBeEqual(t, true, actualOk2)
}
