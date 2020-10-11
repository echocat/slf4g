package fields

import (
	"errors"
	"reflect"
)

// DefaultEntryEqualityFunction is the default instance of an
// EntryEqualityFunction. The initial initialization of this global variable
// should be able to deal with the majority of the cases.
var DefaultEntryEqualityFunction EntryEqualityFunction = func(key string, leftValue, rightValue interface{}) (bool, error) {
	if isFunction(leftValue) {
		lV, rV := reflect.ValueOf(leftValue), reflect.ValueOf(rightValue)
		if lV.Kind() == reflect.Func {
			return rV.Kind() == reflect.Func && lV.Pointer() == rV.Pointer(), nil
		}
	}

	return reflect.DeepEqual(leftValue, rightValue), nil
}

// IsEqual is comparing two fields if all of its key values pairs are equal to
// each other by using the DefaultEntryEqualityFunction.
func IsEqual(left, right Fields) (bool, error) {
	return IsEqualCustom(DefaultEntryEqualityFunction, left, right)
}

// CustomIsEqual is comparing two fields if all of its key values pairs are
// equal to each other by using the given EntryEqualityFunction. If for the
// function nil is provided the DefaultEntryEqualityFunction is used.
func IsEqualCustom(eef EntryEqualityFunction, left, right Fields) (bool, error) {
	if eef == nil {
		eef = DefaultEntryEqualityFunction
	}
	if eef == nil {
		return false, nil
	}

	if left == nil && right == nil {
		return true, nil
	}
	if left == nil {
		left = Empty()
	}
	if right == nil {
		right = Empty()
	}

	if left.Len() != right.Len() {
		return false, nil
	}
	if err := left.ForEach(func(key string, lValue interface{}) error {
		rValue, rExists := right.Get(key)
		if !rExists {
			return entriesNotEqualV
		}
		if equal, err := eef(key, lValue, rValue); err != nil {
			return err
		} else if !equal {
			return entriesNotEqualV
		} else {
			return nil
		}
	}); err == entriesNotEqualV {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// EntryEqualityFunction is comparing two values (of the same key) with each
// other and check if both are equal.
type EntryEqualityFunction func(key string, leftValue, rightValue interface{}) (bool, error)

var (
	entriesNotEqualV = errors.New(reflect.TypeOf(EntryEqualityFunction(nil)).PkgPath() + "/both entries are not equal")
)

func isFunction(arg interface{}) bool {
	return arg != nil && reflect.TypeOf(arg).Kind() == reflect.Func
}
