package formatter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/internal/support"
	"github.com/echocat/slf4g/level"
)

// SimpleTextValue is a simple implementation of TextValue.
type SimpleTextValue struct {
	// QuoteType defines how values are quoted.
	QuoteType QuoteType
}

// NewSimpleTextValue creates a new instance of SimpleTextValue which is ready
// to use.
func NewSimpleTextValue(customizer ...func(*SimpleTextValue)) *SimpleTextValue {
	result := &SimpleTextValue{
		QuoteType: QuoteTypeMinimal,
	}
	for _, c := range customizer {
		c(result)
	}
	return result
}

// FormatTextValue implements TextValue.FormatTextValue().
func (instance *SimpleTextValue) FormatTextValue(v interface{}, _ log.Provider) ([]byte, error) {
	if lazy, ok := v.(fields.Lazy); ok {
		if support.IsNil(lazy) {
			v = nil
		} else {
			v = lazy.Get()
		}
	}

	switch vs := v.(type) {
	case nil:
		v = ""
	case *string:
		if vs == nil {
			v = ""
		} else {
			v = *vs
		}
	case json.Number:
		v = vs.String()
	case time.Time:
		v = vs.String()
	case time.Duration:
		v = vs.String()
	case fmt.Stringer:
		if support.IsNil(vs) {
			v = ""
		} else {
			v = vs.String()
		}
	case fmt.Formatter:
		if support.IsNil(vs) {
			v = ""
		} else {
			v = fmt.Sprint(vs)
		}
	case error:
		if support.IsNil(vs) {
			v = ""
		} else {
			v = vs.Error()
		}
	case string, bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64, complex64, complex128,
		[]byte, []string, []interface{},
		map[string]string, map[string]interface{},
		level.Level:
		// Common values do not require reflection.
	default:
		value := reflect.ValueOf(v)
		if value.Kind() == reflect.Pointer && value.IsNil() {
			v = ""
		}
	}
	switch instance.QuoteType {
	case QuoteTypeMinimal:
		if vs, ok := v.(string); ok && !stringNeedsQuoting(vs) {
			return []byte(vs), nil
		}
		return json.Marshal(v)
	case QuoteTypeEverything:
		return json.Marshal(fmt.Sprint(v))
	default:
		return json.Marshal(v)
	}
}
