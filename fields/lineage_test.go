package fields

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_newLineage_withTargetAndParent(t *testing.T) {
	givenTarget := With("a", 1)
	givenParent := With("b", 2)

	actual := newLineage(givenTarget, givenParent)

	assert.ToBeOfType(t, &lineage{}, actual)
	assert.ToBeSame(t, givenTarget, actual.(*lineage).target)
	assert.ToBeSame(t, givenParent, actual.(*lineage).parent)
}

func Test_newLineage_withTargetAndNilParent(t *testing.T) {
	givenTarget := With("a", 1)

	actual := newLineage(givenTarget, nil)

	assert.ToBeSame(t, givenTarget, actual)
}

func Test_newLineage_withTargetAndEmptyParent(t *testing.T) {
	givenTarget := With("a", 1)
	givenParent := Empty()

	actual := newLineage(givenTarget, givenParent)

	assert.ToBeSame(t, givenTarget, actual)
}

func Test_newLineage_withTargetAndEmptyMapParent(t *testing.T) {
	givenTarget := With("a", 1)
	givenParent := mapped{}

	actual := newLineage(givenTarget, givenParent)

	assert.ToBeSame(t, givenTarget, actual)
}

func Test_newLineage_withNilTargetAndParent(t *testing.T) {
	givenParent := With("b", 2)

	actual := newLineage(nil, givenParent)

	assert.ToBeSame(t, givenParent, actual)
}

func Test_newLineage_withEmptyTargetAndParent(t *testing.T) {
	givenTarget := Empty()
	givenParent := With("b", 2)

	actual := newLineage(givenTarget, givenParent)

	assert.ToBeSame(t, givenParent, actual)
}

func Test_newLineage_withEmptyMapTargetAndParent(t *testing.T) {
	givenTarget := mapped{}
	givenParent := With("b", 2)

	actual := newLineage(givenTarget, givenParent)

	assert.ToBeSame(t, givenParent, actual)
}

