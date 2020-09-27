package log

import (
	"fmt"

	"github.com/echocat/slf4g/fields"
)

func GetLogger(name string) Logger {
	return GetProvider().GetLogger(name)
}

var globalLogger = GetLogger(GlobalLoggerName)

func GetGlobalLogger() Logger {
	return globalLogger
}

func log(level Level, args ...interface{}) {
	if !IsLevelEnabled(level) {
		return
	}

	var f fields.Fields
	if len(args) == 1 {
		f = fields.With(GetProvider().GetFieldKeySpec().GetMessage(), args[0])
	} else if len(args) > 1 {
		f = fields.With(GetProvider().GetFieldKeySpec().GetMessage(), fields.LazyFunc(func() interface{} {
			return fmt.Sprint(args...)
		}))
	}

	GetGlobalLogger().Log(NewEvent(level, f, 2))
}

func logf(level Level, format string, args ...interface{}) {
	if !IsLevelEnabled(level) {
		return
	}

	f := fields.With(GetProvider().GetFieldKeySpec().GetMessage(), fields.LazyFormat(format, args...))

	GetGlobalLogger().Log(NewEvent(level, f, 2))
}

func IsLevelEnabled(level Level) bool {
	return GetGlobalLogger().IsLevelEnabled(level)
}

func Trace(args ...interface{}) {
	log(LevelTrace, args...)
}

func Tracef(format string, args ...interface{}) {
	logf(LevelTrace, format, args...)
}

func IsTraceEnabled() bool {
	return IsLevelEnabled(LevelTrace)
}

func Debug(args ...interface{}) {
	log(LevelDebug, args...)
}

func Debugf(format string, args ...interface{}) {
	logf(LevelDebug, format, args...)
}

func IsDebugEnabled() bool {
	return IsLevelEnabled(LevelDebug)
}

func Info(args ...interface{}) {
	log(LevelInfo, args...)
}

func Infof(format string, args ...interface{}) {
	logf(LevelInfo, format, args...)
}

func IsInfoEnabled() bool {
	return IsLevelEnabled(LevelInfo)
}

func Warn(args ...interface{}) {
	log(LevelWarn, args...)
}

func Warnf(format string, args ...interface{}) {
	logf(LevelWarn, format, args...)
}

func IsWarnEnabled() bool {
	return IsLevelEnabled(LevelWarn)
}

func Error(args ...interface{}) {
	log(LevelError, args...)
}

func Errorf(format string, args ...interface{}) {
	logf(LevelError, format, args...)
}

func IsErrorEnabled() bool {
	return IsLevelEnabled(LevelError)
}

func Fatal(args ...interface{}) {
	log(LevelFatal, args...)
}

func Fatalf(format string, args ...interface{}) {
	logf(LevelFatal, format, args...)
}

func IsFatalEnabled() bool {
	return IsLevelEnabled(LevelFatal)
}

func With(name string, value interface{}) Logger {
	return GetGlobalLogger().With(name, value)
}

func Withf(name string, format string, args ...interface{}) Logger {
	return GetGlobalLogger().Withf(name, format, args...)
}

func WithError(err error) Logger {
	return GetGlobalLogger().WithError(err)
}

func WithAll(of map[string]interface{}) Logger {
	return GetGlobalLogger().WithAll(of)
}
