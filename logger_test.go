package log

import (
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

func newMockLogger(name string) *mockLogger {
	return &mockLogger{NewLogger(newMockCoreLogger(name)).(*loggerImpl)}
}

type mockLogger struct {
	*loggerImpl
}

func (instance *mockLogger) setLevel(in level.Level) {
	instance.coreProvider().(*mockCoreLogger).level = in
}

func (instance *mockLogger) getProvider() *mockProvider {
	return instance.coreProvider().(*mockCoreLogger).provider
}

func (instance *mockLogger) getFieldKeysSpec() fields.KeysSpec {
	return instance.GetProvider().GetFieldKeysSpec()
}

func (instance *mockLogger) initLoggedEvents() {
	instance.coreProvider().(*mockCoreLogger).initLoggedEvents()
}

func (instance *mockLogger) loggedEvents() []Event {
	if v := instance.coreProvider().(*mockCoreLogger).loggedEvents; v != nil {
		return *v
	}
	return nil
}

func (instance *mockLogger) loggedEvent(i int) Event {
	return instance.loggedEvents()[i]
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
