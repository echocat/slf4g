package log

import (
	"github.com/echocat/slf4g/level"
)

type testLogger struct {
	testCoreLogger
}

func (instance *testLogger) Log(Event) {
	panic("not implemented in tests")
}

func (instance *testLogger) IsLevelEnabled(level.Level) bool {
	panic("not implemented in tests")
}

func (instance *testLogger) Trace(...interface{}) {
	panic("not implemented in tests")
}

func (instance *testLogger) Tracef(string, ...interface{}) {
	panic("not implemented in tests")
}

func (instance *testLogger) IsTraceEnabled() bool {
	panic("not implemented in tests")
}

func (instance *testLogger) Debug(...interface{}) {
	panic("not implemented in tests")
}

func (instance *testLogger) Debugf(string, ...interface{}) {
	panic("not implemented in tests")
}

func (instance *testLogger) IsDebugEnabled() bool {
	panic("not implemented in tests")
}

func (instance *testLogger) Info(...interface{}) {
	panic("not implemented in tests")
}

func (instance *testLogger) Infof(string, ...interface{}) {
	panic("not implemented in tests")
}

func (instance *testLogger) IsInfoEnabled() bool {
	panic("not implemented in tests")
}

func (instance *testLogger) Warn(...interface{}) {
	panic("not implemented in tests")
}

func (instance *testLogger) Warnf(string, ...interface{}) {
	panic("not implemented in tests")
}

func (instance *testLogger) IsWarnEnabled() bool {
	panic("not implemented in tests")
}

func (instance *testLogger) Error(...interface{}) {
	panic("not implemented in tests")
}

func (instance *testLogger) Errorf(string, ...interface{}) {
	panic("not implemented in tests")
}

func (instance *testLogger) IsErrorEnabled() bool {
	panic("not implemented in tests")
}

func (instance *testLogger) Fatal(...interface{}) {
	panic("not implemented in tests")
}

func (instance *testLogger) Fatalf(string, ...interface{}) {
	panic("not implemented in tests")
}

func (instance *testLogger) IsFatalEnabled() bool {
	panic("not implemented in tests")
}

func (instance *testLogger) With(string, interface{}) Logger {
	panic("not implemented in tests")
}

func (instance *testLogger) Withf(string, string, ...interface{}) Logger {
	panic("not implemented in tests")
}

func (instance *testLogger) WithError(error) Logger {
	panic("not implemented in tests")
}

func (instance *testLogger) WithAll(map[string]interface{}) Logger {
	panic("not implemented in tests")
}

func (instance *testLogger) Without(...string) Logger {
	panic("not implemented in tests")
}

type wrappingTestLogger struct {
	Logger
}

func (instance *wrappingTestLogger) Unwrap() Logger {
	return instance.Logger
}
