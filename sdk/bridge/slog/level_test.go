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
		{LevelTrace, level.Trace, ""},
		{LevelDebug, level.Debug, ""},
		{LevelInfo, level.Info, ""},
		{LevelWarn, level.Warn, ""},
		{LevelError, level.Error, ""},
		{LevelFatal, level.Fatal, ""},
		{666, 0, "unknown slog level: 666"},
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
		{level.Trace, LevelTrace, ""},
		{level.Debug, LevelDebug, ""},
		{level.Info, LevelInfo, ""},
		{level.Warn, LevelWarn, ""},
		{level.Error, LevelError, ""},
		{level.Fatal, LevelFatal, ""},
		{666, 0, "unknown log level: 666"},
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

	actualB, actualErrB := instance.FromSdk(sdk.Level(666))
	assert.ToBeMatching(t, "unknown slog level: 666", actualErrB)
	assert.ToBeEqual(t, level.Level(0), actualB)
}

func TestLevelMapperFacade_ToSdk(t *testing.T) {
	instance := NewLevelMapperFacade(NewLevelMapper)

	actualA, actualErrA := instance.ToSdk(level.Debug)
	assert.ToBeNoError(t, actualErrA)
	assert.ToBeEqual(t, sdk.LevelDebug, actualA)

	actualB, actualErrB := instance.ToSdk(level.Level(666))
	assert.ToBeMatching(t, "unknown log level: 666", actualErrB)
	assert.ToBeEqual(t, sdk.Level(0), actualB)
}
