package color

import (
	"io"
	"os"
	"syscall"
)

var (
	kernel32Dll    = syscall.NewLazyDLL("Kernel32.dll")
	setConsoleMode = kernel32Dll.NewProc("SetConsoleMode")
	getConsoleMode = syscall.GetConsoleMode
)

func enableVirtualTerminalProcessing(w io.Writer) (bool, error) {
	switch v := w.(type) {
	case *os.File:
		handle := syscall.Handle(v.Fd())
		var mode uint32
		if err := getConsoleMode(handle, &mode); err != nil {
			return false, err
		}

		if ret, _, err := setConsoleMode.Call(uintptr(handle), uintptr(mode|0x4)); ret == 0 {
			return false, err
		}

		return true, nil
	default:
		return false, nil
	}
}
