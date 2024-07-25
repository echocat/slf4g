package testlog

import (
	"testing"

	log "github.com/echocat/slf4g"

	tlevel "github.com/echocat/slf4g/sdk/testlog/level"

	"github.com/echocat/slf4g/fields"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/internal/test/assert"
)

func TestProvider_GetName_default(t *testing.T) {
	instance := NewProvider(t)
	assert.ToBeEqual(t, t.Name(), instance.GetName())
}

func TestProvider_GetName_specific(t *testing.T) {
	instance := NewProvider(t, Name("foo"))
	assert.ToBeEqual(t, "foo", instance.GetName())
}

func TestProvider_GetAllLevels_default(t *testing.T) {
	instance := NewProvider(t)
	assert.ToBeEqual(t, level.GetProvider().GetLevels(), instance.GetAllLevels())
}

func TestProvider_GetAllLevels_specific(t *testing.T) {
	given := level.Levels{level.Info, level.Level(666)}
	instance := NewProvider(t, AllLevels(given))
	assert.ToBeEqual(t, given, instance.GetAllLevels())
}

func TestProvider_GetFieldKeysSpec_default(t *testing.T) {
	instance := NewProvider(t)
	assert.ToBeEqual(t, &fields.KeysSpecImpl{}, instance.GetFieldKeysSpec())
}

func TestProvider_GetFieldKeysSpec_specific(t *testing.T) {
	given := &mockFieldKeysSpec{}
	instance := NewProvider(t, FieldKeysSpec(given))
	assert.ToBeEqual(t, given, instance.GetFieldKeysSpec())
}

func TestProvider_GetLevel_default(t *testing.T) {
	instance := NewProvider(t)
	assert.ToBeEqual(t, DefaultLevel, instance.GetLevel())
}

func TestProvider_GetLevel_specific(t *testing.T) {
	given := level.Level(666)
	instance := NewProvider(t, Level(given))
	assert.ToBeEqual(t, given, instance.GetLevel())
}

func TestProvider_SetLevel(t *testing.T) {
	given := level.Level(666)

	instance := NewProvider(t)
	assert.ToBeEqual(t, DefaultLevel, instance.GetLevel())

	instance.SetLevel(given)
	assert.ToBeEqual(t, given, instance.GetLevel())
}

func TestProvider_getFailAtLevel_default(t *testing.T) {
	instance := NewProvider(t)
	assert.ToBeEqual(t, DefaultFailAtLevel, instance.getFailAtLevel())
}

func TestProvider_getFailAtLevel_specific(t *testing.T) {
	given := level.Level(666)
	instance := NewProvider(t, FailAtLevel(given))
	assert.ToBeEqual(t, given, instance.getFailAtLevel())
}

func TestProvider_getFailNowAtLevel_default(t *testing.T) {
	instance := NewProvider(t)
	assert.ToBeEqual(t, DefaultFailNowAtLevel, instance.getFailNowAtLevel())
}

func TestProvider_getFailNowAtLevel_specific(t *testing.T) {
	given := level.Level(666)
	instance := NewProvider(t, FailNowAtLevel(given))
	assert.ToBeEqual(t, given, instance.getFailNowAtLevel())
}

func TestProvider_getTimeFormat_default(t *testing.T) {
	instance := NewProvider(t)
	assert.ToBeEqual(t, DefaultTimeFormat, instance.getTimeFormat())
}

func TestProvider_getTimeFormat_specific(t *testing.T) {
	given := "15:04:05"
	instance := NewProvider(t, TimeFormat(given))
	assert.ToBeEqual(t, given, instance.getTimeFormat())
}

func TestProvider_getLevelFormatter_default(t *testing.T) {
	instance := NewProvider(t)
	assert.ToBeEqual(t, tlevel.DefaultFormatter, instance.getLevelFormatter())
}

func TestProvider_getLevelFormatter_specific(t *testing.T) {
	given := tlevel.FormatterFunc(nil)
	instance := NewProvider(t, LevelFormatter(given))
	assert.ToBeEqual(t, given, instance.getLevelFormatter())
}

func TestProvider_GetRootLogger(t *testing.T) {
	instance := NewProvider(t)

	actual := instance.GetRootLogger()
	assert.ToBeNotNil(t, actual)

	actualCoreLogger := log.UnwrapCoreLogger(actual)
	assert.ToBeOfType(t, &coreLogger{}, actualCoreLogger)

	assert.ToBeEqual(t, RootLoggerName, actualCoreLogger.GetName())
}

func TestProvider_GetLogger(t *testing.T) {
	instance := NewProvider(t)

	actualRootLogger := instance.GetRootLogger()
	assert.ToBeNotNil(t, actualRootLogger)

	actualRootCoreLogger := log.UnwrapCoreLogger(actualRootLogger)
	assert.ToBeOfType(t, &coreLogger{}, actualRootCoreLogger)

	actual := instance.GetLogger("foo")
	assert.ToBeNotNil(t, actualRootLogger)

	actualCoreLogger := log.UnwrapCoreLogger(actual)
	assert.ToBeOfType(t, &coreLoggerRenamed{}, actualCoreLogger)

	assert.ToBeEqual(t, "foo", actualCoreLogger.GetName())
	assert.ToBeSame(t, actualRootCoreLogger, actualCoreLogger.(*coreLoggerRenamed).coreLogger)
}

func TestProvider_GetLogger_rootLogger(t *testing.T) {
	instance := NewProvider(t)

	actualRootLogger := instance.GetRootLogger()
	assert.ToBeNotNil(t, actualRootLogger)

	actual := instance.GetLogger(RootLoggerName)
	assert.ToBeSame(t, actualRootLogger, actual)
}

type mockFieldKeysSpec struct{}

func (instance *mockFieldKeysSpec) GetTimestamp() string {
	panic("not implemented in tests")
}

func (instance *mockFieldKeysSpec) GetMessage() string {
	panic("not implemented in tests")
}

func (instance *mockFieldKeysSpec) GetError() string {
	panic("not implemented in tests")
}

func (instance *mockFieldKeysSpec) GetLogger() string {
	panic("not implemented in tests")
}
