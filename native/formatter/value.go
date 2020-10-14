package formatter

import (
	log "github.com/echocat/slf4g"
)

var (
	DefaultValueFormatter ValueFormatter = &SimpleValueFormatter{}
)

// ValueFormatter formats a given value
type ValueFormatter interface {
	FormatValue(interface{}, log.Provider) ([]byte, error)
}

// ValueFormatterFunc is wrapping the given function into a ValueFormatter.
type ValueFormatterFunc func(interface{}, log.Provider) ([]byte, error)

// FormatValue implements ValueFormatter.FormatValue().
func (instance ValueFormatterFunc) FormatValue(value interface{}, provider log.Provider) ([]byte, error) {
	return instance(value, provider)
}

// NewValueFormatterFacade creates a new facade instance of ValueFormatter using
// the given provider.
func NewValueFormatterFacade(provider func() ValueFormatter) ValueFormatter {
	return valueFacade(provider)
}

type valueFacade func() ValueFormatter

func (instance valueFacade) FormatValue(value interface{}, provider log.Provider) ([]byte, error) {
	return instance().FormatValue(value, provider)
}
