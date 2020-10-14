package fields

import (
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_NewKeysSpecFacade(t *testing.T) {
	givenSpec := &KeysSpecImpl{Message: "aMessage"}
	givenProvider := func() KeysSpec { return givenSpec }

	actual := NewKeysSpecFacade(givenProvider)

	assert.ToBeSame(t, keysSpecFacade(givenProvider), actual)
}

func Test_keysSpecFacade_GetTimestamp(t *testing.T) {
	givenSpec := &KeysSpecImpl{Timestamp: "xyz"}
	instance := keysSpecFacade(func() KeysSpec { return givenSpec })

	actual := instance.GetTimestamp()

	assert.ToBeEqual(t, "xyz", actual)
}

func Test_keysSpecFacade_GetMessage(t *testing.T) {
	givenSpec := &KeysSpecImpl{Message: "xyz"}
	instance := keysSpecFacade(func() KeysSpec { return givenSpec })

	actual := instance.GetMessage()

	assert.ToBeEqual(t, "xyz", actual)
}

func Test_keysSpecFacade_GetError(t *testing.T) {
	givenSpec := &KeysSpecImpl{Error: "xyz"}
	instance := keysSpecFacade(func() KeysSpec { return givenSpec })

	actual := instance.GetError()

	assert.ToBeEqual(t, "xyz", actual)
}

func Test_keysSpecFacade_GetLogger(t *testing.T) {
	givenSpec := &KeysSpecImpl{Logger: "xyz"}
	instance := keysSpecFacade(func() KeysSpec { return givenSpec })

	actual := instance.GetLogger()

	assert.ToBeEqual(t, "xyz", actual)
}
