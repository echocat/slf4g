package log

import (
	"fmt"
	"github.com/echocat/slf4g/fields"
	"time"
)

func GetMessageOf(fields fields.Fields, using Provider) *string {
	if fields == nil {
		return nil
	}
	pv := fields.Get(using.GetFieldKeys().GetMessage())
	if pv == nil {
		return nil
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

func GetErrorOf(fields fields.Fields, using Provider) error {
	if fields == nil {
		return nil
	}
	pv := fields.Get(using.GetFieldKeys().GetError())
	switch v := pv.(type) {
	case error:
		return v
	default:
		return nil
	}
}

func GetTimestampOf(fields fields.Fields, using Provider) *time.Time {
	if fields == nil {
		return nil
	}
	pv := fields.Get(using.GetFieldKeys().GetTimestamp())
	switch v := pv.(type) {
	case time.Time:
		return &v
	case *time.Time:
		return v
	default:
		return nil
	}
}

func GetLoggerOf(fields fields.Fields, using Provider) *string {
	type getNameAware interface {
		GetName() string
	}
	if fields == nil {
		return nil
	}
	pv := fields.Get(using.GetFieldKeys().GetLogger())
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
