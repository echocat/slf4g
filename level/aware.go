package level

import (
	"reflect"
)

// Aware describes an object that is aware of a Level and exports its current
// state.
type Aware interface {
	// GetLevel returns the current level.
	GetLevel() Level
}

// Get returns the [Level] of the given object.
//
// It returns its [Level] together with `true` if either the direct given object or
// wrapped objects implements [Aware]. It will call recursively `Unwrap` on each
// object until it reaches a matching candidate or returns `0` along with `false` if
// nothing can be found.
func Get(of interface{}) (Level, bool) {
	ofVal := unwrap(of, awareType)
	if ofVal.IsValid() {
		return ofVal.Interface().(Aware).GetLevel(), true
	}
	return 0, false
}

// MutableAware is similar to [Aware] but additionally is able to modify the
// Level by calling SetLevel(Level).
type MutableAware interface {
	Aware

	// SetLevel modifies the current level to the given one.
	SetLevel(Level)
}

// Set sets the [Level] on the given object.
//
// It returns `true` if either the direct given object or wrapped objects implements
// [MutableAware] and it was possible to set the given [Level]. It will call
// recursively `Unwrap` on each object until it reaches a matching candidate or
// `false` if nothing can be found.
func Set(of interface{}, target Level) bool {
	ofVal := unwrap(of, mutableAwareType)
	if ofVal.IsValid() {
		ofVal.Interface().(MutableAware).SetLevel(target)
		return true
	}
	return false
}

var (
	awareType        = reflect.TypeOf((*Aware)(nil)).Elem()
	mutableAwareType = reflect.TypeOf((*MutableAware)(nil)).Elem()
)

func unwrap(v interface{}, targetType reflect.Type) reflect.Value {
	val := reflect.ValueOf(v)
	for {
		if !val.IsValid() {
			return reflect.Value{}
		}

		kind := val.Kind()
		switch kind {
		case reflect.Interface,
			reflect.Pointer,
			reflect.Chan,
			reflect.Func,
			reflect.Map,
			reflect.UnsafePointer,
			reflect.Slice:
			if val.IsNil() {
				return reflect.Value{}
			}
		default:
			// No nil checks
		}

		candidate, ok := checkUnwrapAssignment(val, targetType)
		if ok {
			return candidate
		}
		if !candidate.IsValid() {
			switch kind {
			case reflect.Interface,
				reflect.Pointer:
				candidate, ok = checkUnwrapAssignment(val.Elem(), targetType)
				if ok {
					return candidate
				}
				if candidate.IsValid() {
					val = candidate
					continue
				}
			default:
				// No more alternatives
			}

			return reflect.Value{}
		}
		val = candidate
	}
}

func checkUnwrapAssignment(val reflect.Value, targetType reflect.Type) (reflect.Value, bool) {
	typ := val.Type()
	if typ.AssignableTo(targetType) {
		return val, true
	}

	umt, ok := typ.MethodByName("Unwrap")
	if !ok ||
		!umt.IsExported() ||
		umt.Type.NumIn() != 1 ||
		umt.Type.NumOut() != 1 ||
		!umt.Type.In(0).AssignableTo(typ) {
		return reflect.Value{}, false
	}

	rts := umt.Func.Call([]reflect.Value{val})
	return rts[0], false
}
