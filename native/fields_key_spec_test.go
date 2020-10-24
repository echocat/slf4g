package native

import (
	"testing"

	"github.com/echocat/slf4g/fields"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_FieldKeysSpecImpl_GetLocation_specified(t *testing.T) {
	instance := &FieldKeysSpecImpl{Location: "foo"}

	actual := instance.GetLocation()

	assert.ToBeEqual(t, "foo", actual)
}

func Test_FieldKeysSpecImpl_GetLocation_default(t *testing.T) {
	instance := &FieldKeysSpecImpl{}

	actual := instance.GetLocation()

	assert.ToBeEqual(t, "location", actual)
}

func Test_NewFieldKeysSpecFacade(t *testing.T) {
	delegate := &FieldKeysSpecImpl{
		KeysSpecImpl: fields.KeysSpecImpl{
			Timestamp: "a",
			Message:   "b",
			Logger:    "c",
			Error:     "d",
		},
		Location: "e",
	}

	actual := NewFieldKeysSpecFacade(func() FieldKeysSpec {
		return delegate
	})

	assert.ToBeEqual(t, "a", actual.GetTimestamp())
	assert.ToBeEqual(t, "b", actual.GetMessage())
	assert.ToBeEqual(t, "c", actual.GetLogger())
	assert.ToBeEqual(t, "d", actual.GetError())
	assert.ToBeEqual(t, "e", actual.GetLocation())
}
