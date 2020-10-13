package fields

import "reflect"

// DefaultValueEquality is the default instance of a ValueEquality. The initial
// initialization of this global variable should be able to deal with
// the majority of the cases.
var DefaultValueEquality ValueEquality = ValueEqualityFunc(func(key string, leftValue, rightValue interface{}) (bool, error) {
	if v, ok := leftValue.(Lazy); ok {
		leftValue = v.Get()
	}
	if v, ok := rightValue.(Lazy); ok {
		rightValue = v.Get()
	}

	if isFunction(leftValue) {
		lV, rV := reflect.ValueOf(leftValue), reflect.ValueOf(rightValue)
		if lV.Kind() == reflect.Func {
			return rV.Kind() == reflect.Func && lV.Pointer() == rV.Pointer(), nil
		}
	}

	return reflect.DeepEqual(leftValue, rightValue), nil
})

// ValueEquality is comparing two values (of the same key) with each
// other and check if both are equal.
type ValueEquality interface {
	// AreValuesEqual compares the two given values for their equality for
	// the given key.
	AreValuesEqual(key string, left, right interface{}) (bool, error)
}

// ValueEqualityFunc is wrapping a given func into ValueEquality.
type ValueEqualityFunc func(key string, left, right interface{}) (bool, error)

// AreValuesEqual implements ValueEquality.AreValuesEqual().
func (instance ValueEqualityFunc) AreValuesEqual(key string, left, right interface{}) (bool, error) {
	return instance(key, left, right)
}

// NewValueEqualityFacade creates a re-implementation of ValueEquality which
// uses the given provider to retrieve the actual instance of ValueEquality in
// the moment when it is used. This is useful especially in cases where you want
// to deal with concurrency while creation of objects that need to hold a
// reference to an ValueEquality.
func NewValueEqualityFacade(provider func() ValueEquality) ValueEquality {
	return valueEqualityFacade(provider)
}

type valueEqualityFacade func() ValueEquality

func (instance valueEqualityFacade) AreValuesEqual(name string, left, right interface{}) (bool, error) {
	return instance().AreValuesEqual(name, left, right)
}

func isFunction(arg interface{}) bool {
	return arg != nil && reflect.TypeOf(arg).Kind() == reflect.Func
}
