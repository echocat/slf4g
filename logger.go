package log

import "github.com/echocat/slf4g/fields"

const (
	GlobalLoggerName = "global"
)

type CoreLogger interface {
	GetName() string
	LogEvent(Event)
	IsLevelEnabled(Level) bool
	GetProvider() Provider
}

type Logger interface {
	CoreLogger

	Log(Level, ...interface{})
	Logf(Level, string, ...interface{})

	Trace(...interface{})
	Tracef(string, ...interface{})
	IsTraceEnabled() bool

	Debug(...interface{})
	Debugf(string, ...interface{})
	IsDebugEnabled() bool

	Info(...interface{})
	Infof(string, ...interface{})
	IsInfoEnabled() bool

	Warn(...interface{})
	Warnf(string, ...interface{})
	IsWarnEnabled() bool

	Error(...interface{})
	Errorf(string, ...interface{})
	IsErrorEnabled() bool

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	IsFatalEnabled() bool

	Panic(...interface{})
	Panicf(string, ...interface{})
	IsPanicEnabled() bool

	With(name string, value interface{}) Logger
	WithError(error) Logger
	WithFields(fields.Fields) Logger
}

func AsLogger(in CoreLogger) Logger {
	if v, ok := in.(Logger); ok {
		return v
	}
	return &loggerImpl{CoreLogger: in}
}
