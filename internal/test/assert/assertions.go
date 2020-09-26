package assert

import (
	"reflect"
)

type T interface {
	Errorf(string, ...interface{})
}

func ToBeSame(t T, expected, actual interface{}) {
	if !isSame(expected, actual) {
		Failf(t, "Expected to be same as: <%+v>; but got: <%+v>", expected, actual)
	}
}

func ToBeEqual(t T, expected, actual interface{}) {
	if !isEqual(expected, actual) {
		Failf(t, "Expected to equal to: <%+v>; but got: <%+v>", expected, actual)
	}
}

func ToBeOfType(t T, expectedType, actual interface{}) {
	if !isType(expectedType, actual) {
		Failf(t, "Expected to be type of: <%+v>; but got: <%+v>", reflect.TypeOf(expectedType), reflect.TypeOf(actual))
	}
}

func ToBeNoError(t T, actual error) {
	if actual != nil {
		Failf(t, "Expected to be no error; but got: <%+v>", reflect.TypeOf(actual))
	}
}

func ToBeNil(t T, actual interface{}) {
	ToBeEqual(t, nil, actual)
}

func Failf(t T, fmt string, args ...interface{}) {
	t.Errorf(fmt, args)
}

func isSame(expected, actual interface{}) bool {
	if expected == nil && actual == nil {
		return true
	}
	expectedV, actualV := reflect.ValueOf(expected), reflect.ValueOf(actual)
	if expectedV.Kind() != reflect.Ptr && actualV.Kind() != reflect.Ptr {
		return expected == actual
	}
	if expectedV.Kind() != reflect.Ptr || actualV.Kind() != reflect.Ptr {
		return false
	}
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		return false
	}
	return expected == actual
}

func isEqual(expected, actual interface{}) bool {
	if isFunction(expected) {
		return isSame(expected, actual)
	}
	return reflect.DeepEqual(expected, actual)
}

func isType(expectedType, actual interface{}) bool {
	return reflect.TypeOf(expectedType) == reflect.TypeOf(actual)
}

func isFunction(arg interface{}) bool {
	if arg == nil {
		return false
	}
	return reflect.TypeOf(arg).Kind() == reflect.Func
}
