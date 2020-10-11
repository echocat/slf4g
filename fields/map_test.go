package fields

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_WithAll(t *testing.T) {
	actual := WithAll(map[string]interface{}{"foo": 1, "bar": 2})

	assert.ToBeEqual(t, mapped{"foo": 1, "bar": 2}, actual)
}

func Test_mapped_ForEach(t *testing.T) {
	instance := mapped{"foo": 1, "bar": 2}

	actualConsumed := map[string]interface{}{}
	actualErr := instance.ForEach(func(k string, v interface{}) error {
		actualConsumed[k] = v
		return nil
	})

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, map[string]interface{}{
		"foo": 1,
		"bar": 2,
	}, actualConsumed)
}

func Test_mapped_ForEach_isForwardingErrors(t *testing.T) {
	expectedErr := errors.New("foo")
	instance := mapped{"foo": 1, "bar": 2}

	actualErr := instance.ForEach(func(string, interface{}) error {
		return expectedErr
	})

	assert.ToBeEqual(t, expectedErr, actualErr)
}

//goland:noinspection GoNilness
func Test_mapped_ForEach_withNilInstance(t *testing.T) {
	var instance mapped

	actualErr := instance.ForEach(func(k string, v interface{}) error {
		assert.Fail(t, "should never be called")
		return nil
	})

	assert.ToBeNoError(t, actualErr)
}

func Test_mapped_ForEach_withNilConsumer(t *testing.T) {
	instance := mapped{}

	actualErr := instance.ForEach(nil)

	assert.ToBeNoError(t, actualErr)
}

func Test_mapped_Get(t *testing.T) {
	instance := mapped{"foo": 1, "bar": 2}

	actual1, actual1Exists := instance.Get("foo")
	assert.ToBeEqual(t, 1, actual1)
	assert.ToBeEqual(t, true, actual1Exists)

	actual2, actual2Exists := instance.Get("bar")
	assert.ToBeEqual(t, 2, actual2)
	assert.ToBeEqual(t, true, actual2Exists)

	actual3, actual3Exists := instance.Get("xyz")
	assert.ToBeEqual(t, nil, actual3)
	assert.ToBeEqual(t, false, actual3Exists)
}

//goland:noinspection GoNilness
func Test_mapped_Get_withNilInstance(t *testing.T) {
	var instance mapped

	actual1, actual1Exists := instance.Get("foo")
	assert.ToBeEqual(t, nil, actual1)
	assert.ToBeEqual(t, false, actual1Exists)

	actual2, actual2Exists := instance.Get("bar")
	assert.ToBeEqual(t, nil, actual2)
	assert.ToBeEqual(t, false, actual2Exists)
}

func Test_mapped_With(t *testing.T) {
	instance := mapped{"foo": 1}

	actual := instance.With("bar", 2)
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": 2}, mustAsMap(actual))
}

func Test_mapped_With_overwrites(t *testing.T) {
	instance := mapped{"foo": 1}

	actual := instance.With("foo", 2)
	assert.ToBeEqual(t, mapped{"foo": 2}, mustAsMap(actual))
}

//goland:noinspection GoNilness
func Test_mapped_With_withNilInstance(t *testing.T) {
	var instance mapped

	actual := instance.With("bar", 2)
	assert.ToBeEqual(t, mapped{"bar": 2}, mustAsMap(actual))
}

func Test_mapped_Withf(t *testing.T) {
	instance := mapped{"foo": 1}

	actual := instance.Withf("bar", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": LazyFormat("hello %d", 2)}, mustAsMap(actual))
}

func Test_mapped_Withf_overwrites(t *testing.T) {
	instance := mapped{"foo": 1}

	actual := instance.Withf("foo", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"foo": LazyFormat("hello %d", 2)}, mustAsMap(actual))
}

//goland:noinspection GoNilness
func Test_mapped_Withf_withNilInstance(t *testing.T) {
	var instance mapped

	actual := instance.Withf("bar", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"bar": LazyFormat("hello %d", 2)}, mustAsMap(actual))
}

func Test_mapped_WithAll(t *testing.T) {
	instance := mapped{"foo": 1, "bar": 2}

	actual := instance.WithAll(map[string]interface{}{"bar": 66, "xyz": 3})
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": 66, "xyz": 3}, mustAsMap(actual))
}

//goland:noinspection GoNilness
func Test_mapped_WithAll_withNilInstance(t *testing.T) {
	var instance mapped

	actual := instance.WithAll(map[string]interface{}{"bar": 66, "xyz": 3})
	assert.ToBeEqual(t, mapped{"bar": 66, "xyz": 3}, mustAsMap(actual))
}

func Test_mapped_Without(t *testing.T) {
	instance := mapped{"foo": 1, "bar": 2, "xyz": 3, "abc": 4}

	actual := instance.Without("bar", "xyz", "notExisting")
	assert.ToBeEqual(t, mapped{"foo": 1, "abc": 4}, mustAsMap(actual))
}

//goland:noinspection GoNilness
func Test_mapped_Without_withNilInstance(t *testing.T) {
	var instance mapped

	actual := instance.Without("bar", "foo")
	assert.ToBeEqual(t, mapped{}, mustAsMap(actual))
}

func Test_mapped_Len(t *testing.T) {
	instance := mapped{"foo": 1, "bar": 2, "xyz": 3, "abc": 4}

	actual := instance.Len()

	assert.ToBeEqual(t, 4, actual)
}

//goland:noinspection GoNilness
func Test_mapped_Len_withNilInstance(t *testing.T) {
	var instance mapped

	actual := instance.Len()

	assert.ToBeEqual(t, 0, actual)
}
