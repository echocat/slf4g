package log

import (
	"fmt"
	"github.com/echocat/slf4g/fields"
)

func NewLogger(cl CoreLogger) Logger {
	return &loggerImpl{
		CoreLogger: cl,
		fields:     fields.Empty(),
	}
}

type loggerImpl struct {
	CoreLogger
	fields fields.Fields
}

func (instance *loggerImpl) logM(level Level, message *string) {
	f := instance.fields
	if message != nil {
		f = f.With(instance.GetProvider().GetFieldKeys().GetMessage(), *message)
	}
	instance.LogEvent(NewEvent(level, f, 3))
}

func (instance *loggerImpl) log(level Level, args ...interface{}) {
	var message *string
	if len(args) > 0 {
		str := fmt.Sprint(args...)
		message = &str
	}

	instance.logM(level, message)
}

func (instance *loggerImpl) logf(level Level, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)

	instance.logM(level, &message)
}

func (instance *loggerImpl) Log(level Level, args ...interface{}) {
	instance.log(level, args...)
}

func (instance *loggerImpl) Logf(level Level, format string, args ...interface{}) {
	instance.logf(level, format, args...)
}

func (instance *loggerImpl) getMessageKey() string {
	return instance.GetProvider().GetFieldKeys().GetMessage()
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

func (instance *loggerImpl) Panic(args ...interface{}) {
	instance.log(LevelPanic, args...)
}

func (instance *loggerImpl) Panicf(format string, args ...interface{}) {
	instance.logf(LevelPanic, format, args...)
}

func (instance *loggerImpl) IsPanicEnabled() bool {
	return instance.IsLevelEnabled(LevelPanic)
}

func (instance *loggerImpl) Print(args ...interface{}) {
	instance.log(LevelInfo, args...)
}

func (instance *loggerImpl) Printf(format string, args ...interface{}) {
	instance.logf(LevelInfo, format, args...)
}

func (instance *loggerImpl) Println(args ...interface{}) {
	instance.log(LevelInfo, args...)
}

func (instance *loggerImpl) Fatalln(args ...interface{}) {
	instance.log(LevelFatal, args...)
}

func (instance *loggerImpl) Panicln(args ...interface{}) {
	instance.log(LevelPanic, args...)
}

func (instance *loggerImpl) With(name string, value interface{}) Logger {
	targetFields := instance.fields
	if targetFields != nil {
		targetFields = targetFields.With(name, value)
	} else {
		targetFields = fields.With(name, value)
	}
	return &loggerImpl{
		CoreLogger: instance.CoreLogger,
		fields:     targetFields,
	}
}

func (instance *loggerImpl) WithError(err error) Logger {
	return instance.With(instance.GetProvider().GetFieldKeys().GetError(), err)
}

func (instance *loggerImpl) WithFields(fields fields.Fields) Logger {
	targetFields := instance.fields
	if targetFields != nil {
		targetFields = targetFields.WithFields(fields)
	} else {
		targetFields = fields
	}
	return &loggerImpl{
		CoreLogger: instance.CoreLogger,
		fields:     targetFields,
	}
}
