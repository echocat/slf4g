package consumer

import (
	"testing"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/testing/recording"
)

func Test_Func_Consume(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Warn, map[string]interface{}{
		"foo": "bar",
	})

	wasCalled := false
	instance := Func(func(actualEvent log.Event, actualLogger log.CoreLogger) {
		assert.ToBeSame(t, givenEvent, actualEvent)
		assert.ToBeSame(t, givenLogger, actualLogger)
		wasCalled = true
	})

	instance.Consume(givenEvent, givenLogger)

	assert.ToBeEqual(t, true, wasCalled)
}

func Test_noopV_Consume(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Warn, map[string]interface{}{
		"foo": "bar",
	})

	noopV.Consume(givenEvent, givenLogger)

	// Yeah.. it is a noop... nothing happened :-D
}

func Test_Noop(t *testing.T) {
	actual := Noop()

	assert.ToBeSame(t, noopV, actual)
}

func Test_NewFacade(t *testing.T) {
	givenDelegate := Func(func(actualEvent log.Event, actualLogger log.CoreLogger) {
		assert.Fail(t, "should never be called.")
	})

	instance := NewFacade(func() Consumer {
		return givenDelegate
	})

	assert.ToBeSame(t, givenDelegate, instance.(facade)())
}

func Test_facade_Consume(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Warn, map[string]interface{}{
		"foo": "bar",
	})
	wasCalled := false
	givenDelegate := Func(func(actualEvent log.Event, actualLogger log.CoreLogger) {
		assert.ToBeSame(t, givenEvent, actualEvent)
		assert.ToBeSame(t, givenLogger, actualLogger)
		wasCalled = true
	})

	instance := NewFacade(func() Consumer {
		return givenDelegate
	})

	instance.Consume(givenEvent, givenLogger)

	assert.ToBeEqual(t, true, wasCalled)
}
