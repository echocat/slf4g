package log

import (
	"github.com/echocat/slf4g/level"
)

func newMockCoreLogger(name string) *mockCoreLogger {
	return &mockCoreLogger{name: name}
}

type mockCoreLogger struct {
	name     string
	provider Provider
}

func (instance *mockCoreLogger) Log(Event) {
	panic("not implemented in tests")
}

func (instance *mockCoreLogger) IsLevelEnabled(level.Level) bool {
	panic("not implemented in tests")
}

func (instance *mockCoreLogger) GetName() string {
	if v := instance.name; v != "" {
		return v
	}
	panic("not implemented in tests")
}

func (instance *mockCoreLogger) GetProvider() Provider {
	if v := instance.provider; v != nil {
		return v
	}
	panic("not implemented in tests")
}

func newWrappingCoreLogger(in CoreLogger) *wrappingCoreLogger {
	return &wrappingCoreLogger{in}
}

type wrappingCoreLogger struct {
	CoreLogger
}

func (instance *wrappingCoreLogger) Unwrap() CoreLogger {
	return instance.CoreLogger
}
