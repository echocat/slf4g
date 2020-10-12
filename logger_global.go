package log

import (
	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/fields"
)

// GetLogger returns a logger for the given name from the global Provider.
func GetLogger(name string) Logger {
	return GetProvider().GetLogger(name)
}

// GetLogger returns the ROOT logger from the global Provider.
func GetRootLogger() Logger {
	return GetProvider().GetRootLogger()
}

// IsLevelEnabled checks if the given Level is enabled at the current root
// Logger.
func IsLevelEnabled(level level.Level) bool {
	return GetRootLogger().IsLevelEnabled(level)
}

// Trace logs the provided arguments on Trace at the current root Logger.
func Trace(args ...interface{}) {
	log(level.Trace, args...)
}

// Tracef is like Trace but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
func Tracef(format string, args ...interface{}) {
	logf(level.Trace, format, args...)
}

// IsTraceEnabled checks if Trace is enabled at the current root Logger.
func IsTraceEnabled() bool {
	return IsLevelEnabled(level.Trace)
}

// Debug logs the provided arguments on Debug at the current root Logger.
func Debug(args ...interface{}) {
	log(level.Debug, args...)
}

// Debugf is like Debug but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
func Debugf(format string, args ...interface{}) {
	logf(level.Debug, format, args...)
}

// IsDebugEnabled checks if Debug is enabled at the current root Logger.
func IsDebugEnabled() bool {
	return IsLevelEnabled(level.Debug)
}

// Info logs the provided arguments on Info at the current root Logger.
func Info(args ...interface{}) {
	log(level.Info, args...)
}

// Infof is like Info but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
func Infof(format string, args ...interface{}) {
	logf(level.Info, format, args...)
}

// IsInfoEnabled checks if Info is enabled at the current root Logger.
func IsInfoEnabled() bool {
	return IsLevelEnabled(level.Info)
}

// Warn logs the provided arguments on Warn at the current root Logger.
func Warn(args ...interface{}) {
	log(level.Warn, args...)
}

// Warnf is like Warn but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
func Warnf(format string, args ...interface{}) {
	logf(level.Warn, format, args...)
}

// IsWarnEnabled checks if Warn is enabled at the current root Logger.
func IsWarnEnabled() bool {
	return IsLevelEnabled(level.Warn)
}

// Error logs the provided arguments on Error at the current root Logger.
func Error(args ...interface{}) {
	log(level.Error, args...)
}

// Errorf is like Error but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
func Errorf(format string, args ...interface{}) {
	logf(level.Error, format, args...)
}

// IsErrorEnabled checks if Error is enabled at the current root Logger.
func IsErrorEnabled() bool {
	return IsLevelEnabled(level.Error)
}

// Fatal logs the provided arguments on Fatal at the current root Logger.
//
// IMPORTANT! In contrast to many other log Golang frameworks this logging Fatal
// with slf4g does not lead to an os.Exit() by default. By contract the
// application can do that but it is doing that always GRACEFUL. All processes
// should be always able to do shutdown operations if needed AND possible.
func Fatal(args ...interface{}) {
	log(level.Fatal, args...)
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
	logf(level.Fatal, format, args...)
}

// IsFatalEnabled checks if Fatal is enabled at the current root Logger.
func IsFatalEnabled() bool {
	return IsLevelEnabled(level.Fatal)
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

func log(l level.Level, args ...interface{}) {
	if !IsLevelEnabled(l) {
		return
	}

	p := GetProvider()
	e := NewEvent(p, l, 2)

	if len(args) == 1 {
		e = e.With(p.GetFieldKeysSpec().GetMessage(), args[0])
	} else if len(args) > 0 {
		e = e.With(p.GetFieldKeysSpec().GetMessage(), args)
	}

	GetRootLogger().Log(e)
}

func logf(l level.Level, format string, args ...interface{}) {
	if !IsLevelEnabled(l) {
		return
	}

	p := GetProvider()
	e := NewEvent(p, l, 2).
		With(GetProvider().GetFieldKeysSpec().GetMessage(), fields.LazyFormat(format, args...))

	GetRootLogger().Log(e)
}
