package formatter

import (
	"fmt"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_stringNeedsQuoting(t *testing.T) {
	cases := []struct {
		given    string
		expected bool
	}{
		{"abc", false},
		{"aBc", false},
		{"{aBc}", false},
		{"{[]}", false},
		{"foo@bar.com", false},
		{"-.,_:;!/\\@^+#()[]{}", false},
		{"hello\"", true},
		{"hello$", true},
		{"hello&", true},
		{"hello*", true},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, c.given), func(t *testing.T) {
			actual := stringNeedsQuoting(c.given)

			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}
