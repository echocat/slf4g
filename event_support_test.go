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
	givenProvider := newMockProvider("test").withRootLogger()

	actual := GetMessageOf(nil, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetMessageOf_withNilValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info)

	actual := GetMessageOf(givenEvent, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetMessageOf_withStringValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetMessage(), "value")

	actual := GetMessageOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("value"), actual)
}

func Test_GetMessageOf_withPStringValue(t *testing.T) {
	givenValue := support.PString("value")
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetMessage(), givenValue)

	actual := GetMessageOf(givenEvent, givenProvider)

	assert.ToBeSame(t, givenValue, actual)
}

func Test_GetMessageOf_withStringerValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetMessage(), stringerMock("value"))

	actual := GetMessageOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("value"), actual)
}

func Test_GetMessageOf_withFmtValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetMessage(), 666)

	actual := GetMessageOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("666"), actual)
}

func Test_GetMessageOf_withStringSliceValue(t *testing.T) {
	cases := []struct {
		given    []string
		expected string
	}{
		{[]string{"foo", "bar", "xyz"}, "foo bar xyz"},
		{[]string{"foo", "bar"}, "foo bar"},
		{[]string{"foo"}, "foo"},
		{[]string{}, ""},
	}
	givenProvider := newMockProvider("test").withRootLogger()

	for _, c := range cases {
		t.Run(c.expected, func(t *testing.T) {
			givenEvent := givenProvider.newEvent(level.Info).
				With(givenProvider.fieldKeysSpec.GetMessage(), c.given)

			actual := GetMessageOf(givenEvent, givenProvider)

			assert.ToBeEqual(t, support.PString(c.expected), actual)
		})
	}
}

func Test_GetMessageOf_withAnySliceValue(t *testing.T) {
	cases := []struct {
		given    []any
		expected string
	}{
		{[]any{"foo", "bar", "xyz"}, "foo bar xyz"},
		{[]any{"foo", "bar"}, "foo bar"},
		{[]any{"foo"}, "foo"},
		{[]any{}, ""},
		{[]any{"foo", 1, "bar"}, "foo 1 bar"},
	}
	givenProvider := newMockProvider("test").withRootLogger()

	for _, c := range cases {
		t.Run(c.expected, func(t *testing.T) {
			givenEvent := givenProvider.newEvent(level.Info).
				With(givenProvider.fieldKeysSpec.GetMessage(), c.given)

			actual := GetMessageOf(givenEvent, givenProvider)

			assert.ToBeEqual(t, support.PString(c.expected), actual)
		})
	}
}

func Test_GetMessageOf_withLazyValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetMessage(), lazyMock(666))

	actual := GetMessageOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("666"), actual)
}

func Test_GetMessageOf_withFilteredValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()

	cases := []struct {
		name     string
		filtered fields.Filtered
		expected *string
	}{
		{"respected", fields.RequireMaximalLevel(level.Info, "value"), support.PString("value")},
		{"ignored", fields.RequireMaximalLevel(level.Debug, "secret"), nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			givenEvent := givenProvider.newEvent(level.Info).
				With(givenProvider.fieldKeysSpec.GetMessage(), c.filtered)

			actual := GetMessageOf(givenEvent, givenProvider)

			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}

func Test_GetMessageOf_withFilteredLazyValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()

	cases := []struct {
		name          string
		maximalLevel  level.Level
		expected      *string
		expectedCalls int
	}{
		{"respected", level.Info, support.PString("value"), 1},
		{"ignored", level.Debug, nil, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			calls := 0
			filtered := fields.RequireMaximalLevelLazy(c.maximalLevel, fields.LazyFunc(func() any {
				calls++
				return "value"
			}))
			givenEvent := givenProvider.newEvent(level.Info).
				With(givenProvider.fieldKeysSpec.GetMessage(), filtered)

			actual := GetMessageOf(givenEvent, givenProvider)

			assert.ToBeEqual(t, c.expected, actual)
			assert.ToBeEqual(t, c.expectedCalls, calls)
		})
	}
}

func Test_EventSupport_withTypedNilValues(t *testing.T) {
	var (
		message   *nilLazyMock
		givenErr  *nilErrorMock
		timestamp *time.Time
		logger    *nilNamedMock
	)
	givenProvider := newMockProvider("test").withRootLogger()
	givenSpec := givenProvider.GetFieldKeysSpec()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenSpec.GetMessage(), message).
		With(givenSpec.GetError(), givenErr).
		With(givenSpec.GetTimestamp(), timestamp).
		With(givenSpec.GetLogger(), logger)

	assert.ToBeNil(t, GetMessageOf(givenEvent, givenProvider))
	assert.ToBeNil(t, GetErrorOf(givenEvent, givenProvider))
	assert.ToBeNil(t, GetTimestampOf(givenEvent, givenProvider))
	assert.ToBeNil(t, GetLoggerOf(givenEvent, givenProvider))
}

func Test_GetMessageOf_withTypedNilFilteredValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetMessage(), (*nilFilteredMock)(nil))

	actual := GetMessageOf(givenEvent, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetMessageOf_withLazyTypedNilResult(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetMessage(), fields.LazyFunc(func() any {
			var result *string
			return result
		}))

	actual := GetMessageOf(givenEvent, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetErrorOf_withNilEvent(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()

	actual := GetErrorOf(nil, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetErrorOf_withNilValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info)

	actual := GetErrorOf(givenEvent, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetErrorOf_withErrorValue(t *testing.T) {
	givenError := errors.New("test")
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetError(), givenError)

	actual := GetErrorOf(givenEvent, givenProvider)

	assert.ToBeSame(t, givenError, actual)
}

func Test_GetErrorOf_withStringValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetError(), "value")

	actual := GetErrorOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, stringError("value"), actual)
}

