package level

import (
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_defaultProvider_GetName(t *testing.T) {
	actual := defaultProviderV.GetName()

	assert.ToBeEqual(t, "default", actual)
}

func Test_defaultProvider_GetLevels(t *testing.T) {
	actual := defaultProviderV.GetLevels()

	assert.ToBeEqual(t, Levels{Trace, Debug, Info, Warn, Error, Fatal}, actual)
}
