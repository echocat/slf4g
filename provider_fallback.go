package log

import (
	"io"
	"os"
)

type fallbackProvider struct {
	Provider

	level Level
	out   io.Writer
}

var simpleProviderV = func() *fallbackProvider {
	result := &fallbackProvider{
		level: LevelInfo,
		out:   os.Stderr,
	}
	result.Provider = NewProvider("fallback", result.factory, DefaultLevelProvider)
	return result
}()

func (instance *fallbackProvider) factory(name string) Logger {
	prefix := name
	if prefix == GlobalLoggerName {
		prefix = ""
	}
	cl := &fallbackCoreLogger{
		fallbackProvider: instance,
		name:             name,
	}
	return NewLogger(cl)
}
