package testlog

import (
	"testing"

	log "github.com/echocat/slf4g"
)

// NewLogger creates a new instance of log.Logger ready to use. If you want
// to use a direct instance of a logger, this is the easiest way to get it.
//
// This is a shortcut for NewProvider(..).GetRootLogger().
func NewLogger(tb testing.TB, customizer ...func(*Provider)) log.Logger {
	provider := NewProvider(tb, customizer...)
	return provider.GetRootLogger()
}

// NewNamedLogger creates a new instance of log.Logger ready to use. If you want
// to use a direct instance of a logger with a specific name, this is the
// easiest way to get it.
//
// This is a shortcut for NewProvider(..).GetLogger(...).
func NewNamedLogger(tb testing.TB, name string, customizer ...func(*Provider)) log.Logger {
	provider := NewProvider(tb, customizer...)
	return provider.GetLogger(name)
}
