package formatter

import (
	"errors"
	"fmt"
	"testing"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/testing/recording"
)

func Test_NewNamesBasedLevel(t *testing.T) {
	actual := NewNamesBasedLevel(mockedNames{})

	actualFormatted, actualErr := actual.FormatLevel(level.Level(666), nil)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, "level-666", actualFormatted)
}

func Test_NewNamesBasedLevel_handlesErrors(t *testing.T) {
	givenError := errors.New("expected")
	actual := NewNamesBasedLevel(mockedNames{givenError})

	actualFormatted, actualErr := actual.FormatLevel(level.Level(666), nil)

	assert.ToBeSame(t, givenError, actualErr)
	assert.ToBeEqual(t, "", actualFormatted)
}
func Test_NewOrdinalBasedLevel(t *testing.T) {
	actual := NewOrdinalBasedLevel()

	actualFormatted, actualErr := actual.FormatLevel(level.Level(666), nil)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, uint16(666), actualFormatted)
}

func Test_LevelFunc_FormatLevel(t *testing.T) {
	givenProvider := recording.NewProvider()
	given := level.Level(666)

	wasCalled := false
	instance := LevelFunc(func(actual level.Level, actualProvider log.Provider) (interface{}, error) {
		assert.ToBeEqual(t, given, actual)
		assert.ToBeSame(t, givenProvider, actualProvider)
		wasCalled = true
		return "expected", nil
	})

	actual, actualErr := instance.FormatLevel(given, givenProvider)

	assert.ToBeEqual(t, true, wasCalled)
	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, "expected", actual)
}

func Test_LevelFunc_FormatLevel_errors(t *testing.T) {
	givenProvider := recording.NewProvider()
	given := level.Level(666)
	givenError := errors.New("expected")

	wasCalled := false
	instance := LevelFunc(func(actual level.Level, actualProvider log.Provider) (interface{}, error) {
		assert.ToBeEqual(t, given, actual)
		assert.ToBeSame(t, givenProvider, actualProvider)
		wasCalled = true
		return nil, givenError
	})

	actual, actualErr := instance.FormatLevel(given, givenProvider)

	assert.ToBeEqual(t, true, wasCalled)
	assert.ToBeSame(t, givenError, actualErr)
	assert.ToBeNil(t, actual)
}

func Test_NewLevelFacade(t *testing.T) {
	givenDelegate := LevelFunc(func(actual level.Level, actualProvider log.Provider) (interface{}, error) {
		assert.Fail(t, "should never be called.")
		return nil, nil
	})

	instance := NewLevelFacade(func() Level {
		return givenDelegate
	})

	assert.ToBeSame(t, givenDelegate, instance.(levelFacade)())
}

func Test_levelFacade_FormatLevel(t *testing.T) {
	givenProvider := recording.NewProvider()
	given := level.Level(666)

	wasCalled := false
	instance := NewLevelFacade(func() Level {
		return LevelFunc(func(actual level.Level, actualProvider log.Provider) (interface{}, error) {
			assert.ToBeEqual(t, given, actual)
			assert.ToBeSame(t, givenProvider, actualProvider)
			wasCalled = true
			return "expected", nil
		})
	})

	actual, actualErr := instance.FormatLevel(given, givenProvider)

	assert.ToBeEqual(t, true, wasCalled)
	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, "expected", actual)
}

func Test_textValueFacade_FormatLevel_errors(t *testing.T) {
	givenProvider := recording.NewProvider()
	given := level.Level(666)
	givenError := errors.New("expected")

	wasCalled := false
	instance := NewLevelFacade(func() Level {
		return LevelFunc(func(actual level.Level, actualProvider log.Provider) (interface{}, error) {
			assert.ToBeEqual(t, given, actual)
			assert.ToBeSame(t, givenProvider, actualProvider)
			wasCalled = true
			return nil, givenError
		})
	})

	actual, actualErr := instance.FormatLevel(given, givenProvider)

	assert.ToBeEqual(t, true, wasCalled)
	assert.ToBeSame(t, givenError, actualErr)
	assert.ToBeNil(t, actual)
}

func Test_textValueNoopV_FormatLevel(t *testing.T) {
	givenProvider := recording.NewProvider()
	given := level.Level(666)

	actual, actualErr := noopLevelV.FormatLevel(given, givenProvider)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, given, actual)
}

func Test_NoopLevel(t *testing.T) {
	actual := NoopLevel()

	assert.ToBeSame(t, noopLevelV, actual)
}

type mockedNames struct {
	err error
}

func (instance mockedNames) ToName(l level.Level) (string, error) {
	if err := instance.err; err != nil {
		return "", err
	}
	return fmt.Sprintf("level-%d", l), nil
}

func (instance mockedNames) ToLevel(string) (level.Level, error) {
	panic("should never be called")
}
