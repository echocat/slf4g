//go:build go1.21

package sdk

import (
	"fmt"
	sdk "log/slog"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
)

func TestNewLevelMapper(t *testing.T) {
	actualA := NewLevelMapper()
	assert.ToBeNotNil(t, actualA)

	actualB := NewLevelMapper()
	assert.ToBeNotNil(t, actualB)

	assert.ToBeOfType(t, (*defaultLevelMapper)(nil), actualA)
	assert.ToBeOfType(t, (*defaultLevelMapper)(nil), actualB)
}

func TestDefaultLevelMapper_FromSdk(t *testing.T) {
	instance := &defaultLevelMapper{}

	cases := []struct {
		input       sdk.Level
		expected    level.Level
		expectedErr string
	}{
		{LevelTrace - 1, level.Trace, ""},
		{LevelTrace, level.Trace, ""},
		{LevelDebug - 1, level.Trace, ""},
		{LevelDebug, level.Debug, ""},
		{LevelInfo - 1, level.Debug, ""},
		{LevelInfo, level.Info, ""},
		{LevelWarn - 1, level.Info, ""},
		{LevelWarn, level.Warn, ""},
		{LevelError - 1, level.Warn, ""},
		{LevelError, level.Error, ""},
		{LevelFatal - 1, level.Error, ""},
		{LevelFatal, level.Fatal, ""},
		{LevelFatal + 1, level.Fatal, ""},
	}

	for _, c := range cases {
		t.Run(c.input.String(), func(t *testing.T) {
			actual, actualErr := instance.FromSdk(c.input)
			if c.expectedErr == "" {
				assert.ToBeNoError(t, actualErr)
				assert.ToBeEqual(t, c.expected, actual)
			} else {
				assert.ToBeMatching(t, c.expectedErr, actualErr)
				assert.ToBeEqual(t, level.Level(0), actual)
			}
		})
	}
}

func TestDefaultLevelMapper_ToSdk(t *testing.T) {
	instance := &defaultLevelMapper{}

	cases := []struct {
		input       level.Level
		expected    sdk.Level
		expectedErr string
	}{
		{level.Trace - 1, LevelTrace, ""},
		{level.Trace, LevelTrace, ""},
		{level.Debug - 1, LevelTrace, ""},
		{level.Debug, LevelDebug, ""},
		{level.Info - 1, LevelDebug, ""},
		{level.Info, LevelInfo, ""},
		{level.Warn - 1, LevelInfo, ""},
		{level.Warn, LevelWarn, ""},
		{level.Error - 1, LevelWarn, ""},
		{level.Error, LevelError, ""},
		{level.Fatal - 1, LevelError, ""},
		{level.Fatal, LevelFatal, ""},
		{level.Fatal + 1, LevelFatal, ""},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("l%d", c.input), func(t *testing.T) {
			actual, actualErr := instance.ToSdk(c.input)
			if c.expectedErr == "" {
				assert.ToBeNoError(t, actualErr)
				assert.ToBeEqual(t, c.expected, actual)
			} else {
				assert.ToBeMatching(t, c.expectedErr, actualErr)
				assert.ToBeEqual(t, sdk.Level(0), actual)
			}
		})
	}
}

func TestLevelMapperFacade_FromSdk(t *testing.T) {
	instance := NewLevelMapperFacade(NewLevelMapper)

	actualA, actualErrA := instance.FromSdk(sdk.LevelDebug)
	assert.ToBeNoError(t, actualErrA)
	assert.ToBeEqual(t, level.Debug, actualA)

	actualB, actualErrB := instance.FromSdk(LevelWarn + 1)
	assert.ToBeNoError(t, actualErrB)
	assert.ToBeEqual(t, level.Warn, actualB)
}

func TestLevelMapperFacade_ToSdk(t *testing.T) {
	instance := NewLevelMapperFacade(NewLevelMapper)

	actualA, actualErrA := instance.ToSdk(level.Debug)
	assert.ToBeNoError(t, actualErrA)
	assert.ToBeEqual(t, sdk.LevelDebug, actualA)

	actualB, actualErrB := instance.ToSdk(level.Warn + 1)
	assert.ToBeNoError(t, actualErrB)
	assert.ToBeEqual(t, LevelWarn, actualB)
}

type testingLevelMapper struct {
	fromSdk func(sdk.Level) (level.Level, error)
	toSdk   func(level.Level) (sdk.Level, error)
}

func (instance *testingLevelMapper) FromSdk(s sdk.Level) (level.Level, error) {
	return instance.fromSdk(s)
}

func (instance *testingLevelMapper) ToSdk(l level.Level) (sdk.Level, error) {
	return instance.toSdk(l)
}
