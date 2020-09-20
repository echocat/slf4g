package log

import (
	"fmt"
	"github.com/echocat/slf4g/fields"
	"time"
)

func GetMessageOf(e Event, using Provider) *string {
	if e == nil {
		return nil
	}
	pv := e.GetFields().Get(using.GetFieldKeySpec().GetMessage())
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
	default:
		result := fmt.Sprint(pv)
		return &result
	}
}

func GetErrorOf(e Event, using Provider) error {
	if e == nil {
		return nil
	}
	pv := e.GetFields().Get(using.GetFieldKeySpec().GetError())
	if lv, ok := pv.(fields.Lazy); ok {
		pv = lv.Get()
	}
	switch v := pv.(type) {
	case error:
		return v
	default:
		return nil
	}
}

func GetTimestampOf(e Event, using Provider) *time.Time {
	if e == nil {
		return nil
	}
	pv := e.GetFields().Get(using.GetFieldKeySpec().GetTimestamp())
	if lv, ok := pv.(fields.Lazy); ok {
		pv = lv.Get()
	}
	switch v := pv.(type) {
	case time.Time:
		return &v
	case *time.Time:
		return v
	default:
		return nil
	}
}

func GetLoggerOf(e Event, using Provider) *string {
	type getName interface {
		GetName() string
	}
	if e == nil {
		return nil
	}
	pv := e.GetFields().Get(using.GetFieldKeySpec().GetLogger())
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
