//go:build !mock && !appengine && !js && !windows && !nacl && !plan9
// +build !mock,!appengine,!js,!windows,!nacl,!plan9

package color

import (
	"io"
	"os"
)

func prepareForColors(w io.Writer) (bool, error) {
	switch v := w.(type) {
	case *os.File:
		return isTerminal(int(v.Fd()))
	default:
		return false, nil
	}
}
