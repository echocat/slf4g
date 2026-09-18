//go:build go1.21

package sdk

import (
	"errors"
	"fmt"
	sdk "log/slog"
	"strconv"
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

func TestAttrs_ForEach_resolvesLogValuer(t *testing.T) {
	instance := attrs{
		sdk.Any("secret", redactingLogValuer{"not-for-the-log"}),
		sdk.Group("group", sdk.Any("nested", redactingLogValuer{"also-not-for-the-log"})),
	}

	actual, actualErr := fields.AsMap(instance)
	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, map[string]interface{}{
		"secret": "[REDACTED]",
		"group": []sdk.Attr{
			sdk.String("nested", "[REDACTED]"),
		},
	}, actual)
	assert.ToBeEqual(t, sdk.KindLogValuer, instance[1].Value.Group()[0].Value.Kind())
}

func TestResolvedValueOf_resolvesDeepLogValuerGroups(t *testing.T) {
	const depth = 20_000
	actual := resolvedValueOf(sdk.AnyValue(nestedGroupLogValuer(depth)))

	value := actual
	for i := 0; i < depth; i++ {
		group := value.Group()
		if len(group) != 1 || group[0].Key != "nested" {
			t.Fatalf("unexpected group at depth %d: %v", i, group)
		}
		value = group[0].Value
	}
	assert.ToBeEqual(t, sdk.KindString, value.Kind())
	assert.ToBeEqual(t, "resolved", value.String())
}

func TestResolvedValueOf_preservesNestedGroupOrder(t *testing.T) {
	instance := sdk.GroupValue(
		sdk.String("before", "a"),
		sdk.Group("group",
			sdk.String("nestedBefore", "b"),
			sdk.Any("secret", redactingLogValuer{"not-for-the-log"}),
			sdk.Any("empty", emptyGroupLogValuer{}),
			sdk.String("nestedAfter", "c"),
		),
		sdk.String("after", "d"),
	)

	actual := resolvedValueOf(instance).Group()

	if len(actual) != 3 {
		t.Fatalf("unexpected root group: %v", actual)
	}
	assert.ToBeEqual(t, "before", actual[0].Key)
	assert.ToBeEqual(t, "group", actual[1].Key)
	assert.ToBeEqual(t, "after", actual[2].Key)
	nested := actual[1].Value.Group()
	if len(nested) != 3 {
		t.Fatalf("unexpected nested group: %v", nested)
	}
	assert.ToBeEqual(t, "nestedBefore", nested[0].Key)
	assert.ToBeEqual(t, "secret", nested[1].Key)
	assert.ToBeEqual(t, "[REDACTED]", nested[1].Value.String())
	assert.ToBeEqual(t, "nestedAfter", nested[2].Key)
	assert.ToBeEqual(t, sdk.KindLogValuer, instance.Group()[1].Value.Group()[1].Value.Kind())
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
		sdk.Any("secret", redactingLogValuer{"not-for-the-log"}),
	}

	cases := []struct {
		key      string
		expected interface{}
	}{
		{"foo", int64(1)},
		{"bar", int64(2)},
		{"secret", "[REDACTED]"},
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

type redactingLogValuer struct {
	Secret string
}

func (instance redactingLogValuer) LogValue() sdk.Value {
	return sdk.StringValue("[REDACTED]")
}

type emptyGroupLogValuer struct{}

func (emptyGroupLogValuer) LogValue() sdk.Value {
	return sdk.GroupValue()
}

type nestedGroupLogValuer int

func (instance nestedGroupLogValuer) LogValue() sdk.Value {
	if instance == 0 {
		return sdk.StringValue("resolved")
	}
	return sdk.GroupValue(sdk.Any("nested", instance-1))
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

func TestAttrs_addManyPreservesOrderAndLastValue(t *testing.T) {
	instance := attrs{
		sdk.Int("first", 1),
		sdk.Int("prefix.replaced", 2),
		sdk.Int("prefix.replaced", 99),
	}
	values := make([]sdk.Attr, indexedAttrsThreshold)
	for i := range values {
		values[i] = sdk.Int(strconv.Itoa(i), i)
	}
	values[1] = sdk.Int("replaced", 3)
	values[len(values)-1] = sdk.Int("0", 4)

	instance.add("prefix.", values...)

	assert.ToBeEqual(t, sdk.Int("first", 1), instance[0])
	assert.ToBeEqual(t, sdk.Int("prefix.replaced", 3), instance[1])
	assert.ToBeEqual(t, sdk.Int("prefix.replaced", 99), instance[2])
	assert.ToBeEqual(t, sdk.Int("prefix.0", 4), instance[3])
	assert.ToBeEqual(t, sdk.Int("prefix.2", 2), instance[4])
	assert.ToBeEqual(t, indexedAttrsThreshold+1, len(instance))
}
