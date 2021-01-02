package fields

import "sort"

// DefaultKeySorter is the default instance of an KeySorter. By default this is
// sorting in increasing order.
var DefaultKeySorter KeySorter = func(what []string) {
	sort.Strings(what)
}

// SortedForEach is calling the consumer for all entries of ForEachEnabled but
// in the order ensured by KeySorter. If KeySorter is nil DefaultKeySorter is
// used.
func SortedForEach(input ForEachEnabled, sorter KeySorter, consumer func(key string, value interface{}) error) error {
	if input == nil {
		return nil
	}
	if sorter == nil {
		sorter = DefaultKeySorter
	}
	if sorter == nil || isNoopKeySorter(sorter) {
		return input.ForEach(consumer)
	}

	m, err := asMap(input)
	if err != nil {
		return err
	}

	keys := make([]string, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}

	sorter(keys)

	for _, key := range keys {
		if err := consumer(key, m[key]); err != nil {
			return err
		}
	}
	return nil
}

// KeySorter is used to sort all keys. See Sort() for more details.
type KeySorter func(keys []string)

// NoopKeySorter provides a noop implementation of KeySorter.
func NoopKeySorter() KeySorter {
	return noopKeySorterV
}

var noopKeySorterV = func(what []string) {}

func isNoopKeySorter(v KeySorter) bool {
	return interface{}(noopKeySorterV) == interface{}(v)
}
