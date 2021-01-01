package formatter

import (
	log "github.com/echocat/slf4g"
)

// DefaultTextValue is the default instance of TextValue
// which should cover the most of the cases.
var DefaultTextValue TextValue = NewSimpleTextValue()

// TextValue formats a given value to be printed in the text.
type TextValue interface {
	// FormatValue formats the given value to a readable format.
	FormatTextValue(value interface{}, provider log.Provider) ([]byte, error)
}

// TextValueFunc is wrapping the given function into a
// TextValue.
type TextValueFunc func(interface{}, log.Provider) ([]byte, error)

// FormatTextValue implements TextValue.FormatTextValue().
func (instance TextValueFunc) FormatTextValue(value interface{}, provider log.Provider) ([]byte, error) {
	return instance(value, provider)
}

// NewTextValueFacade creates a new facade instance of
// TextValue using the given provider.
func NewTextValueFacade(provider func() TextValue) TextValue {
	return textValueFacade(provider)
}

type textValueFacade func() TextValue

func (instance textValueFacade) FormatTextValue(value interface{}, provider log.Provider) ([]byte, error) {
	return instance().FormatTextValue(value, provider)
}

// NoopTextValue provides a noop implementation of TextValue.
func NoopTextValue() TextValue {
	return noopTextValueV
}

var noopTextValueV = TextValueFunc(func(interface{}, log.Provider) ([]byte, error) {
	return []byte{}, nil
})
