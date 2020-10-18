package fields

// WithAll wraps a map into Fields with all its functions. This can be better in
// case of performance because it could safe several Fields.With() calls.
//
// WARNING! This type is potentially dangerous. On the one hand anybody can cast
// to this type and can directly modify the contents. This applies to the
// initial creator of this mapped instance, too. Both cases might result in
// breaking the basic contracts of Fields: Be immutable. So it is only recommend
// to use this when it really makes sense out of readability or performance.
func WithAll(of map[string]interface{}) Fields {
	if of == nil || len(of) == 0 {
		return Empty()
	}
	return mapped(of)
}

type mapped map[string]interface{}

func (instance mapped) ForEach(consumer func(key string, value interface{}) error) error {
	if instance == nil || consumer == nil {
		return nil
	}
	for k, v := range instance {
		if err := consumer(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (instance mapped) Get(key string) (interface{}, bool) {
	if instance == nil {
		return nil, false
	}
	v, exists := instance[key]
	return v, exists
}

func (instance mapped) With(key string, value interface{}) Fields {
	return instance.asParentOf(With(key, value))
}

func (instance mapped) Withf(key string, format string, args ...interface{}) Fields {
	return instance.asParentOf(Withf(key, format, args...))
}

func (instance mapped) WithAll(of map[string]interface{}) Fields {
	return instance.asParentOf(WithAll(of))
}

func (instance mapped) Without(keys ...string) Fields {
	return newWithout(instance, keys...)
}

func (instance mapped) asParentOf(fields Fields) Fields {
	return newLineage(fields, instance)
}

func (instance mapped) Len() (result int) {
	if instance == nil {
		return
	}
	return len(instance)
}
