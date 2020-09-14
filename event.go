package log

import "github.com/echocat/slf4g/fields"

type Event interface {
	fields.Fields

	GetLevel() Level
	GetCallDepth() int
	GetContext() interface{}

	WithField(key string, value interface{}) Event
}

func NewEvent(level Level, fields fields.Fields, callDepth int) *EventImpl {
	return &EventImpl{
		Level:     level,
		Fields:    fields,
		CallDepth: callDepth,
	}
}

type EventImpl struct {
	fields.Fields
	Level     Level
	CallDepth int
	Context   interface{}
}

func (instance *EventImpl) GetLevel() Level {
	return instance.Level
}

func (instance *EventImpl) GetCallDepth() int {
	return instance.CallDepth
}

func (instance *EventImpl) GetContext() interface{} {
	return instance.Context
}

func (instance *EventImpl) WithField(key string, value interface{}) Event {
	return &EventImpl{
		Fields:    instance.Fields.With(key, value),
		Level:     instance.Level,
		CallDepth: instance.CallDepth,
		Context:   instance.Context,
	}
}
