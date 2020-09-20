package log

import (
	"io"
	"os"
)

type simpleProvider struct {
	Provider

	Level Level
	Out   io.Writer
}

var simpleProviderV = func() *simpleProvider {
	result := &simpleProvider{
		Level: LevelInfo,
		Out:   os.Stderr,
	}
	result.Provider = NewProvider("simple", result.factory, DefaultLevelProvider)
	return result
}()

func (instance *simpleProvider) factory(name string) Logger {
	prefix := name
	if prefix == GlobalLoggerName {
		prefix = ""
	}
	cl := &simpleCoreLogger{
		simpleProvider: instance,
		name:           name,
	}
	return NewLogger(cl)
}

func (instance *simpleProvider) SetLevel(level Level) {
	instance.Level = level
}

func (instance *simpleProvider) GetLevel() Level {
	return instance.Level
}
