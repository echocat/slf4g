// +build linux aix

package color

import "golang.org/x/sys/unix"

func isTerminal(fd int) (bool, error) {
	_, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	return err == nil, err
}
