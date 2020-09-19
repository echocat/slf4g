package log

import (
	"fmt"
	"github.com/echocat/slf4g/fields"
)

func GetLogger(name string) Logger {
	return GetProvider().GetLogger(name)
}

func GetGlobalLogger() Logger {
	return globalLoggerFacadeV
}

func logM(level Level, message *string) {
	f := fields.Empty()
	if message != nil {
		f = f.With(GetProvider().GetFieldKeySpec().GetMessage(), *message)
	}

	GetGlobalLogger().LogEvent(NewEvent(level, f, 3))
}

func log(level Level, args ...interface{}) {
	var message *string
	if len(args) > 0 {
		str := fmt.Sprint(args...)
		message = &str
	}

	logM(level, message)
}

func logf(level Level, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)

	logM(level, &message)
}

func Log(level Level, args ...interface{}) {
	log(level, args...)
}

func Logf(level Level, format string, args ...interface{}) {
	logf(level, format, args...)
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

func Panic(args ...interface{}) {
	log(LevelPanic, args...)
}

func Panicf(format string, args ...interface{}) {
	logf(LevelPanic, format, args...)
}

func IsPanicEnabled() bool {
	return IsLevelEnabled(LevelPanic)
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

func WithFields(fields fields.Fields) Logger {
	return GetGlobalLogger().WithFields(fields)
}
