package level

import (
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_providerFacade_GetName(t *testing.T) {
	inner := &testProviderForFacade{}
	called := false

	instance := providerFacade(func() Provider {
		called = true
		return inner
	})

	actual := instance.GetName()

	assert.ToBeEqual(t, "foo", actual)
	assert.ToBeEqual(t, true, called)
}

func Test_providerFacade_GetLevels(t *testing.T) {
	inner := &testProviderForFacade{}
	called := false

	instance := providerFacade(func() Provider {
		called = true
		return inner
	})

	actual := instance.GetLevels()

	assert.ToBeEqual(t, Levels{Info, Error}, actual)
	assert.ToBeEqual(t, true, called)
}

type testProviderForFacade struct{}

func (instance *testProviderForFacade) GetName() string {
	return "foo"
}

func (instance *testProviderForFacade) GetLevels() Levels {
	return Levels{Info, Error}
}
