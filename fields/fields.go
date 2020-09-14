package fields

type Fields interface {
	ForEach(consumer Consumer) error
	Get(key string) interface{}

	With(key string, value interface{}) Fields
	Without(keys ...string) Fields
	WithFields(Fields) Fields
}

type Consumer func(key string, value interface{}) error
