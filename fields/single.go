package fields

func With(key string, value interface{}) Fields {
	return &single{key: key, value: value}
}

func Withf(key string, format string, args ...interface{}) Fields {
	return With(key, LazyFormat(format, args...))
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
	return instance.asParentOf(With(key, value))
}

func (instance *single) Withf(key string, format string, args ...interface{}) Fields {
	return instance.asParentOf(Withf(key, format, args...))
}

func (instance *single) WithFields(fields Fields) Fields {
	return instance.asParentOf(fields)
}

func (instance *single) asParentOf(fields Fields) Fields {
	return newLineage(fields, instance)
}

func (instance *single) Without(keys ...string) Fields {
	return Without(instance, keys...)
}
