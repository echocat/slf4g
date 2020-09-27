package fields

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_newWithout(t *testing.T) {
	given := mapped{"a": 1, "b": 2, "c": 3}

	actual := newWithout(given, "b")

	assert.ToBeOfType(t, &without{}, actual)
	assert.ToBeEqual(t, given, actual.(*without).fields)
	assert.ToBeEqual(t, withoutKeys{"b": withoutPresent}, actual.(*without).excludedKeys)
}

func Test_newWithout_returnsEmptyOnEmpty(t *testing.T) {
	given := mapped{}

	actual := newWithout(given, "b")

	assert.ToBeEqual(t, Empty(), actual)
}

func Test_newWithout_returnsSameOnNoKeys(t *testing.T) {
	given := mapped{"a": 1, "b": 2, "c": 3}

	actual := newWithout(given)

	assert.ToBeEqual(t, given, actual)
}

func Test_without_ForEach(t *testing.T) {
	instance := newWithout(mapped{"a": 1, "b": 2, "c": 3}, "b")

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
	instance := newWithout(mapped{"a": 1, "b": 2, "c": 3}, "b")

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
	instance := newWithout(mapped{"foo": 1, "bar": 2, "xyz": 3}, "xyz")

	assert.ToBeEqual(t, 1, instance.Get("foo"))
	assert.ToBeEqual(t, 2, instance.Get("bar"))
	assert.ToBeEqual(t, nil, instance.Get("xyz"))
}

//goland:noinspection GoNilness
func Test_without_Get_withNilInstance(t *testing.T) {
	var instance *without

	assert.ToBeEqual(t, nil, instance.Get("foo"))
	assert.ToBeEqual(t, nil, instance.Get("bar"))
}

//goland:noinspection GoNilness
func Test_without_Get_withNilFields(t *testing.T) {
	instance := &without{}

	assert.ToBeEqual(t, nil, instance.Get("foo"))
	assert.ToBeEqual(t, nil, instance.Get("bar"))
}

func Test_without_With(t *testing.T) {
	instance := newWithout(mapped{"foo": 1, "xyz": 3}, "xyz")

	actual := instance.With("bar", 2)
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": 2}, asMap(actual))
}

func Test_without_With_overwrites(t *testing.T) {
	instance := newWithout(mapped{"foo": 1, "xyz": 3}, "xyz")

	actual := instance.With("foo", 2)
	assert.ToBeEqual(t, mapped{"foo": 2}, asMap(actual))
}

//goland:noinspection GoNilness
func Test_without_With_withNilInstance(t *testing.T) {
	var instance *sorted

	actual := instance.With("bar", 2)
	assert.ToBeEqual(t, mapped{"bar": 2}, asMap(actual))
}

func Test_without_Withf(t *testing.T) {
	instance := newWithout(mapped{"foo": 1, "xyz": 3}, "xyz")

	actual := instance.Withf("bar", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": LazyFormat("hello %d", 2)}, asMap(actual))
}

func Test_without_Withf_overwrites(t *testing.T) {
	instance := newWithout(mapped{"foo": 1, "xyz": 3}, "xyz")

	actual := instance.Withf("foo", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"foo": LazyFormat("hello %d", 2)}, asMap(actual))
}

//goland:noinspection GoNilness
func Test_without_Withf_withNilInstance(t *testing.T) {
	var instance *without

	actual := instance.Withf("bar", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"bar": LazyFormat("hello %d", 2)}, asMap(actual))
}

func Test_without_WithAll(t *testing.T) {
	instance := newWithout(mapped{"foo": 1, "xyz": 3}, "xyz")

	actual := instance.WithAll(map[string]interface{}{"bar": 2, "xyz": 3})
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": 2, "xyz": 3}, asMap(actual))
}

func Test_without_WithAll_overwrites(t *testing.T) {
	instance := newWithout(mapped{"foo": 1, "xyz": 3}, "xyz")

	actual := instance.WithAll(map[string]interface{}{"foo": 66, "xyz": 3})
	assert.ToBeEqual(t, mapped{"foo": 66, "xyz": 3}, asMap(actual))
}

//goland:noinspection GoNilness
func Test_without_WithAll_withNilInstance(t *testing.T) {
	var instance *without

	actual := instance.WithAll(map[string]interface{}{"bar": 66, "xyz": 3})
	assert.ToBeEqual(t, mapped{"bar": 66, "xyz": 3}, asMap(actual))
}

func Test_without_Without(t *testing.T) {
	instance := newWithout(mapped{"foo": 1, "xyz": 3}, "xyz")

	actual := instance.Without("foo", "notExisting")
	assert.ToBeEqual(t, mapped{}, asMap(actual))
}

//goland:noinspection GoNilness
func Test_without_Without_withNilInstance(t *testing.T) {
	var instance *without

	actual := instance.Without("bar", "foo")
	assert.ToBeEqual(t, mapped{}, asMap(actual))
}
