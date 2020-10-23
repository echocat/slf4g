package level

import (
	"sort"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_Level_CompareTo(t *testing.T) {
	assert.ToBeEqual(t, 1000, Info.CompareTo(Debug))
	assert.ToBeEqual(t, -1000, Debug.CompareTo(Info))
	assert.ToBeEqual(t, 0, Info.CompareTo(Info))
	assert.ToBeEqual(t, -5000, Level(0).CompareTo(Error))
}

func Test_Levels_Sorting(t *testing.T) {
	instance := Levels{Fatal, Info, Debug, Warn, Error}

	sort.Sort(instance)

	assert.ToBeEqual(t, Levels{Debug, Info, Warn, Error, Fatal}, instance)
}

func Test_Levels_ToProvider(t *testing.T) {
	instance := Levels{Fatal, Info, Debug, Warn, Error}

	actual := instance.ToProvider("foo")

	assert.ToBeEqual(t, levelsAsProvider{instance, "foo"}, actual)
}

func Test_levelsAsProvider_GetName(t *testing.T) {
	instance := levelsAsProvider{Levels{}, "foo"}

	actual := instance.GetName()

	assert.ToBeEqual(t, "foo", actual)
}

func Test_levelsAsProvider_GetLevels(t *testing.T) {
	givenLevels := Levels{Fatal, Info, Debug, Warn, Error}
	instance := levelsAsProvider{givenLevels, "foo"}

	actual := instance.GetLevels()

	assert.ToBeEqual(t, givenLevels, actual)
}
