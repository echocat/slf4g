package fields

// Empty returns an empty instance of Fields.
func Empty() Fields {
	return emptyV
}

type empty struct{}

var emptyV = &empty{}

func (instance *empty) ForEach(func(key string, value any) error) error {
	return nil
}

func (instance *empty) Get(string) (any, bool) {
	return nil, false
}

func (instance *empty) With(key string, value any) Fields {
	return With(key, value)
}

func (instance *empty) Withf(key string, format string, args ...any) Fields {
	return Withf(key, format, args...)
}

func (instance *empty) WithAll(of map[string]any) Fields {
	return WithAll(of)
}

func (instance *empty) Without(...string) Fields {
	return instance
}

func (instance *empty) Len() int {
	return 0
}
