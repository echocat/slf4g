package level

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_NewNamesFacade(t *testing.T) {
	givenNames := &mockNames{}

	actual := NewNamesFacade(func() Names {
		return givenNames
	})

	assert.ToBeSame(t, givenNames, actual.(namesFacade)())
}

func Test_namesFacade_ToName(t *testing.T) {
	givenLevel := Warn
	givenError := errors.New("foo")
	givenNames := &mockNames{onToName: func(actualLevel Level) (string, error) {
		assert.ToBeEqual(t, givenLevel, actualLevel)
		return "bar", givenError
	}}
	instance := namesFacade(func() Names { return givenNames })

	actual, actualErr := instance.ToName(givenLevel)

	assert.ToBeEqual(t, "bar", actual)
	assert.ToBeSame(t, givenError, actualErr)
}

func Test_namesFacade_ToLevel(t *testing.T) {
	givenLevel := Warn
	givenError := errors.New("foo")
	givenNames := &mockNames{onToLevel: func(actualName string) (Level, error) {
		assert.ToBeEqual(t, "bar", actualName)
		return givenLevel, givenError
	}}
	instance := namesFacade(func() Names { return givenNames })

	actual, actualErr := instance.ToLevel("bar")

	assert.ToBeEqual(t, givenLevel, actual)
	assert.ToBeSame(t, givenError, actualErr)
}

type mockNames struct {
	onToName  func(lvl Level) (string, error)
	onToLevel func(name string) (Level, error)
}

func (instance *mockNames) ToName(l Level) (string, error) {
	if v := instance.onToName; v != nil {
		return v(l)
	}
	panic("not implemented")
}

func (instance *mockNames) ToLevel(s string) (Level, error) {
	if v := instance.onToLevel; v != nil {
		return v(s)
	}
	panic("not implemented")
}
