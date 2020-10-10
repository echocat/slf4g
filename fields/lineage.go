package fields

type lineage struct {
	target Fields
	parent Fields
}

func newLineage(target Fields, parent Fields) Fields {
	if isEmpty(parent) {
		return target
	}
	if isEmpty(target) {
		return parent
	}
	return &lineage{target, parent}
}

func (instance *lineage) ForEach(consumer func(key string, value interface{}) error) error {
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

	if f := instance.target; f != nil {
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
	if f := instance.target; f != nil {
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
