package fields

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_SortedForEach(t *testing.T) {
	given := mapped{"f": "value_f", "h": "value_h", "z": "value_z", "a": "value_a"}
	expected := entries{{"a", "value_a"}, {"f", "value_f"}, {"h", "value_h"}, {"z", "value_z"}}

	actualEntries := entries{}
	actualErr := SortedForEach(given, DefaultKeySorter, func(key string, value interface{}) error {
		actualEntries.add(key, value)
		return nil
	})

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, expected, actualEntries)
}

func Test_SortedForEach_withAsMapError(t *testing.T) {
	givenErr := errors.New("expected")
	given := ForEachFunc(func(func(key string, value interface{}) error) error {
		return givenErr
	})

	actualErr := SortedForEach(given, DefaultKeySorter, func(key string, value interface{}) error {
		panic("should never be called")
	})

	assert.ToBeEqual(t, givenErr, actualErr)
}

func Test_SortedForEach_withError(t *testing.T) {
	given := mapped{"f": "value_f", "h": "value_h", "z": "value_z", "a": "value_a"}
	giveError := errors.New("expected")
	expected := entries{{"a", "value_a"}}

	actualEntries := entries{}
	actualErr := SortedForEach(given, DefaultKeySorter, func(key string, value interface{}) error {
		if key == "f" {
			return giveError
		}
		actualEntries.add(key, value)
		return nil
	})

	assert.ToBeSame(t, giveError, actualErr)
	assert.ToBeEqual(t, expected, actualEntries)
}

func Test_SortedForEach_withNilSorter(t *testing.T) {
	given := mapped{"f": "value_f", "h": "value_h", "z": "value_z", "a": "value_a"}
	expected := entries{{"a", "value_a"}, {"f", "value_f"}, {"h", "value_h"}, {"z", "value_z"}}

	actualEntries := entries{}
	actualErr := SortedForEach(given, DefaultKeySorter, func(key string, value interface{}) error {
		actualEntries.add(key, value)
		return nil
	})

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, expected, actualEntries)
}

func Test_SortedForEach_withNilSorterAndDefaultIsNilToo(t *testing.T) {
	v := DefaultKeySorter
	defer func() {
		DefaultKeySorter = v
	}()
	DefaultKeySorter = nil

	given := mapped{"f": "value_f", "h": "value_h", "z": "value_z", "a": "value_a"}
	expected := mapped{"f": "value_f", "h": "value_h", "z": "value_z", "a": "value_a"}

	actualEntries := mapped{}
	actualErr := SortedForEach(given, nil, func(key string, value interface{}) error {
		actualEntries[key] = value
		return nil
	})

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, expected, actualEntries)
}

func Test_SortedForEach_withEmptyFields(t *testing.T) {
	given := mapped{}
	expected := entries{}

	actualEntries := entries{}
	actualErr := SortedForEach(given, nil, func(key string, value interface{}) error {
		actualEntries.add(key, value)
		return nil
	})

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, expected, actualEntries)
}

type entries []entry

func (instance *entries) add(k string, v interface{}) {
	*instance = append(*instance, entry{k, v})
}

type entry struct {
	key   string
	value interface{}
}
