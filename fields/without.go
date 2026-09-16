package fields

// NewWithout creates a new instance of target where the given keys are no longer included.
func NewWithout(target Fields, keys ...string) Fields {
	if isEmpty(target) {
		return Empty()
	}
	if len(keys) == 0 {
		return target
	}
	result := &without{
		fields: target,
	}
	result.excludedKeys = make(keySet, len(keys))
	for _, key := range keys {
		result.excludedKeys[key] = keyPresent
	}
	return result
}

type without struct {
	fields       Fields
	excludedKeys keySet
}

var keyPresent = struct{}{}

func (instance *without) ForEach(consumer func(key string, value interface{}) error) error {
	return forEachField(instance, consumer, false)
}

func (instance *without) Get(key string) (interface{}, bool) {
	return getField(instance, key)
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
	return NewWithout(instance, keys...)
}

func (instance *without) asParentOf(fields Fields) Fields {
	return NewLineage(fields, instance)
}

func (instance *without) Len() (result int) {
	return fieldCount(instance)
}
