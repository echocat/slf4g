package log

import "github.com/echocat/slf4g/level"

// CoreLogger defines the base functions of all Logger of the slf4g framework.
//
// This needs to be usually implemented by loggers that interfaces with the
// slf4g framework and can be elevated to a full instance of Logger by calling
// NewLogger().
type CoreLogger interface {
	// Log is called to log the given Event. It depends on the implementation
	// if this action will be synchronous or asynchronous.
	Log(Event)

	// IsLevelEnabled returns <true> if the given Level is enabled to be logged
	// with this (Core)Logger.
	IsLevelEnabled(level.Level) bool

	// GetName returns the name of this (Core)Logger instance.
	GetName() string

	// GetProvider will return the Provider where this (Core)Logger belongs to.
	// This is for example used to access the AllLevels or fields.KeysSpec used
	// by this (Core)Logger.
	GetProvider() Provider
}
