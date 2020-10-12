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

func ToBeNotSame(t testing.TB, expected, actual interface{}) {
	t.Helper()
	if isSame(expected, actual) {
		Failf(t, "Expected to be not same as: <%+v>; but got: <%+v>", expected, actual)
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
		Failf(t, "Expected to be not equal to: <%+v>; but got: <%+v>", expected, actual)
	}
}

func ToBeEqualUsing(t testing.TB, expected, actual interface{}, comparator interface{}) {
	t.Helper()
	isEqual, err := callComparator(expected, actual, comparator)
	if err != nil {
		Failf(t, "Expected to be no error; but got: <%+v>", err)
	} else if !isEqual {
		Failf(t, "Expected to be equal to: <%+v>; but got: <%+v>", expected, actual)
	}
}

func ToBeNotEqualUsing(t testing.TB, expected, actual interface{}, comparator interface{}) {
	t.Helper()
	isEqual, err := callComparator(expected, actual, comparator)
	if err != nil {
		Failf(t, "Expected to be no error; but got: <%+v>", err)
	} else if !isEqual {
		Failf(t, "Expected to be not equal to: <%+v>; but got: <%+v>", expected, actual)
	}
}

func ToBeMatching(t testing.TB, expectedPattern string, actual interface{}) {
	t.Helper()
	if actual == nil {
		Failf(t, "Expected to be matching: <%s>; but got: <%+v>", expectedPattern, actual)
	} else {
		expectedRegexp := regexp.MustCompile(expectedPattern)
		actualStr := fmt.Sprint(actual)
		if !expectedRegexp.MatchString(actualStr) {
			Failf(t, "Expected to be matching: <%s>; but got: <%+v>", expectedPattern, actual)
		}
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
	if !isNil(actual) {
		Failf(t, "Expected to be nil; but got: <%+v>", actual)
	}
}

func ToBeNotNil(t testing.TB, actual interface{}) {
	t.Helper()
	if isNil(actual) {
		Failf(t, "Expected to be not nil; but got: <%+v>", actual)
	}
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

func callComparator(expected, actual interface{}, comparator interface{}) (equal bool, err error) {
	if comparator == nil {
		panic("comparator of kind func expected; but got: <nil>")
	}
	cv := reflect.ValueOf(comparator)
	if cv.Kind() != reflect.Func {
		panic(fmt.Sprintf("comparator of kind func expected; but got: <%+v>", comparator))
	}
	ct := cv.Type()
	if ct.NumIn() != 2 ||
		ct.NumOut() != 2 ||
		ct.Out(0) != reflect.TypeOf(true) ||
		ct.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		panic(fmt.Sprintf("comparator of signature <func(interface{}, interface{}) (bool, error)>; but got: <%+v>", ct))
	}
	if ct.In(0) != ct.In(1) {
		panic(fmt.Sprintf("comparator of signature with same expected and actual type expected; but got: <%+v>", ct))
	}
	result := cv.Call([]reflect.Value{
		reflect.ValueOf(expected),
		reflect.ValueOf(actual),
	})
	if v, ok := result[0].Interface().(bool); ok {
		equal = v
	}
	if v, ok := result[1].Interface().(error); ok {
		err = v
	}
	return
}

func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	isNilableKind := containsKind(
		[]reflect.Kind{
			reflect.Chan, reflect.Func,
			reflect.Interface, reflect.Map,
			reflect.Ptr, reflect.Slice},
		kind)

	if isNilableKind && value.IsNil() {
		return true
	}

	return false
}

func containsKind(kinds []reflect.Kind, kind reflect.Kind) bool {
	for i := 0; i < len(kinds); i++ {
		if kind == kinds[i] {
			return true
		}
	}

	return false
}
