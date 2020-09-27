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
}

func Test_LazyFunc_callsItselfOnGet(t *testing.T) {
	expected := struct{ foo string }{foo: "bar"}
	givenProvider := func() interface{} { return expected }

	actualInstance := LazyFunc(givenProvider)
	actual := actualInstance.Get()

	assert.ToBeEqual(t, expected, actual)
}

func ExampleLazyFormat() {
	lazy := LazyFormat("Hello, %s!", "world")

	fmt.Println(lazy.Get())
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

var someVariable = &someStruct{}

type someStruct struct {
}

func (instance *someStruct) someResourceIntensiveMethod() string {
	return "foobar"
}
