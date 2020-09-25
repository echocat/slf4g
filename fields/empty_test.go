package fields

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Empty(t *testing.T) {
	actual := Empty()

	assert.IsType(t, &empty{}, actual)
}

func Test_Empty_alwaysTheSame(t *testing.T) {
	actual1 := Empty()
	actual2 := Empty()

	assert.Same(t, actual1, actual2)
}
