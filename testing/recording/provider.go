package recording

import (
	"sync/atomic"
	"unsafe"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

// DefaultProviderName specifies the default name an instance of Provider which
// will be used if no other name was defined.
const DefaultProviderName = "recording"

// DefaultLevel specifies the default level.Level of an instance of Provider
// which be used if no other level was defined.
const DefaultLevel = level.Info

// Provider is an implementation of log.Provider which is simply recording all
// events which are logged with it's loggers. They can be simply received using
// GetAll(), GetAllRoot() or GetAllOf()...
type Provider struct {
	// Name specifies the name of this Provider. If empty this Provider will
	// use DefaultProviderName.
	Name string

	// Level specifies the level of this Provider which will be also inherited
	// by all of its loggers. If 0 this Provider will use DefaultLevel.
	Level level.Level

	// AllLevels specifies the levels which are supported by this Provider and all
	// of its loggers. If nil this Provider will use
	// level.GetProvider()#GetLevels().
	AllLevels level.Levels

	// FieldKeysSpec specifies the spec of the fields are supported by this
	// Provider and all of its loggers. If nil this Provider will use the
	// default instance of fields.KeysSpecImpl.
	FieldKeysSpec fields.KeysSpec

	cachePointer unsafe.Pointer
}

// NewProvider creates a new instance of Provider which is ready to use.
func NewProvider() *Provider {
	return &Provider{}
}

// HookGlobally hooks itself into the global provider registry and forces this
// provider with all its loggers to be used by all libraries. It also returns
// a function pointer which should be executed when you're done. It will set
// everything back as it was before. Usually this should be done by using a
// defer statement.
//
// This is quite useful in tests where you want to record also test outputs
// and want to reset everything afterwards into a clean state.
func (instance *Provider) HookGlobally() func() {
	previous := log.SetProvider(instance)
	return func() {
		log.SetProvider(previous)
	}
}

// Contains checks if the given log.Event was recorded by at least of one
// CoreLogger of this Provider. It will use the
// fields.DefaultEntryEqualityFunction to checks the equality but will ignore
// the timestamp, log.Event#GetCallDepth() and log.Event#GetContext().
func (instance *Provider) Contains(expected log.Event) (bool, error) {
	return instance.ContainsCustom(instance.defaultEventEquality(), expected)
}

// MustContains is like Contains but will panic on errors instead returning
// them.
func (instance *Provider) MustContains(expected log.Event) bool {
	result, err := instance.Contains(expected)
	if err != nil {
		panic(err)
	}
	return result
}

// ContainsCustom checks if the given log.Event was recorded by at least of one
// CoreLogger of this Provider. It will use the given log.EventEquality to
// checks the equality.
func (instance *Provider) ContainsCustom(eef log.EventEquality, expected log.Event) (bool, error) {
	if found, err := instance.getRootLogger().ContainsCustom(eef, expected); err != nil {
		return false, err
	} else if found {
		return true, nil
	}

	for _, name := range instance.getCache().GetNames() {
		if found, err := instance.getLogger(name).ContainsCustom(eef, expected); err != nil {
			return false, err
		} else if found {
			return true, nil
		}
	}

	return false, nil
}

// MustContainsCustom is like ContainsCustom but will panic on errors instead
// returning them.
func (instance *Provider) MustContainsCustom(eef log.EventEquality, expected log.Event) bool {
	result, err := instance.ContainsCustom(eef, expected)
	if err != nil {
		panic(err)
	}
	return result
}

// Len the number of all events recorded by all of the CoreLoggers of this
// Provider.
func (instance *Provider) Len() int {
	result := instance.getRootLogger().Len()
	for _, name := range instance.getCache().GetNames() {
		result += instance.getLogger(name).Len()
	}
	return result
}

// GetAll returns all instances of log.Event which where recorded so far by all
// instances of CoreLogger which are associated to this instance of Provider.
func (instance *Provider) GetAll() []log.Event {
	result := instance.GetAllRoot()
	for _, name := range instance.getCache().GetNames() {
		result = append(result, instance.GetAllOf(name)...)
	}
	return result
}

// GetAllRoot returns all instances of log.Event which where recorded so far by
// the instance of root CoreLogger which is associated to this instance of
// Provider.
func (instance *Provider) GetAllRoot() []log.Event {
	return instance.getRootLogger().GetAll()
}

// GetAllOf returns all instances of log.Event which where recorded so far by
// the instance of named CoreLogger which are associated to this instance of
// Provider.
func (instance *Provider) GetAllOf(name string) []log.Event {
	return instance.getLogger(name).GetAll()
}

// ResetAll removes all recorded entries from all instances of CoreLogger which
// are associated to this instance of Provider.
func (instance *Provider) ResetAll() {
	instance.ResetRoot()
	for _, name := range instance.getCache().GetNames() {
		instance.Reset(name)
	}
}

// ResetRoot removes all recorded entries from the instance of root CoreLogger
// which is associated to this instance of Provider.
func (instance *Provider) ResetRoot() {
	instance.getRootLogger().Reset()
}

// Reset removes all recorded entries from the instance of named CoreLogger
// which is associated to this instance of Provider.
func (instance *Provider) Reset(name string) {
	instance.getLogger(name).Reset()
}

func (instance *Provider) getRootLogger() *Logger {
	return instance.GetRootLogger().(*Logger)
}

func (instance *Provider) getLogger(name string) *Logger {
	return instance.GetLogger(name).(*Logger)
}

func (instance *Provider) rootFactory() log.Logger {
	return instance.factory(RootLoggerName)
}

func (instance *Provider) factory(name string) log.Logger {
	result := NewLogger()
	result.Name = name
	result.Provider = instance
	return result
}

// GetRootLogger implements log.Provider#GetRootLogger()
func (instance *Provider) GetRootLogger() log.Logger {
	return instance.getCache().GetRootLogger()
}

// GetLogger implements log.Provider#GetLogger()
func (instance *Provider) GetLogger(name string) log.Logger {
	return instance.getCache().GetLogger(name)
}

// GetName implements log.Provider#GetName()
func (instance *Provider) GetName() string {
	if v := instance.Name; v != "" {
		return v
	}
	return DefaultProviderName
}

// GetAllLevels implements log.Provider#GetAllLevels()
func (instance *Provider) GetAllLevels() level.Levels {
	if v := instance.AllLevels; v != nil {
		return v
	}
	return level.GetProvider().GetLevels()
}

// GetFieldKeysSpec implements log.Provider#GetFieldKeysSpec()
func (instance *Provider) GetFieldKeysSpec() fields.KeysSpec {
	if v := instance.FieldKeysSpec; v != nil {
		return v
	}
	return &fields.KeysSpecImpl{}
}

// GetLevel returns the current level.Level where this log.Provider is set to.
func (instance *Provider) GetLevel() level.Level {
	if v := instance.Level; v != 0 {
		return v
	}
	return DefaultLevel
}

// SetLevel changes the current level.Level of this log.Provider. If set to
// 0 it will force this Provider to use DefaultLevel.
func (instance *Provider) SetLevel(v level.Level) {
	instance.Level = v
}

func (instance *Provider) getCache() log.LoggerCache {
	for {
		v := (*log.LoggerCache)(atomic.LoadPointer(&instance.cachePointer))
		if v != nil && *v != nil {
			return *v
		}

		c := log.NewLoggerCache(instance.rootFactory, instance.factory)

		if atomic.CompareAndSwapPointer(&instance.cachePointer, unsafe.Pointer(v), unsafe.Pointer(&c)) {
			return c
		}
	}
}

func (instance *Provider) defaultEventEquality() log.EventEquality {
	spec := instance.GetFieldKeysSpec()
	return log.DefaultEventEquality.WithIgnoringKeys(spec.GetTimestamp(), spec.GetLogger())
}
