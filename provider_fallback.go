package log

import (
	"github.com/echocat/slf4g/fields"
	"io"
	"os"
)

type fallbackProvider struct {
	cache LoggerCache
	level Level
	out   io.Writer
}

var simpleProviderV = func() *fallbackProvider {
	result := &fallbackProvider{
		level: LevelInfo,
		out:   os.Stderr,
	}
	result.cache = NewLoggerCache(result.factory)
	return result
}()

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

func (instance *fallbackProvider) GetLogger(name string) Logger {
	return instance.cache.GetLogger(name)
}

func (instance *fallbackProvider) GetAllLevels() []Level {
	return DefaultLevelProvider()
}

func (instance *fallbackProvider) GetFieldKeySpec() fields.KeysSpec {
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
