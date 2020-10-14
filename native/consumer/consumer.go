// Package consumer provides the functionally to print log events either to
// console, files, ...
package consumer

import (
	"os"

	log "github.com/echocat/slf4g"
)

// Default is the default instance of Consumer which should cover the majority
// of all cases.
var Default = NewWriter(os.Stderr)

// Consumer consumes instances of log.Event of a log.CoreLogger and for example
// print them to the console, to files, ...
type Consumer interface {
	// Consume consumes the event.
	Consume(event log.Event, source log.CoreLogger)
}

// Func is wrapping a given function into an instance of Consumer.
type Func func(event log.Event, source log.CoreLogger)

// Consume implements Consumer.Consume()
func (instance Func) Consume(event log.Event, source log.CoreLogger) {
	instance(event, source)
}

var noopV = Func(func(log.Event, log.CoreLogger) {})

// Noop provides a noop implementation of Consumer.
func Noop() Consumer {
	return noopV
}

// NewFacade creates a new facade of Consumer with the given function that
// provides the actual Consumer to use.
func NewFacade(provider func() Consumer) Consumer {
	return facade(provider)
}

type facade func() Consumer

func (instance facade) Consume(event log.Event, source log.CoreLogger) {
	instance().Consume(event, source)
}
