package fields

func Without(fields Fields, keys ...string) *without {
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
	f := instance.fields
	if instance == nil || f == nil {
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
	f := instance.fields
	if instance == nil || f == nil {
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

func (instance *without) WithFields(fields Fields) Fields {
	return instance.asParentOf(fields)
}

func (instance *without) asParentOf(fields Fields) Fields {
	return newLineage(fields, instance)
}

func (instance *without) Without(keys ...string) Fields {
	return Without(instance, keys...)
}
