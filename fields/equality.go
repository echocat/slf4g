package fields

import (
	"errors"
	"reflect"
)

// AreEqual is comparing two given Fields using DefaultEquality.
func AreEqual(left, right Fields) (bool, error) {
	if v := DefaultEquality; v != nil {
		return v.AreFieldsEqual(left, right)
	}
	return false, nil
}

// DefaultEquality is the default instance of a Equality. The initial
// initialization of this global variable should be able to deal with the
// majority of the cases. There is also a shortcut function:
// AreEqual(left,right)
var DefaultEquality Equality = &privateEqualityImpl{&EqualityImpl{
	ValueEquality: NewValueEqualityFacade(func() ValueEquality {
		return DefaultValueEquality
	}),
}}

// Equality is comparing two Fields with each other and check if both are equal.
type Equality interface {
	// AreFieldsEqual compares the two given Fields for their equality.
	AreFieldsEqual(left, right Fields) (bool, error)
}

// EqualityFunc is wrapping a given func into Equality.
type EqualityFunc func(left, right Fields) (bool, error)

// AreFieldsEqual implements Equality.AreFieldsEqual().
func (instance EqualityFunc) AreFieldsEqual(left, right Fields) (bool, error) {
	return instance(left, right)
}

// EqualityImpl is a default implementation of Equality which compares all its
// values using the configured ValueEquality.
type EqualityImpl struct {
	// ValueEquality is used to compare the values of the given fields. If nil
	// DefaultValueEquality is used. If this is nil too, there is always false
	// returned.
	ValueEquality ValueEquality
}

// AreFieldsEqual implements Equality.AreFieldsEqual().
func (instance *EqualityImpl) AreFieldsEqual(left, right Fields) (bool, error) {
	if instance == nil {
		return false, nil
	}

	ve := instance.ValueEquality

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

	if ve != nil {
		if err := left.ForEach(func(key string, lValue interface{}) error {
			rValue, rExists := right.Get(key)
			if !rExists {
				return entriesNotEqualV
			}
			if equal, err := ve.AreValuesEqual(key, lValue, rValue); err != nil {
				return err
			} else if !equal {
				return entriesNotEqualV
			}
			return nil
		}); err == entriesNotEqualV {
			return false, nil
		} else if err != nil {
			return false, err
		}
	}

	return true, nil
}

// NewEqualityFacade creates a re-implementation of Equality which uses the
// given provider to retrieve the actual instance of Equality in the moment when
// it is used. This is useful especially in cases where you want to deal with
// concurrency while creation of objects that need to hold a reference to an
// Equality.
func NewEqualityFacade(provider func() Equality) Equality {
	return equalityFacade(provider)
}

type equalityFacade func() Equality

func (instance equalityFacade) AreFieldsEqual(left, right Fields) (bool, error) {
	return instance().AreFieldsEqual(left, right)
}

type privateEqualityImpl struct {
	inner *EqualityImpl
}

func (instance *privateEqualityImpl) AreFieldsEqual(left, right Fields) (bool, error) {
	return instance.inner.AreFieldsEqual(left, right)
}

var (
	entriesNotEqualV = errors.New(reflect.TypeOf((*Equality)(nil)).PkgPath() + "/both entries are not equal")
)
