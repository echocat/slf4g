package fields

type lineage struct {
	fields Fields
	parent Fields
}

func newLineage(fields Fields, parent Fields) Fields {
	if parent == nil {
		return fields
	}
	if _, ok := parent.(*empty); ok {
		return fields
	}
	return &lineage{fields, parent}
}

func (instance *lineage) ForEach(consumer Consumer) error {
	if instance == nil || consumer == nil {
		return nil
	}

	handledKeys := map[string]bool{}
	duplicatePreventingConsumer := func(key string, value interface{}) error {
		if handledKeys[key] {
			return nil
		} else {
			handledKeys[key] = true
			return consumer(key, value)
		}
	}

	if f := instance.fields; f != nil {
		if err := f.ForEach(duplicatePreventingConsumer); err != nil {
			return err
		}
	}
	if f := instance.parent; f != nil {
		if err := f.ForEach(duplicatePreventingConsumer); err != nil {
			return err
		}
	}
	return nil
}

func (instance *lineage) Get(key string) interface{} {
	if instance == nil {
		return nil
	}
	if f := instance.fields; f != nil {
		if v := f.Get(key); v != nil {
			return v
		}
	}
	if f := instance.parent; f != nil {
		if v := f.Get(key); v != nil {
			return v
		}
	}
	return nil
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
	return newLineage(fields, instance)
}

func (instance *lineage) Without(keys ...string) Fields {
	return newWithout(instance, keys...)
}
