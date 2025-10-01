package log

import (
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

func (instance *loggerImpl) Log(event Event, skipFrames uint16) {
	instance.Helper()()
	instance.Unwrap().Log(event, skipFrames+1)
}

func (instance *loggerImpl) NewEvent(l level.Level, values map[string]interface{}) Event {
	return NewEvent(instance.Unwrap(), l, values)
}

func (instance *loggerImpl) NewEventWithFields(l level.Level, f fields.ForEachEnabled) Event {
	return NewEventWithFields(instance.Unwrap(), l, f)
}

func (instance *loggerImpl) Accepts(event Event) bool {
	return instance.Unwrap().Accepts(event)
}

func (instance *loggerImpl) IsLevelEnabled(level level.Level) bool {
	return instance.Unwrap().IsLevelEnabled(level)
}

func (instance *loggerImpl) GetProvider() Provider {
	return instance.Unwrap().GetProvider()
}

func (instance *loggerImpl) log(level level.Level, args ...interface{}) (doLog, helper func()) {
	return instance.doLog(level, 2, args...)
}

func (instance *loggerImpl) logf(level level.Level, format string, args ...interface{}) (doLog, helper func()) {
	return instance.doLogf(level, 2, format, args...)
}

func (instance *loggerImpl) DoLog(level level.Level, skipFrames uint16, args ...interface{}) {
	l, helper := instance.doLog(level, skipFrames+1, args...)
	helper()
	l()
}

func (instance *loggerImpl) doLog(level level.Level, skipFrames uint16, args ...interface{}) (doLog, helper func()) {
	delegate := instance.Unwrap()
	helper = helperOf(delegate)
	if !delegate.IsLevelEnabled(level) {
		return func() {}, helper
	}
	provider := instance.GetProvider()
	e := instance.NewEventWithFields(level, instance.fields)

	if len(args) == 1 {
		e = e.With(provider.GetFieldKeysSpec().GetMessage(), args[0])
	} else if len(args) > 1 {
		e = e.With(provider.GetFieldKeysSpec().GetMessage(), args)
	}

	return func() {
		helper()
		instance.Unwrap().Log(e, skipFrames+1)
	}, helper

}

func (instance *loggerImpl) DoLogf(level level.Level, skipFrames uint16, format string, args ...interface{}) {
	l, helper := instance.doLogf(level, skipFrames+1, format, args...)
	helper()
	l()
}

func (instance *loggerImpl) doLogf(level level.Level, skipFrames uint16, format string, args ...interface{}) (doLog, helper func()) {
	delegate := instance.Unwrap()
	helper = helperOf(delegate)
	if !delegate.IsLevelEnabled(level) {
		return func() {}, helper
	}
	provider := instance.GetProvider()
	e := instance.NewEventWithFields(level, instance.fields).
		Withf(provider.GetFieldKeysSpec().GetMessage(), format, args...)

	return func() {
		helper()
		instance.Unwrap().Log(e, skipFrames+1)
	}, helper
}

func (instance *loggerImpl) Trace(args ...interface{}) {
	l, helper := instance.log(level.Trace, args...)
	helper()
	l()
}

func (instance *loggerImpl) Tracef(format string, args ...interface{}) {
	l, helper := instance.logf(level.Trace, format, args...)
	helper()
	l()
}

func (instance *loggerImpl) IsTraceEnabled() bool {
	return instance.IsLevelEnabled(level.Trace)
}

func (instance *loggerImpl) Debug(args ...interface{}) {
	l, helper := instance.log(level.Debug, args...)
	helper()
	l()
}

func (instance *loggerImpl) Debugf(format string, args ...interface{}) {
	l, helper := instance.logf(level.Debug, format, args...)
	helper()
	l()
}

func (instance *loggerImpl) IsDebugEnabled() bool {
	return instance.IsLevelEnabled(level.Debug)
}

func (instance *loggerImpl) Info(args ...interface{}) {
	l, helper := instance.log(level.Info, args...)
	helper()
	l()
}

func (instance *loggerImpl) Infof(format string, args ...interface{}) {
	l, helper := instance.logf(level.Info, format, args...)
	helper()
	l()
}

func (instance *loggerImpl) IsInfoEnabled() bool {
	return instance.IsLevelEnabled(level.Info)
}

func (instance *loggerImpl) Warn(args ...interface{}) {
	l, helper := instance.log(level.Warn, args...)
	helper()
	l()
}

func (instance *loggerImpl) Warnf(format string, args ...interface{}) {
	l, helper := instance.logf(level.Warn, format, args...)
	helper()
	l()
}

func (instance *loggerImpl) IsWarnEnabled() bool {
	return instance.IsLevelEnabled(level.Warn)
}

func (instance *loggerImpl) Error(args ...interface{}) {
	l, helper := instance.log(level.Error, args...)
	helper()
	l()
}

func (instance *loggerImpl) Errorf(format string, args ...interface{}) {
	l, helper := instance.logf(level.Error, format, args...)
	helper()
	l()
}

func (instance *loggerImpl) IsErrorEnabled() bool {
	return instance.IsLevelEnabled(level.Error)
}

func (instance *loggerImpl) Fatal(args ...interface{}) {
	l, helper := instance.log(level.Fatal, args...)
	helper()
	l()
}

func (instance *loggerImpl) Fatalf(format string, args ...interface{}) {
	l, helper := instance.logf(level.Fatal, format, args...)
	helper()
	l()
}

func (instance *loggerImpl) IsFatalEnabled() bool {
	return instance.IsLevelEnabled(level.Fatal)
}

func (instance *loggerImpl) With(name string, value interface{}) Logger {
	return &loggerImpl{
		coreProvider: instance.coreProvider,
		fields:       instance.fields.With(name, value),
	}
}

func (instance *loggerImpl) Withf(name string, format string, args ...interface{}) Logger {
	return instance.With(name, fields.LazyFormat(format, args...))
}

func (instance *loggerImpl) WithError(err error) Logger {
	return instance.With(instance.GetProvider().GetFieldKeysSpec().GetError(), err)
}

func (instance *loggerImpl) WithAll(of map[string]interface{}) Logger {
	return &loggerImpl{
		coreProvider: instance.coreProvider,
		fields:       instance.fields.WithAll(of),
	}
}

func (instance *loggerImpl) Without(keys ...string) Logger {
	return &loggerImpl{
		coreProvider: instance.coreProvider,
		fields:       instance.fields.Without(keys...),
	}
}

func (instance *loggerImpl) Helper() func() {
	return helperOf(instance.Unwrap())
}

func helperOf(instance CoreLogger) func() {
	if wh, ok := instance.(interface {
		Helper() func()
	}); ok {
		return wh.Helper()
	}
	return func() {}
}
