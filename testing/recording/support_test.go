package recording

import (
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_BeTrue(t *testing.T) {
	actual1 := BeTrue()
	actual2 := BeTrue()

	assert.ToBeEqual(t, true, *actual1)
	assert.ToBeEqual(t, true, *actual2)
	assert.ToBeEqual(t, actual1, actual2)
	assert.ToBeNotSame(t, actual1, actual2)
}

func Test_BeFalse(t *testing.T) {
	actual1 := BeFalse()
	actual2 := BeFalse()

	assert.ToBeEqual(t, false, *actual1)
	assert.ToBeEqual(t, false, *actual2)
	assert.ToBeEqual(t, actual1, actual2)
	assert.ToBeNotSame(t, actual1, actual2)
}
