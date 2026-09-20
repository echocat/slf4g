package support

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// EscapeNonGraphic replaces non-graphic runes with visible Go escape sequences.
func EscapeNonGraphic(value string, allowNewlines bool) string {
	firstUnsafe := -1
	for offset := 0; offset < len(value); {
		r, size := utf8.DecodeRuneInString(value[offset:])
		if r == utf8.RuneError && size == 1 || !unicode.IsGraphic(r) && (!allowNewlines || r != '\n') {
			firstUnsafe = offset
			break
		}
		offset += size
	}
	if firstUnsafe < 0 {
		return value
	}

	var result strings.Builder
	result.Grow(len(value))
	result.WriteString(value[:firstUnsafe])

	const hex = "0123456789abcdef"
	for offset := firstUnsafe; offset < len(value); {
		r, size := utf8.DecodeRuneInString(value[offset:])
		if r == utf8.RuneError && size == 1 {
			b := value[offset]
			result.WriteString(`\x`)
			result.WriteByte(hex[b>>4])
			result.WriteByte(hex[b&0x0f])
			offset++
			continue
		}
		if unicode.IsGraphic(r) || allowNewlines && r == '\n' {
			result.WriteRune(r)
		} else {
			quoted := strconv.QuoteRuneToGraphic(r)
			result.WriteString(quoted[1 : len(quoted)-1])
		}
		offset += size
	}

	return result.String()
}
