package level

import (
	"fmt"
	"strings"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
)

func Test_DefaultFormatter(t *testing.T) {
	cases := []struct {
		given    level.Level
		expected string
	}{
		{level.Trace, "TRACE"},
		{level.Debug, "DEBUG"},
		{level.Info, " INFO"},
		{level.Warn, " WARN"},
		{level.Error, "ERROR"},
		{level.Fatal, "FATAL"},
		{level.Level(666), "  666"},
	}

	for _, c := range cases {
		t.Run(strings.TrimSpace(c.expected), func(t *testing.T) {
			actual := DefaultFormatter.Format(c.given)

			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}

func Test_FormatterFunc_Format(t *testing.T) {
	given := level.Level(666)

	instance := FormatterFunc(func(l level.Level) string {
		return fmt.Sprintf("x%dx", l)
	})

	actual := instance.Format(given)

	assert.ToBeEqual(t, "x666x", actual)
}
