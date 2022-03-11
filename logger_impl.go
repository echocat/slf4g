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

func (instance *loggerImpl) Log(event Event, skipFrames uint16) {
	instance.Unwrap().Log(event, skipFrames+1)
}

func (instance *loggerImpl) NewEvent(l level.Level, values map[string]interface{}) Event {
	if v := instance.fields; v != nil {
		return instance.NewEventWithFields(l, v.WithAll(values))
	}
	return instance.Unwrap().NewEvent(l, values)
}

func (instance *loggerImpl) NewEventWithFields(l level.Level, f fields.ForEachEnabled) Event {
	target := instance.Unwrap()
	if wf, ok := target.(interface {
		NewEventWithFields(l level.Level, f fields.ForEachEnabled) Event
	}); ok {
		return wf.NewEventWithFields(l, f)
	}
	asMap, err := fields.AsMap(f)
	if err != nil {
		panic(fmt.Errorf("cannot make a map out of %v: %w", f, err))
	}
	return target.NewEvent(l, asMap)
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

func (instance *loggerImpl) log(level level.Level, args ...interface{}) {
	if !instance.IsLevelEnabled(level) {
		return
	}
	provider := instance.GetProvider()
	e := instance.NewEventWithFields(level, instance.fields)

	if len(args) == 1 {
		e = e.With(provider.GetFieldKeysSpec().GetMessage(), args[0])
	} else if len(args) > 1 {
		e = e.With(provider.GetFieldKeysSpec().GetMessage(), args)
	}

	instance.Unwrap().Log(e, 2)
}

func (instance *loggerImpl) logf(level level.Level, format string, args ...interface{}) {
	if !instance.IsLevelEnabled(level) {
		return
	}
	provider := instance.GetProvider()
	e := instance.NewEventWithFields(level, instance.fields).
		Withf(provider.GetFieldKeysSpec().GetMessage(), format, args...)

	instance.Unwrap().Log(e, 2)
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
