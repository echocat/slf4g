package fields

func newWithout(fields Fields, keys ...string) Fields {
	if isEmpty(fields) {
		return Empty()
	}
	if len(keys) == 0 {
		return fields
	}
	result := &without{
		fields: fields,
	}
	result.excludedKeys = make(withoutKeys, len(keys))
	for _, key := range keys {
		result.excludedKeys[key] = withoutPresent
	}
	return result
}

type without struct {
	fields       Fields
	excludedKeys withoutKeys
}

type withoutKeys map[string]struct{}

var withoutPresent = struct{}{}

func (instance *without) ForEach(consumer func(key string, value interface{}) error) error {
	if instance == nil || consumer == nil {
		return nil
	}
	f := instance.fields
	if f == nil {
		return nil
	}

	excludedKeys := instance.excludedKeys
	filteringConsumer := func(key string, value interface{}) error {
		if _, ok := excludedKeys[key]; ok {
			return nil
		} else {
			return consumer(key, value)
		}
	}

	return f.ForEach(filteringConsumer)
}

func (instance *without) Get(key string) interface{} {
	if instance == nil {
		return nil
	}
	f := instance.fields
	if f == nil {
		return nil
	}

	excludedKeys := instance.excludedKeys
	if _, ok := excludedKeys[key]; ok {
		return nil
	}

	return f.Get(key)
}

func (instance *without) With(key string, value interface{}) Fields {
	return instance.asParentOf(With(key, value))
}

func (instance *without) Withf(key string, format string, args ...interface{}) Fields {
	return instance.asParentOf(Withf(key, format, args...))
}

func (instance *without) WithAll(of map[string]interface{}) Fields {
	return instance.asParentOf(WithAll(of))
}

func (instance *without) Without(keys ...string) Fields {
	return newWithout(instance, keys...)
}

func (instance *without) asParentOf(fields Fields) Fields {
	return newLineage(fields, instance)
}
