package fields

import (
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_asMap_withMapped(t *testing.T) {
	given := mapped{"foo": 1, "bar": 2}

	actual := asMap(given)

	assert.ToBeEqual(t, given, actual)
}

func Test_asMap_withPMapped(t *testing.T) {
	given := mapped{"foo": 1, "bar": 2}

	actual := asMap(&given)

	assert.ToBeEqual(t, given, actual)
}

func Test_asMap_copies(t *testing.T) {
	given := &lineage{With("a", 1), With("b", 2)}

	actual := asMap(given)

	assert.ToBeEqual(t, mapped{"a": 1, "b": 2}, actual)
}
