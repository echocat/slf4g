//go:build go1.21

package sdk

import log "github.com/echocat/slf4g"

type eventWithProgramCounter struct {
	log.Event
	programCounter uintptr
}

func (instance eventWithProgramCounter) GetProgramCounter() uintptr {
	return instance.programCounter
}

func (instance eventWithProgramCounter) With(key string, value interface{}) log.Event {
	instance.Event = instance.Event.With(key, value)
	return instance
}

func (instance eventWithProgramCounter) Withf(key string, format string, args ...interface{}) log.Event {
	instance.Event = instance.Event.Withf(key, format, args...)
	return instance
}

func (instance eventWithProgramCounter) WithError(err error) log.Event {
	instance.Event = instance.Event.WithError(err)
	return instance
}

func (instance eventWithProgramCounter) WithAll(values map[string]interface{}) log.Event {
	instance.Event = instance.Event.WithAll(values)
	return instance
}

func (instance eventWithProgramCounter) Without(keys ...string) log.Event {
	instance.Event = instance.Event.Without(keys...)
	return instance
}
