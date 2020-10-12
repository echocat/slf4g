package log

import (
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

func newMockProvider(name string) *mockProvider {
	return &mockProvider{
		name: name,
		fieldKeysSpec: &mockFieldKeysSpec{
			timestamp: "aTimestamp",
			message:   "aMessage",
			error:     "anError",
			logger:    "aLogger",
		},
	}
}

type mockProvider struct {
	rootProvider  func() Logger
	provider      func(name string) Logger
	name          string
	fieldKeysSpec fields.KeysSpec
	levels        level.Levels
}

func (instance *mockProvider) GetName() string {
	if v := instance.name; v != "" {
		return v
	}
	panic("not implemented in tests")
}

func (instance *mockProvider) GetRootLogger() Logger {
	if v := instance.rootProvider; v != nil {
		return v()
	}
	panic("not implemented in tests")
}

func (instance *mockProvider) GetLogger(name string) Logger {
	if v := instance.provider; v != nil {
		return v(name)
	}
	panic("not implemented in tests")
}

func (instance *mockProvider) GetAllLevels() level.Levels {
	if v := instance.levels; v != nil {
		return v
	}
	panic("not implemented in tests")
}

func (instance *mockProvider) GetFieldKeysSpec() fields.KeysSpec {
	if v := instance.fieldKeysSpec; v != nil {
		return v
	}
	panic("not implemented in tests")
}

func newWrappingProvider(in Provider) *wrappingProvider {
	return &wrappingProvider{in}
}

type wrappingProvider struct {
	Provider
}

func (instance *wrappingProvider) Unwrap() Provider {
	return instance.Provider
}

type mockFieldKeysSpec struct {
	timestamp string
	message   string
	error     string
	logger    string
}

func (instance *mockFieldKeysSpec) GetTimestamp() string {
	if v := instance.timestamp; v != "" {
		return v
	}
	panic("not implemented in tests")
}

func (instance *mockFieldKeysSpec) GetMessage() string {
	if v := instance.message; v != "" {
		return v
	}
	panic("not implemented in tests")
}

func (instance *mockFieldKeysSpec) GetError() string {
	if v := instance.error; v != "" {
		return v
	}
	panic("not implemented in tests")
}

func (instance *mockFieldKeysSpec) GetLogger() string {
	if v := instance.logger; v != "" {
		return v
	}
	panic("not implemented in tests")
}
