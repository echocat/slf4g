package log

import (
	"fmt"

	"github.com/echocat/slf4g/fields"
)

type loggerImpl struct {
	coreProvider func() CoreLogger
	fields       fields.Fields
}

func (instance *loggerImpl) UnwrapCore() CoreLogger {
	return instance.coreProvider()
}

func (instance *loggerImpl) GetName() string {
	return instance.UnwrapCore().GetName()
}

func (instance *loggerImpl) Log(event Event) {
	if !instance.IsLevelEnabled(event.GetLevel()) {
		return
	}
	instance.UnwrapCore().Log(event.WithCallDepth(1))
}

func (instance *loggerImpl) IsLevelEnabled(level Level) bool {
	return instance.UnwrapCore().IsLevelEnabled(level)
}

func (instance *loggerImpl) GetProvider() Provider {
	return instance.UnwrapCore().GetProvider()
}

func (instance *loggerImpl) log(level Level, args ...interface{}) {
	f := instance.fields
	if len(args) == 1 {
		f = f.With(instance.GetProvider().GetFieldKeySpec().GetMessage(), args[0])
	} else if len(args) > 1 {
		f = f.With(instance.GetProvider().GetFieldKeySpec().GetMessage(), fields.LazyFunc(func() interface{} {
			return fmt.Sprint(args...)
		}))
	}
	instance.Log(NewEvent(level, f, 2))
}

func (instance *loggerImpl) logf(level Level, format string, args ...interface{}) {
	f := instance.fields.
		Withf(instance.GetProvider().GetFieldKeySpec().GetMessage(), format, args...)
	instance.Log(NewEvent(level, f, 2))
}

func (instance *loggerImpl) Trace(args ...interface{}) {
	instance.log(LevelTrace, args...)
}

func (instance *loggerImpl) Tracef(format string, args ...interface{}) {
	instance.logf(LevelTrace, format, args...)
}

func (instance *loggerImpl) IsTraceEnabled() bool {
	return instance.IsLevelEnabled(LevelTrace)
}

func (instance *loggerImpl) Debug(args ...interface{}) {
	instance.log(LevelDebug, args...)
}

func (instance *loggerImpl) Debugf(format string, args ...interface{}) {
	instance.logf(LevelDebug, format, args...)
}

func (instance *loggerImpl) IsDebugEnabled() bool {
	return instance.IsLevelEnabled(LevelDebug)
}

func (instance *loggerImpl) Info(args ...interface{}) {
	instance.log(LevelInfo, args...)
}

func (instance *loggerImpl) Infof(format string, args ...interface{}) {
	instance.logf(LevelInfo, format, args...)
}

func (instance *loggerImpl) IsInfoEnabled() bool {
	return instance.IsLevelEnabled(LevelInfo)
}

func (instance *loggerImpl) Warn(args ...interface{}) {
	instance.log(LevelWarn, args...)
}

func (instance *loggerImpl) Warnf(format string, args ...interface{}) {
	instance.logf(LevelWarn, format, args...)
}

func (instance *loggerImpl) IsWarnEnabled() bool {
	return instance.IsLevelEnabled(LevelWarn)
}

func (instance *loggerImpl) Error(args ...interface{}) {
	instance.log(LevelError, args...)
}

func (instance *loggerImpl) Errorf(format string, args ...interface{}) {
	instance.logf(LevelError, format, args...)
}

func (instance *loggerImpl) IsErrorEnabled() bool {
	return instance.IsLevelEnabled(LevelError)
}

func (instance *loggerImpl) Fatal(args ...interface{}) {
	instance.log(LevelFatal, args...)
}

func (instance *loggerImpl) Fatalf(format string, args ...interface{}) {
	instance.logf(LevelFatal, format, args...)
}

func (instance *loggerImpl) IsFatalEnabled() bool {
	return instance.IsLevelEnabled(LevelFatal)
}

func (instance *loggerImpl) With(name string, value interface{}) Logger {
	targetFields := instance.fields
	if targetFields != nil {
		targetFields = targetFields.With(name, value)
	} else {
		targetFields = fields.With(name, value)
	}
	return &loggerImpl{
		coreProvider: instance.coreProvider,
		fields:       targetFields,
	}
}

func (instance *loggerImpl) Withf(name string, format string, args ...interface{}) Logger {
	return instance.With(name, fields.LazyFormat(format, args...))
}

func (instance *loggerImpl) WithError(err error) Logger {
	return instance.With(instance.GetProvider().GetFieldKeySpec().GetError(), err)
}

func (instance *loggerImpl) WithAll(of map[string]interface{}) Logger {
	targetFields := instance.fields
	if targetFields != nil {
		targetFields = targetFields.WithAll(of)
	} else {
		targetFields = fields.WithAll(of)
	}
	return &loggerImpl{
		coreProvider: instance.coreProvider,
		fields:       targetFields,
	}
}
