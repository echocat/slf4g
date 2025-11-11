//go:build !mock && appengine

package color

import (
	"io"
)

func prepareForColors(w io.Writer) (bool, error) {
	return true, nil
}
