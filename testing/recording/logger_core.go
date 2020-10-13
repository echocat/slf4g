package recording

import (
	"fmt"
	"sync"
	"time"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/level"
)

// RootLoggerName specifies the name of the root version of CoreLogger
// instances which are managed by Provider.
const RootLoggerName = "ROOT"

// CoreLogger implements log.CoreLogger and records simply every logged event.
// Each of these events can be received using GetAll() or Get(index).
//
// Mutation
//
// This log.CoreLogger will mutate the recorded instances of log.Event if
// one of the following fields is absent: logger, timestamp. This behavior can
// be disabled using SetIfAbsent field.
//
// Mutation will also happen if you use this logger as a global registered one
// by calling methods like log.Info(..), log.Warn(..), ... because they are
// adding in any way their context information.
type CoreLogger struct {
	// Provider which are managing this instance. If this value is nil
	// log.GetProvider() will be used instead. This default behaviour could lead
	// in some cases to unintended consequences.
	Provider log.Provider

	// Name represents the name of this instance of CoreLogger. If this value is
	// empty RootLoggerName will be used instead.
	Name string

	// Level represents the level of this instance of CoreLogger. If this value
	// is empty the level of its Provider will be used. If it is not possible
	// to receive the level of its Provider it will use DefaultLevel instead.
	Level level.Level

	recorded []log.Event
	mutex    sync.RWMutex
}

// NewCoreLogger creates a new instance of CoreLogger which is ready to use.
func NewCoreLogger() *CoreLogger {
	return &CoreLogger{}
}

// Contains checks if the given log.Event was recorded by this CoreLogger.
// It will use the fields.DefaultEntryEqualityFunction to checks the equality
// but will ignore the timestamp, log.Event#GetCallDepth() and
// log.Event#GetContext().
func (instance *CoreLogger) Contains(expected log.Event) (bool, error) {
	return instance.ContainsCustom(instance.defaultEventEquality(), expected)
}

// MustContains is like Contains but will panic on errors instead returning
// them.
func (instance *CoreLogger) MustContains(expected log.Event) bool {
	result, err := instance.Contains(expected)
	if err != nil {
		panic(err)
	}
	return result
}

// ContainsCustom checks if the given log.Event was recorded by this CoreLogger.
// It will use the given log.EventEquality to checks the equality.
func (instance *CoreLogger) ContainsCustom(eef log.EventEquality, expected log.Event) (bool, error) {
	instance.mutex.RLock()
	defer instance.mutex.RUnlock()

	for _, candidate := range instance.recorded {
		if matches, err := eef.AreEventsEqual(expected, candidate); err != nil {
			return false, err
		} else if matches {
			return true, nil
		}
	}

	return false, nil
}

// MustContainsCustom is like ContainsCustom but will panic on errors instead
// returning them.
func (instance *CoreLogger) MustContainsCustom(eef log.EventEquality, expected log.Event) bool {
	result, err := instance.ContainsCustom(eef, expected)
	if err != nil {
		panic(err)
	}
	return result
}

// Len returns the length of all recorded events.
func (instance *CoreLogger) Len() int {
	instance.mutex.RLock()
	defer instance.mutex.RUnlock()

	return len(instance.recorded)
}

// GetAll returns all recorded events.
func (instance *CoreLogger) GetAll() []log.Event {
	instance.mutex.RLock()
	defer instance.mutex.RUnlock()

	result := make([]log.Event, len(instance.recorded))
	copy(result, instance.recorded)

	return result
}

// Get return an event at the given index. If this index does not exists this
// method will panic.
func (instance *CoreLogger) Get(index uint) log.Event {
	instance.mutex.RLock()
	defer instance.mutex.RUnlock()

	if uint(len(instance.recorded)) <= index {
		panic(fmt.Sprintf("Index %d requested but the amount of recorded events is only %d", index, len(instance.recorded)))
	}

	return instance.recorded[index]
}

// Reset will remove all recorded events of this CoreLogger.
func (instance *CoreLogger) Reset() {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()

	instance.recorded = []log.Event{}
}

// Log implements log.CoreLogger#Log(event).
func (instance *CoreLogger) Log(event log.Event) {
	if !instance.IsLevelEnabled(event.GetLevel()) {
		return
	}

	instance.mutex.Lock()
	defer instance.mutex.Unlock()

	provider := instance.GetProvider()
	if v := log.GetTimestampOf(event, provider); v == nil {
		event = event.With(provider.GetFieldKeysSpec().GetTimestamp(), time.Now())
	}
	if v := log.GetLoggerOf(event, provider); v == nil {
		event = event.With(provider.GetFieldKeysSpec().GetLogger(), instance)
	}

	instance.recorded = append(instance.recorded, event)
}

// GetLevel returns the current level.Level where this log.CoreLogger is set to.
func (instance *CoreLogger) GetLevel() level.Level {
	if v := instance.Level; v != 0 {
		return v
	}
	if la, ok := instance.GetProvider().(level.Aware); ok {
		return la.GetLevel()
	}
	return DefaultLevel
}

// SetLevel changes the current level.Level of this log.CoreLogger. If set to
// 0 it will force this CoreLogger to use DefaultLevel.
func (instance *CoreLogger) SetLevel(v level.Level) {
	instance.Level = v
}

// IsLevelEnabled implements log.CoreLogger#IsLevelEnabled()
func (instance *CoreLogger) IsLevelEnabled(v level.Level) bool {
	return instance.GetLevel().CompareTo(v) <= 0
}

// GetName implements log.CoreLogger#GetName()
func (instance *CoreLogger) GetName() string {
	if v := instance.Name; v != "" {
		return v
	}
	return RootLoggerName
}

// GetProvider implements log.CoreLogger#GetProvider()
func (instance *CoreLogger) GetProvider() log.Provider {
	if v := instance.Provider; v != nil {
		return v
	}
	return log.GetProvider()
}

func (instance *CoreLogger) defaultEventEquality() log.EventEquality {
	spec := instance.GetProvider().GetFieldKeysSpec()
	return log.DefaultEventEquality.WithIgnoringKeys(spec.GetTimestamp(), spec.GetLogger())
}
