package fields

import "github.com/echocat/slf4g/value"

func With(key string, value interface{}) Fields {
	return &single{key: key, value: value}
}

func Withf(key string, format string, args ...interface{}) Fields {
	return With(key, value.Format(format, args...))
}

type single struct {
	key   string
	value interface{}
}

func (instance *single) ForEach(consumer Consumer) error {
	if instance == nil {
		return nil
	}
	return consumer(instance.key, instance.value)
}

func (instance *single) Get(key string) interface{} {
	if instance != nil && key == instance.key {
		return instance.value
	}
	return nil
}

func (instance *single) With(key string, value interface{}) Fields {
	return &lineage{
		fields: &single{key: key, value: value},
		parent: instance,
	}
}

func (instance *single) Withf(key string, format string, args ...interface{}) Fields {
	return &lineage{
		fields: Withf(key, format, args...),
		parent: instance,
	}
}

func (instance *single) WithFields(fields Fields) Fields {
	return &lineage{
		fields: fields,
		parent: instance,
	}
}

func (instance *single) Without(keys ...string) Fields {
	return Without(instance, keys...)
}
