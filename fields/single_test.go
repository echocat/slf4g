package fields

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_With(t *testing.T) {
	actual := With("foo", 2)

	assert.ToBeEqual(t, mapped{"foo": 2}, mustAsMap(actual))
}

func Test_Withf(t *testing.T) {
	actual := Withf("foo", "hello %d", 2)

	assert.ToBeEqual(t, mapped{"foo": LazyFormat("hello %d", 2)}, mustAsMap(actual))
}

func Test_single_ForEach(t *testing.T) {
	instance := &single{"foo", 1}

	actualConsumed := map[string]interface{}{}
	actualErr := instance.ForEach(func(k string, v interface{}) error {
		actualConsumed[k] = v
		return nil
	})

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, map[string]interface{}{
		"foo": 1,
	}, actualConsumed)
}

func Test_single_ForEach_isForwardingErrors(t *testing.T) {
	expectedErr := errors.New("foo")
	instance := &single{"foo", 1}

	actualErr := instance.ForEach(func(string, interface{}) error {
		return expectedErr
	})

	assert.ToBeEqual(t, expectedErr, actualErr)
}

//goland:noinspection GoNilness
func Test_single_ForEach_withNilInstance(t *testing.T) {
	var instance *single

	actualErr := instance.ForEach(func(k string, v interface{}) error {
		assert.Fail(t, "should never be called")
		return nil
	})

	assert.ToBeNoError(t, actualErr)
}

func Test_single_ForEach_withNilConsumer(t *testing.T) {
	instance := &single{"foo", 1}

	actualErr := instance.ForEach(nil)

	assert.ToBeNoError(t, actualErr)
}

func Test_single_Get(t *testing.T) {
	instance := &single{"foo", 1}

	actual1, actual1Exists := instance.Get("foo")
	assert.ToBeEqual(t, 1, actual1)
	assert.ToBeEqual(t, true, actual1Exists)

	actual2, actual2Exists := instance.Get("bar")
	assert.ToBeEqual(t, nil, actual2)
	assert.ToBeEqual(t, false, actual2Exists)

	actual3, actual3Exists := instance.Get("xyz")
	assert.ToBeEqual(t, nil, actual3)
	assert.ToBeEqual(t, false, actual3Exists)
}

//goland:noinspection GoNilness
func Test_single_Get_withNilInstance(t *testing.T) {
	var instance *single

	actual1, actual1Exists := instance.Get("foo")
	assert.ToBeEqual(t, nil, actual1)
	assert.ToBeEqual(t, false, actual1Exists)

	actual2, actual2Exists := instance.Get("bar")
	assert.ToBeEqual(t, nil, actual2)
	assert.ToBeEqual(t, false, actual2Exists)
}

func Test_single_With(t *testing.T) {
	instance := &single{"foo", 1}

	actual := instance.With("bar", 2)
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": 2}, mustAsMap(actual))
}

func Test_single_With_overwrites(t *testing.T) {
	instance := &single{"foo", 1}

	actual := instance.With("foo", 2)
	assert.ToBeEqual(t, mapped{"foo": 2}, mustAsMap(actual))
}

//goland:noinspection GoNilness
func Test_single_With_withNilInstance(t *testing.T) {
	var instance *single

	actual := instance.With("bar", 2)
	assert.ToBeEqual(t, mapped{"bar": 2}, mustAsMap(actual))
}

func Test_single_Withf(t *testing.T) {
	instance := &single{"foo", 1}

	actual := instance.Withf("bar", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": LazyFormat("hello %d", 2)}, mustAsMap(actual))
}

func Test_single_Withf_overwrites(t *testing.T) {
	instance := &single{"foo", 1}

	actual := instance.Withf("foo", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"foo": LazyFormat("hello %d", 2)}, mustAsMap(actual))
}

//goland:noinspection GoNilness
func Test_single_Withf_withNilInstance(t *testing.T) {
	var instance *single

	actual := instance.Withf("bar", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"bar": LazyFormat("hello %d", 2)}, mustAsMap(actual))
}

func Test_single_WithAll(t *testing.T) {
	instance := &single{"foo", 1}

	actual := instance.WithAll(map[string]interface{}{"bar": 2, "xyz": 3})
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": 2, "xyz": 3}, mustAsMap(actual))
}

func Test_single_WithAll_overwrites(t *testing.T) {
	instance := &single{"foo", 1}

	actual := instance.WithAll(map[string]interface{}{"foo": 66, "xyz": 3})
	assert.ToBeEqual(t, mapped{"foo": 66, "xyz": 3}, mustAsMap(actual))
}

//goland:noinspection GoNilness
func Test_single_WithAll_withNilInstance(t *testing.T) {
	var instance *single

	actual := instance.WithAll(map[string]interface{}{"bar": 66, "xyz": 3})
	assert.ToBeEqual(t, mapped{"bar": 66, "xyz": 3}, mustAsMap(actual))
}

func Test_single_Without(t *testing.T) {
	instance := &single{"foo", 1}

	actual := instance.Without("foo", "notExisting")
	assert.ToBeEqual(t, mapped{}, mustAsMap(actual))
}

//goland:noinspection GoNilness
func Test_single_Without_withNilInstance(t *testing.T) {
	var instance *single

	actual := instance.Without("bar", "foo")
	assert.ToBeEqual(t, mapped{}, mustAsMap(actual))
}

func Test_single_Len(t *testing.T) {
	instance := &single{"foo", 1}

	actual := instance.Len()

	assert.ToBeEqual(t, 1, actual)
}

func Test_single_Len_withNilValue(t *testing.T) {
	instance := &single{"foo", nil}

	actual := instance.Len()

	assert.ToBeEqual(t, 1, actual)
}

//goland:noinspection GoNilness
func Test_single_Len_withNilInstance(t *testing.T) {
	var instance *single

	actual := instance.Len()

	assert.ToBeEqual(t, 0, actual)
}
