package log

import (
	"github.com/echocat/slf4g/level"
)

func newMockCoreLogger(name string) *mockCoreLogger {
	return &mockCoreLogger{
		name:     name,
		provider: newMockProvider("mocked"),
	}
}

type mockCoreLogger struct {
	name         string
	provider     *mockProvider
	level        level.Level
	loggedEvents *[]Event
}

func (instance *mockCoreLogger) initLoggedEvents() {
	instance.loggedEvents = new([]Event)
}

func (instance *mockCoreLogger) loggedEvent(i int) Event {
	return (*instance.loggedEvents)[i]
}

func (instance *mockCoreLogger) Log(e Event) {
	if v := instance.loggedEvents; v != nil {
		*v = append(*v, e)
		return
	}
	panic("not implemented in tests")
}

func (instance *mockCoreLogger) IsLevelEnabled(l level.Level) bool {
	if v := instance.level; v != 0 {
		return instance.level.CompareTo(l) <= 0
	}
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
