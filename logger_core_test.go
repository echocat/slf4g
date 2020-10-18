package log

import (
	"github.com/echocat/slf4g/fields"
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

func (instance *mockCoreLogger) Log(e Event, _ uint16) {
	if v := instance.loggedEvents; v != nil {
		*v = append(*v, e)
		return
	}
	panic("not implemented in tests")
}

func (instance *mockCoreLogger) NewEvent(level level.Level, values map[string]interface{}) Event {
	return &fallbackEvent{
		provider: instance.provider,
		level:    level,
		fields:   fields.WithAll(values),
	}
}

func (instance *mockCoreLogger) Accepts(e Event) bool {
	return e != nil
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
