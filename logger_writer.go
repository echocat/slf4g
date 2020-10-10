package log

import (
	"github.com/echocat/slf4g/level"
)

type LoggingWriter struct {
	CoreLogger
	LogAs level.Level
}

func (instance *LoggingWriter) Write(p []byte) (n int, err error) {
	provider := GetProvider()
	instance.Log(NewEvent(provider, instance.LogAs, 3).
		With(provider.GetFieldKeysSpec().GetMessage(), string(p)),
	)
	return len(p), nil
}
