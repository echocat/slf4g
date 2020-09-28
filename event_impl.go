package log

import "github.com/echocat/slf4g/fields"

type eventImpl struct {
	Fields    fields.Fields
	Level     Level
	CallDepth int
	Context   interface{}
}

func (instance *eventImpl) GetFields() fields.Fields {
	return instance.Fields
}

func (instance *eventImpl) GetLevel() Level {
	return instance.Level
}

func (instance *eventImpl) GetCallDepth() int {
	return instance.CallDepth
}

func (instance *eventImpl) GetContext() interface{} {
	return instance.Context
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
		Fields:    instance.Fields,
		Level:     instance.Level,
		CallDepth: instance.CallDepth,
		Context:   ctx,
	}
}

func (instance *eventImpl) WithCallDepth(add int) Event {
	return &eventImpl{
		Fields:    instance.Fields,
		Level:     instance.Level,
		CallDepth: instance.CallDepth + add,
		Context:   instance.Context,
	}
}

func (instance *eventImpl) with(mod func(fields.Fields) fields.Fields) Event {
	return &eventImpl{
		Fields:    mod(instance.Fields),
		Level:     instance.Level,
		CallDepth: instance.CallDepth,
		Context:   instance.Context,
	}
}
