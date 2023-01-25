package native

import (
	"sync/atomic"
	"unsafe"

	"github.com/echocat/slf4g/level"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/native/consumer"
	nlevel "github.com/echocat/slf4g/native/level"
	"github.com/echocat/slf4g/native/location"
)

var DefaultProvider = &Provider{}

// Provider implements log.Provider of the slf4g framework for the "native"
// implementation.
//
// Usually you should not be required to create by your self. Either use simply
// log.GetProvider() (which will return this provider once you imported this
// package at least one time somewhere) or if you want to customize its behavior
// simply modify DefaultProvider.
type Provider struct {
	// Name represents the name of this Provider. If empty it will be "native"
	// by default.
	Name string

	// Level represents the level.Level of this Provider that is at least
	// required that the loggers managed by this Provider will respect logged
	// events. This can be overwritten by individual loggers. If this value is
	// not set it will be log.Info by default.
	Level level.Level

	// LevelNames is used to format the levels as human-readable
	// representations. If this is not set it will be level.DefaultNames by
	// default.
	LevelNames level.Names

	// LevelProvider is used to determine the log.Levels support by this
	// Provider and all of its managed loggers. If this is not set it will be
	// level.GetProvider() by default.
	LevelProvider level.Provider

	// Consumer is used to handle the logged events with. If this is not set it
	// will be consumer.Default by default.
	Consumer consumer.Consumer

	// LocationDiscovery is used to discover the location.Location where events
	// are happen. If this is not set it will be location.DefaultDiscovery by
	// default.
	LocationDiscovery location.Discovery

	// FieldKeysSpec defines what are the keys of the major fields managed by
	// this Provider and its managed loggers. If this is not set it will be
	// DefaultFieldKeysSpec by default.
	FieldKeysSpec FieldKeysSpec

	// CoreLoggerCustomizer will be called in every moment a logger instance
	// needs to be created (if configured).
	CoreLoggerCustomizer CoreLoggerCustomizer

	cachePointer unsafe.Pointer
}

// GetName implements log.Provider#GetName()
func (instance *Provider) GetName() string {
	if v := instance.Name; v != "" {
		return v
	}
	return "native"
}

// GetRootLogger implements log.Provider#GetRootLogger()
func (instance *Provider) GetRootLogger() log.Logger {
	return instance.getCache().GetRootLogger()
}

// GetLogger implements log.Provider#GetLogger()
func (instance *Provider) GetLogger(name string) log.Logger {
	return instance.getCache().GetLogger(name)
}

// SetLevel changes the current level.Level of this log.Provider. If set to
// 0 it will force this Provider to use log.Info.
func (instance *Provider) SetLevel(v level.Level) {
	instance.Level = v
}

// GetLevel returns the current level.Level where this log.Provider is set to.
func (instance *Provider) GetLevel() level.Level {
	if v := instance.Level; v != 0 {
		return v
	}
	return level.Info
}

// SetConsumer changes the current consumer.Consumer of this log.Provider. If set
// to nil consumer.Default will be used.
func (instance *Provider) SetConsumer(v consumer.Consumer) {
	instance.Consumer = v
}

// GetConsumer returns the current consumer.Consumer where this log.Provider is set to.
func (instance *Provider) GetConsumer() consumer.Consumer {
	if v := instance.Consumer; v != nil {
		return v
	}
	if v := consumer.Default; v != nil {
		return v
	}
	return consumer.Noop()
}

// GetFieldKeysSpec implements log.Provider#GetFieldKeysSpec()
func (instance *Provider) GetFieldKeysSpec() fields.KeysSpec {
	return instance.getFieldKeysSpec()
}

// GetLevelNames returns an instance of level.Names that support by formatting
// level.Level managed by this Provider.
func (instance *Provider) GetLevelNames() level.Names {
	if v := instance.LevelNames; v != nil {
		return v
	}
	if v := nlevel.DefaultNames; v != nil {
		return v
	}
	return nlevel.NewNames()
}

// GetAllLevels implements log.Provider#GetAllLevels()
func (instance *Provider) GetAllLevels() level.Levels {
	p := instance.LevelProvider
	if p == nil {
		p = level.GetProvider()
	}
	return p.GetLevels()
}

func (instance *Provider) getFieldKeysSpec() FieldKeysSpec {
	if v := instance.FieldKeysSpec; v != nil {
		return v
	}
	if v := DefaultFieldKeysSpec; v != nil {
		return v
	}
	return &FieldKeysSpecImpl{}
}

func (instance *Provider) rootFactory() log.Logger {
	return instance.factory(rootLoggerName)
}

func (instance *Provider) factory(name string) log.Logger {
	cl := &CoreLogger{
		provider: instance,
		name:     name,
	}
	if c := instance.CoreLoggerCustomizer; c != nil {
		return log.NewLogger(c(instance, cl))
	}
	return log.NewLogger(cl)
}

func (instance *Provider) getLocationDiscovery() location.Discovery {
	if v := instance.LocationDiscovery; v != nil {
		return v
	}
	if v := location.DefaultDiscovery; v != nil {
		return v
	}
	return location.NoopDiscovery()
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

// CoreLoggerCustomizer can be used by the Provider to customize created
// instances of CoreLogger. See Provider.CoreLoggerCustomizer
type CoreLoggerCustomizer func(*Provider, *CoreLogger) log.CoreLogger

func init() {
	log.RegisterProvider(DefaultProvider)
}
