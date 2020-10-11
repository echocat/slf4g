package log

import (
	"fmt"
	"time"

	"github.com/echocat/slf4g/fields"
)

// GetMessageOf returns for the given Event the contained message (if exists).
func GetMessageOf(e Event, using Provider) *string {
	if e == nil {
		return nil
	}
	pv, _ := e.Get(using.GetFieldKeysSpec().GetMessage())
	if pv == nil {
		return nil
	}
	if lv, ok := pv.(fields.Lazy); ok {
		pv = lv.Get()
	}
	switch v := pv.(type) {
	case string:
		return &v
	case *string:
		return v
	case fmt.Stringer:
		s := v.String()
		return &s
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
	if lv, ok := pv.(fields.Lazy); ok {
		pv = lv.Get()
	}
	switch v := pv.(type) {
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
	if lv, ok := pv.(fields.Lazy); ok {
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
	type getName interface {
		GetName() string
	}
	if e == nil {
		return nil
	}
	pv, _ := e.Get(using.GetFieldKeysSpec().GetLogger())
	if lv, ok := pv.(fields.Lazy); ok {
		pv = lv.Get()
	}
	if pv == nil {
		return nil
	}
	switch v := pv.(type) {
	case string:
		return &v
	case *string:
		return v
	case getName:
		result := v.GetName()
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
