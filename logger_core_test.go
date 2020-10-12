package log

import (
	"github.com/echocat/slf4g/level"
)

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

type wrappingMockCoreLogger struct {
	CoreLogger
}

func (instance *wrappingMockCoreLogger) Unwrap() CoreLogger {
	return instance.CoreLogger
}
