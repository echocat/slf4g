package log

import (
	"fmt"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

// Event represents an event which can be logged using a Logger (or CoreLogger).
//
// # Contents
//
// Events containing always present content provided by GetLevel().
//
// They are providing additionally dynamic content (messages, errors, ...)
// which are accessible via ForEach() and Get(). None of these fields are
// required to exist by contract. The keys of these fields are defined by
// Provider.GetFieldKeysSpec(). For example using fields.KeysSpec.GetMessage()
// it might be possible to get the key of the message.
//
// The keys are always of type string and should be only printable characters
// which can be printed in any context. Recommended are everything that matches:
//
//	^[a-zA-Z0-9._-]+$
//
// The values could be everything including nils.
//
// # Immutability
//
// Fields are defined as immutable. Calling the methods With, Withf, WithAll
// and  Without always results in a new instance of Event that could be either
// brand new, a copy of the source one or do inherit some stuff of the original
// called one; but it never modifies the called instance.
type Event interface {
	// GetLevel returns the Level of this event.
	GetLevel() level.Level

	// ForEach will call the provided consumer for each field which is provided
	// by this Fields instance.
	ForEach(consumer func(key string, value interface{}) error) error

	// Get will return for the given key the corresponding value if exists.
	// Otherwise, it will return nil.
	Get(key string) (value interface{}, exists bool)

	// Len returns the len of all key value pairs contained in this event which
	// can be received by using ForEach() or Get().
	Len() int

	// With returns a variant of this Event with the given key
	// value pair contained inside. If the given key already exists in the
	// current instance this means it will be overwritten.
	With(key string, value interface{}) Event

	// Withf is similar to With, but it adds classic fmt.Printf functions to it.
	// It is defined that the format itself will not be executed before the
	// consumption of the value. (See fields.Fields.ForEach() and
	// fields.Fields.Get())
	Withf(key string, format string, args ...interface{}) Event

	// WithError is similar to With, but it adds an error as field.
	WithError(error) Event

	// WithAll is similar to With, but it can consume more than one field at
	// once. Be aware: There is neither a guarantee that this instance will be
	// copied or not.
	WithAll(map[string]interface{}) Event

	// Without returns a variant of this Event without the given
	// key contained inside. In other words: If someone afterwards tries to
	// call either fields.Fields.ForEach() or fields.Fields.Get() nothing with
	// this key(s) will be returned.
	Without(keys ...string) Event
}

// EventFactory creates a new instance of [Event] using the provided [level.Level]
// and values.
type EventFactory interface {
	NewEvent(l level.Level, values map[string]interface{}) Event
}

// EventFactoryWithFields creates a new instance of [Event] using the provided
// [level.Level] and fields.
type EventFactoryWithFields interface {
	NewEventWithFields(level.Level, fields.ForEachEnabled) Event
}

// NewEvent creates a new [Event] using the provided [EventFactory], [level.Level]
// and values.
func NewEvent(factory EventFactory, l level.Level, values map[string]interface{}) Event {
	return factory.NewEvent(l, values)
}

// NewEventWithFields creates a new [Event] using the provided [EventFactory], [level.Level]
// and fields.
func NewEventWithFields(factory EventFactory, l level.Level, f fields.ForEachEnabled) Event {
	if wf, ok := factory.(EventFactoryWithFields); ok {
		return wf.NewEventWithFields(l, f)
	}
	asMap, err := fields.AsMap(f)
	if err != nil {
		panic(fmt.Errorf("cannot make a map out of %v: %w", f, err))
	}
	return factory.NewEvent(l, asMap)
}
