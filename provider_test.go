package log

import (
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

type testProvider struct {
	name          string
	fieldKeysSpec fields.KeysSpec
}

func (instance *testProvider) GetName() string {
	if v := instance.name; v != "" {
		return v
	}
	panic("not implemented in tests")
}

func (instance *testProvider) GetRootLogger() Logger {
	panic("not implemented in tests")
}

func (instance *testProvider) GetLogger(string) Logger {
	panic("not implemented in tests")
}

func (instance *testProvider) GetAllLevels() level.Levels {
	panic("not implemented in tests")
}

func (instance *testProvider) GetFieldKeysSpec() fields.KeysSpec {
	if v := instance.fieldKeysSpec; v != nil {
		return v
	}
	panic("not implemented in tests")
}

type wrappingTestProvider struct {
	Provider
}

func (instance *wrappingTestProvider) Unwrap() Provider {
	return instance.Provider
}

type testFieldKeysSpec struct {
	timestamp string
	message   string
	error     string
	logger    string
}

func (instance *testFieldKeysSpec) GetTimestamp() string {
	if v := instance.timestamp; v != "" {
		return v
	}
	panic("not implemented in tests")
}

func (instance *testFieldKeysSpec) GetMessage() string {
	if v := instance.message; v != "" {
		return v
	}
	panic("not implemented in tests")
}

func (instance *testFieldKeysSpec) GetError() string {
	if v := instance.error; v != "" {
		return v
	}
	panic("not implemented in tests")
}

func (instance *testFieldKeysSpec) GetLogger() string {
	if v := instance.logger; v != "" {
		return v
	}
	panic("not implemented in tests")
}
