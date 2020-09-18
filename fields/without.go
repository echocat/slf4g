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
	return &lineage{
		fields: With(key, value),
		parent: instance,
	}
}

func (instance *without) Withf(key string, format string, args ...interface{}) Fields {
	return &lineage{
		fields: Withf(key, format, args...),
		parent: instance,
	}
}

func (instance *without) WithFields(fields Fields) Fields {
	return &lineage{
		fields: fields,
		parent: instance,
	}
}

func (instance *without) Without(keys ...string) Fields {
	return Without(instance, keys...)
}
