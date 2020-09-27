package log

import "github.com/echocat/slf4g/fields"

type Event interface {
	GetLevel() Level
	GetCallDepth() int
	GetContext() interface{}
	GetFields() fields.Fields

	With(key string, value interface{}) Event
	Withf(key string, format string, args ...interface{}) Event
	WithAll(map[string]interface{}) Event
	Without(keys ...string) Event
	WithCallDepth(int) Event
}

func NewEvent(level Level, f fields.Fields, callDepth int) *EventImpl {
	if f == nil {
		f = fields.Empty()
	}
	return &EventImpl{
		Level:     level,
		Fields:    f,
		CallDepth: callDepth,
	}
}

type EventImpl struct {
	Fields    fields.Fields
	Level     Level
	CallDepth int
	Context   interface{}
}

func (instance *EventImpl) GetFields() fields.Fields {
	return instance.Fields
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

func (instance *EventImpl) With(key string, value interface{}) Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.With(key, value)
	})
}

func (instance *EventImpl) Withf(key string, format string, args ...interface{}) Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.Withf(key, format, args...)
	})
}

func (instance *EventImpl) WithAll(of map[string]interface{}) Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.WithAll(of)
	})
}

func (instance *EventImpl) Without(keys ...string) Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.Without(keys...)
	})
}

func (instance *EventImpl) WithCallDepth(add int) Event {
	return &EventImpl{
		Fields:    instance.Fields,
		Level:     instance.Level,
		CallDepth: instance.CallDepth + add,
		Context:   instance.Context,
	}
}

func (instance *EventImpl) with(mod func(fields.Fields) fields.Fields) Event {
	return &EventImpl{
		Fields:    mod(instance.Fields),
		Level:     instance.Level,
		CallDepth: instance.CallDepth,
		Context:   instance.Context,
	}
}
