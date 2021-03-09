package functions

import (
	"fmt"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_EnsureWidth(t *testing.T) {
	cases := []struct {
		given                string
		givenMaxWidth        int32
		givenCutOffIfTooLong bool
		expected             string
	}{{
		given:         "hello, world",
		givenMaxWidth: 0,
		expected:      "hello, world",
	}, {
		given:                "hello, world",
		givenMaxWidth:        12,
		givenCutOffIfTooLong: true,
		expected:             "hello, world",
	}, {
		given:                "hello, world",
		givenMaxWidth:        -12,
		givenCutOffIfTooLong: true,
		expected:             "hello, world",
	}, {
		given:                "hello, world",
		givenMaxWidth:        11,
		givenCutOffIfTooLong: true,
		expected:             "hello, worl",
	}, {
		given:                "hello, world",
		givenMaxWidth:        11,
		givenCutOffIfTooLong: false,
		expected:             "hello, world",
	}, {
		given:                "hello, world",
		givenMaxWidth:        -11,
		givenCutOffIfTooLong: true,
		expected:             "hello, worl",
	}, {
		given:                "hello, world",
		givenMaxWidth:        -11,
		givenCutOffIfTooLong: false,
		expected:             "hello, world",
	}, {
		given:                "hello, world",
		givenMaxWidth:        13,
		givenCutOffIfTooLong: true,
		expected:             "hello, world ",
	}, {
		given:                "hello, world",
		givenMaxWidth:        -13,
		givenCutOffIfTooLong: true,
		expected:             " hello, world",
	}}

	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			actual := EnsureWidth(c.givenMaxWidth, c.givenCutOffIfTooLong, c.given)

			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}
