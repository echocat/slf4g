package log

import (
	"io"
	"os"
)

type SimpleProvider struct {
	Provider

	Level Level
	Out   io.Writer
}

func NewSimpleProvider(name string) *SimpleProvider {
	result := &SimpleProvider{
		Level: LevelInfo,
		Out:   os.Stderr,
	}
	result.Provider = NewProvider(name, result.factory, DefaultLevelProvider)
	return result
}

func (instance *SimpleProvider) factory(name string) Logger {
	prefix := name
	if prefix == GlobalLoggerName {
		prefix = ""
	}
	cl := &simpleCoreLogger{
		SimpleProvider: instance,
		name:           name,
	}
	return NewLogger(cl)
}

func (instance *SimpleProvider) SetLevel(level Level) {
	instance.Level = level
}

func (instance *SimpleProvider) GetLevel() Level {
	return instance.Level
}
