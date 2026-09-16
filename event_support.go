package log

import (
	"fmt"
	"time"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/internal/support"
)

// GetMessageOf returns for the given Event the contained message (if exists).
func GetMessageOf(e Event, using Provider) *string {
	if e == nil {
		return nil
	}
	pv, _ := e.Get(using.GetFieldKeysSpec().GetMessage())
	pv = resolveEventValue(e, pv)
	switch v := pv.(type) {
	case nil:
		return nil
	case string:
		return &v
	case *string:
		if v == nil {
			return nil
		}
		return v
	case fmt.Stringer:
		if support.IsNil(v) {
			return nil
		}
		s := v.String()
		return &s
	case []string:
		result := formatStrSlice(v)
		return &result
	case []interface{}:
		result := formatAnySlice(v)
		return &result
	default:
		result := fmt.Sprint(pv)
		return &result
	}
}

// GetErrorOf returns for the given Event the contained error (if exists).
func GetErrorOf(e Event, using Provider) error {
	if e == nil {
		return nil
	}
	pv, _ := e.Get(using.GetFieldKeysSpec().GetError())
	pv = resolveEventValue(e, pv)
	switch v := pv.(type) {
	case nil:
		return nil
	case error:
		if support.IsNil(v) {
			return nil
		}
		return v
	case string:
		return stringError(v)
	case *string:
		if v == nil {
			return nil
		}
		return stringError(*v)
	case fmt.Stringer:
		if support.IsNil(v) {
			return nil
		}
		return stringError(v.String())
	default:
		return stringError(fmt.Sprint(pv))
	}
}

// GetTimestampOf returns for the given Event the contained timestamp
// (if exists).
func GetTimestampOf(e Event, using Provider) *time.Time {
	if e == nil {
		return nil
	}
	pv, _ := e.Get(using.GetFieldKeysSpec().GetTimestamp())
	pv = resolveEventValue(e, pv)
	switch v := pv.(type) {
	case time.Time:
		if v.IsZero() {
			return nil
		}
		return &v
	case *time.Time:
		if v == nil {
			return nil
		}
		if v.IsZero() {
			return nil
		}
		return v
	default:
		return nil
	}
}

// GetLoggerOf returns for the given Event the contained logger (name)
// (if exists).
func GetLoggerOf(e Event, using Provider) *string {
	if e == nil {
		return nil
	}
	pv, _ := e.Get(using.GetFieldKeysSpec().GetLogger())
	pv = resolveEventValue(e, pv)
	switch v := pv.(type) {
	case nil:
		return nil
	case string:
		return &v
	case *string:
		if v == nil {
			return nil
		}
		return v
	case Logger:
		if support.IsNil(v) {
			return nil
		}
		result := v.GetName()
		return &result
	case interface {
		GetName() string
	}:
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
		result := fmt.Sprint(pv)
		return &result
	}
}

func resolveEventValue(event Event, value interface{}) interface{} {
	if value == nil {
		return nil
	}
	if filtered, ok := value.(fields.Filtered); ok {
		if support.IsNil(filtered) {
			return nil
		}
		resolved, respected := filtered.Filter(event)
		if !respected {
			return nil
		}
		return resolved
	}
	if lazy, ok := value.(fields.Lazy); ok {
		if support.IsNil(lazy) {
			return nil
		}
		return lazy.Get()
	}
	return value
}

type stringError string

func (instance stringError) Error() string {
	return string(instance)
}

func (instance stringError) String() string {
	return string(instance)
}
