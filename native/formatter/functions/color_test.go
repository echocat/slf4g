package functions

import (
	"fmt"
	"testing"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/native/hints"

	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/native/color"
	nlevel "github.com/echocat/slf4g/native/level"
)

func Test_ColorizeByLevel_enabled(t *testing.T) {
	colorizer := nlevel.DefaultColorizer

	for _, l := range level.GetProvider().GetLevels() {
		t.Run(fmt.Sprint(l), func(t *testing.T) {
			expected := colorizer.ColorizeByLevel(l, "foo")
			actual := ColorizeByLevel(l, mockHintColorMode(color.ModeAlways), "foo")

			assert.ToBeEqual(t, expected, actual)
		})
	}
}

func Test_ColorizeByLevel_disabled(t *testing.T) {
	for _, l := range level.GetProvider().GetLevels() {
		t.Run(fmt.Sprint(l), func(t *testing.T) {
			actual := ColorizeByLevel(l, mockHintColorMode(color.ModeNever), "foo")

			assert.ToBeEqual(t, "foo", actual)
		})
	}
}

func Test_Colorize(t *testing.T) {
	cases := []struct {
		givenColorCode string
		givenText      string
		shouldColorize bool
		expected       string
	}{{
		givenColorCode: "15;1",
		givenText:      "hello, world",
		shouldColorize: true,
		expected:       `[15;1mhello, world[0m`,
	}, {
		givenColorCode: "1",
		givenText:      "hello, world",
		shouldColorize: true,
		expected:       `[1mhello, world[0m`,
	}, {
		givenColorCode: "15;1",
		givenText:      "hello, world",
		shouldColorize: false,
		expected:       `hello, world`,
	}}

	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var h hints.Hints
			if c.shouldColorize {
				h = mockHintColorMode(color.ModeAlways)
			}
			actual := Colorize(c.givenColorCode, h, c.givenText)

			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}

func Test_ShouldColorize(t *testing.T) {
	cases := []struct {
		name     string
		given    hints.Hints
		expected bool
	}{{
		name:     "nothing given",
		given:    nil,
		expected: false,
	}, {
		name:     "supported assumed given",
		given:    mockHintColorsSupport(color.SupportedAssumed),
		expected: true,
	}, {
		name:     "mode always given",
		given:    mockHintColorMode(color.ModeAlways),
		expected: true,
	}, {
		name:     "supported assumed and mode never given",
		given:    mockHintColorsCombined{mockHintColorsSupport(color.SupportedAssumed), mockHintColorMode(color.ModeNever)},
		expected: false,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := ShouldColorize(c.given)

			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}

func Test_LevelColorizer_fromHints(t *testing.T) {
	givenLevelColorizer := nlevel.NewColorizerFacade(func() nlevel.Colorizer {
		panic("should never be called")
	})

	actual := LevelColorizer(mockHintLevelColorizer{givenLevelColorizer})

	assert.ToBeSame(t, givenLevelColorizer, actual)
}

func Test_LevelColorizer_default(t *testing.T) {
	actual := LevelColorizer(nil)

	assert.ToBeEqual(t, nlevel.DefaultColorizer, actual)
}

func Test_LevelColorizer_noop(t *testing.T) {
	old := nlevel.DefaultColorizer
	defer func() {
		nlevel.DefaultColorizer = old
	}()
	nlevel.DefaultColorizer = nil

	actual := LevelColorizer(nil)

	assert.ToBeSame(t, nlevel.NoopColorizer(), actual)
}

type mockHintColorsCombined struct {
	mockHintColorsSupport
	mockHintColorMode
}

type mockHintColorsSupport color.Supported

func (instance mockHintColorsSupport) IsColorSupported() color.Supported {
	return color.Supported(instance)
}

type mockHintColorMode color.Mode

func (instance mockHintColorMode) ColorMode() color.Mode {
	return color.Mode(instance)
}

type mockHintLevelColorizer struct {
	value nlevel.Colorizer
}

func (instance mockHintLevelColorizer) LevelColorizer() nlevel.Colorizer {
	return instance.value
}
