package fields

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_newWithout(t *testing.T) {
	given := mapped{"a": 1, "b": 2, "c": 3}

	actual := NewWithout(given, "b")

	assert.ToBeOfType(t, &without{}, actual)
	assert.ToBeEqual(t, given, actual.(*without).fields)
	assert.ToBeEqual(t, keySet{"b": keyPresent}, actual.(*without).excludedKeys)
}

func Test_newWithout_returnsEmptyOnEmpty(t *testing.T) {
	given := mapped{}

	actual := NewWithout(given, "b")

	assert.ToBeEqual(t, Empty(), actual)
}

func Test_newWithout_returnsSameOnNoKeys(t *testing.T) {
	given := mapped{"a": 1, "b": 2, "c": 3}

	actual := NewWithout(given)

	assert.ToBeEqual(t, given, actual)
}

func Test_without_ForEach(t *testing.T) {
	instance := NewWithout(mapped{"a": 1, "b": 2, "c": 3}, "b")

	actualConsumed := mapped{}
	actualErr := instance.ForEach(func(k string, v interface{}) error {
		actualConsumed[k] = v
		return nil
	})

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, mapped{"a": 1, "c": 3}, actualConsumed)
}

func Test_without_ForEach_isForwardingErrors(t *testing.T) {
	expectedErr := errors.New("foo")
	instance := NewWithout(mapped{"a": 1, "b": 2, "c": 3}, "b")

	actualErr := instance.ForEach(func(string, interface{}) error {
		return expectedErr
	})

	assert.ToBeEqual(t, expectedErr, actualErr)
}

//goland:noinspection GoNilness
func Test_without_ForEach_withNilInstance(t *testing.T) {
	var instance *without

	actualErr := instance.ForEach(func(k string, v interface{}) error {
		assert.Fail(t, "should never be called")
		return nil
	})

	assert.ToBeNoError(t, actualErr)
}

func Test_without_ForEach_withNilConsumer(t *testing.T) {
	instance := &without{}

	actualErr := instance.ForEach(nil)

	assert.ToBeNoError(t, actualErr)
}
func Test_without_ForEach_withNilFields(t *testing.T) {
	instance := &without{}

	actualErr := instance.ForEach(func(k string, v interface{}) error {
		assert.Fail(t, "should never be called")
		return nil
	})

	assert.ToBeNoError(t, actualErr)
}

func Test_without_Get(t *testing.T) {
	instance := NewWithout(mapped{"foo": 1, "bar": 2, "xyz1": 3}, "xyz1", "xyz2", "xyz3")

	actual1, actual1Exists := instance.Get("foo")
	assert.ToBeEqual(t, 1, actual1)
	assert.ToBeEqual(t, true, actual1Exists)

	actual2, actual2Exists := instance.Get("bar")
	assert.ToBeEqual(t, 2, actual2)
	assert.ToBeEqual(t, true, actual2Exists)

	actual3, actual3Exists := instance.Get("xyz1")
	assert.ToBeEqual(t, nil, actual3)
	assert.ToBeEqual(t, false, actual3Exists)

	actual4, actual4Exists := instance.Get("xyz2")
	assert.ToBeEqual(t, nil, actual4)
	assert.ToBeEqual(t, false, actual4Exists)
}

//goland:noinspection GoNilness
func Test_without_Get_withNilInstance(t *testing.T) {
	var instance *without

	actual1, actual1Exists := instance.Get("foo")
	assert.ToBeEqual(t, nil, actual1)
	assert.ToBeEqual(t, false, actual1Exists)

	actual2, actual2Exists := instance.Get("bar")
	assert.ToBeEqual(t, nil, actual2)
	assert.ToBeEqual(t, false, actual2Exists)
}

//goland:noinspection GoNilness
func Test_without_Get_withNilFields(t *testing.T) {
	instance := &without{}

	actual1, actual1Exists := instance.Get("foo")
	assert.ToBeEqual(t, nil, actual1)
	assert.ToBeEqual(t, false, actual1Exists)

	actual2, actual2Exists := instance.Get("bar")
	assert.ToBeEqual(t, nil, actual2)
	assert.ToBeEqual(t, false, actual2Exists)
}

