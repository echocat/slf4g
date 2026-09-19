// Package support provides helper methods for internal usage.
package support

import (
	"reflect"
	"time"
)

// IsNil detects nil values wrapped in an interface. Callers should prefer a
// direct nil check and only use this as a fallback before calling a method.
func IsNil(v any) bool {
	if v == nil {
		return true
	}

	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return value.IsNil()
	default:
		return false
	}
}

func PString(v string) *string {
	return &v
}

func PTime(v time.Time) *time.Time {
	return &v
}
