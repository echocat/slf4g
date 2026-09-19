package support

import (
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func TestIsNil(t *testing.T) {
	var (
		pointer  *int
		function func()
		mapping  map[string]string
		slice    []string
		channel  chan string
	)
	cases := []struct {
		name     string
		given    any
		expected bool
	}{
		{"nil", nil, true},
		{"pointer", pointer, true},
		{"function", function, true},
		{"map", mapping, true},
		{"slice", slice, true},
		{"channel", channel, true},
		{"value", 42, false},
		{"non-nil pointer", new(int), false},
		{"non-nil map", map[string]string{}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.ToBeEqual(t, c.expected, IsNil(c.given))
		})
	}
}
