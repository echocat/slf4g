package log

import (
	"fmt"
	"time"
)

// GetMessageOf returns for the given Event the contained message (if exists).
func GetMessageOf(e Event, using Provider) *string {
	if e == nil {
		return nil
	}
	pv, _ := e.Get(using.GetFieldKeysSpec().GetMessage())
	if lv, ok := pv.(interface {
		Get() interface{}
	}); ok {
		pv = lv.Get()
	}
	switch v := pv.(type) {
	case nil:
		return nil
	case string:
		return &v
	case *string:
		return v
	case fmt.Stringer:
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

func formatStrSlice(in []string) string {
	var result []byte
	for i, v := range in {
		if i > 0 {
			result = append(result, ' ')
		}
		result = fmt.Append(result, v)
	}
	return string(result)
}

func formatAnySlice(in []interface{}) string {
	var result []byte
	for i, v := range in {
		if i > 0 {
			result = append(result, ' ')
		}
		result = fmt.Append(result, v)
	}
	return string(result)
}

// GetErrorOf returns for the given Event the contained error (if exists).
func GetErrorOf(e Event, using Provider) error {
	if e == nil {
		return nil
	}
	pv, _ := e.Get(using.GetFieldKeysSpec().GetError())
	if lv, ok := pv.(interface {
		Get() interface{}
	}); ok {
		pv = lv.Get()
	}
	switch v := pv.(type) {
	case nil:
		return nil
	case error:
		return v
	case string:
		return stringError(v)
	case *string:
		return stringError(*v)
	case fmt.Stringer:
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
	if lv, ok := pv.(interface {
		Get() interface{}
	}); ok {
		pv = lv.Get()
	}
	switch v := pv.(type) {
	case time.Time:
		if v.IsZero() {
			return nil
		}
		return &v
	case *time.Time:
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
	if lv, ok := pv.(interface {
		Get() interface{}
	}); ok {
		pv = lv.Get()
	}
	switch v := pv.(type) {
	case nil:
		return nil
	case string:
		return &v
	case *string:
		return v
	case Logger:
		result := v.GetName()
		return &result
	case interface {
		GetName() string
	}:
		result := v.GetName()
		return &result
	case fmt.Stringer:
		result := v.String()
		return &result
	default:
		result := fmt.Sprint(pv)
		return &result
	}
}

type stringError string

func (instance stringError) Error() string {
	return string(instance)
}

func (instance stringError) String() string {
	return string(instance)
}