func Test_GetErrorOf_withPStringValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetError(), support.PString("value"))

	actual := GetErrorOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, stringError("value"), actual)
}

func Test_GetErrorOf_withStringerValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetError(), stringerMock("value"))

	actual := GetErrorOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, stringError("value"), actual)
}

func Test_GetErrorOf_withFmtValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetError(), 666)

	actual := GetErrorOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, stringError("666"), actual)
}

func Test_GetErrorOf_withLazyValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetError(), lazyMock(666))

	actual := GetErrorOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, stringError("666"), actual)
}

func Test_GetErrorOf_withFilteredValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenError := errors.New("value")

	cases := []struct {
		name     string
		filtered fields.Filtered
		expected error
	}{
		{"respected", fields.RequireMaximalLevel(level.Info, givenError), givenError},
		{"ignored", fields.RequireMaximalLevel(level.Debug, errors.New("secret")), nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			givenEvent := givenProvider.newEvent(level.Info).
				With(givenProvider.fieldKeysSpec.GetError(), c.filtered)

			actual := GetErrorOf(givenEvent, givenProvider)

			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}

func Test_GetTimestampOf_withNilEvent(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()

	actual := GetTimestampOf(nil, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetTimestampOf_withNilValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info)

	actual := GetTimestampOf(givenEvent, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetTimestampOf_withTimeValue(t *testing.T) {
	givenTimestamp := time.Now()
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetTimestamp(), givenTimestamp)

	actual := GetTimestampOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, &givenTimestamp, actual)
}

func Test_GetTimestampOf_withZeroTimeValue(t *testing.T) {
	givenTimestamp := time.Time{}
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetTimestamp(), givenTimestamp)

	actual := GetTimestampOf(givenEvent, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetTimestampOf_withPTimeValue(t *testing.T) {
	givenTimestamp := support.PTime(time.Now())
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetTimestamp(), givenTimestamp)

	actual := GetTimestampOf(givenEvent, givenProvider)

	assert.ToBeSame(t, givenTimestamp, actual)
}

func Test_GetTimestampOf_withPZeroTimeValue(t *testing.T) {
	givenTimestamp := support.PTime(time.Time{})
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetTimestamp(), givenTimestamp)

	actual := GetTimestampOf(givenEvent, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetTimestampOf_withLazyValue(t *testing.T) {
	givenTimestamp := support.PTime(time.Now())
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetTimestamp(), fields.LazyFunc(func() any {
			return givenTimestamp
		}))

	actual := GetTimestampOf(givenEvent, givenProvider)

	assert.ToBeSame(t, givenTimestamp, actual)
}

func Test_GetTimestampOf_withFilteredValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenTimestamp := time.Now()

	cases := []struct {
		name     string
		filtered fields.Filtered
		expected *time.Time
	}{
		{"respected", fields.RequireMaximalLevel(level.Info, givenTimestamp), &givenTimestamp},
		{"ignored", fields.RequireMaximalLevel(level.Debug, givenTimestamp), nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			givenEvent := givenProvider.newEvent(level.Info).
				With(givenProvider.fieldKeysSpec.GetTimestamp(), c.filtered)

			actual := GetTimestampOf(givenEvent, givenProvider)

			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}

func Test_GetLoggerOf_withNilEvent(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()

	actual := GetLoggerOf(nil, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetLoggerOf_withNilValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info)

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeNil(t, actual)
}

func Test_GetLoggerOf_withLoggerValue(t *testing.T) {
	givenLogger := newMockLogger("foo")
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetLogger(), givenLogger)

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("foo"), actual)
}

func Test_GetLoggerOf_withStringValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetLogger(), "foo")

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("foo"), actual)
}

func Test_GetLoggerOf_withPStringValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetLogger(), support.PString("foo"))

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("foo"), actual)
}

func Test_GetLoggerOf_withStringerValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetLogger(), stringerMock("value"))

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("value"), actual)
}

func Test_GetLoggerOf_withNamedValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetLogger(), namedMock("value"))

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("value"), actual)
}

func Test_GetLoggerOf_withFmtValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetLogger(), 666)

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("666"), actual)
}

func Test_GetLoggerOf_withLazyValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()
	givenEvent := givenProvider.newEvent(level.Info).
		With(givenProvider.fieldKeysSpec.GetLogger(), lazyMock(666))

	actual := GetLoggerOf(givenEvent, givenProvider)

	assert.ToBeEqual(t, support.PString("666"), actual)
}

func Test_GetLoggerOf_withFilteredValue(t *testing.T) {
	givenProvider := newMockProvider("test").withRootLogger()

	cases := []struct {
		name     string
		filtered fields.Filtered
		expected *string
	}{
		{"respected", fields.RequireMaximalLevel(level.Info, "value"), support.PString("value")},
		{"ignored", fields.RequireMaximalLevel(level.Debug, "secret"), nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			givenEvent := givenProvider.newEvent(level.Info).
				With(givenProvider.fieldKeysSpec.GetLogger(), c.filtered)

			actual := GetLoggerOf(givenEvent, givenProvider)

			assert.ToBeEqual(t, c.expected, actual)
		})
	}
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

func (instance lazyMock) Get() any {
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

type nilLazyMock struct{}

func (*nilLazyMock) Get() any {
	panic("must not be called")
}

type nilErrorMock struct{}

func (*nilErrorMock) Error() string {
	panic("must not be called")
}

type nilFilteredMock struct{}

func (*nilFilteredMock) Get() any {
	panic("must not be called")
}

func (*nilFilteredMock) Filter(fields.FilterContext) (any, bool) {
	panic("must not be called")
}

type nilNamedMock struct{}

func (*nilNamedMock) GetName() string {
	panic("must not be called")
}
