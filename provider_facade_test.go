package log

import (
	"testing"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_providerFacade_GetName(t *testing.T) {
	givenProvider := newMockProvider("foo")
	instance := providerFacade(func() Provider { return givenProvider })

	actual := instance.GetName()

	assert.ToBeEqual(t, "foo", actual)
}

func Test_providerFacade_GetRootLogger(t *testing.T) {
	givenRootLogger := newMockLogger("root")
	givenProvider := newMockProvider("foo")
	givenProvider.rootProvider = func() Logger { return givenRootLogger }
	instance := providerFacade(func() Provider { return givenProvider })

	actual := instance.GetRootLogger()

	assert.ToBeOfType(t, &loggerImpl{}, actual)
	assert.ToBeSame(t, givenRootLogger, actual.(*loggerImpl).Unwrap())
}

func Test_providerFacade_GetLogger(t *testing.T) {
	givenLogger1 := newMockLogger("1")
	givenLogger2 := newMockLogger("2")
	givenProvider := newMockProvider("foo")
	givenProvider.provider = func(name string) Logger {
		switch name {
		case "1":
			return givenLogger1
		case "2":
			return givenLogger2
		default:
			panic(name)
		}
	}
	instance := providerFacade(func() Provider { return givenProvider })

	actual1 := instance.GetLogger("1")
	actual2 := instance.GetLogger("2")

	assert.ToBeOfType(t, &loggerImpl{}, actual1)
	assert.ToBeOfType(t, &loggerImpl{}, actual2)

	assert.ToBeSame(t, givenLogger1, actual1.(*loggerImpl).Unwrap())
	assert.ToBeSame(t, givenLogger2, actual2.(*loggerImpl).Unwrap())
}

func Test_providerFacade_GetAllLevels(t *testing.T) {
	givenProvider := newMockProvider("foo")
	givenProvider.levels = []level.Level{level.Info, level.Fatal}
	instance := providerFacade(func() Provider { return givenProvider })

	actual := instance.GetAllLevels()

	assert.ToBeEqual(t, givenProvider.levels, actual)
}

func Test_providerFacade_GetFieldKeysSpec(t *testing.T) {
	givenProvider := newMockProvider("foo")
	givenProvider.fieldKeysSpec = &mockFieldKeysSpec{"1", "2", "3", "4"}
	givenProvider.levels = []level.Level{level.Info, level.Fatal}
	instance := providerFacade(func() Provider { return givenProvider })

	actual := instance.GetFieldKeysSpec()

	assert.ToBeEqual(t, givenProvider.fieldKeysSpec, actual)
}
