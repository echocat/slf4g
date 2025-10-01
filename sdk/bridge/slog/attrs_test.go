//go:build go1.21
// +build go1.21

package sdk

import (
	"errors"
	"fmt"
	sdk "log/slog"
	"testing"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/internal/test/assert"
)

func TestAttrs_ForEach_success(t *testing.T) {
	instance := attrs{
		{Key: "foo", Value: sdk.IntValue(1)},
		{Key: "bar", Value: sdk.IntValue(2)},
	}

	actual, actualErr := fields.AsMap(instance)
	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, map[string]interface{}{
		"foo": int64(1),
		"bar": int64(2),
	}, actual)
}

func TestAttrs_ForEach_empty(t *testing.T) {
	instance := attrs{}

	actualErr := instance.ForEach(func(string, interface{}) error {
		return fmt.Errorf("should never be called")
	})
	assert.ToBeNoError(t, actualErr)
}

func TestAttrs_ForEach_nilConsumer(t *testing.T) {
	instance := attrs{}

	actualErr := instance.ForEach(nil)
	assert.ToBeNoError(t, actualErr)
}

func TestAttrs_ForEach_error(t *testing.T) {
	instance := attrs{
		{Key: "foo", Value: sdk.IntValue(1)},
	}
	anError := errors.New("just an error")

	actualErr := instance.ForEach(func(string, interface{}) error {
		return anError
	})
	assert.ToBeSame(t, anError, actualErr)
}

func TestAttrs_Get(t *testing.T) {
	instance := attrs{
		{Key: "foo", Value: sdk.IntValue(1)},
		{Key: "bar", Value: sdk.IntValue(2)},
	}

	cases := []struct {
		key      string
		expected interface{}
	}{
		{"foo", int64(1)},
		{"bar", int64(2)},
		{"xyz", nil},
	}
	for _, c := range cases {
		t.Run(c.key, func(t *testing.T) {
			actual, actualOk := instance.Get(c.key)
			if c.expected != nil {
				assert.ToBeEqual(t, true, actualOk)
				assert.ToBeEqual(t, c.expected, actual)
			} else {
				assert.ToBeEqual(t, false, actualOk)
				assert.ToBeNil(t, actual)
			}
		})
	}
}

func TestAttrs_With(t *testing.T) {
	instance := attrs{
		{Key: "foo", Value: sdk.IntValue(1)},
		{Key: "bar", Value: sdk.IntValue(2)},
	}

	actual := instance.With("xyz", int64(3))

	actualAsMap, actualErr := fields.AsMap(actual)
	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, map[string]interface{}{
		"foo": int64(1),
		"bar": int64(2),
		"xyz": int64(3),
	}, actualAsMap)
}

func TestAttrs_Withf(t *testing.T) {
	instance := attrs{
		{Key: "foo", Value: sdk.IntValue(1)},
		{Key: "bar", Value: sdk.IntValue(2)},
	}

	actual := instance.Withf("xyz", "[%d]", 3)

	actualAsMap, actualErr := fields.AsMap(actual)
	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, map[string]interface{}{
		"foo": int64(1),
		"bar": int64(2),
		"xyz": fields.LazyFormat("[%d]", 3),
	}, actualAsMap)
}

func TestAttrs_WithAll(t *testing.T) {
	instance := attrs{
		{Key: "foo", Value: sdk.IntValue(1)},
		{Key: "bar", Value: sdk.IntValue(2)},
	}

	actual := instance.WithAll(map[string]interface{}{
		"xyz": int64(3),
		"abc": int64(4),
	})

	actualAsMap, actualErr := fields.AsMap(actual)
	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, map[string]interface{}{
		"foo": int64(1),
		"bar": int64(2),
		"xyz": int64(3),
		"abc": int64(4),
	}, actualAsMap)
}

func TestAttrs_Without(t *testing.T) {
	instance := attrs{
		{Key: "foo", Value: sdk.IntValue(1)},
		{Key: "bar", Value: sdk.IntValue(2)},
	}

	actual := instance.Without("foo")

	actualAsMap, actualErr := fields.AsMap(actual)
	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, map[string]interface{}{
		"bar": int64(2),
	}, actualAsMap)
}

func TestAttrs_Len(t *testing.T) {
	instance := attrs{
		{Key: "foo", Value: sdk.IntValue(1)},
		{Key: "bar", Value: sdk.IntValue(2)},
	}

	actual := instance.Len()
	assert.ToBeEqual(t, 2, actual)
}

func TestAttrs_clone(t *testing.T) {
	instance := attrs{
		{Key: "foo", Value: sdk.IntValue(1)},
		{Key: "bar", Value: sdk.IntValue(2)},
	}

	actual := instance.clone()

	assert.ToBeEqual(t, instance, actual)
	assert.ToBeNotSame(t, instance, actual)
}

func TestAttrs_add(t *testing.T) {
	instance := attrs{
		{Key: "foo", Value: sdk.IntValue(1)},
		{Key: "bar", Value: sdk.IntValue(2)},
	}

	instance.add("aPrefix.",
		sdk.Attr{Key: "xyz", Value: sdk.IntValue(3)},
		sdk.Attr{Key: "abc", Value: sdk.IntValue(4)},
	)
	instance.add("",
		sdk.Attr{Key: "xyz", Value: sdk.IntValue(5)},
		sdk.Attr{Key: "bar", Value: sdk.IntValue(123)},
		sdk.Attr{Key: "abc", Value: sdk.IntValue(6)},
	)

	actual, actualErr := fields.AsMap(instance)
	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, map[string]interface{}{
		"foo":         int64(1),
		"bar":         int64(123),
		"aPrefix.xyz": int64(3),
		"aPrefix.abc": int64(4),
		"xyz":         int64(5),
		"abc":         int64(6),
	}, actual)
}
