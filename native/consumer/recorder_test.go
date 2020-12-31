package consumer

import (
	"testing"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/testing/recording"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_NewRecorder(t *testing.T) {
	instance := NewRecorder()

	assert.ToBeEqual(t, 0, len(instance.recorded))
	assert.ToBeEqual(t, true, instance.Synchronized)
}

func Test_NewRecorder_withCustomization(t *testing.T) {
	instance := NewRecorder(func(recorder *Recorder) {
		recorder.Synchronized = false
	})

	assert.ToBeEqual(t, 0, len(instance.recorded))
	assert.ToBeEqual(t, false, instance.Synchronized)
}

func Test_NewRecorder_Consume(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent1 := givenLogger.NewEvent(level.Info, nil)
	givenEvent2 := givenLogger.NewEvent(level.Warn, nil)
	instance := NewRecorder()

	instance.Consume(givenEvent1, givenLogger)
	instance.Consume(givenEvent2, givenLogger)

	assert.ToBeEqual(t, []log.Event{givenEvent1, givenEvent2}, instance.recorded)
}

func Test_NewRecorder_Len(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent1 := givenLogger.NewEvent(level.Info, nil)
	givenEvent2 := givenLogger.NewEvent(level.Warn, nil)
	instance := NewRecorder(func(recorder *Recorder) {
		recorder.recorded = []log.Event{givenEvent1, givenEvent2}
	})

	actual := instance.Len()

	assert.ToBeEqual(t, 2, actual)
}

func Test_NewRecorder_Get(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent1 := givenLogger.NewEvent(level.Info, nil)
	givenEvent2 := givenLogger.NewEvent(level.Warn, nil)
	instance := NewRecorder(func(recorder *Recorder) {
		recorder.recorded = []log.Event{givenEvent1, givenEvent2}
	})

	actual1 := instance.Get(0)
	actual2 := instance.Get(1)

	assert.ToBeSame(t, givenEvent1, actual1)
	assert.ToBeSame(t, givenEvent2, actual2)
	assert.Execution(t, func() {
		instance.Get(3)
	}).WillPanicWith("Index 3 requested but the amount of recorded events is only 2")
}

func Test_NewRecorder_GetAll(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent1 := givenLogger.NewEvent(level.Info, nil)
	givenEvent2 := givenLogger.NewEvent(level.Warn, nil)
	instance := NewRecorder(func(recorder *Recorder) {
		recorder.recorded = []log.Event{givenEvent1, givenEvent2}
	})

	actual := instance.GetAll()

	assert.ToBeEqual(t, []log.Event{givenEvent1, givenEvent2}, actual)
}

func Test_NewRecorder_Reset(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent1 := givenLogger.NewEvent(level.Info, nil)
	givenEvent2 := givenLogger.NewEvent(level.Warn, nil)
	instance := NewRecorder(func(recorder *Recorder) {
		recorder.recorded = []log.Event{givenEvent1, givenEvent2}
	})

	instance.Reset()

	assert.ToBeEqual(t, []log.Event{}, instance.recorded)
}
