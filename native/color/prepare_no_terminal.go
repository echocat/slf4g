//go:build (!mock && js) || nacl || plan9
// +build !mock,js nacl plan9

package color

import (
	"io"
)

func prepareForColors(w io.Writer) (bool, error) {
	return false, nil
}
