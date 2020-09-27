package fields

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_Sort(t *testing.T) {
	given := mapped{"f": "value_f", "h": "value_h", "z": "value_z", "a": "value_a"}
	expected := consumedEntries{{"a", "value_a"}, {"f", "value_f"}, {"h", "value_h"}, {"z", "value_z"}}

	var actualEntries consumedEntries
	actualErr := Sort(given, DefaultKeySorter).ForEach(func(key string, value interface{}) error {
		actualEntries.add(key, value)
		return nil
	})

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, expected, actualEntries)
}

func Test_Sort_withNilSorter(t *testing.T) {
	given := mapped{"f": "value_f", "h": "value_h", "z": "value_z", "a": "value_a"}

	actual := Sort(given, nil)

	assert.ToBeEqual(t, given, actual)
}

func Test_Sort_withEmptyFields(t *testing.T) {
	given := mapped{}

	actual := Sort(given, DefaultKeySorter)

	assert.ToBeEqual(t, given, actual)
}

func Test_sorted_ForEach(t *testing.T) {
	instance := Sort(mapped{"a": 1, "b": 2}, DefaultKeySorter)

	var actualConsumed consumedEntries
	actualErr := instance.ForEach(func(k string, v interface{}) error {
		actualConsumed.add(k, v)
		return nil
	})

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, consumedEntries{{"a", 1}, {"b", 2}}, actualConsumed)
}

func Test_sorted_ForEach_isForwardingErrors(t *testing.T) {
	expectedErr := errors.New("foo")
	instance := Sort(mapped{"a": 1, "b": 2}, DefaultKeySorter)

	actualErr := instance.ForEach(func(string, interface{}) error {
		return expectedErr
	})

	assert.ToBeEqual(t, expectedErr, actualErr)
}

//goland:noinspection GoNilness
func Test_sorted_ForEach_withNilInstance(t *testing.T) {
	var instance *sorted

	actualErr := instance.ForEach(func(k string, v interface{}) error {
		assert.Fail(t, "should never be called")
		return nil
	})

	assert.ToBeNoError(t, actualErr)
}

func Test_sorted_ForEach_withNilConsumer(t *testing.T) {
	instance := &sorted{}

	actualErr := instance.ForEach(nil)

	assert.ToBeNoError(t, actualErr)
}

func Test_sorted_Get(t *testing.T) {
	instance := Sort(mapped{"foo": 1, "bar": 2}, DefaultKeySorter)

	assert.ToBeEqual(t, 1, instance.Get("foo"))
	assert.ToBeEqual(t, 2, instance.Get("bar"))
	assert.ToBeEqual(t, nil, instance.Get("xyz"))
}

//goland:noinspection GoNilness
func Test_sorted_Get_withNilInstance(t *testing.T) {
	var instance *sorted

	assert.ToBeEqual(t, nil, instance.Get("foo"))
	assert.ToBeEqual(t, nil, instance.Get("bar"))
}

func Test_sorted_With(t *testing.T) {
	instance := Sort(mapped{"foo": 1}, DefaultKeySorter)

	actual := instance.With("bar", 2)
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": 2}, asMap(actual))
}

func Test_sorted_With_overwrites(t *testing.T) {
	instance := Sort(mapped{"foo": 1}, DefaultKeySorter)

	actual := instance.With("foo", 2)
	assert.ToBeEqual(t, mapped{"foo": 2}, asMap(actual))
}

//goland:noinspection GoNilness
func Test_sorted_With_withNilInstance(t *testing.T) {
	var instance *sorted

	actual := instance.With("bar", 2)
	assert.ToBeEqual(t, mapped{"bar": 2}, asMap(actual))
}

func Test_sorted_Withf(t *testing.T) {
	instance := Sort(mapped{"foo": 1}, DefaultKeySorter)

	actual := instance.Withf("bar", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": LazyFormat("hello %d", 2)}, asMap(actual))
}

func Test_sorted_Withf_overwrites(t *testing.T) {
	instance := Sort(mapped{"foo": 1}, DefaultKeySorter)

	actual := instance.Withf("foo", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"foo": LazyFormat("hello %d", 2)}, asMap(actual))
}

//goland:noinspection GoNilness
func Test_sorted_Withf_withNilInstance(t *testing.T) {
	var instance *sorted

	actual := instance.Withf("bar", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"bar": LazyFormat("hello %d", 2)}, asMap(actual))
}

func Test_sorted_WithAll(t *testing.T) {
	instance := Sort(mapped{"foo": 1}, DefaultKeySorter)

	actual := instance.WithAll(map[string]interface{}{"bar": 2, "xyz": 3})
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": 2, "xyz": 3}, asMap(actual))
}

func Test_sorted_WithAll_overwrites(t *testing.T) {
	instance := Sort(mapped{"foo": 1}, DefaultKeySorter)

	actual := instance.WithAll(map[string]interface{}{"foo": 66, "xyz": 3})
	assert.ToBeEqual(t, mapped{"foo": 66, "xyz": 3}, asMap(actual))
}

//goland:noinspection GoNilness
func Test_sorted_WithAll_withNilInstance(t *testing.T) {
	var instance *sorted

	actual := instance.WithAll(map[string]interface{}{"bar": 66, "xyz": 3})
	assert.ToBeEqual(t, mapped{"bar": 66, "xyz": 3}, asMap(actual))
}

func Test_sorted_Without(t *testing.T) {
	instance := Sort(mapped{"foo": 1}, DefaultKeySorter)

	actual := instance.Without("foo", "notExisting")
	assert.ToBeEqual(t, mapped{}, asMap(actual))
}

//goland:noinspection GoNilness
func Test_sorted_Without_withNilInstance(t *testing.T) {
	var instance *sorted

	actual := instance.Without("bar", "foo")
	assert.ToBeEqual(t, mapped{}, asMap(actual))
}

type consumedEntries []consumedEntry

func (instance *consumedEntries) add(k string, v interface{}) {
	*instance = append(*instance, consumedEntry{k, v})
}

type consumedEntry struct {
	key   string
	value interface{}
}
