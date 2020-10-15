// +build !mock

package color

import (
	"golang.org/x/sys/unix"
)

// IsTerminal returns true if the given file descriptor is a terminal.
func isTerminal(fd int) (bool, error) {
	_, err := unix.IoctlGetTermios(fd, unix.TCGETA)
	if err == unix.ENOTTY {
		return false, nil
	}
	return err == nil, err
}
