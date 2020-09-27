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
	if sorter == nil || isEmpty(fields) {
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
	if instance == nil {
		return nil
	}
	return instance.m.Get(key)
}

func (instance *sorted) With(key string, value interface{}) Fields {
	if instance == nil {
		return With(key, value)
	}
	return instance.m.With(key, value)
}

func (instance *sorted) Withf(key string, format string, args ...interface{}) Fields {
	if instance == nil {
		return Withf(key, format, args...)
	}
	return instance.m.Withf(key, format, args...)
}

func (instance *sorted) WithAll(of map[string]interface{}) Fields {
	if instance == nil {
		return WithAll(of)
	}
	return instance.m.WithAll(of)
}

func (instance *sorted) Without(keys ...string) Fields {
	if instance == nil {
		return Empty()
	}
	return instance.m.Without(keys...)
}
