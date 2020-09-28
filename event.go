package log

import "github.com/echocat/slf4g/fields"

// Event represents an event which can be logged using a Logger (or CoreLogger).
//
// Contents
//
// Events are haven dynamic contents (such as messages, errors, ...) provided
// by GetFields(). Only be contract always provided information are provided by
// GetLevel() and GetCalDepth().
//
// Immutability
//
// Fields are defined as immutable. Calling the methods With, Withf, WithAll,
// Without, WithCallDepth and WithContext always results in a new instance of
// Event that could be either brand new, a copy of the source one or do inherit
// some stuff of the original called one; but it never modifies the called
// instance.
type Event interface {
	// GetLevel returns the Level of this event.
	GetLevel() Level

	// GetCallDepth returns the call depth inside of the call stack that should
	// be ignored before capturing the caller position (if required). This could
	// be increased (if delegating from instance to instance) by calling
	// WithCallDepth().
	GetCallDepth() int

	// GetContext returns an optional context of this event. This is stuff which
	// should not be represented and/or reported and/or could contain hints for
	// the actual logger. Therefore and can be <nil>. This can altered by
	// WithContext().
	GetContext() interface{}

	// GetFields will return all fields which are associated with this Event.
	// This could contain a message, timestamp, error and so on. None of this
	// fields is required to exists by contract. The keys of these fields is
	// defined by Provider.GetFieldKeysSpec(). For example using
	// fields.KeysSpec.GetMessage() it might be possible to get the
	// key of the message.
	GetFields() fields.Fields

	// With returns an variant of this Event with the given key
	// value pair contained inside. If the given key already exists in the
	// current instance this means it will be overwritten.
	With(key string, value interface{}) Event

	// Withf is similar to With but it adds classic fmt.Printf functions to it.
	// It is defined that the format itself will not be executed before the
	// consumption of the value. (See fields.Fields.ForEach() and
	// fields.Fields.Get())
	Withf(key string, format string, args ...interface{}) Event

	// WithAll is similar to With but it can consume more than one field at
	// once. Be aware: There is neither a guarantee that this instance will be
	// copied or not.
	WithAll(map[string]interface{}) Event

	// Without returns a variant of this Event without the given
	// key contained inside. In other words: If someone afterwards tries to
	// call either fields.Fields.ForEach() or fields.Fields.Get() nothing with
	// this key(s) will be returned.
	Without(keys ...string) Event

	// WithCallDepth returns an variant of this Event with the given
	// call depth is added to the existing one of this Event. All other values
	// remaining the same.
	WithCallDepth(int) Event

	// WithContext returns an variant of this Event with the given
	// context is replaced with the existing one of this Event. All other values
	// remaining the same.
	WithContext(ctx interface{}) Event
}

// NewEvent creates a new instance of Event from the given Level, fields.Fields
// and the given callDepth.
func NewEvent(level Level, f fields.Fields, callDepth int) Event {
	if f == nil {
		f = fields.Empty()
	}
	return &eventImpl{
		Level:     level,
		Fields:    f,
		CallDepth: callDepth,
	}
}
