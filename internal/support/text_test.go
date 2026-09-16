package support

import (
	"testing"

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
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := EscapeNonGraphic(c.given, c.allowNewlines)

			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}
