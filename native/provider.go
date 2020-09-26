package native

import (
	"fmt"
	"os"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/native/consumer"
	"github.com/echocat/slf4g/native/formatter"
	"github.com/echocat/slf4g/native/interceptor"
	"github.com/echocat/slf4g/native/level"
	"github.com/echocat/slf4g/native/location"
)

var DefaultProvider = NewProvider("native")
var _ = log.RegisterProvider(DefaultProvider)

type Provider struct {
	Cache log.LoggerCache
	Name  string

	Level         log.Level
	LevelNames    level.Names
	LevelProvider log.LevelProvider

	Formatter       formatter.Formatter
	Interceptor     interceptor.Interceptor
	Consumer        consumer.Consumer
	LocationFactory location.Factory
	FieldsKeysSpec  FieldsKeysSpec
}

func NewProvider(name string) *Provider {
	result := &Provider{
		Name:            name,
		Level:           log.LevelInfo,
		LevelNames:      level.DefaultLevelNamesFacade,
		LocationFactory: location.DefaultFactoryFacade,
		FieldsKeysSpec:  DefaultFieldsKeySpecFacade,
		Formatter:       formatter.DefaultFacade,
	}
	result.Cache = log.NewLoggerCache(result.factory)
	result.Consumer = consumer.NewWritingConsumer(result, os.Stderr)
	return result
}

func (instance *Provider) GetName() string {
	return instance.Name
}

func (instance *Provider) GetLogger(name string) log.Logger {
	return instance.Cache.GetLogger(name)
}

func (instance *Provider) factory(name string) log.Logger {
	cl := &CoreLogger{
		provider: instance,
		name:     name,
	}
	return log.NewLogger(cl)
}

func (instance *Provider) SetLevel(level log.Level) {
	instance.Level = level
}

func (instance *Provider) GetLevel() log.Level {
	return instance.Level
}

func (instance *Provider) GetInterceptor() interceptor.Interceptor {
	return instance.Interceptor
}

func (instance *Provider) SetInterceptor(v interceptor.Interceptor) {
	instance.Interceptor = v
}

func (instance *Provider) GetConsumer() consumer.Consumer {
	return instance.Consumer
}

func (instance *Provider) SetConsumer(v consumer.Consumer) {
	if v == nil {
		panic(fmt.Sprintf("Provider %s cannot handle a consumer of nil.", instance.GetName()))
	}
	instance.Consumer = v
}

func (instance *Provider) getConsumer() consumer.Consumer {
	if c := instance.GetConsumer(); c != nil {
		return c
	}
	panic(fmt.Sprintf("There is no consume for provider %s configured.", instance.GetName()))
}

func (instance *Provider) GetFormatter() formatter.Formatter {
	return instance.Formatter
}

func (instance *Provider) SetFormatter(v formatter.Formatter) {
	instance.Formatter = v
}

func (instance *Provider) GetAllLevels() []log.Level {
	if p := instance.LevelProvider; p != nil {
		return p()
	}
	return log.DefaultLevelProvider()
}

func (instance *Provider) GetLocationFactory() location.Factory {
	return instance.LocationFactory
}

func (instance *Provider) SetLocationFactory(v location.Factory) {
	instance.LocationFactory = v
}

func (instance *Provider) getLocationFactory() location.Factory {
	if f := instance.GetLocationFactory(); f != nil {
		return f
	}
	return location.NoopFactory
}

func (instance *Provider) GetFieldKeySpec() fields.KeysSpec {
	return instance.FieldsKeysSpec
}

func (instance *Provider) GetLevelNames() level.Names {
	return instance.LevelNames
}
