package value

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/echocat/slf4g/native/formatter"
)

// DefaultFormatterCodec is the default instance of FormatterCodec which should cover the
// most of the cases.
var DefaultFormatterCodec = MappingFormatterCodec{
	"text": func() (formatter.Formatter, error) {
		return formatter.NewText(), nil
	},
	"json": func() (formatter.Formatter, error) {
		return formatter.NewJson(), nil
	},
}

// FormatterCodec transforms strings to formatter.Formatter and other way around.
type FormatterCodec interface {
	// Parse takes a string and creates out of it an instance of formatter.Formatter.
	Parse(plain string) (formatter.Formatter, error)

	// Formatter takes an instance of formatter.Formatter and formats it as a string.
	Format(what formatter.Formatter) (string, error)
}

// MappingFormatterCodec is a default implementation of FormatterCodec which handles
// the most common cases by default.
type MappingFormatterCodec map[string]func() (formatter.Formatter, error)

// Parse implements FormatterCodec.Parse
func (instance MappingFormatterCodec) Parse(plain string) (formatter.Formatter, error) {
	if plain == "" {
		if d := formatter.Default; d != nil {
			return d, nil
		}
	}
	for n, f := range instance {
		if strings.EqualFold(n, plain) {
			return f()
		}
	}
	return nil, fmt.Errorf("unknown log format: %s", plain)
}

// Format implements FormatterCodec.Format
func (instance MappingFormatterCodec) Format(what formatter.Formatter) (string, error) {
	for n, f := range instance {
		candidate, err := f()
		if err == nil && reflect.DeepEqual(what, candidate) {
			return n, nil
		}
	}
	return "", fmt.Errorf("unknown log formatter: %v", reflect.TypeOf(what))
}

type noopFormatterCodec struct{}

func (instance *noopFormatterCodec) Parse(plain string) (formatter.Formatter, error) {
	return nil, fmt.Errorf("unknown log format: %s", plain)
}

func (instance *noopFormatterCodec) Format(what formatter.Formatter) (string, error) {
	return "", fmt.Errorf("unknown log formatter: %v", reflect.TypeOf(what))
}

var noopFormatterCodecV = &noopFormatterCodec{}

// NoopFormatterCodec provides a noop implementation of FormatterCodec.
func NoopFormatterCodec() FormatterCodec {
	return noopFormatterCodecV
}

// NewFormatterCodecFacade creates a new facade of FormatterCodec with the given function
// that provides the actual FormatterCodec to use.
func NewFormatterCodecFacade(provider func() FormatterCodec) FormatterCodec {
	return formatterCodecFacade(provider)
}

type formatterCodecFacade func() FormatterCodec

func (instance formatterCodecFacade) Parse(plain string) (formatter.Formatter, error) {
	return instance().Parse(plain)
}

func (instance formatterCodecFacade) Format(what formatter.Formatter) (string, error) {
	return instance().Format(what)
}
