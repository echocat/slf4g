package fields

func newWithout(fields Fields, keys ...string) Fields {
	if fields == nil {
		return Empty()
	}
	result := &without{
		fields: fields,
	}
	result.excludedKeys = make(map[string]bool, len(keys))
	for _, key := range keys {
		result.excludedKeys[key] = true
	}
	return result
}

type without struct {
	fields       Fields
	excludedKeys map[string]bool
}

func (instance *without) ForEach(consumer Consumer) error {
	if instance == nil || consumer == nil {
		return nil
	}
	f := instance.fields
	if f == nil {
		return nil
	}

	excludedKeys := instance.excludedKeys
	if excludedKeys == nil {
		return f.ForEach(consumer)
	}

	filteringConsumer := func(key string, value interface{}) error {
		if excludedKeys[key] {
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
	if excludedKeys == nil {
		return f.Get(key)
	}

	if excludedKeys[key] {
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

func (instance *without) asParentOf(fields Fields) Fields {
	return newLineage(fields, instance)
}

func (instance *without) Without(keys ...string) Fields {
	return newWithout(instance, keys...)
}
