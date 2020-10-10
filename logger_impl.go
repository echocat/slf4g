package log

import (
	"fmt"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/fields"
)

type loggerImpl struct {
	coreProvider func() CoreLogger
	fields       fields.Fields
}

func (instance *loggerImpl) Unwrap() CoreLogger {
	return instance.coreProvider()
}

func (instance *loggerImpl) GetName() string {
	return instance.Unwrap().GetName()
}

func (instance *loggerImpl) Log(event Event) {
	if !instance.IsLevelEnabled(event.GetLevel()) {
		return
	}
	instance.Unwrap().Log(event.WithCallDepth(1))
}

func (instance *loggerImpl) IsLevelEnabled(level level.Level) bool {
	return instance.Unwrap().IsLevelEnabled(level)
}

func (instance *loggerImpl) GetProvider() Provider {
	return instance.Unwrap().GetProvider()
}

func (instance *loggerImpl) log(level level.Level, args ...interface{}) {
	e := NewEvent(instance.GetProvider(), level, 2, instance.fields)

	if len(args) == 1 {
		e = e.With(instance.GetProvider().GetFieldKeysSpec().GetMessage(), args[0])
	} else if len(args) > 1 {
		e = e.With(instance.GetProvider().GetFieldKeysSpec().GetMessage(), fields.LazyFunc(func() interface{} {
			return fmt.Sprint(args...)
		}))
	}

	instance.Log(e)
}

func (instance *loggerImpl) logf(level level.Level, format string, args ...interface{}) {
	e := NewEvent(instance.GetProvider(), level, 2, instance.fields).
		Withf(instance.GetProvider().GetFieldKeysSpec().GetMessage(), format, args...)

	instance.Log(e)
}

func (instance *loggerImpl) Trace(args ...interface{}) {
	instance.log(level.Trace, args...)
}

func (instance *loggerImpl) Tracef(format string, args ...interface{}) {
	instance.logf(level.Trace, format, args...)
}

func (instance *loggerImpl) IsTraceEnabled() bool {
	return instance.IsLevelEnabled(level.Trace)
}

func (instance *loggerImpl) Debug(args ...interface{}) {
	instance.log(level.Debug, args...)
}

func (instance *loggerImpl) Debugf(format string, args ...interface{}) {
	instance.logf(level.Debug, format, args...)
}

func (instance *loggerImpl) IsDebugEnabled() bool {
	return instance.IsLevelEnabled(level.Debug)
}

func (instance *loggerImpl) Info(args ...interface{}) {
	instance.log(level.Info, args...)
}

func (instance *loggerImpl) Infof(format string, args ...interface{}) {
	instance.logf(level.Info, format, args...)
}

func (instance *loggerImpl) IsInfoEnabled() bool {
	return instance.IsLevelEnabled(level.Info)
}

func (instance *loggerImpl) Warn(args ...interface{}) {
	instance.log(level.Warn, args...)
}

func (instance *loggerImpl) Warnf(format string, args ...interface{}) {
	instance.logf(level.Warn, format, args...)
}

func (instance *loggerImpl) IsWarnEnabled() bool {
	return instance.IsLevelEnabled(level.Warn)
}

func (instance *loggerImpl) Error(args ...interface{}) {
	instance.log(level.Error, args...)
}

func (instance *loggerImpl) Errorf(format string, args ...interface{}) {
	instance.logf(level.Error, format, args...)
}

func (instance *loggerImpl) IsErrorEnabled() bool {
	return instance.IsLevelEnabled(level.Error)
}

func (instance *loggerImpl) Fatal(args ...interface{}) {
	instance.log(level.Fatal, args...)
}

func (instance *loggerImpl) Fatalf(format string, args ...interface{}) {
	instance.logf(level.Fatal, format, args...)
}

func (instance *loggerImpl) IsFatalEnabled() bool {
	return instance.IsLevelEnabled(level.Fatal)
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
	return instance.With(instance.GetProvider().GetFieldKeysSpec().GetError(), err)
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

func (instance *loggerImpl) Without(keys ...string) Logger {
	targetFields := instance.fields
	if targetFields != nil {
		targetFields = targetFields.Without(keys...)
	} else {
		targetFields = fields.Empty()
	}
	return &loggerImpl{
		coreProvider: instance.coreProvider,
		fields:       targetFields,
	}
}
