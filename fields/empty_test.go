package fields

import (
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_Empty(t *testing.T) {
	actual := Empty()

	assert.ToBeOfType(t, &empty{}, actual)
}

func Test_Empty_alwaysTheSame(t *testing.T) {
	actual1 := Empty()
	actual2 := Empty()

	assert.ToBeSame(t, actual1, actual2)
}

func Test_empty_ForEach_isNeverConsumingSomething(t *testing.T) {
	instance := Empty()

	actualErr := instance.ForEach(func(key string, value interface{}) error {
		assert.Failf(t, "Expected to be never call; but was called with <%+v>=<%+v>", key, value)
		return nil
	})

	assert.ToBeNoError(t, actualErr)
}

func Test_empty_Get_isAlwaysNil(t *testing.T) {
	instance := Empty()

	actual1, actual1Exists := instance.Get("foo")
	assert.ToBeNil(t, actual1)
	assert.ToBeEqual(t, false, actual1Exists)

	actual2, actual2Exists := instance.Get("foo")
	assert.ToBeNil(t, actual2)
	assert.ToBeEqual(t, false, actual2Exists)
}

func Test_empty_With_isAlwaysWithResult(t *testing.T) {
	instance := Empty()

	actual := instance.With("foo", "bar")

	assert.ToBeEqual(t, With("foo", "bar"), actual)
}

func Test_empty_Withf_isAlwaysWithfResult(t *testing.T) {
	instance := Empty()

	actual := instance.Withf("foo", "bar %s", "xyz")

	assert.ToBeEqual(t, Withf("foo", "bar %s", "xyz"), actual)
}

func Test_empty_WithAll_isAlwaysTheWrappedInput(t *testing.T) {
	instance := Empty()
	given := map[string]interface{}{
		"foo": 1,
		"bar": 2,
	}

	actual := instance.WithAll(given)

	assert.ToBeEqual(t, WithAll(given), actual)
}

func Test_empty_Without_isAlwaysSameEmptyInstance(t *testing.T) {
	instance := Empty()

	actual := instance.Without("a", "b")

	assert.ToBeSame(t, instance, actual)
}

func Test_empty_Len(t *testing.T) {
	instance := Empty()

	actual := instance.Len()

	assert.ToBeEqual(t, 0, actual)
}
