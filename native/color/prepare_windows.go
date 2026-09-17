//go:build !mock

package color

import (
	"io"
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

func prepareForColors(w io.Writer) (bool, error) {
	switch v := w.(type) {
	case *os.File:
		var mode uint32
		if err := syscall.GetConsoleMode(syscall.Handle(v.Fd()), &mode); err == windows.ERROR_INVALID_HANDLE {
			return false, nil
		} else if err != nil {
			return false, err
		}
		return enableVirtualTerminalProcessing(w)
	default:
		return false, nil
	}
}
