package formatter

import (
	"encoding/json"
	"fmt"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
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
	if vl, ok := v.(fields.Lazy); ok {
		v = vl.Get()
	}

	switch vs := v.(type) {
	case *string:
		v = *vs
	case fmt.Stringer:
		v = vs.String()
	case fmt.Formatter:
		v = fmt.Sprint(vs)
	case error:
		v = vs.Error()
	}

	if v == nil {
		v = ""
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
