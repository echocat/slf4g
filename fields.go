package log

import (
	"fmt"
	"github.com/echocat/slf4g/fields"
	"time"
)

func GetMessageOf(f fields.Fields, using Provider) *string {
	if f == nil {
		return nil
	}
	pv := f.Get(using.GetFieldKeySpec().GetMessage())
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

func GetErrorOf(f fields.Fields, using Provider) error {
	if f == nil {
		return nil
	}
	pv := f.Get(using.GetFieldKeySpec().GetError())
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

func GetTimestampOf(f fields.Fields, using Provider) *time.Time {
	if f == nil {
		return nil
	}
	pv := f.Get(using.GetFieldKeySpec().GetTimestamp())
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

func GetLoggerOf(f fields.Fields, using Provider) *string {
	type getNameAware interface {
		GetName() string
	}
	if f == nil {
		return nil
	}
	pv := f.Get(using.GetFieldKeySpec().GetLogger())
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
	case getNameAware:
		result := v.GetName()
		return &result
	default:
		result := fmt.Sprint(pv)
		return &result
	}
}
