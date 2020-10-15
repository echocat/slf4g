package formatter

import (
	"encoding/json"
	"fmt"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
)

// SimpleTextValueFormatter is a simple implementation of
// TextValueFormatter.
type SimpleTextValueFormatter struct {
	// QuoteType defines how values are quoted.
	QuoteType QuoteType
}

// NewSimpleTextValueFormatter creates a new instance of
// SimpleTextValueFormatter which is ready to use.
func NewSimpleTextValueFormatter(customizer ...func(*SimpleTextValueFormatter)) *SimpleTextValueFormatter {
	result := &SimpleTextValueFormatter{
		QuoteType: QuoteTypeMinimal,
	}
	for _, c := range customizer {
		c(result)
	}
	return result
}

// FormatValue implements TextValueFormatter.FormatValue().
func (instance *SimpleTextValueFormatter) FormatValue(v interface{}, _ log.Provider) ([]byte, error) {
	if vl, ok := v.(fields.Lazy); ok {
		v = vl.Get()
	}

	if instance.QuoteType == QuoteTypeMinimal {
		switch vs := v.(type) {
		case string:
			if stringNeedsQuoting(vs) {
				return json.Marshal(vs)
			}
			return []byte(vs), nil
		case *string:
			if stringNeedsQuoting(*vs) {
				return json.Marshal(*vs)
			}
			return []byte(*vs), nil
		case fmt.Stringer:
			str := vs.String()
			if stringNeedsQuoting(str) {
				return json.Marshal(str)
			}
			return []byte(str), nil
		case fmt.Formatter:
			str := fmt.Sprint(vs)
			if stringNeedsQuoting(str) {
				return json.Marshal(str)
			}
			return []byte(str), nil
		case error:
			str := vs.Error()
			if stringNeedsQuoting(str) {
				return json.Marshal(str)
			}
			return []byte(str), nil
		case *error:
			str := (*vs).Error()
			if stringNeedsQuoting(str) {
				return json.Marshal(str)
			}
			return []byte(str), nil
		}
		return json.Marshal(v)
	}

	if instance.QuoteType == QuoteTypeEverything {
		return json.Marshal(fmt.Sprint(v))
	}

	return json.Marshal(v)
}
