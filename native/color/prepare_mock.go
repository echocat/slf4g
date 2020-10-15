// +build mock

package color

import "io"

var prepareCallback func(w io.Writer) (bool, error)

func prepareForColors(w io.Writer) (bool, error) {
	if cb := prepareCallback; cb != nil {
		return cb(w)
	}
	return false, nil
}
