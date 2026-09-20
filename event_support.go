package log

import (
	"fmt"
	"time"

	"github.com/echocat/slf4g/internal/eventvalue"
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
	case []any:
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
	pv, _ = eventvalue.Resolve(e, pv)
	return eventvalue.AsTimestamp(pv)
}

// GetLoggerOf returns for the given Event the contained logger (name)
// (if exists).
func GetLoggerOf(e Event, using Provider) *string {
	if e == nil {
		return nil
	}
	pv, _ := e.Get(using.GetFieldKeysSpec().GetLogger())
	pv, _ = eventvalue.Resolve(e, pv)
	return eventvalue.AsLogger(pv)
}

func resolveEventValue(event Event, value any) any {
	resolved, _ := eventvalue.Resolve(event, value)
	return resolved
}

type stringError string

func (instance stringError) Error() string {
	return string(instance)
}

func (instance stringError) String() string {
	return string(instance)
}
