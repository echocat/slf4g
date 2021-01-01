package formatter

import (
	"errors"
	"testing"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/testing/recording"
)

func Test_TextValueFunc_Format(t *testing.T) {
	givenProvider := recording.NewProvider()
	given := &struct{}{}

	wasCalled := false
	instance := TextValueFunc(func(actual interface{}, actualProvider log.Provider) ([]byte, error) {
		assert.ToBeSame(t, given, actual)
		assert.ToBeSame(t, givenProvider, actualProvider)
		wasCalled = true
		return []byte("expected"), nil
	})

	actual, actualErr := instance.FormatTextValue(given, givenProvider)

	assert.ToBeEqual(t, true, wasCalled)
	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, "expected", string(actual))
}

func Test_TextValueFunc_Format_errors(t *testing.T) {
	givenProvider := recording.NewProvider()
	given := &struct{}{}
	givenError := errors.New("expected")

	wasCalled := false
	instance := TextValueFunc(func(actual interface{}, actualProvider log.Provider) ([]byte, error) {
		assert.ToBeSame(t, given, actual)
		assert.ToBeSame(t, givenProvider, actualProvider)
		wasCalled = true
		return nil, givenError
	})

	actual, actualErr := instance.FormatTextValue(given, givenProvider)

	assert.ToBeEqual(t, true, wasCalled)
	assert.ToBeSame(t, givenError, actualErr)
	assert.ToBeNil(t, actual)
}

func Test_NewTextValueFacade(t *testing.T) {
	givenDelegate := TextValueFunc(func(actual interface{}, actualProvider log.Provider) ([]byte, error) {
		assert.Fail(t, "should never be called.")
		return nil, nil
	})

	instance := NewTextValueFacade(func() TextValue {
		return givenDelegate
	})

	assert.ToBeSame(t, givenDelegate, instance.(textValueFacade)())
}

func Test_textValueFacade_Format(t *testing.T) {
	givenProvider := recording.NewProvider()
	given := &struct{}{}

	wasCalled := false
	instance := NewTextValueFacade(func() TextValue {
		return TextValueFunc(func(actual interface{}, actualProvider log.Provider) ([]byte, error) {
			assert.ToBeSame(t, given, actual)
			assert.ToBeSame(t, givenProvider, actualProvider)
			wasCalled = true
			return []byte("expected"), nil
		})
	})

	actual, actualErr := instance.FormatTextValue(given, givenProvider)

	assert.ToBeEqual(t, true, wasCalled)
	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, "expected", string(actual))
}

func Test_textValueFacade_FormatTextValue_errors(t *testing.T) {
	givenProvider := recording.NewProvider()
	given := &struct{}{}
	givenError := errors.New("expected")

	wasCalled := false
	instance := NewTextValueFacade(func() TextValue {
		return TextValueFunc(func(actual interface{}, actualProvider log.Provider) ([]byte, error) {
			assert.ToBeSame(t, given, actual)
			assert.ToBeSame(t, givenProvider, actualProvider)
			wasCalled = true
			return nil, givenError
		})
	})

	actual, actualErr := instance.FormatTextValue(given, givenProvider)

	assert.ToBeEqual(t, true, wasCalled)
	assert.ToBeSame(t, givenError, actualErr)
	assert.ToBeNil(t, actual)
}

func Test_textValueNoopV_FormatTextValue(t *testing.T) {
	givenProvider := recording.NewProvider()
	given := &struct{}{}

	actual, actualErr := noopTextValueV.FormatTextValue(given, givenProvider)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, []byte{}, actual)
}

func Test_NoopTextValue(t *testing.T) {
	actual := NoopTextValue()

	assert.ToBeSame(t, noopTextValueV, actual)
}
