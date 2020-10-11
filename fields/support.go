package fields

// ForEachEnabled defines a type that handles iterates over all its
// key value pairs and providing it to the the given consumer.
type ForEachEnabled interface {
	ForEach(consumer func(key string, value interface{}) error) error
}

// ForEachFunc is a utility type to wrapping simple functions into
// ForEachEnabled.
type ForEachFunc func(consumer func(key string, value interface{}) error) error

func (instance ForEachFunc) ForEach(consumer func(key string, value interface{}) error) error {
	return instance(consumer)
}

func asMap(f ForEachEnabled) (mapped, error) {
	if f == nil {
		return mapped{}, nil
	}

	switch v := f.(type) {
	case mapped:
		return v, nil
	case *mapped:
		return *v, nil
	}

	result := mapped{}
	if err := f.ForEach(func(key string, value interface{}) error {
		result[key] = value
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func mustAsMap(f ForEachEnabled) mapped {
	if m, err := asMap(f); err != nil {
		panic(err)
	} else {
		return m
	}
}

func isEmpty(given ForEachEnabled) bool {
	if given == nil {
		return true
	}
	if _, ok := given.(*empty); ok {
		return true
	}
	if v, ok := given.(mapped); ok && len(v) == 0 {
		return true
	}
	return false
}

type keySet map[string]struct{}
