package log

import (
	"github.com/echocat/slf4g/fields"
)

type LoggingWriter struct {
	CoreLogger
	LogAs Level
}

func (instance *LoggingWriter) Write(p []byte) (n int, err error) {
	instance.Log(NewEvent(
		instance.LogAs,
		fields.With(GetProvider().GetFieldKeySpec().GetMessage(), string(p)),
		3,
	))
	return len(p), nil
}
