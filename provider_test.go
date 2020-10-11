package log

import (
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

type testProvider struct {
	name string
}

func (instance *testProvider) GetName() string {
	return instance.name
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
	panic("not implemented in tests")
}

type wrappingTestProvider struct {
	Provider
}

func (instance *wrappingTestProvider) Unwrap() Provider {
	return instance.Provider
}
