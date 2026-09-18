package fields

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func ExampleLazyFunc() {
	lazy := LazyFunc(func() interface{} {
		return someVariable.someResourceIntensiveMethod()
	})

	fmt.Println(lazy.Get())

	// Output:
	// foobar
}

func Test_LazyFunc_callsItselfOnGet(t *testing.T) {
	expected := struct{ foo string }{foo: "bar"}
	givenProvider := func() interface{} { return expected }

	actualInstance := LazyFunc(givenProvider)
	actual := actualInstance.Get()

	assert.ToBeEqual(t, expected, actual)
}

func Test_LazyFunc_nil(t *testing.T) {
	actual := LazyFunc(nil).Get()

	assert.ToBeNil(t, actual)
}

func ExampleLazyFormat() {
	lazy := LazyFormat("Hello, %s!", "world")

	fmt.Println(lazy.Get())

	// Output:
	//Hello, world!
}

func Test_LazyFormat_formats(t *testing.T) {
	actualCallAmount := uint64(0)

	instance := LazyFormat("foo%s", LazyFunc(func() interface{} {
		atomic.AddUint64(&actualCallAmount, 1)
		return "bar"
	}))
	assert.ToBeEqual(t, uint64(0), atomic.LoadUint64(&actualCallAmount))

	actual := instance.Get()

	assert.ToBeEqual(t, uint64(1), atomic.LoadUint64(&actualCallAmount))
	assert.ToBeEqual(t, "foobar", actual)
}

func Test_LazyFormat_copiesArguments(t *testing.T) {
	args := []interface{}{"before"}
	instance := LazyFormat("%s", args...)

	args[0] = "after"

	assert.ToBeEqual(t, "before", instance.Get())
}

func Test_LazyFormat_withTypedNilLazy(t *testing.T) {
	var value *nilLazy

	actual := LazyFormat("%v", value).Get()

	assert.ToBeEqual(t, "<nil>", actual)
}

var someVariable = &someStruct{}

type someStruct struct {
}

type nilLazy struct{}

func (*nilLazy) Get() interface{} {
	panic("must not be called")
}

func (instance *someStruct) someResourceIntensiveMethod() string {
	return "foobar"
}
