package log

import (
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

type eventImpl struct {
	provider  Provider
	fields    fields.Fields
	level     level.Level
	callDepth int
	context   interface{}
}

func (instance *eventImpl) ForEach(consumer func(key string, value interface{}) error) error {
	return instance.fields.ForEach(consumer)
}

func (instance *eventImpl) Get(key string) (interface{}, bool) {
	return instance.fields.Get(key)
}

func (instance *eventImpl) Len() int {
	return instance.fields.Len()
}

func (instance *eventImpl) GetLevel() level.Level {
	return instance.level
}

func (instance *eventImpl) GetCallDepth() int {
	return instance.callDepth
}

func (instance *eventImpl) GetContext() interface{} {
	return instance.context
}

func (instance *eventImpl) With(key string, value interface{}) Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.With(key, value)
	})
}

func (instance *eventImpl) Withf(key string, format string, args ...interface{}) Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.Withf(key, format, args...)
	})
}

func (instance *eventImpl) WithError(err error) Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.With(instance.provider.GetFieldKeysSpec().GetError(), err)
	})
}

func (instance *eventImpl) WithAll(of map[string]interface{}) Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.WithAll(of)
	})
}

func (instance *eventImpl) Without(keys ...string) Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.Without(keys...)
	})
}

func (instance *eventImpl) WithContext(ctx interface{}) Event {
	return &eventImpl{
		provider:  instance.provider,
		fields:    instance.fields,
		level:     instance.level,
		callDepth: instance.callDepth,
		context:   ctx,
	}
}

func (instance *eventImpl) WithCallDepth(add int) Event {
	return &eventImpl{
		provider:  instance.provider,
		fields:    instance.fields,
		level:     instance.level,
		callDepth: instance.callDepth + add,
		context:   instance.context,
	}
}

func (instance *eventImpl) with(mod func(fields.Fields) fields.Fields) Event {
	return &eventImpl{
		provider:  instance.provider,
		fields:    mod(instance.fields),
		level:     instance.level,
		callDepth: instance.callDepth,
		context:   instance.context,
	}
}
