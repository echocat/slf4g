package support

import (
	"testing"
	"unicode/utf8"

	"github.com/echocat/slf4g/internal/test/assert"
)

func TestEscapeNonGraphic(t *testing.T) {
	cases := []struct {
		name          string
		given         string
		allowNewlines bool
		expected      string
	}{
		{"plain", "hello, world", false, "hello, world"},
		{"controls", "a\r\n\tb\x1b[2J\u202ec", false, `a\r\n\tb\x1b[2J\u202ec`},
		{"allowed newline", "a\r\nb\x1b[2J", true, "a\\r\nb\\x1b[2J"},
		{"invalid byte", "a\xffb", false, `a\xffb`},
		{"truncated sequence", "a\xe2\x82", false, `a\xe2\x82`},
		{"overlong sequence", "a\xc0\xafb", false, `a\xc0\xafb`},
		{"surrogate sequence", "a\xed\xa0\x80b", false, `a\xed\xa0\x80b`},
		{"control and invalid byte", "a\x00\xff\n", false, `a\x00\xff\n`},
		{"invalid byte with allowed newline", "a\n\xffb", true, "a\n\\xffb"},
		{"replacement character", "a\ufffdb", false, "a\ufffdb"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := EscapeNonGraphic(c.given, c.allowNewlines)

			assert.ToBeEqual(t, c.expected, actual)
			assert.ToBeEqual(t, true, utf8.ValidString(actual))
		})
	}
}

func TestEscapeNonGraphic_doesNotAllocateForSafeValue(t *testing.T) {
	const given = "hello, 世界"
	var actual string

	allocations := testing.AllocsPerRun(100, func() {
		actual = EscapeNonGraphic(given, false)
	})

	assert.ToBeEqual(t, given, actual)
	assert.ToBeEqual(t, float64(0), allocations)
}
