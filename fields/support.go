package fields

func asMap(f Fields) mapped {
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

func isEmpty(given Fields) bool {
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
