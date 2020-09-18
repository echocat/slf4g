package fields

type lineage struct {
	fields Fields
	parent Fields
}

func (instance *lineage) ForEach(consumer Consumer) error {
	if instance == nil {
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
	return &lineage{
		fields: With(key, value),
		parent: instance,
	}
}

func (instance *lineage) Withf(key string, format string, args ...interface{}) Fields {
	return &lineage{
		fields: Withf(key, format, args...),
		parent: instance,
	}
}

func (instance *lineage) WithFields(fields Fields) Fields {
	return &lineage{
		fields: fields,
		parent: instance,
	}
}

func (instance *lineage) Without(keys ...string) Fields {
	return Without(instance, keys...)
}
