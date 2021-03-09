package execution

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_Execute(t *testing.T) {
	expectedErr := errors.New("expected")
	executed1, executed2 := false, false
	execution1 := func() error {
		executed1 = true
		return nil
	}
	execution2 := func() error {
		executed2 = true
		return expectedErr
	}

	actualErr := Execute(execution1, execution2)

	assert.ToBeSame(t, expectedErr, actualErr)
	assert.ToBeEqual(t, true, executed1)
	assert.ToBeEqual(t, true, executed2)
}

func Test_Join(t *testing.T) {
	expectedErr := errors.New("expected")
	executed1, executed2 := false, false
	execution1 := func() error {
		executed1 = true
		return nil
	}
	execution2 := func() error {
		executed2 = true
		return expectedErr
	}

	actual := Join(execution1, execution2)
	actualErr := actual()

	assert.ToBeSame(t, expectedErr, actualErr)
	assert.ToBeEqual(t, true, executed1)
	assert.ToBeEqual(t, true, executed2)
}
