package fields

type Map map[string]interface{}

func (instance Map) ForEach(consumer Consumer) error {
	if instance == nil {
		return nil
	}
	for k, v := range instance {
		if err := consumer(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (instance Map) Get(key string) interface{} {
	if instance == nil {
		return nil
	}
	return instance[key]
}

func (instance Map) With(key string, value interface{}) Fields {
	return &lineage{
		fields: &single{key: key, value: value},
		parent: instance,
	}
}

func (instance Map) Withf(key string, format string, args ...interface{}) Fields {
	return &lineage{
		fields: Withf(key, format, args...),
		parent: instance,
	}
}

func (instance Map) WithFields(fields Fields) Fields {
	return &lineage{
		fields: fields,
		parent: instance,
	}
}

func (instance Map) Without(keys ...string) Fields {
	return Without(instance, keys...)
}
