package log

import (
	"fmt"

	"github.com/echocat/slf4g/fields"
)

// GetLogger returns a logger for the given name from the global Provider.
func GetLogger(name string) Logger {
	return GetProvider().GetLogger(name)
}

// GetLogger returns the ROOT logger from the global Provider.
func GetRootLogger() Logger {
	return rootLogger
}

// IsLevelEnabled checks if the given Level is enabled at the current root
// Logger.
func IsLevelEnabled(level Level) bool {
	return GetRootLogger().IsLevelEnabled(level)
}

// Trace logs the provided arguments on LevelTrace at the current root Logger.
func Trace(args ...interface{}) {
	log(LevelTrace, args...)
}

// Tracef is like Trace but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
func Tracef(format string, args ...interface{}) {
	logf(LevelTrace, format, args...)
}

// IsTraceEnabled checks if LevelTrace is enabled at the current root Logger.
func IsTraceEnabled() bool {
	return IsLevelEnabled(LevelTrace)
}

// Debug logs the provided arguments on LevelDebug at the current root Logger.
func Debug(args ...interface{}) {
	log(LevelDebug, args...)
}

// Debugf is like Debug but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
func Debugf(format string, args ...interface{}) {
	logf(LevelDebug, format, args...)
}

// IsDebugEnabled checks if LevelDebug is enabled at the current root Logger.
func IsDebugEnabled() bool {
	return IsLevelEnabled(LevelDebug)
}

// Info logs the provided arguments on LevelInfo at the current root Logger.
func Info(args ...interface{}) {
	log(LevelInfo, args...)
}

// Infof is like Info but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
func Infof(format string, args ...interface{}) {
	logf(LevelInfo, format, args...)
}

// IsInfoEnabled checks if LevelInfo is enabled at the current root Logger.
func IsInfoEnabled() bool {
	return IsLevelEnabled(LevelInfo)
}

// Warn logs the provided arguments on LevelWarn at the current root Logger.
func Warn(args ...interface{}) {
	log(LevelWarn, args...)
}

// Warnf is like Warn but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
func Warnf(format string, args ...interface{}) {
	logf(LevelWarn, format, args...)
}

// IsWarnEnabled checks if LevelWarn is enabled at the current root Logger.
func IsWarnEnabled() bool {
	return IsLevelEnabled(LevelWarn)
}

// Error logs the provided arguments on LevelError at the current root Logger.
func Error(args ...interface{}) {
	log(LevelError, args...)
}

// Errorf is like Error but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
func Errorf(format string, args ...interface{}) {
	logf(LevelError, format, args...)
}

// IsErrorEnabled checks if LevelError is enabled at the current root Logger.
func IsErrorEnabled() bool {
	return IsLevelEnabled(LevelError)
}

// Fatal logs the provided arguments on LevelFatal at the current root Logger.
//
// IMPORTANT! In contrast to many other log Golang frameworks this logging Fatal
// with slf4g does not lead to an os.Exit() by default. By contract the
// application can do that but it is doing that always GRACEFUL. All processes
// should be always able to do shutdown operations if needed AND possible.
func Fatal(args ...interface{}) {
	log(LevelFatal, args...)
}

// Fatalf is like Fatal but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
//
// IMPORTANT! In contrast to many other log Golang frameworks this logging Fatal
// with slf4g does not lead to an os.Exit() by default. By contract the
// application can do that but it is doing that always GRACEFUL. All processes
// should be always able to do shutdown operations if needed AND possible.
func Fatalf(format string, args ...interface{}) {
	logf(LevelFatal, format, args...)
}

// IsFatalEnabled checks if LevelFatal is enabled at the current root Logger.
func IsFatalEnabled() bool {
	return IsLevelEnabled(LevelFatal)
}

// With returns a root Logger which will contain the provided field.
func With(name string, value interface{}) Logger {
	return GetRootLogger().With(name, value)
}

// Withf is similar to With but it adds classic fmt.Printf functions to it.
// It is defined that the format itself will not be executed before the
// consumption of the value.
func Withf(name string, format string, args ...interface{}) Logger {
	return GetRootLogger().Withf(name, format, args...)
}

// WithError is similar to With but it adds specially an error field.
func WithError(err error) Logger {
	return GetRootLogger().WithError(err)
}

// WithAll is similar to With but it can consume more than one field at
// once. Be aware: There is neither a guarantee that this instance will be
// copied or not.
func WithAll(of map[string]interface{}) Logger {
	return GetRootLogger().WithAll(of)
}

var rootLogger = GetLogger(RootLoggerName)

func log(level Level, args ...interface{}) {
	if !IsLevelEnabled(level) {
		return
	}

	var f fields.Fields
	if len(args) == 1 {
		f = fields.With(GetProvider().GetFieldKeysSpec().GetMessage(), args[0])
	} else if len(args) > 1 {
		f = fields.With(GetProvider().GetFieldKeysSpec().GetMessage(), fields.LazyFunc(func() interface{} {
			return fmt.Sprint(args...)
		}))
	}

	GetRootLogger().Log(NewEvent(level, f, 2))
}

func logf(level Level, format string, args ...interface{}) {
	if !IsLevelEnabled(level) {
		return
	}

	f := fields.With(GetProvider().GetFieldKeysSpec().GetMessage(), fields.LazyFormat(format, args...))

	GetRootLogger().Log(NewEvent(level, f, 2))
}
