package functions

import (
	"strings"
	"unicode/utf8"
)

// EnsureWidth will ensure that the given width is applied to the given string `of`.
// If the given string `of` is too short whitespaces will be added. If the string is too
// long and cutOffIfTooLong is true it will be cut of.
// If a negative value of `width` has nearly the same effect as it is positive. The only
// difference is, that in case of positive: whitespaces will be suffixed; in case
// of negative: whitespaces will be prefixed. If `width` is `0` this method will do nothing.
func EnsureWidth(width int32, cutOffIfTooLong bool, of string) string {
	l2r := true
	if width < 0 {
		width *= -1
		l2r = false
	}
	if width == 0 {
		return of
	}
	l := utf8.RuneCountInString(of)
	if l >= int(width) {
		if cutOffIfTooLong {
			return of[:width]
		}
		return of
	}
	if l2r {
		return of + strings.Repeat(" ", int(width)-l)
	}
	return strings.Repeat(" ", int(width)-l) + of
}
