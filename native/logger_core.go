package native

import (
	"time"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native/consumer"
	"github.com/echocat/slf4g/native/location"
)

type CoreLogger struct {
	Level           *log.Level
	Consumer        consumer.Consumer
	LocationFactory location.Factory

	provider *Provider
	name     string
}

func (instance *CoreLogger) Log(event log.Event) {
	if event == nil {
		return
	}
	if !instance.IsLevelEnabled(event.GetLevel()) {
		return
	}

	if v := log.GetTimestampOf(event, instance.provider); v == nil {
		event = event.With(instance.provider.GetFieldKeysSpec().GetTimestamp(), time.Now())
	}
	if v := log.GetLoggerOf(event, instance.provider); v == nil {
		event = event.With(instance.provider.GetFieldKeysSpec().GetLogger(), instance.name)
	}

	if v := instance.getLocationFactory()(event, event.GetCallDepth()+1); v != nil {
		event = event.With(instance.provider.FieldsKeysSpec.GetLocation(), v)
	}

	instance.getConsumer().Consume(event, instance)
}

func (instance *CoreLogger) IsLevelEnabled(level log.Level) bool {
	return instance.GetLevel().CompareTo(level) <= 0
}

func (instance *CoreLogger) SetLevel(level log.Level) {
	instance.Level = &level
}

func (instance *CoreLogger) GetLevel() log.Level {
	if v := instance.Level; v != nil {
		return *v
	}
	return instance.provider.GetLevel()
}

func (instance *CoreLogger) GetName() string {
	return instance.name
}

func (instance *CoreLogger) GetProvider() log.Provider {
	return instance.provider
}

func (instance *CoreLogger) GetConsumer() consumer.Consumer {
	return instance.Consumer
}

func (instance *CoreLogger) SetConsumer(v consumer.Consumer) {
	instance.Consumer = v
}

func (instance *CoreLogger) getConsumer() consumer.Consumer {
	if c := instance.GetConsumer(); c != nil {
		return c
	}
	return instance.provider.getConsumer()
}

func (instance *CoreLogger) GetLocationFactory() location.Factory {
	return instance.LocationFactory
}

func (instance *CoreLogger) SetLocationFactory(v location.Factory) {
	instance.LocationFactory = v
}

func (instance *CoreLogger) getLocationFactory() location.Factory {
	if f := instance.GetLocationFactory(); f != nil {
		return f
	}
	return instance.provider.getLocationFactory()
}
