package fields

import (
	"fmt"
	"testing"
)

func Test_LazyFunc_callsItselfOnGet(t *testing.T) {

}

func ExampleLazyFunc() {
	lazy := LazyFunc(func() interface{} {
		return someVariable.someResourceIntensiveMethod()
	})

	fmt.Println(lazy.Get())
}

var someVariable = &someStruct{}

type someStruct struct {
}

func (instance *someStruct) someResourceIntensiveMethod() string {
	return "foobar"
}
