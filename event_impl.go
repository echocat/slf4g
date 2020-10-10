package log

import (
	"github.com/echocat/slf4g/fields"
	level2 "github.com/echocat/slf4g/level"
)

type eventImpl struct {
	provider  Provider
	fields    fields.Fields
	level     level2.Level
	callDepth int
	context   interface{}
}

func (instance *eventImpl) ForEach(consumer func(key string, value interface{}) error) error {
	if instance == nil {
		return nil
	}
	if v := instance.fields; v != nil {
		return v.ForEach(consumer)
	}
	return nil
}

func (instance *eventImpl) Get(key string) interface{} {
	if instance == nil {
		return nil
	}
	if v := instance.fields; v != nil {
		return v.Get(key)
	}
	return nil
}

func (instance *eventImpl) GetLevel() level2.Level {
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
