package log

import (
	"io"
	"os"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/fields"
)

type fallbackProvider struct {
	cache LoggerCache
	level level.Level
	out   io.Writer
}

var fallbackProviderV = func() *fallbackProvider {
	result := &fallbackProvider{
		level: level.Info,
		out:   os.Stderr,
	}
	result.cache = NewLoggerCache(result.rootFactory, result.factory)
	return result
}()

func (instance *fallbackProvider) rootFactory() Logger {
	return instance.factory(fallbackRootLoggerName)
}

func (instance *fallbackProvider) factory(name string) Logger {
	cl := &fallbackCoreLogger{
		fallbackProvider: instance,
		name:             name,
	}
	return NewLogger(cl)
}

func (instance *fallbackProvider) GetName() string {
	return "fallback"
}

func (instance *fallbackProvider) GetRootLogger() Logger {
	return instance.cache.GetRootLogger()
}

func (instance *fallbackProvider) GetLogger(name string) Logger {
	return instance.cache.GetLogger(name)
}

func (instance *fallbackProvider) GetAllLevels() level.Levels {
	return level.GetProvider().GetLevels()
}

func (instance *fallbackProvider) GetFieldKeysSpec() fields.KeysSpec {
	return fallbackFieldKeysSpecV
}

var fallbackFieldKeysSpecV = &fallbackFieldKeysSpec{}

type fallbackFieldKeysSpec struct{}

func (instance *fallbackFieldKeysSpec) GetTimestamp() string {
	return "timestamp"
}

func (instance *fallbackFieldKeysSpec) GetMessage() string {
	return "message"
}

func (instance *fallbackFieldKeysSpec) GetError() string {
	return "error"
}

func (instance *fallbackFieldKeysSpec) GetLogger() string {
	return "logger"
}
