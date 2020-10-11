package log

import (
	"github.com/echocat/slf4g/level"
)

type testCoreLogger struct {
	name     string
	provider Provider
}

func (instance *testCoreLogger) Log(Event) {
	panic("not implemented in tests")
}

func (instance *testCoreLogger) IsLevelEnabled(level.Level) bool {
	panic("not implemented in tests")
}

func (instance *testCoreLogger) GetName() string {
	if v := instance.name; v != "" {
		return v
	}
	panic("not implemented in tests")
}

func (instance *testCoreLogger) GetProvider() Provider {
	if v := instance.provider; v != nil {
		return v
	}
	panic("not implemented in tests")
}

type wrappingCoreTestLogger struct {
	CoreLogger
}

func (instance *wrappingCoreTestLogger) Unwrap() CoreLogger {
	return instance.CoreLogger
}
