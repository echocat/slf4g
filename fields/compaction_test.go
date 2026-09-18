package fields

import (
	"errors"
	"strconv"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_compaction_releasesOverwrittenHistory(t *testing.T) {
	original := With("key", "original")
	instance := original
	for i := 1; i <= minimumCompactionCost; i++ {
		instance = instance.With("key", i)
	}

	compacted, ok := instance.(*compactedFields)
	assert.ToBeEqual(t, true, ok)
	assert.ToBeEqual(t, map[string]interface{}{"key": minimumCompactionCost}, compacted.values)
	assert.ToBeEqual(t, []string{"key"}, compacted.keys)
	assert.ToBeNil(t, compacted.base)
	assert.ToBeEqual(t, minimumCompactionCost, compacted.compactAt)

	actual, actualExists := instance.Get("key")
	assert.ToBeEqual(t, minimumCompactionCost, actual)
	assert.ToBeEqual(t, true, actualExists)
	originalValue, originalExists := original.Get("key")
	assert.ToBeEqual(t, "original", originalValue)
	assert.ToBeEqual(t, true, originalExists)
}

func Test_compaction_releasesHistoryCreatedWithPublicConstructor(t *testing.T) {
	instance := With("key", "original")
	for i := 1; i <= minimumCompactionCost; i++ {
		instance = NewLineage(With("key", i), instance)
	}

	compacted, ok := instance.(*compactedFields)
	assert.ToBeEqual(t, true, ok)
	assert.ToBeEqual(t, map[string]interface{}{"key": minimumCompactionCost}, compacted.values)
	assert.ToBeNil(t, compacted.base)
}

func Test_compaction_preservesOrderAndGrowsThreshold(t *testing.T) {
	instance := With("base", -1)
	for i := 0; i < minimumCompactionCost; i++ {
		instance = instance.With(strconv.Itoa(i), i)
	}

	compacted := instance.(*compactedFields)
	assert.ToBeEqual(t, minimumCompactionCost+1, compacted.compactAt)
	assert.ToBeEqual(t, minimumCompactionCost+1, compacted.Len())

	var actualKeys []string
	actualErr := compacted.ForEach(func(key string, _ interface{}) error {
		actualKeys = append(actualKeys, key)
		return nil
	})
	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, strconv.Itoa(minimumCompactionCost-1), actualKeys[0])
	assert.ToBeEqual(t, "base", actualKeys[len(actualKeys)-1])

	derived := compacted.With("after", true)
	assert.ToBeOfType(t, &lineage{}, derived)
	assert.ToBeEqual(t, 1, derived.(*lineage).pendingCompactionCost)
	assert.ToBeEqual(t, minimumCompactionCost+1, derived.(*lineage).compactAt)
}

func Test_compaction_releasesRemovedHistory(t *testing.T) {
	instance := With("dead", &struct{}{})
	for i := 0; i < minimumCompactionCost; i++ {
		instance = instance.Without("dead")
	}

	assert.ToBeSame(t, Empty(), instance)
}

func Test_compaction_releasesRemovedHistoryCreatedWithPublicConstructor(t *testing.T) {
	instance := With("dead", &struct{}{})
	for i := 0; i < minimumCompactionCost; i++ {
		instance = NewWithout(instance, "dead")
	}

	assert.ToBeSame(t, Empty(), instance)
}

func Test_compaction_preservesReaddedValue(t *testing.T) {
	instance := With("key", "old").Without("key")
	for i := 0; i < minimumCompactionCost-2; i++ {
		instance = instance.With("other", i)
	}
	instance = instance.With("key", "new")

	actual, actualExists := instance.Get("key")
	assert.ToBeEqual(t, "new", actual)
	assert.ToBeEqual(t, true, actualExists)
	assert.ToBeEqual(t, 2, instance.Len())
}

func Test_compaction_preservesLocalWithoutScope(t *testing.T) {
	hiddenTarget := NewWithout(
		NewLineage(With("shared", "target"), With("shared", "target-parent")),
		"shared",
	)
	instance := NewLineage(hiddenTarget, With("shared", "parent"))
	for i := 0; i < minimumCompactionCost; i++ {
		instance = instance.With("overlay", i)
	}

	actual, actualErr := asMap(instance)

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, mapped{"overlay": minimumCompactionCost - 1, "shared": "parent"}, actual)
}

func Test_compaction_doesNotEvaluateLazyOrFilteredValues(t *testing.T) {
	lazyCalls := 0
	newLazy := func(value int) Lazy {
		return LazyFunc(func() interface{} {
			lazyCalls++
			return value
		})
	}
	instance := With("lazy", newLazy(0)).With("filtered", RequireMaximalLevelLazy(0, newLazy(0)))
	for i := 1; i <= minimumCompactionCost; i++ {
		instance = instance.With("lazy", newLazy(i)).With("filtered", RequireMaximalLevelLazy(0, newLazy(i)))
	}

	assert.ToBeEqual(t, 0, lazyCalls)
	actualLazy, actualLazyExists := instance.Get("lazy")
	assert.ToBeEqual(t, true, actualLazyExists)
	assert.ToBeEqual(t, minimumCompactionCost, actualLazy.(Lazy).Get())
	assert.ToBeEqual(t, 1, lazyCalls)
	actualFiltered, actualFilteredExists := instance.Get("filtered")
	assert.ToBeEqual(t, true, actualFilteredExists)
	assert.ToBeOfType(t, RequireMaximalLevelLazy(0, nil), actualFiltered)
	assert.ToBeEqual(t, 1, lazyCalls)
}

