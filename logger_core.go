package log

import "github.com/echocat/slf4g/level"

// CoreLogger defines the base functions of all Logger of the slf4g framework.
//
// This needs to be usually implemented by loggers that interfaces with the
// slf4g framework and can be elevated to a full instance of Logger by calling
// NewLogger().
type CoreLogger interface {
	// Log is called to log the given Event. It depends on the implementation
	// if this action will be synchronous or asynchronous. An Event handed over
	// to this (Core)Logger should usually be created by the NewEvent() method
	// of this instance. An implementation of CoreLogger is allowed to either
	// accept events bt foreigners or reject them with a panic. You can check
	// this using Accepts().
	//
	// skipFrames defines how many frame should be skipped to determine the real
	// caller of the log event from the call stack. In cse of delegating do not
	// forget to increase.
	Log(event Event, skipFrames uint16)

	// IsLevelEnabled returns <true> if the given Level is enabled to be logged
	// with this (Core)Logger.
	IsLevelEnabled(level.Level) bool

	// GetName returns the name of this (Core)Logger instance.
	GetName() string

	// NewEvent creates a new instance of an Event which can be modified before
	// provided back to this CoreLogger using Log().
	NewEvent(level level.Level, values map[string]interface{}) Event

	// Accepts is used to determine of this instance of a (Core)Logger might
	// accept also Events created by foreigners. Events creates by NewEvent()
	// of this instance should always return true.
	Accepts(Event) bool

	// GetProvider will return the Provider where this (Core)Logger belongs to.
	// This is for example used to access the AllLevels or fields.KeysSpec used
	// by this (Core)Logger.
	GetProvider() Provider
}
