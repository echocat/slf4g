package eventlog

import (
	"testing"
	"time"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/native/consumer"
	"github.com/echocat/slf4g/native/location"
)

func Test_CoreLogger_Log_doesNotRestoreFilteredReservedFields(t *testing.T) {
	recorder := consumer.NewRecorder()
	provider := &Provider{
		Consumer:          recorder,
		LocationDiscovery: location.NoopDiscovery(),
	}
	instance := &CoreLogger{
		provider: provider,
		name:     "test",
	}
	timestamp := &countingFilteredValue{value: time.Now()}
	logger := &countingFilteredValue{value: "secret"}
	givenEvent := instance.NewEvent(level.Info, nil).
		With("timestamp", timestamp).
		With("logger", logger)

	instance.Log(givenEvent, 0)

	assert.ToBeEqual(t, 1, timestamp.calls)
	assert.ToBeEqual(t, 1, logger.calls)
	actual := recorder.Get(0)
	assert.ToBeNil(t, log.GetTimestampOf(actual, provider))
	assert.ToBeNil(t, log.GetLoggerOf(actual, provider))
}

func Test_CoreLogger_Log_appliesDefaultsForInvalidRespectedReservedFields(t *testing.T) {
	recorder := consumer.NewRecorder()
	provider := &Provider{
		Consumer:          recorder,
		LocationDiscovery: location.NoopDiscovery(),
	}
	instance := &CoreLogger{
		provider: provider,
		name:     "test",
	}
	timestamp := &countingFilteredValue{value: time.Time{}, respected: true}
	logger := &countingFilteredValue{value: "different", respected: true}
	givenEvent := instance.NewEvent(level.Info, nil).
		With("timestamp", timestamp).
		With("logger", logger)

	instance.Log(givenEvent, 0)

	assert.ToBeEqual(t, 1, timestamp.calls)
	assert.ToBeEqual(t, 1, logger.calls)
	actual := recorder.Get(0)
	assert.ToBeNotNil(t, log.GetTimestampOf(actual, provider))
	expectedLogger := "test"
	assert.ToBeEqual(t, &expectedLogger, log.GetLoggerOf(actual, provider))
}

type countingFilteredValue struct {
	value     any
	respected bool
	calls     int
}

func (instance *countingFilteredValue) Filter(fields.FilterContext) (any, bool) {
	instance.calls++
	return instance.value, instance.respected
}

func (instance *countingFilteredValue) Get() any {
	return instance.value
}
