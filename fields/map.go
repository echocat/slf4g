package fields

func AsMap(f Fields) Map {
	switch v := f.(type) {
	case Map:
		return v
	case *Map:
		return *v
	}

	result := Map{}
	if err := f.ForEach(func(key string, value interface{}) error {
		result[key] = value
		return nil
	}); err != nil {
		panic(err)
	}

	return result
}

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
	return instance.asParentOf(With(key, value))
}

func (instance Map) Withf(key string, format string, args ...interface{}) Fields {
	return instance.asParentOf(Withf(key, format, args...))
}

func (instance Map) WithFields(fields Fields) Fields {
	return instance.asParentOf(fields)
}

func (instance Map) asParentOf(fields Fields) Fields {
	return newLineage(fields, instance)
}

func (instance Map) Without(keys ...string) Fields {
	return Without(instance, keys...)
}
