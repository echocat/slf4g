package assert

import (
	"fmt"
	"reflect"
	"regexp"
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
		Failf(t, "Expected to be equal to: <%+v>; but got: <%+v>", expected, actual)
	}
}

func ToBeNotEqual(t testing.TB, expected, actual interface{}) {
	t.Helper()
	if isEqual(expected, actual) {
		Failf(t, "Expected to be no equal to: <%+v>; but got: <%+v>", expected, actual)
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

func ToBeNotNil(t testing.TB, actual interface{}) {
	t.Helper()
	ToBeNotEqual(t, nil, actual)
}

func Execution(t testing.TB, f func()) *ExecutionT {
	return &ExecutionT{t, f}
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
	if expectedV.Kind() == reflect.Func {
		return actualV.Kind() == reflect.Func && expectedV.Pointer() == actualV.Pointer()
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

type ExecutionT struct {
	testing.TB
	what func()
}

func (instance *ExecutionT) WillPanicWith(pattern string) {
	instance.TB.Helper()
	instance.WillPanicWithRegexp(regexp.MustCompile(pattern))
}

func (instance *ExecutionT) WillPanicWithRegexp(pattern *regexp.Regexp) {
	instance.TB.Helper()
	defer func() {
		instance.TB.Helper()
		if p := recover(); p != nil {
			if !pattern.MatchString(fmt.Sprint(p)) {
				instance.TB.Errorf("Expected to panics with: <%v>; but got: <%+v>", pattern.String(), p)
			}
		} else {
			instance.TB.Errorf("Expected to panics with: <%v>; but does not", pattern.String())
		}
	}()
	instance.what()
}
