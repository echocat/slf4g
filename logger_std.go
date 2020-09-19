package log

import (
	"github.com/echocat/slf4g/fields"
	stdlog "log"
)

type StdLogger interface {
	CoreLogger

	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}

func AsStdInterfacingLogger(in CoreLogger) StdLogger {
	if v, ok := in.(StdLogger); ok {
		return v
	}
	return &loggerImpl{getCoreLogger: func() CoreLogger {
		return in
	}}
}

func AsStdLogger(in CoreLogger, logAs Level) *stdlog.Logger {
	return stdlog.New(&Writer{
		CoreLogger: in,
		LogAs:      logAs,
	}, "", 0)
}

func ConfigureStd() {
	ConfigureStdWith(GetGlobalLogger(), LevelInfo)
}

func ConfigureStdWith(in CoreLogger, logAs Level) {
	w := &Writer{
		CoreLogger: in,
		LogAs:      logAs,
	}
	stdlog.SetOutput(w)
	stdlog.SetPrefix("")
	stdlog.SetFlags(0)
}

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
