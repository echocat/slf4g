package support

import (
	"strconv"
	"strings"
	"unicode"
)

// EscapeNonGraphic replaces non-graphic runes with visible Go escape sequences.
func EscapeNonGraphic(value string, allowNewlines bool) string {
	firstUnsafe := strings.IndexFunc(value, func(r rune) bool {
		return !unicode.IsGraphic(r) && (!allowNewlines || r != '\n')
	})
	if firstUnsafe < 0 {
		return value
	}

	var result strings.Builder
	result.Grow(len(value))
	result.WriteString(value[:firstUnsafe])

	for _, r := range value[firstUnsafe:] {
		if unicode.IsGraphic(r) || allowNewlines && r == '\n' {
			result.WriteRune(r)
			continue
		}

		quoted := strconv.QuoteRuneToGraphic(r)
		result.WriteString(quoted[1 : len(quoted)-1])
	}

	return result.String()
}
