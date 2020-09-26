package fields

// AsMap will assumes every given Fields instance as a Map instance.
// If the given argument is already of Map this will be simply returned. In all
// other cases a new instance of Map with all values will be created.
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

// Map is a simple map which interfaces Fields with all its functions. This can
// be better in case of performance because it could safe several Fields.With()
// calls.
//
// WARNING! This type is potentially dangerous. On the one hand anybody can cast
// to this type and can directly modify the contents. This applies to the
// initial creator of this Map instance, too. Both cases might result in
// breaking the basic contracts of Fields: Be immutable. So it is only recommend
// to use this when it really makes sense out of readability or performance.
type Map map[string]interface{}

// See Fields.ForEach()
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

// See Fields.Get()
func (instance Map) Get(key string) interface{} {
	if instance == nil {
		return nil
	}
	return instance[key]
}

// See Fields.With()
func (instance Map) With(key string, value interface{}) Fields {
	return instance.asParentOf(With(key, value))
}

// See Fields.Withf()
func (instance Map) Withf(key string, format string, args ...interface{}) Fields {
	return instance.asParentOf(Withf(key, format, args...))
}

// See Fields.WithFields()
func (instance Map) WithFields(fields Fields) Fields {
	return instance.asParentOf(fields)
}

// See Fields.Without()
func (instance Map) Without(keys ...string) Fields {
	return newWithout(instance, keys...)
}

func (instance Map) asParentOf(fields Fields) Fields {
	return newLineage(fields, instance)
}
