package fields

type ForEachEnabled interface {
	ForEach(func(key string, value interface{}) error) error
}

func asMap(f ForEachEnabled) mapped {
	switch v := f.(type) {
	case mapped:
		return v
	case *mapped:
		return *v
	}

	result := mapped{}
	if err := f.ForEach(func(key string, value interface{}) error {
		result[key] = value
		return nil
	}); err != nil {
		panic(err)
	}

	return result
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
