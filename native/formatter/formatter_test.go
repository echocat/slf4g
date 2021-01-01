package formatter

import (
	"errors"
	"testing"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/native/hints"
	"github.com/echocat/slf4g/testing/recording"
)

func Test_Func_Format(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Warn, map[string]interface{}{
		"foo": "bar",
	})
	givenHints := &struct{}{}

	wasCalled := false
	instance := Func(func(actualEvent log.Event, actualProvider log.Provider, actualHints hints.Hints) ([]byte, error) {
		assert.ToBeSame(t, givenEvent, actualEvent)
		assert.ToBeSame(t, givenLogger.GetProvider(), actualProvider)
		assert.ToBeSame(t, givenHints, actualHints)
		wasCalled = true
		return []byte("expected"), nil
	})

	actual, actualErr := instance.Format(givenEvent, givenLogger.GetProvider(), givenHints)

	assert.ToBeEqual(t, true, wasCalled)
	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, "expected", string(actual))
}

func Test_Func_Format_errors(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Warn, map[string]interface{}{
		"foo": "bar",
	})
	givenHints := &struct{}{}
	givenError := errors.New("expected")

	wasCalled := false
	instance := Func(func(actualEvent log.Event, actualProvider log.Provider, actualHints hints.Hints) ([]byte, error) {
		assert.ToBeSame(t, givenEvent, actualEvent)
		assert.ToBeSame(t, givenLogger.GetProvider(), actualProvider)
		assert.ToBeSame(t, givenHints, actualHints)
		wasCalled = true
		return nil, givenError
	})

	actual, actualErr := instance.Format(givenEvent, givenLogger.GetProvider(), givenHints)

	assert.ToBeEqual(t, true, wasCalled)
	assert.ToBeSame(t, givenError, actualErr)
	assert.ToBeNil(t, actual)
}

func Test_NewFacade(t *testing.T) {
	givenDelegate := Func(func(log.Event, log.Provider, hints.Hints) ([]byte, error) {
		assert.Fail(t, "should never be called.")
		return nil, nil
	})

	instance := NewFacade(func() Formatter {
		return givenDelegate
	})

	assert.ToBeSame(t, givenDelegate, instance.(facade)())
}

func Test_facade_Format(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Warn, map[string]interface{}{
		"foo": "bar",
	})
	givenHints := &struct{}{}

	wasCalled := false
	instance := NewFacade(func() Formatter {
		return Func(func(actualEvent log.Event, actualProvider log.Provider, actualHints hints.Hints) ([]byte, error) {
			assert.ToBeSame(t, givenEvent, actualEvent)
			assert.ToBeSame(t, givenLogger.GetProvider(), actualProvider)
			assert.ToBeSame(t, givenHints, actualHints)
			wasCalled = true
			return []byte("expected"), nil
		})
	})

	actual, actualErr := instance.Format(givenEvent, givenLogger.GetProvider(), givenHints)

	assert.ToBeEqual(t, true, wasCalled)
	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, "expected", string(actual))
}

func Test_facade_Format_errors(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Warn, map[string]interface{}{
		"foo": "bar",
	})
	givenHints := &struct{}{}
	givenError := errors.New("expected")

	wasCalled := false
	instance := NewFacade(func() Formatter {
		return Func(func(actualEvent log.Event, actualProvider log.Provider, actualHints hints.Hints) ([]byte, error) {
			assert.ToBeSame(t, givenEvent, actualEvent)
			assert.ToBeSame(t, givenLogger.GetProvider(), actualProvider)
			assert.ToBeSame(t, givenHints, actualHints)
			wasCalled = true
			return nil, givenError
		})
	})

	actual, actualErr := instance.Format(givenEvent, givenLogger.GetProvider(), givenHints)

	assert.ToBeEqual(t, true, wasCalled)
	assert.ToBeSame(t, givenError, actualErr)
	assert.ToBeNil(t, actual)
}

func Test_noopV_Format(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Warn, map[string]interface{}{
		"foo": "bar",
	})
	givenHints := &struct{}{}

	actual, actualErr := noopV.Format(givenEvent, givenLogger.GetProvider(), givenHints)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, []byte{}, actual)
}

func Test_Noop(t *testing.T) {
	actual := Noop()

	assert.ToBeSame(t, noopV, actual)
}
