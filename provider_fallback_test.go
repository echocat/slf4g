package log

import (
	"testing"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_IsFallbackProvider(t *testing.T) {
	previous := SetProvider(nil)
	defer SetProvider(previous)

	assert.ToBeEqual(t, true, IsFallbackProvider(fallbackProviderV))
	assert.ToBeEqual(t, true, IsFallbackProvider(GetProvider()))
	assert.ToBeEqual(t, []Provider{}, GetAllProviders())

	SetProvider(newMockProvider("foo"))

	assert.ToBeEqual(t, false, IsFallbackProvider(GetProvider()))

	SetProvider(nil)

	assert.ToBeEqual(t, true, IsFallbackProvider(GetProvider()))
}

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

func Test_fallbackProvider_GetLevel(t *testing.T) {
	instance := &fallbackProvider{}

	assert.ToBeEqual(t, level.Info, instance.GetLevel())

	for _, l := range instance.GetAllLevels() {
		instance.level = l
		assert.ToBeEqual(t, l, instance.GetLevel())
	}

	instance.level = 0
	assert.ToBeEqual(t, level.Info, instance.GetLevel())
}

func Test_fallbackProvider_SetLevel(t *testing.T) {
	instance := &fallbackProvider{}

	assert.ToBeEqual(t, level.Level(0), instance.level)

	for _, l := range instance.GetAllLevels() {
		instance.SetLevel(l)
		assert.ToBeEqual(t, l, instance.level)
	}

	instance.SetLevel(0)
	assert.ToBeEqual(t, level.Level(0), instance.level)
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

func Test_fallbackProvider_levelAware(t *testing.T) {
	defer func() {
		fallbackProviderV.level = 0
	}()

	actual, actualOk := level.Get(fallbackProviderV)
	assert.ToBeEqual(t, level.Info, actual)
	assert.ToBeEqual(t, true, actualOk)
	assert.ToBeEqual(t, true, level.Set(fallbackProviderV, level.Level(666)))

	actual2, actualOk2 := level.Get(fallbackProviderV)
	assert.ToBeEqual(t, level.Level(666), actual2)
	assert.ToBeEqual(t, true, actualOk2)
}

func Test_fallbackProvider_levelAwareLogger(t *testing.T) {
	defer func() {
		fallbackProviderV.level = 0
	}()

	fooLogger := fallbackProviderV.GetLogger("foo")
	actual, actualOk := level.Get(fooLogger)
	assert.ToBeEqual(t, level.Info, actual)
	assert.ToBeEqual(t, true, actualOk)
	assert.ToBeEqual(t, true, level.Set(fooLogger, level.Level(666)))

	fooLogger2 := fallbackProviderV.GetLogger("foo")
	actual2, actualOk2 := level.Get(fooLogger2)
	assert.ToBeEqual(t, level.Level(666), actual2)
	assert.ToBeEqual(t, true, actualOk2)
}
