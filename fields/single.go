package fields

// With creates an instance of Fields for the given key value pair.
func With(key string, value any) Fields {
	return &single{key: key, value: value}
}

// Withf creates an instance of Fields for the given key and a Lazy fmt.Sprintf
// value from the given format and args.
func Withf(key string, format string, args ...any) Fields {
	return With(key, LazyFormat(format, args...))
}

type single struct {
	key   string
	value any
}

func (instance *single) ForEach(consumer func(key string, value any) error) error {
	if instance == nil || consumer == nil {
		return nil
	}
	return consumer(instance.key, instance.value)
}

func (instance *single) Get(key string) (any, bool) {
	if instance != nil && key == instance.key {
		return instance.value, true
	}
	return nil, false
}

func (instance *single) With(key string, value any) Fields {
	return instance.asParentOf(With(key, value))
}

func (instance *single) Withf(key string, format string, args ...any) Fields {
	return instance.asParentOf(Withf(key, format, args...))
}

func (instance *single) WithAll(of map[string]any) Fields {
	return instance.asParentOf(WithAll(of))
}

func (instance *single) asParentOf(fields Fields) Fields {
	return newDerivedLineage(fields, instance)
}

func (instance *single) Without(keys ...string) Fields {
	return newDerivedWithout(instance, keys...)
}

func (instance *single) Len() (result int) {
	if instance == nil {
		return 0
	}
	return 1
}
