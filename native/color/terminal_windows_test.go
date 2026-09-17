package color

import (
	"errors"
	"os"
	"syscall"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_enableVirtualTerminalProcessing_readsModeFromTarget(t *testing.T) {
	oldGetConsoleMode := getConsoleMode
	defer func() {
		getConsoleMode = oldGetConsoleMode
	}()

	expectedHandle := syscall.Handle(1234)
	expectedError := errors.New("expected")
	var actualHandle syscall.Handle
	getConsoleMode = func(handle syscall.Handle, _ *uint32) error {
		actualHandle = handle
		return expectedError
	}

	actual, actualErr := enableVirtualTerminalProcessing(os.NewFile(uintptr(expectedHandle), "expected"))

	assert.ToBeEqual(t, false, actual)
	assert.ToBeSame(t, expectedError, actualErr)
	assert.ToBeEqual(t, expectedHandle, actualHandle)
}
