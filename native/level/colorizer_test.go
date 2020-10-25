package level

import (
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
)

func Test_ColorizerMap_ColorizeByLevel(t *testing.T) {
	instance := ColorizerMap{
		level.Level(1): `[31;1m`,
		level.Level(2): `[32;1m`,
		level.Level(3): `[33;1m`,
	}

	cases := []struct {
		givenLevel level.Level
		expected   string
	}{
		{level.Level(1), `[31;1mfoo[0m`},
		{level.Level(2), `[32;1mfoo[0m`},
		{level.Level(3), `[33;1mfoo[0m`},
		{level.Level(666), `[37;1mfoo[0m`},
	}

	for _, c := range cases {
		name, _ := DefaultNames.ToName(c.givenLevel)
		t.Run(name, func(t *testing.T) {
			actual := instance.ColorizeByLevel(c.givenLevel, "foo")

			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}

func Test_NewColorizerFacade(t *testing.T) {
	givenColorizer := ColorizerMap{level.Level(666): "foo"}

	actual := NewColorizerFacade(func() Colorizer {
		return givenColorizer
	})

	assert.ToBeEqual(t, givenColorizer, actual.(colorizerFacade)())
}

func Test_colorizerFacade_ColorizeByLevel(t *testing.T) {
	givenLevel := level.Warn
	givenColorizer := mockColorizer(func(actualLevel level.Level, actualInput string) string {
		assert.ToBeEqual(t, givenLevel, actualLevel)
		assert.ToBeEqual(t, "foo", actualInput)
		return "bar"
	})
	instance := colorizerFacade(func() Colorizer {
		return givenColorizer
	})

	actual := instance.ColorizeByLevel(givenLevel, "foo")

	assert.ToBeEqual(t, "bar", actual)
}

func Test_NoopColorizer(t *testing.T) {
	actual := NoopColorizer()

	assert.ToBeSame(t, noopColorizerV, actual)
}

func Test_noopColorizer_ColorizeByLevel(t *testing.T) {
	instance := &noopColorizer{}

	actual := instance.ColorizeByLevel(level.Warn, "foo")

	assert.ToBeEqual(t, "foo", actual)
}

type mockColorizer func(lvl level.Level, input string) string

func (instance mockColorizer) ColorizeByLevel(lvl level.Level, input string) string {
	return instance(lvl, input)
}
