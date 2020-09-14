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
	return &lineage{
		fields: &single{key: key, value: value},
	}
}

func (instance *empty) WithFields(fields Fields) Fields {
	return &lineage{
		fields: fields,
	}
}

func (instance *empty) Without(keys ...string) Fields {
	return Without(instance, keys...)
}
