package log

import (
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

type fallbackEvent struct {
	provider Provider
	fields   fields.Fields
	level    level.Level
}

func (instance *fallbackEvent) ForEach(consumer func(key string, value interface{}) error) error {
	return instance.fields.ForEach(consumer)
}

func (instance *fallbackEvent) Get(key string) (interface{}, bool) {
	return instance.fields.Get(key)
}

func (instance *fallbackEvent) Len() int {
	return instance.fields.Len()
}

func (instance *fallbackEvent) GetLevel() level.Level {
	return instance.level
}

func (instance *fallbackEvent) With(key string, value interface{}) Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.With(key, value)
	})
}

func (instance *fallbackEvent) Withf(key string, format string, args ...interface{}) Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.Withf(key, format, args...)
	})
}

func (instance *fallbackEvent) WithError(err error) Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.With(instance.provider.GetFieldKeysSpec().GetError(), err)
	})
}

func (instance *fallbackEvent) WithAll(of map[string]interface{}) Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.WithAll(of)
	})
}

func (instance *fallbackEvent) Without(keys ...string) Event {
	return instance.with(func(s fields.Fields) fields.Fields {
		return s.Without(keys...)
	})
}

func (instance *fallbackEvent) with(mod func(fields.Fields) fields.Fields) Event {
	return &fallbackEvent{
		provider: instance.provider,
		fields:   mod(instance.fields),
		level:    instance.level,
	}
}
