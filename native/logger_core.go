package native

import (
	"fmt"
	"reflect"
	"time"

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
	Level             *level.Level
	Consumer          consumer.Consumer
	LocationDiscovery location.Discovery

	provider *Provider
	name     string
}

// Log implements log.CoreLogger#Log()
func (instance *CoreLogger) Log(event log.Event) {
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
	if v := log.GetLoggerOf(event, provider); v == nil {
		event = event.With(fieldKeysSpec.GetLogger(), instance.name)
	}
	if v := instance.getLocationDiscovery().DiscoveryLocation(event, event.GetCallDepth()+1); v != nil {
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
	instance.Level = &level
}

// GetLevel returns the current level.Level where this log.CoreLogger is set to.
func (instance *CoreLogger) GetLevel() level.Level {
	if v := instance.Level; v != nil {
		return *v
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

func (instance *CoreLogger) getConsumer() consumer.Consumer {
	if c := instance.Consumer; c != nil {
		return c
	}
	return instance.getProvider().getConsumer()
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
