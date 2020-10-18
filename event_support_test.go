package log

import (
	"errors"
	"testing"
	"time"

	"github.com/echocat/slf4g/fields"

	"github.com/echocat/slf4g/internal/support"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_GetMessageOf_withNilEvent(t *testing.T) {
	givenProvider := newMockProvider("test")

	actual := GetMessageOf(nil, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetMessageOf_withNilValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info)

	actual := GetMessageOf(givenEvent, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetMessageOf_withStringValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetMessage(), "value")

	actual := GetMessageOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("value"), actual)
}

func Test_GetMessageOf_withPStringValue(t *testing.T) {
	givenValue := support.PString("value")
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetMessage(), givenValue)

	actual := GetMessageOf(givenEvent, givenProvider)

	assert.ToBeSame(t, givenValue, actual)
}

func Test_GetMessageOf_withStringerValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetMessage(), stringerMock("value"))

	actual := GetMessageOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("value"), actual)
}

func Test_GetMessageOf_withFmtValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetMessage(), 666)

	actual := GetMessageOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("666"), actual)
}

func Test_GetMessageOf_withLazyValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetMessage(), lazyMock(666))

	actual := GetMessageOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("666"), actual)
}

func Test_GetErrorOf_withNilEvent(t *testing.T) {
	givenProvider := newMockProvider("test")

	actual := GetErrorOf(nil, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetErrorOf_withNilValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info)

	actual := GetErrorOf(givenEvent, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetErrorOf_withErrorValue(t *testing.T) {
	givenError := errors.New("test")
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetError(), givenError)

	actual := GetErrorOf(givenEvent, givenProvider)

	assert.ToBeSame(t, givenError, actual)
}

func Test_GetErrorOf_withStringValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetError(), "value")

	actual := GetErrorOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, stringError("value"), actual)
}

func Test_GetErrorOf_withPStringValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetError(), support.PString("value"))

	actual := GetErrorOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, stringError("value"), actual)
}

func Test_GetErrorOf_withStringerValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetError(), stringerMock("value"))

	actual := GetErrorOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, stringError("value"), actual)
}

func Test_GetErrorOf_withFmtValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetError(), 666)

	actual := GetErrorOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, stringError("666"), actual)
}

func Test_GetErrorOf_withLazyValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetError(), lazyMock(666))

	actual := GetErrorOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, stringError("666"), actual)
}

func Test_GetTimestampOf_withNilEvent(t *testing.T) {
	givenProvider := newMockProvider("test")

	actual := GetTimestampOf(nil, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetTimestampOf_withNilValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info)

	actual := GetTimestampOf(givenEvent, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetTimestampOf_withTimeValue(t *testing.T) {
	givenTimestamp := time.Now()
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetTimestamp(), givenTimestamp)

	actual := GetTimestampOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, &givenTimestamp, actual)
}

func Test_GetTimestampOf_withZeroTimeValue(t *testing.T) {
	givenTimestamp := time.Time{}
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetTimestamp(), givenTimestamp)

	actual := GetTimestampOf(givenEvent, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetTimestampOf_withPTimeValue(t *testing.T) {
	givenTimestamp := support.PTime(time.Now())
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetTimestamp(), givenTimestamp)

	actual := GetTimestampOf(givenEvent, givenProvider)

	assert.ToBeSame(t, givenTimestamp, actual)
}

func Test_GetTimestampOf_withPZeroTimeValue(t *testing.T) {
	givenTimestamp := support.PTime(time.Time{})
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetTimestamp(), givenTimestamp)

	actual := GetTimestampOf(givenEvent, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetTimestampOf_withLazyValue(t *testing.T) {
	givenTimestamp := support.PTime(time.Now())
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetTimestamp(), fields.LazyFunc(func() interface{} {
			return givenTimestamp
		}))

	actual := GetTimestampOf(givenEvent, givenProvider)

	assert.ToBeSame(t, givenTimestamp, actual)
}

func Test_GetLoggerOf_withNilEvent(t *testing.T) {
	givenProvider := newMockProvider("test")

	actual := GetLoggerOf(nil, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetLoggerOf_withNilValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info)

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetLoggerOf_withLoggerValue(t *testing.T) {
	givenLogger := newMockLogger("foo")
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetLogger(), givenLogger)

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("foo"), actual)
}

func Test_GetLoggerOf_withStringValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetLogger(), "foo")

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("foo"), actual)
}

func Test_GetLoggerOf_withPStringValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetLogger(), support.PString("foo"))

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("foo"), actual)
}

func Test_GetLoggerOf_withStringerValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetLogger(), stringerMock("value"))

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("value"), actual)
}

func Test_GetLoggerOf_withNamedValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetLogger(), namedMock("value"))

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("value"), actual)
}

func Test_GetLoggerOf_withFmtValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetLogger(), 666)

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("666"), actual)
}

func Test_GetLoggerOf_withLazyValue(t *testing.T) {
	givenProvider := newMockProvider("test")
	givenEvent := NewEvent(givenProvider, level.Info).
		With(givenProvider.fieldKeysSpec.GetLogger(), lazyMock(666))

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("666"), actual)
}

func Test_stringError_Error(t *testing.T) {
	instance := stringError("foo")

	actual := instance.Error()

	assert.ToBeEqual(t, "foo", actual)
}

func Test_stringError_String(t *testing.T) {
	instance := stringError("foo")

	actual := instance.String()

	assert.ToBeEqual(t, "foo", actual)
}

type lazyMock int

func (instance lazyMock) Get() interface{} {
	return int(instance)
}

type stringerMock string

func (instance stringerMock) String() string {
	return string(instance)
}

type namedMock string

func (instance namedMock) GetName() string {
	return string(instance)
}
