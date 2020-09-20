package log

import (
	"github.com/echocat/slf4g/fields"
)

type Writer struct {
	CoreLogger
	LogAs Level
}

func (instance *Writer) Write(p []byte) (n int, err error) {
	instance.LogEvent(NewEvent(
		instance.LogAs,
		fields.With(GetProvider().GetFieldKeySpec().GetMessage(), string(p)),
		3,
	))
	return len(p), nil
}
