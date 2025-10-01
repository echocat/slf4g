//go:build go1.23
// +build go1.23

package fields

import (
	"errors"
	"iter"
)

var errStopForEachNow = errors.New("stop forEach now")

// Iter takes a Fields or ForEachEnabled and makes it iterable.
func Iter(source ForEachEnabled) iter.Seq2[Field, error] {
	return func(yield func(Field, error) bool) {
		if source == nil {
			return
		}

		err := source.ForEach(func(key string, value interface{}) error {
			if !yield(NewField(key, value), nil) {
				return errStopForEachNow
			}
			return nil
		})
		if errors.Is(err, errStopForEachNow) {
			return
		}
		if err != nil {
			yield(nil, err)
		}
	}
}

// Collect collects a sequence of [Field] into [Fields].
func Collect(i iter.Seq[Field]) Fields {
	result := mapped{}
	for f := range i {
		result[f.Key()] = f.Value()
	}
	return result
}

// CollectKeyValue collects a sequence of key and value into [Fields].
func CollectKeyValue(i iter.Seq2[string, interface{}]) Fields {
	result := mapped{}
	for k, v := range i {
		result[k] = v
	}
	return result
}

// CollectErr collects a sequence of [Field] into [Fields].
func CollectErr(i iter.Seq2[Field, error]) (Fields, error) {
	result := mapped{}
	for f, err := range i {
		if err != nil {
			return nil, err
		}
		result[f.Key()] = f.Value()
	}
	return result, nil
}
