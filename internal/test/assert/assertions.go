package assert

import (
	"reflect"
	"testing"
)

func ToBeSame(t testing.TB, expected, actual interface{}) {
	t.Helper()
	if !isSame(expected, actual) {
		Failf(t, "Expected to be same as: <%+v>; but got: <%+v>", expected, actual)
	}
}

func ToBeEqual(t testing.TB, expected, actual interface{}) {
	t.Helper()
	if !isEqual(expected, actual) {
		Failf(t, "Expected to equal to: <%+v>; but got: <%+v>", expected, actual)
	}
}

func ToBeOfType(t testing.TB, expectedType, actual interface{}) {
	t.Helper()
	if !isType(expectedType, actual) {
		Failf(t, "Expected to be type of: <%+v>; but got: <%+v>", reflect.TypeOf(expectedType), reflect.TypeOf(actual))
	}
}

func ToBeNoError(t testing.TB, actual error) {
	t.Helper()
	if actual != nil {
		Failf(t, "Expected to be no error; but got: <%+v>", reflect.TypeOf(actual))
	}
}

func ToBeNil(t testing.TB, actual interface{}) {
	t.Helper()
	ToBeEqual(t, nil, actual)
}

func Fail(t testing.TB, fmt string, args ...interface{}) {
	t.Helper()
	t.Errorf(fmt, args...)
}

func Failf(t testing.TB, fmt string, args ...interface{}) {
	t.Helper()
	t.Errorf(fmt, args...)
}

func isSame(expected, actual interface{}) bool {
	if expected == nil && actual == nil {
		return true
	}
	expectedV, actualV := reflect.ValueOf(expected), reflect.ValueOf(actual)
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
