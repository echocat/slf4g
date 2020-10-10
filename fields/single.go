package fields

// With creates an instance of Fields for the given key value pair.
func With(key string, value interface{}) Fields {
	return &single{key: key, value: value}
}

// With creates an instance of Fields for the given key and a Lazy fmt.Sprintf
// value from the given format and args.
func Withf(key string, format string, args ...interface{}) Fields {
	return With(key, LazyFormat(format, args...))
}

type single struct {
	key   string
	value interface{}
}

func (instance *single) ForEach(consumer func(key string, value interface{}) error) error {
	if instance == nil || consumer == nil {
		return nil
	}
	return consumer(instance.key, instance.value)
}

func (instance *single) Get(key string) interface{} {
	if instance != nil && key == instance.key {
		return instance.value
	}
	return nil
}

func (instance *single) With(key string, value interface{}) Fields {
	return instance.asParentOf(With(key, value))
}

func (instance *single) Withf(key string, format string, args ...interface{}) Fields {
	return instance.asParentOf(Withf(key, format, args...))
}

func (instance *single) WithAll(of map[string]interface{}) Fields {
	return instance.asParentOf(WithAll(of))
}

func (instance *single) asParentOf(fields Fields) Fields {
	return newLineage(fields, instance)
}

func (instance *single) Without(keys ...string) Fields {
	return newWithout(instance, keys...)
}
