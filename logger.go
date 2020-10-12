package log

import "github.com/echocat/slf4g/fields"

// Logger defines an instance which executes log event actions.
//
// Implementation hints
//
// If you considering to implement slf4g you're usually not required to
// implement a full instance of Logger. Usually you just need to implement
// CoreLogger and call NewLogger(with the CoreLogger) to create a full
// implemented instance of a Logger.
type Logger interface {
	CoreLogger

	// Trace logs the provided arguments on LevelTrace with this Logger.
	Trace(...interface{})

	// Tracef is like Trace but wraps the message itself in a fmt.Sprintf action.
	// By contract the actual format action will not be executed before the value
	// will be really consumed.
	Tracef(string, ...interface{})

	// IsTraceEnabled checks if LevelTrace is enabled at this Logger.
	IsTraceEnabled() bool

	// Debug logs the provided arguments on LevelDebug with this Logger.
	Debug(...interface{})

	// Debugf is like Debug but wraps the message itself in a fmt.Sprintf action.
	// By contract the actual format action will not be executed before the value
	// will be really consumed.
	Debugf(string, ...interface{})

	// IsDebugEnabled checks if LevelDebug is enabled at this Logger.
	IsDebugEnabled() bool

	// Info logs the provided arguments on LevelInfo with this Logger.
	Info(...interface{})

	// Infof is like Info but wraps the message itself in a fmt.Sprintf action.
	// By contract the actual format action will not be executed before the value
	// will be really consumed.
	Infof(string, ...interface{})

	// IsInfoEnabled checks if LevelInfo is enabled at this Logger.
	IsInfoEnabled() bool

	// Warn logs the provided arguments on LevelWarn with this Logger.
	Warn(...interface{})

	// Warnf is like Warn but wraps the message itself in a fmt.Sprintf action.
	// By contract the actual format action will not be executed before the value
	// will be really consumed.
	Warnf(string, ...interface{})

	// IsWarnEnabled checks if LevelWarn is enabled at this Logger.
	IsWarnEnabled() bool

	// Error logs the provided arguments on LevelError with this Logger.
	Error(...interface{})

	// Errorf is like Error but wraps the message itself in a fmt.Sprintf action.
	// By contract the actual format action will not be executed before the value
	// will be really consumed.
	Errorf(string, ...interface{})

	// IsErrorEnabled checks if LevelError is enabled at this Logger.
	IsErrorEnabled() bool

	// Fatal logs the provided arguments on LevelFatal with this Logger.
	//
	// IMPORTANT! In contrast to many other log Golang frameworks this logging Fatal
	// with slf4g does not lead to an os.Exit() by default. By contract the
	// application can do that but it is doing that always GRACEFUL. All processes
	// should be always able to do shutdown operations if needed AND possible.
	Fatal(...interface{})

	// Fatalf is like Fatal but wraps the message itself in a fmt.Sprintf action.
	// By contract the actual format action will not be executed before the value
	// will be really consumed.
	//
	// IMPORTANT! In contrast to many other log Golang frameworks this logging Fatal
	// with slf4g does not lead to an os.Exit() by default. By contract the
	// application can do that but it is doing that always GRACEFUL. All processes
	// should be always able to do shutdown operations if needed AND possible.
	Fatalf(string, ...interface{})

	// IsFatalEnabled checks if LevelFatal is enabled at this Logger.
	IsFatalEnabled() bool

	// With returns an variant of this Logger with the given key
	// value pair contained inside. If the given key already exists in the
	// current instance this means it will be overwritten.
	With(name string, value interface{}) Logger

	// Withf is similar to With but it adds classic fmt.Printf functions to it.
	// It is defined that the format itself will not be executed before the
	// consumption of the value.
	Withf(name string, format string, args ...interface{}) Logger

	// WithError is similar to With but it adds an error as field.
	WithError(error) Logger

	// WithAll is similar to With but it can consume more than one field at
	// once. Be aware: There is neither a guarantee that this instance will be
	// copied or not.
	WithAll(map[string]interface{}) Logger

	// Without returns a variant of this Logger without the given
	// key contained inside. In other words: If someone afterwards tries to
	// call either ForEach() or Get() nothing with this key(s) will be returned.
	Without(keys ...string) Logger
}

// NewLogger create a new fully implemented instance of a logger out of a given
// CoreLogger instance.
func NewLogger(cl CoreLogger) Logger {
	return NewLoggerFacade(func() CoreLogger { return cl })
}

// NewLoggerFacade is like NewLogger but takes a provider function that can
// potentially return every time another instance of a CoreLogger. This is
// useful especially in cases where you want to deal with concurrency while
// creation of objects that need to hold a reference to a Logger.
func NewLoggerFacade(provider func() CoreLogger) Logger {
	return &loggerImpl{
		coreProvider: provider,
		fields:       fields.Empty(),
	}
}
