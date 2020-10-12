package log

import (
	"testing"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_fallbackProvider_GetName(t *testing.T) {
	actual := fallbackProviderV.GetName()

	assert.ToBeEqual(t, "fallback", actual)
}

func Test_fallbackProvider_GetRootLogger(t *testing.T) {
	actual1 := fallbackProviderV.GetRootLogger()
	actual2 := fallbackProviderV.GetRootLogger()

	assert.ToBeSame(t, actual1, actual2)
	assert.ToBeOfType(t, &loggerImpl{}, actual1)
	assert.ToBeOfType(t, &loggerImpl{}, actual2)

	assert.ToBeSame(t, fallbackProviderV, actual1.GetProvider())
	assert.ToBeEqual(t, fallbackRootLoggerName, actual1.GetName())
	assert.ToBeOfType(t, &fallbackCoreLogger{}, actual1.(*loggerImpl).coreProvider())
}

func Test_fallbackProvider_GetLogger(t *testing.T) {
	actualFoo1 := fallbackProviderV.GetLogger("foo")
	actualFoo2 := fallbackProviderV.GetLogger("foo")
	actualBar1 := fallbackProviderV.GetLogger("bar")
	actualBar2 := fallbackProviderV.GetLogger("bar")

	assert.ToBeSame(t, actualFoo1, actualFoo2)
	assert.ToBeSame(t, actualBar1, actualBar2)
	assert.ToBeNotSame(t, actualFoo1, actualBar1)
	assert.ToBeOfType(t, &loggerImpl{}, actualFoo1)
	assert.ToBeOfType(t, &loggerImpl{}, actualFoo2)
	assert.ToBeOfType(t, &loggerImpl{}, actualBar1)
	assert.ToBeOfType(t, &loggerImpl{}, actualBar2)

	assert.ToBeSame(t, fallbackProviderV, actualFoo1.GetProvider())
	assert.ToBeEqual(t, "foo", actualFoo1.GetName())
	assert.ToBeOfType(t, &fallbackCoreLogger{}, actualFoo1.(*loggerImpl).coreProvider())

	assert.ToBeSame(t, fallbackProviderV, actualBar1.GetProvider())
	assert.ToBeEqual(t, "bar", actualBar1.GetName())
	assert.ToBeOfType(t, &fallbackCoreLogger{}, actualBar1.(*loggerImpl).coreProvider())
}

func Test_fallbackProvider_GetAllLevels(t *testing.T) {
	actual := fallbackProviderV.GetAllLevels()

	assert.ToBeEqual(t, level.GetProvider().GetLevels(), actual)
}

func Test_fallbackProvider_GetFieldKeysSpec(t *testing.T) {
	actual := fallbackProviderV.GetFieldKeysSpec()

	assert.ToBeEqual(t, "error", actual.GetError())
	assert.ToBeEqual(t, "logger", actual.GetLogger())
	assert.ToBeEqual(t, "message", actual.GetMessage())
	assert.ToBeEqual(t, "timestamp", actual.GetTimestamp())
}
