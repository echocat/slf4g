// +build darwin dragonfly freebsd netbsd openbsd

package color

import "golang.org/x/sys/unix"

func isTerminal(fd int) (bool, error) {
	_, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	return err == nil, err
}
