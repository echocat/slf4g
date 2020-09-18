package fields

import "sort"

type KeySorter func([]string)

var DefaultKeySorter KeySorter = func(what []string) {
	sort.Strings(what)
}

func Sort(f Fields, sorter KeySorter) Fields {
	if sorter == nil {
		return f
	}
	result := sorted{
		m: AsMap(f),
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

type sorted struct {
	m    Map
	keys []string
}

func (instance *sorted) ForEach(consumer Consumer) error {
	if instance == nil {
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

func (instance *sorted) WithFields(fields Fields) Fields {
	return instance.m.WithFields(fields)
}
