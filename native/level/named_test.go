package level

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
)

func Test_AsNamed(t *testing.T) {
	givenLevel := level.Warn
	givenLevelP := &givenLevel
	givenNames := &mockNames{}

	actual := AsNamed(givenLevelP, givenNames)

	assert.ToBeOfType(t, &namedImpl{}, actual)
	assert.ToBeSame(t, givenLevelP, actual.(*namedImpl).level)
	assert.ToBeSame(t, givenNames, actual.(*namedImpl).names)
}

func Test_namedImpl_Get(t *testing.T) {
	givenLevel := level.Warn
	givenLevelP := &givenLevel
	instance := &namedImpl{level: givenLevelP}

	actual := instance.Get()

	assert.ToBeSame(t, givenLevelP, actual)
}

func Test_namedImpl_Unwrap(t *testing.T) {
	givenLevel := level.Warn
	givenLevelP := &givenLevel
	instance := &namedImpl{level: givenLevelP}

	actual := instance.Unwrap()

	assert.ToBeSame(t, givenLevelP, actual)
}

func Test_namedImpl_String_valid(t *testing.T) {
	givenLevel := level.Warn
	givenLevelP := &givenLevel
	instance := &namedImpl{level: givenLevelP, names: &mockNames{onToName: func(actualLevel level.Level) (string, error) {
		assert.ToBeEqual(t, givenLevel, actualLevel)
		return "foo", nil
	}}}

	actual := instance.String()

	assert.ToBeEqual(t, "foo", actual)
}

func Test_namedImpl_String_withError(t *testing.T) {
	givenLevel := level.Warn
	givenLevelP := &givenLevel
	givenError := errors.New("bar")
	instance := &namedImpl{level: givenLevelP, names: &mockNames{onToName: func(actualLevel level.Level) (string, error) {
		assert.ToBeEqual(t, givenLevel, actualLevel)
		return "", givenError
	}}}

	actual := instance.String()

	assert.ToBeEqual(t, "illegal-level-4000", actual)
}

func Test_namedImpl_Set_valid(t *testing.T) {
	givenLevel := level.Warn
	givenLevelP := &givenLevel
	instance := &namedImpl{level: givenLevelP, names: &mockNames{onToLevel: func(actualName string) (level.Level, error) {
		assert.ToBeEqual(t, "foo", actualName)
		return level.Error, nil
	}}}

	actualErr := instance.Set("foo")

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, level.Error, givenLevel)
}

func Test_namedImpl_Set_withError(t *testing.T) {
	givenLevel := level.Warn
	givenLevelP := &givenLevel
	givenError := errors.New("bar")
	instance := &namedImpl{level: givenLevelP, names: &mockNames{onToLevel: func(actualName string) (level.Level, error) {
		assert.ToBeEqual(t, "foo", actualName)
		return 0, givenError
	}}}

	actualErr := instance.Set("foo")

	assert.ToBeSame(t, givenError, actualErr)
	assert.ToBeEqual(t, level.Warn, givenLevel)
}

func Test_namedImpl_getNames_configured(t *testing.T) {
	givenNames := &mockNames{}
	instance := &namedImpl{names: givenNames}

	actual := instance.getNames()

	assert.ToBeSame(t, givenNames, actual)
}

func Test_namedImpl_getNames_default(t *testing.T) {
	before := DefaultNames
	defer func() { DefaultNames = before }()
	givenNames := &mockNames{}
	DefaultNames = givenNames

	instance := &namedImpl{}

	actual := instance.getNames()

	assert.ToBeSame(t, givenNames, actual)
}

func Test_namedImpl_getNames_panics(t *testing.T) {
	before := DefaultNames
	defer func() { DefaultNames = before }()
	DefaultNames = nil
	instance := &namedImpl{}

	assert.Execution(t, func() {
		instance.getNames()
	}).WillPanicWith("^Neither names configured nor level\\.DefaultNames set\\.$")
}
