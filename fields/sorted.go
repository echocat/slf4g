package fields

import "sort"

// DefaultKeySorter is the default instance of an KeySorter. By default this is
// sorting in increasing order.
var DefaultKeySorter KeySorter = func(what []string) {
	sort.Strings(what)
}

// Sort returns an instance Fields which contains the exact same key value pairs
// but on calling Fields.ForEach() the key value pairs are returned ordered by
// all fields sorted by the provided sorter.
func Sort(fields Fields, sorter KeySorter) Fields {
	if sorter == nil {
		return fields
	}
	result := sorted{
		m: asMap(fields),
	}
	result.keys = make([]string, len(result.m))
	var i int
	for k := range result.m {
		result.keys[i] = k
		i++
	}
	sorter(result.keys)
	return &result
}

// KeySorter is used to sort all keys. See Sort() for more details.
type KeySorter func(keys []string)

type sorted struct {
	m    mapped
	keys []string
}

func (instance *sorted) ForEach(consumer Consumer) error {
	if instance == nil || consumer == nil {
		return nil
	}
	for _, k := range instance.keys {
		v := instance.m[k]
		if err := consumer(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (instance *sorted) Get(key string) interface{} {
	return instance.m.Get(key)
}

func (instance *sorted) With(key string, value interface{}) Fields {
	return instance.m.With(key, value)
}

func (instance *sorted) Withf(key string, format string, args ...interface{}) Fields {
	return instance.m.Withf(key, format, args...)
}

func (instance *sorted) Without(keys ...string) Fields {
	return instance.m.Without(keys...)
}

func (instance *sorted) WithAll(of map[string]interface{}) Fields {
	return instance.m.WithAll(of)
}

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
