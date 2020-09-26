package fields

func newDummyFields() Fields {
	return &dummyFields{}
}

type dummyFields struct{}

func (instance *dummyFields) ForEach(Consumer) error {
	panic("should never be called")
}

func (instance *dummyFields) Get(string) interface{} {
	panic("should never be called")
}

func (instance *dummyFields) With(string, interface{}) Fields {
	panic("should never be called")
}

func (instance *dummyFields) Withf(string, string, ...interface{}) Fields {
	panic("should never be called")
}

func (instance *dummyFields) Without(...string) Fields {
	panic("should never be called")
}

func (instance *dummyFields) WithFields(Fields) Fields {
	panic("should never be called")
}
