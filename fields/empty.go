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
	return instance.asParentOf(With(key, value))
}

func (instance *empty) Withf(key string, format string, args ...interface{}) Fields {
	return instance.asParentOf(Withf(key, format, args...))
}

func (instance *empty) WithFields(fields Fields) Fields {
	return instance.asParentOf(fields)
}

func (instance *empty) asParentOf(fields Fields) Fields {
	return newLineage(fields, nil)
}

func (instance *empty) Without(keys ...string) Fields {
	return Without(instance, keys...)
}
