package level

import (
	"errors"
	"fmt"
	"testing"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_NewNames(t *testing.T) {
	actual := NewNames()

	assert.ToBeEqual(t, &defaultNames{}, actual)
}

func Test_NewNamesFacade(t *testing.T) {
	givenNames := &defaultNames{}

	actual := NewNamesFacade(func() Names {
		return givenNames
	})

	assert.ToBeSame(t, givenNames, actual.(namesFacade)())
}

func Test_namesFacade_ToName(t *testing.T) {
	givenLevel := level.Warn
	givenError := errors.New("foo")
	givenNames := &mockNames{onToName: func(actualLevel level.Level) (string, error) {
		assert.ToBeEqual(t, givenLevel, actualLevel)
		return "bar", givenError
	}}
	instance := namesFacade(func() Names { return givenNames })

	actual, actualErr := instance.ToName(givenLevel)

	assert.ToBeEqual(t, "bar", actual)
	assert.ToBeSame(t, givenError, actualErr)
}

func Test_namesFacade_ToLevel(t *testing.T) {
	givenLevel := level.Warn
	givenError := errors.New("foo")
	givenNames := &mockNames{onToLevel: func(actualName string) (level.Level, error) {
		assert.ToBeEqual(t, "bar", actualName)
		return givenLevel, givenError
	}}
	instance := namesFacade(func() Names { return givenNames })

	actual, actualErr := instance.ToLevel("bar")

	assert.ToBeEqual(t, givenLevel, actual)
	assert.ToBeSame(t, givenError, actualErr)
}

func Test_defaultNames_ToName(t *testing.T) {
	instance := &defaultNames{}

	cases := []struct {
		given    level.Level
		expected string
	}{
		{level.Trace, "TRACE"},
		{level.Debug, "DEBUG"},
		{level.Info, "INFO"},
		{level.Warn, "WARN"},
		{level.Error, "ERROR"},
		{level.Fatal, "FATAL"},
		{level.Level(666), "666"},
	}

	for _, c := range cases {
		t.Run(c.expected, func(t *testing.T) {
			actual, actualErr := instance.ToName(c.given)

			assert.ToBeNil(t, actualErr)
			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}

func Test_defaultNames_ToLevel(t *testing.T) {
	instance := &defaultNames{}

	cases := []struct {
		expected level.Level
		given    string
	}{
		{level.Trace, "TRACE"},
		{level.Debug, "DEBUG"},
		{level.Info, "INFO"},
		{level.Warn, "WARN"},
		{level.Error, "ERROR"},
		{level.Fatal, "FATAL"},
		{level.Level(666), "666"},
	}

	for _, c := range cases {
		t.Run(c.given, func(t *testing.T) {
			actual, actualErr := instance.ToLevel(c.given)

			assert.ToBeNil(t, actualErr)
			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}

func Test_defaultNames_ToLevel_withNonNumber(t *testing.T) {
	instance := &defaultNames{}

	actual, actualErr := instance.ToLevel("abc")

	assert.ToBeEqual(t, fmt.Errorf("%w: abc", ErrIllegalLevel), actualErr)
	assert.ToBeEqual(t, level.Level(0), actual)
}

type mockNames struct {
	onToName  func(lvl level.Level) (string, error)
	onToLevel func(name string) (level.Level, error)
}

func (instance *mockNames) ToName(l level.Level) (string, error) {
	if v := instance.onToName; v != nil {
		return v(l)
	}
	panic("not implemented")
}

func (instance *mockNames) ToLevel(s string) (level.Level, error) {
	if v := instance.onToLevel; v != nil {
		return v(s)
	}
	panic("not implemented")
}
