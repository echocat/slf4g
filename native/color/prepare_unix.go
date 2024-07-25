//go:build (!mock && linux) || aix
// +build !mock,linux aix

package color

import "golang.org/x/sys/unix"

func isTerminal(fd int) (bool, error) {
	_, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err == unix.ENOTTY {
		return false, nil
	}
	return err == nil, err
}