func Test_compaction_preservesOpaqueBaseEvaluation(t *testing.T) {
	expectedErr := errors.New("expected")
	base := &observedFields{
		values: map[string]interface{}{"shadowed": "base", "base": "visible"},
		err:    expectedErr,
	}
	var instance Fields = base
	for i := 0; i <= minimumCompactionCost; i++ {
		instance = instance.With("shadowed", i)
	}

	assert.ToBeEqual(t, 0, base.forEachCalls)
	assert.ToBeEqual(t, 0, base.getCalls)
	assert.ToBeEqual(t, 0, base.lenCalls)

	var actualKeys []string
	actualErr := instance.ForEach(func(key string, _ interface{}) error {
		actualKeys = append(actualKeys, key)
		return nil
	})
	assert.ToBeEqual(t, expectedErr, actualErr)
	assert.ToBeEqual(t, []string{"shadowed", "base"}, actualKeys)
	assert.ToBeEqual(t, 1, base.forEachCalls)

	actualShadowed, actualShadowedExists := instance.Get("shadowed")
	assert.ToBeEqual(t, minimumCompactionCost, actualShadowed)
	assert.ToBeEqual(t, true, actualShadowedExists)
	assert.ToBeEqual(t, 0, base.getCalls)
	actualBase, actualBaseExists := instance.Get("base")
	assert.ToBeEqual(t, "visible", actualBase)
	assert.ToBeEqual(t, true, actualBaseExists)
	assert.ToBeEqual(t, 1, base.getCalls)

	assert.ToBeEqual(t, 2, instance.Len())
	assert.ToBeEqual(t, 2, base.forEachCalls)
	assert.ToBeEqual(t, 0, base.lenCalls)
}

func Test_compaction_LenIgnoresOpaqueBaseErrors(t *testing.T) {
	expectedErr := errors.New("expected")
	base := NewLineage(
		&observedFields{values: map[string]interface{}{"target": true}, err: expectedErr},
		With("parent", true),
	)
	instance := base
	for i := 0; i < minimumCompactionCost; i++ {
		instance = instance.With("overlay", i)
	}

	assert.ToBeOfType(t, &compactedFields{}, instance)
	assert.ToBeEqual(t, 3, instance.Len())
	actualErr := instance.ForEach(func(string, interface{}) error { return nil })
	assert.ToBeEqual(t, expectedErr, actualErr)
}

func Test_compaction_matchesUncompactedMixedOperations(t *testing.T) {
	actual := With("base", -1)
	expected := With("base", -1)
	for i := 0; i < 1_000; i++ {
		key := strconv.Itoa(i % 17)
		if i%5 == 0 {
			actual = actual.Without(key)
			expected = &without{fields: expected, excludedKeys: keySet{key: keyPresent}}
		} else {
			actual = actual.With(key, i)
			expected = &lineage{target: With(key, i), parent: expected}
		}
	}

	assert.ToBeEqual(t, expected.Len(), actual.Len())
	for i := 0; i < 17; i++ {
		key := strconv.Itoa(i)
		expectedValue, expectedExists := expected.Get(key)
		actualValue, actualExists := actual.Get(key)
		assert.ToBeEqual(t, expectedValue, actualValue)
		assert.ToBeEqual(t, expectedExists, actualExists)
	}
	var expectedKeys []string
	expectedErr := expected.ForEach(func(key string, _ interface{}) error {
		expectedKeys = append(expectedKeys, key)
		return nil
	})
	var actualKeys []string
	actualErr := actual.ForEach(func(key string, _ interface{}) error {
		actualKeys = append(actualKeys, key)
		return nil
	})
	assert.ToBeNoError(t, expectedErr)
	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, expectedKeys, actualKeys)
}

type observedFields struct {
	values       map[string]interface{}
	err          error
	forEachCalls int
	getCalls     int
	lenCalls     int
}

func (instance *observedFields) ForEach(consumer func(key string, value interface{}) error) error {
	instance.forEachCalls++
	for key, value := range instance.values {
		if err := consumer(key, value); err != nil {
			return err
		}
	}
	return instance.err
}

func (instance *observedFields) Get(key string) (interface{}, bool) {
	instance.getCalls++
	value, exists := instance.values[key]
	return value, exists
}

func (instance *observedFields) Len() int {
	instance.lenCalls++
	return len(instance.values)
}

func (instance *observedFields) With(key string, value interface{}) Fields {
	return NewLineage(With(key, value), instance)
}

func (instance *observedFields) Withf(key string, format string, args ...interface{}) Fields {
	return NewLineage(Withf(key, format, args...), instance)
}

func (instance *observedFields) WithAll(values map[string]interface{}) Fields {
	return NewLineage(WithAll(values), instance)
}

func (instance *observedFields) Without(keys ...string) Fields {
	return NewWithout(instance, keys...)
}
