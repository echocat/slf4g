package fields

import (
	"github.com/echocat/slf4g/internal/test/assert"
	"testing"
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

	assert.ToBeNil(t, instance.Get("foo"))
	assert.ToBeNil(t, instance.Get("bar"))
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

func Test_empty_WithFields_isAlwaysTheInput(t *testing.T) {
	instance := Empty()
	given := newDummyFields()

	actual := instance.WithFields(given)

	assert.ToBeSame(t, given, actual)
}

func Test_empty_Without_isAlwaysSameEmptyInstance(t *testing.T) {
	instance := Empty()

	actual := instance.Without("a", "b")

	assert.ToBeSame(t, instance, actual)
}
