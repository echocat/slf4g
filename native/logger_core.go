package native

import (
	"fmt"
	"reflect"
	"time"

	"github.com/echocat/slf4g/fields"

	"github.com/echocat/slf4g/level"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native/consumer"
	"github.com/echocat/slf4g/native/location"
)

const (
	rootLoggerName = "ROOT"
)

// CoreLogger implements log.CoreLogger of the slf4g framework for the "native"
// implementation.
//
// You cannot create a working instance of this by yourself. It can only be done
// by the Provider instance. If you want to customize it you can use
// Provider.CoreLoggerCustomizer to done this.
type CoreLogger struct {
	Level             level.Level
	Consumer          consumer.Consumer
	LocationDiscovery location.Discovery

	provider *Provider
	name     string
}

// Log implements log.CoreLogger#Log()
func (instance *CoreLogger) Log(event log.Event, skipFrames uint16) {
	if event == nil {
		return
	}
	if !instance.IsLevelEnabled(event.GetLevel()) {
		return
	}
	provider := instance.getProvider()
	fieldKeysSpec := provider.getFieldKeysSpec()

	if v := log.GetTimestampOf(event, provider); v == nil {
		event = event.With(fieldKeysSpec.GetTimestamp(), time.Now())
	}
	if v := log.GetLoggerOf(event, provider); v == nil || *v != instance.name {
		event = event.With(fieldKeysSpec.GetLogger(), instance.name)
	}
	if v := instance.getLocationDiscovery().DiscoverLocation(event, skipFrames+1); v != nil {
		event = event.With(fieldKeysSpec.GetLocation(), v)
	}

	instance.getConsumer().Consume(event, instance)
}

// IsLevelEnabled implements log.CoreLogger#IsLevelEnabled()
func (instance *CoreLogger) IsLevelEnabled(level level.Level) bool {
	return instance.GetLevel().CompareTo(level) <= 0
}

// SetLevel changes the current level.Level of this log.CoreLogger. If set to
// 0 it use the value of Provider.GetLevel().
func (instance *CoreLogger) SetLevel(level level.Level) {
	instance.Level = level
}

// GetLevel returns the current level.Level where this log.CoreLogger is set to.
func (instance *CoreLogger) GetLevel() level.Level {
	if v := instance.Level; v != 0 {
		return v
	}
	return instance.getProvider().GetLevel()
}

// GetName implements log.CoreLogger#GetName()
func (instance *CoreLogger) GetName() string {
	if v := instance.name; v != "" {
		return v
	}
	panic(fmt.Sprintf("This %v was not initiated by a %v.", reflect.TypeOf(*instance), reflect.TypeOf(Provider{})))
}

// GetProvider implements log.CoreLogger#GetProvider()
func (instance *CoreLogger) GetProvider() log.Provider {
	return instance.getProvider()
}

// NewEvent implements log.CoreLogger#NewEvent()
func (instance *CoreLogger) NewEvent(l level.Level, values map[string]interface{}) log.Event {
	return instance.NewEventWithFields(l, fields.WithAll(values))
}

// NewEventWithFields provides a shortcut if an event should directly create
// from fields.
func (instance *CoreLogger) NewEventWithFields(l level.Level, f fields.ForEachEnabled) log.Event {
	asFields, err := fields.AsFields(f)
	if err != nil {
		panic(err)
	}
	return &event{
		provider: instance.provider,
		fields:   asFields,
		level:    l,
	}
}

// Accepts implements log.CoreLogger#Accepts()
func (instance *CoreLogger) Accepts(e log.Event) bool {
	return e != nil
}

func (instance *CoreLogger) getConsumer() consumer.Consumer {
	if c := instance.Consumer; c != nil {
		return c
	}
	return instance.getProvider().GetConsumer()
}

func (instance *CoreLogger) getLocationDiscovery() location.Discovery {
	if f := instance.LocationDiscovery; f != nil {
		return f
	}
	return instance.getProvider().getLocationDiscovery()
}

func (instance *CoreLogger) getProvider() *Provider {
	if v := instance.provider; v != nil {
		return v
	}
	panic(fmt.Sprintf("This %v was not initiated by a %v.", reflect.TypeOf(*instance), reflect.TypeOf(Provider{})))
}
