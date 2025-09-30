package fields

// NewLineage creates a new version of Fields where target is a lineage of parent.
func NewLineage(target Fields, parent Fields) Fields {
	if isEmpty(parent) {
		return target
	}
	if isEmpty(target) {
		return parent
	}
	return &lineage{target, parent}
}

type lineage struct {
	target Fields
	parent Fields
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

func (instance *lineage) Get(key string) (interface{}, bool) {
	if instance == nil {
		return nil, false
	}
	if f := instance.target; f != nil {
		if v, exists := f.Get(key); exists {
			return v, true
		}
	}
	if f := instance.parent; f != nil {
		if v, exists := f.Get(key); exists {
			return v, true
		}
	}
	return nil, false
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
	return NewLineage(fields, instance)
}

func (instance *lineage) Without(keys ...string) Fields {
	return NewWithout(instance, keys...)
}

func (instance *lineage) Len() int {
	if instance == nil {
		return 0
	}
	consumedKeys := keySet{}
	if f := instance.target; f != nil {
		_ = f.ForEach(func(key string, value interface{}) error {
			consumedKeys[key] = keyPresent
			return nil
		})
	}
	if f := instance.parent; f != nil {
		_ = f.ForEach(func(key string, value interface{}) error {
			consumedKeys[key] = keyPresent
			return nil
		})
	}
	return len(consumedKeys)
}
