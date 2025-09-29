package log

import (
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/names"

	"github.com/echocat/slf4g/fields"
)

// GetRootLogger returns the ROOT logger from the global Provider.
func GetRootLogger() Logger {
	return GetProvider().GetRootLogger()
}

// GetLogger returns a logger for the given name from the global Provider.
// If instead of a string another object is given this will be used to create
// a logger name from its package name.
func GetLogger(nameOrReference interface{}) Logger {
	return GetProvider().GetLogger(names.FullLoggerNameGenerator(nameOrReference))
}

// GetLoggerForCurrentPackage is similar to GetLogger(name) but it extracts the
// name automatically from the current call stack. That means: this method
// should be only used for the package this logger should be for.
func GetLoggerForCurrentPackage() Logger {
	return GetProvider().GetLogger(names.CurrentPackageLoggerNameGenerator(1))
}

// IsLevelEnabled checks if the given Level is enabled at the current root
// Logger.
func IsLevelEnabled(level level.Level) bool {
	return GetRootLogger().IsLevelEnabled(level)
}

// Trace logs the provided arguments on Trace at the current root Logger.
func Trace(args ...interface{}) {
	l, helper := log(level.Trace, args...)
	helper()
	l()
}

// Tracef is like Trace but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
func Tracef(format string, args ...interface{}) {
	l, helper := logf(level.Trace, format, args...)
	helper()
	l()
}

// IsTraceEnabled checks if Trace is enabled at the current root Logger.
func IsTraceEnabled() bool {
	return IsLevelEnabled(level.Trace)
}

// Debug logs the provided arguments on Debug at the current root Logger.
func Debug(args ...interface{}) {
	l, helper := log(level.Debug, args...)
	helper()
	l()
}

// Debugf is like Debug but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
func Debugf(format string, args ...interface{}) {
	l, helper := logf(level.Debug, format, args...)
	helper()
	l()
}

// IsDebugEnabled checks if Debug is enabled at the current root Logger.
func IsDebugEnabled() bool {
	return IsLevelEnabled(level.Debug)
}

// Info logs the provided arguments on Info at the current root Logger.
func Info(args ...interface{}) {
	l, helper := log(level.Info, args...)
	helper()
	l()
}

// Infof is like Info but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
func Infof(format string, args ...interface{}) {
	l, helper := logf(level.Info, format, args...)
	helper()
	l()
}

// IsInfoEnabled checks if Info is enabled at the current root Logger.
func IsInfoEnabled() bool {
	return IsLevelEnabled(level.Info)
}

// Warn logs the provided arguments on Warn at the current root Logger.
func Warn(args ...interface{}) {
	l, helper := log(level.Warn, args...)
	helper()
	l()
}

// Warnf is like Warn but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
func Warnf(format string, args ...interface{}) {
	l, helper := logf(level.Warn, format, args...)
	helper()
	l()
}

// IsWarnEnabled checks if Warn is enabled at the current root Logger.
func IsWarnEnabled() bool {
	return IsLevelEnabled(level.Warn)
}

// Error logs the provided arguments on Error at the current root Logger.
func Error(args ...interface{}) {
	l, helper := log(level.Error, args...)
	helper()
	l()
}

// Errorf is like Error but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
func Errorf(format string, args ...interface{}) {
	l, helper := logf(level.Error, format, args...)
	helper()
	l()
}

// IsErrorEnabled checks if Error is enabled at the current root Logger.
func IsErrorEnabled() bool {
	return IsLevelEnabled(level.Error)
}

// Fatal logs the provided arguments on Fatal at the current root Logger.
//
// IMPORTANT! In contrast to many other log Golang frameworks this logging Fatal
// with slf4g does not lead to an os.Exit() by default. By contract the
// application can do that, but it is doing that always GRACEFUL. All processes
// should be always able to do shut down operations if needed AND possible.
func Fatal(args ...interface{}) {
	l, helper := log(level.Fatal, args...)
	helper()
	l()
}

// Fatalf is like Fatal but wraps the message itself in a fmt.Sprintf action.
// By contract the actual format action will not be executed before the value
// will be really consumed.
//
// IMPORTANT! In contrast to many other log Golang frameworks this logging Fatal
// with slf4g does not lead to an os.Exit() by default. By contract the
// application can do that, but it is doing that always GRACEFUL. All processes
// should be always able to do shut down operations if needed AND possible.
func Fatalf(format string, args ...interface{}) {
	l, helper := logf(level.Fatal, format, args...)
	helper()
	l()
}

// IsFatalEnabled checks if Fatal is enabled at the current root Logger.
func IsFatalEnabled() bool {
	return IsLevelEnabled(level.Fatal)
}

// With returns a root Logger which will contain the provided field.
func With(name string, value interface{}) Logger {
	return GetRootLogger().With(name, value)
}

// Withf is similar to With, but it adds classic fmt.Printf functions to it.
// It is defined that the format itself will not be executed before the
// consumption of the value.
func Withf(name string, format string, args ...interface{}) Logger {
	return GetRootLogger().Withf(name, format, args...)
}

// WithError is similar to With, but it adds specially an error field.
func WithError(err error) Logger {
	return GetRootLogger().WithError(err)
}

// WithAll is similar to With, but it can consume more than one field at
// once. Be aware: There is neither a guarantee that this instance will be
// copied or not.
func WithAll(of map[string]interface{}) Logger {
	return GetRootLogger().WithAll(of)
}

func log(l level.Level, args ...interface{}) (doLog, helper func()) {
	p := GetProvider()
	logger := p.GetRootLogger()
	helper = helperOf(logger)
	if !logger.IsLevelEnabled(l) {
		return func() {}, helper
	}

	var values map[string]interface{}

	if len(args) > 0 {
		values = make(map[string]interface{}, 1)
		if len(args) == 1 {
			values[p.GetFieldKeysSpec().GetMessage()] = args[0]
		} else {
			values[p.GetFieldKeysSpec().GetMessage()] = args
		}
	}

	e := logger.NewEvent(l, values)
	return func() {
		helper()
		logger.Log(e, 2)
	}, helper
}

func logf(l level.Level, format string, args ...interface{}) (doLog, helper func()) {
	p := GetProvider()
	logger := p.GetRootLogger()
	helper = helperOf(logger)
	if !logger.IsLevelEnabled(l) {
		return func() {}, helper
	}

	values := map[string]interface{}{
		p.GetFieldKeysSpec().GetMessage(): fields.LazyFormat(format, args...),
	}

	e := logger.NewEvent(l, values)
	return func() {
		helper()
		logger.Log(e, 2)
	}, helper
}
