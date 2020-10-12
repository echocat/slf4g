package log

import (
	"errors"
	"reflect"

	"github.com/echocat/slf4g/fields"
)

// IsEventEqual is comparing two fields if all of its key values pairs are equal
// to each other by using the fields.DefaultEntryEqualityFunction.
//
// Important: The content of Event.GetContext() will not be compared.
func IsEventEqual(left, right Event) (bool, error) {
	return IsEventEqualCustom(fields.DefaultEntryEqualityFunction, left, right)
}

// IsEventEqualCustom is comparing two fields if all of its key values pairs are
// equal to each other by using the given EntryEqualityFunction. If for the
// function nil is provided the fields.DefaultEntryEqualityFunction is used.
//
// Important: The content of Event.GetContext() will not be compared.
func IsEventEqualCustom(eef fields.EntryEqualityFunction, left, right Event) (bool, error) {
	if eef == nil {
		eef = fields.DefaultEntryEqualityFunction
	}
	if eef == nil {
		return false, nil
	}

	if left == nil && right == nil {
		return true, nil
	}
	if left == nil || right == nil {
		return false, nil
	}

	if left.GetLevel() != right.GetLevel() {
		return false, nil
	}
	if left.GetCallDepth() != right.GetCallDepth() {
		return false, nil
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
