package fields

// NewWithout creates a new instance of target where the given keys are no longer included.
func NewWithout(target Fields, keys ...string) Fields {
	return newDerivedWithout(target, keys...)
}

type without struct {
	fields                Fields
	excludedKeys          keySet
	pendingCompactionCost int
	compactAt             int
}

var keyPresent = struct{}{}

func (instance *without) ForEach(consumer func(key string, value any) error) error {
	return forEachField(instance, consumer, false)
}

func (instance *without) Get(key string) (any, bool) {
	return getField(instance, key)
}

func (instance *without) With(key string, value any) Fields {
	return instance.asParentOf(With(key, value))
}

func (instance *without) Withf(key string, format string, args ...any) Fields {
	return instance.asParentOf(Withf(key, format, args...))
}

func (instance *without) WithAll(of map[string]any) Fields {
	return instance.asParentOf(WithAll(of))
}

func (instance *without) Without(keys ...string) Fields {
	return newDerivedWithout(instance, keys...)
}

func (instance *without) asParentOf(fields Fields) Fields {
	return newDerivedLineage(fields, instance)
}

func (instance *without) Len() (result int) {
	return fieldCount(instance)
}
