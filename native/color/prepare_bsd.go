//go:build (!mock && darwin) || dragonfly || freebsd || netbsd || openbsd
// +build !mock,darwin dragonfly freebsd netbsd openbsd

package color

import "golang.org/x/sys/unix"

func isTerminal(fd int) (bool, error) {
	_, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	if err == unix.ENOTTY {
		return false, nil
	}
	return err == nil, err
}
