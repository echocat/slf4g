package native

import (
	"fmt"
	"github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native/consumer"
	"github.com/echocat/slf4g/native/formatter"
	"github.com/echocat/slf4g/native/interceptor"
	"os"
)

var DefaultProvider = NewProvider("native")
var _ = log.RegisterProvider(DefaultProvider)

type Provider struct {
	log.Provider

	Level         log.Level
	LevelProvider log.LevelProvider

	EventFormatter formatter.Formatter
	Interceptor    interceptor.Interceptor
	Consumer       consumer.Consumer
}

func NewProvider(name string) *Provider {
	result := &Provider{
		Level: log.LevelInfo,
	}
	result.Provider = log.NewProvider(name, result.factory, result.provideLevels)
	result.Consumer = consumer.NewWritingConsumer(result, os.Stderr)
	return result
}

func (instance *Provider) factory(name string) log.Logger {
	prefix := name
	if prefix == log.GlobalLoggerName {
		prefix = ""
	}
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

func (instance *Provider) getInterceptor() interceptor.Interceptor {
	if i := instance.GetInterceptor(); i != nil {
		return i
	}
	if i := interceptor.Default; i != nil {
		return i
	}
	return interceptor.Noop()
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
	return instance.EventFormatter
}

func (instance *Provider) SetFormatter(v formatter.Formatter) {
	instance.EventFormatter = v
}

func (instance *Provider) provideLevels() []log.Level {
	if p := instance.LevelProvider; p != nil {
		return p()
	}
	return log.DefaultLevelProvider()
}
