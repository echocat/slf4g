package fields

// NewLineage creates a new version of Fields where target is a lineage of parent.
func NewLineage(target Fields, parent Fields) Fields {
	return newDerivedLineage(target, parent)
}

type lineage struct {
	target                Fields
	parent                Fields
	pendingCompactionCost int
	compactAt             int
}

func (instance *lineage) ForEach(consumer func(key string, value interface{}) error) error {
	return forEachField(instance, consumer, false)
}

func (instance *lineage) Get(key string) (interface{}, bool) {
	return getField(instance, key)
}

func (instance *lineage) With(key string, value interface{}) Fields {
	return instance.asParentOf(With(key, value))
}

func (instance *lineage) Withf(key string, format string, args ...interface{}) Fields {
	return instance.asParentOf(Withf(key, format, args...))
}

func (instance *lineage) WithAll(of map[string]interface{}) Fields {
	return instance.asParentOf(WithAll(of))
}

func (instance *lineage) asParentOf(fields Fields) Fields {
	return newDerivedLineage(fields, instance)
}

func (instance *lineage) Without(keys ...string) Fields {
	return newDerivedWithout(instance, keys...)
}

func (instance *lineage) Len() int {
	return fieldCount(instance)
}