func Test_lineage_ForEach(t *testing.T) {
	instance := &lineage{With("foo", 1), With("bar", 2)}

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

func Test_lineage_ForEach_isForwardingTargetErrors(t *testing.T) {
	expectedErr := errors.New("foo")
	instance := &lineage{With("foo", 1), With("bar", 2)}

	actualErr := instance.ForEach(func(string, interface{}) error {
		return expectedErr
	})

	assert.ToBeEqual(t, expectedErr, actualErr)
}

func Test_lineage_ForEach_isForwardingParentErrors(t *testing.T) {
	expectedErr := errors.New("foo")
	instance := &lineage{With("foo", 1), With("bar", 2)}

	actualErr := instance.ForEach(func(k string, _ interface{}) error {
		if k == "bar" {
			return expectedErr
		}
		return nil
	})

	assert.ToBeEqual(t, expectedErr, actualErr)
}

//goland:noinspection GoNilness
func Test_lineage_ForEach_withNilInstance(t *testing.T) {
	var instance *lineage

	actualErr := instance.ForEach(func(k string, v interface{}) error {
		assert.Fail(t, "should never be called")
		return nil
	})

	assert.ToBeNoError(t, actualErr)
}

func Test_lineage_ForEach_withNilConsumer(t *testing.T) {
	instance := &lineage{}

	actualErr := instance.ForEach(nil)

	assert.ToBeNoError(t, actualErr)
}

func Test_lineage_Get(t *testing.T) {
	instance := &lineage{With("foo", 1), With("bar", 2)}

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
func Test_lineage_Get_withNilInstance(t *testing.T) {
	var instance *lineage

	actual1, actual1Exists := instance.Get("foo")
	assert.ToBeEqual(t, nil, actual1)
	assert.ToBeEqual(t, false, actual1Exists)

	actual2, actual2Exists := instance.Get("bar")
	assert.ToBeEqual(t, nil, actual2)
	assert.ToBeEqual(t, false, actual2Exists)
}

func Test_lineage_With(t *testing.T) {
	instance := &lineage{With("foo", 1), With("bar", 2)}

	actual := instance.With("xyz", 3)
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": 2, "xyz": 3}, mustAsMap(actual))
}

func Test_lineage_With_overwrites(t *testing.T) {
	instance := &lineage{With("foo", 1), With("bar", 2)}

	actual := instance.With("foo", 2)
	assert.ToBeEqual(t, mapped{"foo": 2, "bar": 2}, mustAsMap(actual))
}

//goland:noinspection GoNilness
func Test_lineage_With_withNilInstance(t *testing.T) {
	var instance *lineage

	actual := instance.With("bar", 2)
	assert.ToBeEqual(t, mapped{"bar": 2}, mustAsMap(actual))
}

func Test_lineage_Withf(t *testing.T) {
	instance := &lineage{With("foo", 1), With("bar", 2)}

	actual := instance.Withf("xyz", "hello %d", 3)
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": 2, "xyz": LazyFormat("hello %d", 3)}, mustAsMap(actual))
}

func Test_lineage_Withf_overwrites(t *testing.T) {
	instance := &lineage{With("foo", 1), With("bar", 2)}

	actual := instance.Withf("foo", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"foo": LazyFormat("hello %d", 2), "bar": 2}, mustAsMap(actual))
}

//goland:noinspection GoNilness
func Test_lineage_Withf_withNilInstance(t *testing.T) {
	var instance *lineage

	actual := instance.Withf("bar", "hello %d", 2)
	assert.ToBeEqual(t, mapped{"bar": LazyFormat("hello %d", 2)}, mustAsMap(actual))
}

func Test_lineage_WithAll(t *testing.T) {
	instance := &lineage{With("foo", 1), With("bar", 2)}

	actual := instance.WithAll(map[string]interface{}{"bar": 66, "xyz": 3})
	assert.ToBeEqual(t, mapped{"foo": 1, "bar": 66, "xyz": 3}, mustAsMap(actual))
}

//goland:noinspection GoNilness
func Test_lineage_WithAll_withNilInstance(t *testing.T) {
	var instance *lineage

	actual := instance.WithAll(map[string]interface{}{"bar": 66, "xyz": 3})
	assert.ToBeEqual(t, mapped{"bar": 66, "xyz": 3}, mustAsMap(actual))
}

func Test_lineage_Without(t *testing.T) {
	instance := &lineage{With("foo", 1), With("bar", 2)}

	actual1 := instance.Without("bar", "notExisting")
	assert.ToBeEqual(t, mapped{"foo": 1}, mustAsMap(actual1))
	actual2 := actual1.Without("foo", "notExisting")
	assert.ToBeEqual(t, mapped{}, mustAsMap(actual2))
}

//goland:noinspection GoNilness
func Test_lineage_Without_withNilInstance(t *testing.T) {
	var instance *lineage

	actual := instance.Without("bar", "foo")
	assert.ToBeEqual(t, mapped{}, mustAsMap(actual))
}

func Test_lineage_Len(t *testing.T) {
	instance := &lineage{With("a", 1).With("b", 2), With("a", 1).With("c", 3)}

	actual := instance.Len()

	assert.ToBeEqual(t, 3, actual)
}

func Test_lineage_Len_withNilTarget(t *testing.T) {
	instance := &lineage{nil, With("a", 1)}

	actual := instance.Len()

	assert.ToBeEqual(t, 1, actual)
}

func Test_lineage_Len_withNilParent(t *testing.T) {
	instance := &lineage{With("a", 1), nil}

	actual := instance.Len()

	assert.ToBeEqual(t, 1, actual)
}

func Test_lineage_Len_withNilTargetAndParent(t *testing.T) {
	instance := &lineage{nil, nil}

	actual := instance.Len()

	assert.ToBeEqual(t, 0, actual)
}

//goland:noinspection GoNilness
func Test_lineage_Len_withNilInstance(t *testing.T) {
	var instance *lineage

	actual := instance.Len()

	assert.ToBeEqual(t, 0, actual)
}
