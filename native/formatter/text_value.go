package formatter

import (
	log "github.com/echocat/slf4g"
)

// DefaultTextValueFormatter is the default instance of TextValueFormatter
// which should cover the most of the cases.
var DefaultTextValueFormatter TextValueFormatter = NewSimpleTextValueFormatter()

// TextValueFormatter formats a given value to be printed in the text.
type TextValueFormatter interface {
	// FormatValue formats the given value to a readable format.
	FormatValue(value interface{}, provider log.Provider) ([]byte, error)
}

// TextValueFormatterFunc is wrapping the given function into a
// TextValueFormatter.
type TextValueFormatterFunc func(interface{}, log.Provider) ([]byte, error)

// FormatValue implements TextValueFormatter.FormatValue().
func (instance TextValueFormatterFunc) FormatValue(value interface{}, provider log.Provider) ([]byte, error) {
	return instance(value, provider)
}

// NewTextValueFormatterFacade creates a new facade instance of
// TextValueFormatter using the given provider.
func NewTextValueFormatterFacade(provider func() TextValueFormatter) TextValueFormatter {
	return valueFacade(provider)
}

type valueFacade func() TextValueFormatter

func (instance valueFacade) FormatValue(value interface{}, provider log.Provider) ([]byte, error) {
	return instance().FormatValue(value, provider)
}

// NoopTextValueFormatter provides a noop implementation of TextValueFormatter.
func NoopTextValueFormatter() TextValueFormatter {
	return noopTextValueFormatterV
}

var noopTextValueFormatterV = TextValueFormatterFunc(func(interface{}, log.Provider) ([]byte, error) {
	return []byte{}, nil
})
