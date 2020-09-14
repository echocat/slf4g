package color

import (
	"io"
	"os"
	"syscall"
)

var (
	kernel32Dll    = syscall.NewLazyDLL("Kernel32.dll")
	setConsoleMode = kernel32Dll.NewProc("SetConsoleMode")
)

func enableVirtualTerminalProcessing(w io.Writer) (bool, error) {
	switch v := w.(type) {
	case *os.File:
		var mode uint32
		if err := syscall.GetConsoleMode(syscall.Stdout, &mode); err != nil {
			return false, err
		}

		if ret, _, err := setConsoleMode.Call(v.Fd(), uintptr(mode|0x4)); ret == 0 {
			return false, err
		}

		return true, nil
	default:
		return false, nil
	}
}

func prepareForColors(w io.Writer) (bool, error) {
	switch v := w.(type) {
	case *os.File:
		var mode uint32
		if err := syscall.GetConsoleMode(syscall.Handle(v.Fd()), &mode); err != nil {
			return false, err
		}
		return enableVirtualTerminalProcessing(w)
	default:
		return false, nil
	}
}
