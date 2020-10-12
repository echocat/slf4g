package log

func newMockLogger(name string) *mockLogger {
	return &mockLogger{
		mockCoreLogger: newMockCoreLogger(name),
	}
}

type mockLogger struct {
	*mockCoreLogger
}

func (instance *mockLogger) Trace(...interface{}) {
	panic("not implemented in tests")
}

func (instance *mockLogger) Tracef(string, ...interface{}) {
	panic("not implemented in tests")
}

func (instance *mockLogger) IsTraceEnabled() bool {
	panic("not implemented in tests")
}

func (instance *mockLogger) Debug(...interface{}) {
	panic("not implemented in tests")
}

func (instance *mockLogger) Debugf(string, ...interface{}) {
	panic("not implemented in tests")
}

func (instance *mockLogger) IsDebugEnabled() bool {
	panic("not implemented in tests")
}

func (instance *mockLogger) Info(...interface{}) {
	panic("not implemented in tests")
}

func (instance *mockLogger) Infof(string, ...interface{}) {
	panic("not implemented in tests")
}

func (instance *mockLogger) IsInfoEnabled() bool {
	panic("not implemented in tests")
}

func (instance *mockLogger) Warn(...interface{}) {
	panic("not implemented in tests")
}

func (instance *mockLogger) Warnf(string, ...interface{}) {
	panic("not implemented in tests")
}

func (instance *mockLogger) IsWarnEnabled() bool {
	panic("not implemented in tests")
}

func (instance *mockLogger) Error(...interface{}) {
	panic("not implemented in tests")
}

func (instance *mockLogger) Errorf(string, ...interface{}) {
	panic("not implemented in tests")
}

func (instance *mockLogger) IsErrorEnabled() bool {
	panic("not implemented in tests")
}

func (instance *mockLogger) Fatal(...interface{}) {
	panic("not implemented in tests")
}

func (instance *mockLogger) Fatalf(string, ...interface{}) {
	panic("not implemented in tests")
}

func (instance *mockLogger) IsFatalEnabled() bool {
	panic("not implemented in tests")
}

func (instance *mockLogger) With(string, interface{}) Logger {
	panic("not implemented in tests")
}

func (instance *mockLogger) Withf(string, string, ...interface{}) Logger {
	panic("not implemented in tests")
}

func (instance *mockLogger) WithError(error) Logger {
	panic("not implemented in tests")
}

func (instance *mockLogger) WithAll(map[string]interface{}) Logger {
	panic("not implemented in tests")
}

func (instance *mockLogger) Without(...string) Logger {
	panic("not implemented in tests")
}

func newWrappingLogger(in Logger) *wrappingLogger {
	return &wrappingLogger{in}
}

type wrappingLogger struct {
	Logger
}

func (instance *wrappingLogger) Unwrap() Logger {
	return instance.Logger
}
