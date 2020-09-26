package fields

func Empty() Fields {
	return emptyV
}

type empty struct{}

var emptyV = &empty{}

func (instance *empty) ForEach(Consumer) error {
	return nil
}

func (instance *empty) Get(string) interface{} {
	return nil
}

func (instance *empty) With(key string, value interface{}) Fields {
	return With(key, value)
}

func (instance *empty) Withf(key string, format string, args ...interface{}) Fields {
	return Withf(key, format, args...)
}

func (instance *empty) WithFields(fields Fields) Fields {
	return fields
}

func (instance *empty) Without(...string) Fields {
	return instance
}
