package log

import "github.com/echocat/slf4g/fields"

const (
	GlobalLoggerName = "global"
)

type Logger interface {
	CoreLogger

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

	With(name string, value interface{}) Logger
	Withf(name string, format string, args ...interface{}) Logger
	WithError(error) Logger
	WithFields(fields.Fields) Logger
}

func NewLogger(cl CoreLogger) Logger {
	return NewLoggerFacade(func() CoreLogger { return cl })
}

func NewLoggerFacade(provider func() CoreLogger) Logger {
	return &loggerImpl{
		coreProvider: provider,
		fields:       fields.Empty(),
	}
}

type LoggerFactory func(name string) Logger
