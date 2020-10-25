package log

import (
	"errors"
	"reflect"

	"github.com/echocat/slf4g/fields"
)

// AreEventsEqual is comparing two given Events using DefaultEventEquality.
func AreEventsEqual(left, right Event) (bool, error) {
	if v := DefaultEventEquality; v != nil {
		return v.AreEventsEqual(left, right)
	}
	return false, nil
}

// DefaultEventEquality is the default instance of a Equality. The initial
// initialization of this global variable should be able to deal with the
// majority of the cases. There is also a shortcut function:
// AreEventsEqual(left,right)
var DefaultEventEquality EventEquality = &privateEventEqualityImpl{&EventEqualityImpl{
	CompareLevel: true,
	CompareValuesUsing: fields.NewValueEqualityFacade(func() fields.ValueEquality {
		return fields.DefaultValueEquality
	}),
}}

// EventEquality is comparing two Fields with each other and check if both are
// equal.
type EventEquality interface {
	// AreEventsEqual compares the two given Events for their equality.
	AreEventsEqual(left, right Event) (bool, error)

	// WithIgnoringKeys creates a new instance of this EventEquality which will
	// ignore the given keys while the equality check run.
	WithIgnoringKeys(keys ...string) EventEquality
}

// EventEqualityFunc is wrapping a given func into EventEquality.
type EventEqualityFunc func(left, right Event) (bool, error)

// AreEventsEqual implements EventEquality.AreEventsEqual().
func (instance EventEqualityFunc) AreEventsEqual(left, right Event) (bool, error) {
	return instance(left, right)
}

// WithIgnoringKeys implements EventEquality.WithIgnoringKeys().
func (instance EventEqualityFunc) WithIgnoringKeys(keys ...string) EventEquality {
	return &ignoringKeysEventEquality{instance, keys}
}

type EventEqualityImpl struct {
	// CompareLevel will configure to compare the Event.GetLevel() of both
	// events to be the same.
	CompareLevel bool

	// CompareValuesUsing is used to compare the fields of the given events. If
	// nil the values will not be compared.
	CompareValuesUsing fields.ValueEquality
}

// AreEventsEqual implements EventEquality.AreEventsEqual().
func (instance *EventEqualityImpl) AreEventsEqual(left, right Event) (bool, error) {
	if instance == nil {
		return false, nil
	}

	if left == nil && right == nil {
		return true, nil
	}
	if left == nil || right == nil {
		return false, nil
	}

	if instance.CompareLevel && left.GetLevel() != right.GetLevel() {
		return false, nil
	}

	if ve := instance.CompareValuesUsing; ve != nil {
		if left.Len() != right.Len() {
			return false, nil
		}

		if err := left.ForEach(func(key string, lValue interface{}) error {
			rValue, rExists := right.Get(key)
			if !rExists {
				return entriesNotEqualV
			}
			if equal, err := ve.AreValuesEqual(key, lValue, rValue); err != nil {
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
		}
	}

	return true, nil
}

// WithIgnoringKeys implements EventEquality.WithIgnoringKeys().
func (instance *EventEqualityImpl) WithIgnoringKeys(keys ...string) EventEquality {
	return &ignoringKeysEventEquality{instance, keys}
}

// NewEventEqualityFacade creates a re-implementation of EventEquality which
// uses the given provider to retrieve the actual instance of EventEquality in
// the moment when it is used. This is useful especially in cases where you want
// to deal with concurrency while creation of objects that need to hold a
// reference to a EventEquality.
func NewEventEqualityFacade(provider func() EventEquality) EventEquality {
	return eventEqualityFacade(provider)
}

type eventEqualityFacade func() EventEquality

func (instance eventEqualityFacade) AreEventsEqual(left, right Event) (bool, error) {
	return instance().AreEventsEqual(left, right)
}

func (instance eventEqualityFacade) WithIgnoringKeys(keys ...string) EventEquality {
	return &ignoringKeysEventEquality{instance, keys}
}

type privateEventEqualityImpl struct {
	inner *EventEqualityImpl
}

func (instance *privateEventEqualityImpl) AreEventsEqual(left, right Event) (bool, error) {
	return instance.inner.AreEventsEqual(left, right)
}

func (instance *privateEventEqualityImpl) WithIgnoringKeys(keys ...string) EventEquality {
	return &ignoringKeysEventEquality{instance, keys}
}

type ignoringKeysEventEquality struct {
	parent       EventEquality
	keysToIgnore []string
}

func (instance *ignoringKeysEventEquality) AreEventsEqual(left, right Event) (bool, error) {
	if left == nil && right == nil {
		return true, nil
	}
	if left == nil || right == nil {
		return false, nil
	}

	return instance.parent.AreEventsEqual(
		left.Without(instance.keysToIgnore...),
		right.Without(instance.keysToIgnore...),
	)
}

func (instance *ignoringKeysEventEquality) WithIgnoringKeys(keys ...string) EventEquality {
	return &ignoringKeysEventEquality{
		parent:       instance.parent,
		keysToIgnore: append(instance.keysToIgnore, keys...),
	}
}

var (
	entriesNotEqualV = errors.New(reflect.TypeOf((*EventEquality)(nil)).PkgPath() + "/both entries are not equal")
)