func Test_without_With(t *testing.T) {
	instance := NewWithout(mapped{"foo": 1, "xyz": 3}, "xyz")

	actual := instance.With("bar", 2)
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": 2}, mustAsMap(actual))
}

func Test_without_With_overwrites(t *testing.T) {
	instance := NewWithout(mapped{"foo": 1, "xyz": 3}, "xyz")

	actual := instance.With("foo", 2)
	assert.ToBeEqual(t, mapped{"foo": 2}, mustAsMap(actual))
}

//goland:noinspection GoNilness
func Test_without_With_withNilInstance(t *testing.T) {
	var instance mapped

	actual := instance.With("bar", 2)
	assert.ToBeEqual(t, mapped{"bar": 2}, mustAsMap(actual))
}

func Test_without_Withf(t *testing.T) {
	instance := NewWithout(mapped{"foo": 1, "xyz": 3}, "xyz")

	actual := instance.Withf("bar", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": LazyFormat("hello %d", 2)}, mustAsMap(actual))
}

func Test_without_Withf_overwrites(t *testing.T) {
	instance := NewWithout(mapped{"foo": 1, "xyz": 3}, "xyz")

	actual := instance.Withf("foo", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"foo": LazyFormat("hello %d", 2)}, mustAsMap(actual))
}

//goland:noinspection GoNilness
func Test_without_Withf_withNilInstance(t *testing.T) {
	var instance *without

	actual := instance.Withf("bar", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"bar": LazyFormat("hello %d", 2)}, mustAsMap(actual))
}

func Test_without_WithAll(t *testing.T) {
	instance := NewWithout(mapped{"foo": 1, "xyz": 3}, "xyz")

	actual := instance.WithAll(map[string]interface{}{"bar": 2, "xyz": 3})
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": 2, "xyz": 3}, mustAsMap(actual))
}

func Test_without_WithAll_overwrites(t *testing.T) {
	instance := NewWithout(mapped{"foo": 1, "xyz": 3}, "xyz")

	actual := instance.WithAll(map[string]interface{}{"foo": 66, "xyz": 3})
	assert.ToBeEqual(t, mapped{"foo": 66, "xyz": 3}, mustAsMap(actual))
}

//goland:noinspection GoNilness
func Test_without_WithAll_withNilInstance(t *testing.T) {
	var instance *without

	actual := instance.WithAll(map[string]interface{}{"bar": 66, "xyz": 3})
	assert.ToBeEqual(t, mapped{"bar": 66, "xyz": 3}, mustAsMap(actual))
}

func Test_without_Without(t *testing.T) {
	instance := NewWithout(mapped{"foo": 1, "xyz": 3}, "xyz")

	actual := instance.Without("foo", "notExisting")
	assert.ToBeEqual(t, mapped{}, mustAsMap(actual))
}

//goland:noinspection GoNilness
func Test_without_Without_withNilInstance(t *testing.T) {
	var instance *without

	actual := instance.Without("bar", "foo")
	assert.ToBeEqual(t, mapped{}, mustAsMap(actual))
}

func Test_without_Len(t *testing.T) {
	instance := NewWithout(mapped{"foo": 1, "bar": 2, "xyz1": 3}, "xyz1", "xyz2", "xyz3")

	actual := instance.Len()

	assert.ToBeEqual(t, 2, actual)
}

//goland:noinspection GoNilness
func Test_without_Len_withNilInstance(t *testing.T) {
	var instance *without

	actual := instance.Len()

	assert.ToBeEqual(t, 0, actual)
}

func Test_without_Len_withNilFields(t *testing.T) {
	instance := &without{nil, keySet{}}

	actual := instance.Len()

	assert.ToBeEqual(t, 0, actual)
}

func Test_without_Len_withNilExcludedFields(t *testing.T) {
	instance := &without{Empty(), nil}

	actual := instance.Len()

	assert.ToBeEqual(t, 0, actual)
}
