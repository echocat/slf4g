// Package eventvalue resolves and converts deferred event field values.
package eventvalue

import (
	"fmt"
	"time"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/internal/support"
)

// Resolve evaluates filtered and lazy field values. Respected is false only
// when a non-nil filter explicitly excludes the value.
func Resolve(context fields.FilterContext, value any) (resolved any, respected bool) {
	if value == nil {
		return nil, true
	}
	if filtered, ok := value.(fields.Filtered); ok {
		if support.IsNil(filtered) {
			return nil, true
		}
		resolved, respected = filtered.Filter(context)
		if !respected {
			return nil, false
		}
		return resolved, true
	}
	if lazy, ok := value.(fields.Lazy); ok {
		if support.IsNil(lazy) {
			return nil, true
		}
		return lazy.Get(), true
	}
	return value, true
}

// AsTimestamp converts a resolved field value into a non-zero timestamp.
func AsTimestamp(value any) *time.Time {
	switch v := value.(type) {
	case time.Time:
		if v.IsZero() {
			return nil
		}
		return &v
	case *time.Time:
		if v == nil || v.IsZero() {
			return nil
		}
		return v
	default:
		return nil
	}
}

// AsLogger converts a resolved field value into a logger name.
func AsLogger(value any) *string {
	switch v := value.(type) {
	case nil:
		return nil
	case string:
		return &v
	case *string:
		if v == nil {
			return nil
		}
		return v
	case interface{ GetName() string }:
		if support.IsNil(v) {
			return nil
		}
		result := v.GetName()
		return &result
	case fmt.Stringer:
		if support.IsNil(v) {
			return nil
		}
		result := v.String()
		return &result
	default:
		result := fmt.Sprint(value)
		return &result
	}
}
