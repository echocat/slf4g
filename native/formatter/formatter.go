// Formatter is used to format log events to a format which can logged to a
// console, file, ...
package formatter

import (
	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native/hints"
)

// Default is the default instance of Formatter which should cover the most of
// the cases.
var Default Formatter = NewText()

// Formatter is used to format log events to a format which can logged to a
// console, file, ...
type Formatter interface {
	// Format formats the given event to a format which  can logged to a
	// console, file, ...
	Format(log.Event, log.Provider, hints.Hints) ([]byte, error)
}

// Func is wrapping the given function into a Formatter.
type Func func(log.Event, log.Provider, hints.Hints) ([]byte, error)

// Format implements Formatter.Format()
func (instance Func) Format(event log.Event, provider log.Provider, h hints.Hints) ([]byte, error) {
	return instance(event, provider, h)
}

// NewFacade creates a new facade instance of Formatter using the given
// provider.
func NewFacade(provider func() Formatter) Formatter {
	return facade(provider)
}

type facade func() Formatter

func (instance facade) Format(event log.Event, provider log.Provider, h hints.Hints) ([]byte, error) {
	return instance().Format(event, provider, h)
}

// Noop provides a noop implementation of Formatter.
func Noop() Formatter {
	return noopV
}

var noopV = Func(func(log.Event, log.Provider, hints.Hints) ([]byte, error) {
	return []byte{}, nil
})
