package formatter

import (
	"testing"

	"github.com/echocat/slf4g/fields"

	"github.com/echocat/slf4g/level"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/native/color"
	nlevel "github.com/echocat/slf4g/native/level"
)

func Test_NewText(t *testing.T) {
	actual := NewText()

	assert.ToBeEqual(t, color.ModeAuto, actual.ColorMode)
	assert.ToBeNil(t, actual.LevelColorizer)
	assert.ToBeEqual(t, "", actual.TimeLayout)
	assert.ToBeNil(t, actual.LevelWidth)
	assert.ToBeNil(t, actual.MinMessageWidth)
	assert.ToBeNil(t, actual.MultiLineMessageAfterFields)
	assert.ToBeNil(t, actual.AllowMultiLineMessage)
	assert.ToBeNil(t, actual.PrintRootLogger)
	assert.ToBeNil(t, actual.ValueFormatter)
	assert.ToBeNil(t, actual.KeySorter)
}

func Test_NewText_customized(t *testing.T) {
	givenColorMode := color.ModeAlways
	givenLevelColorzier := nlevel.ColorizerMap{level.Level(666): "-6-6-6-"}
	givenTimeLayout := "abc"
	givenLevelWidth := int8(66)
	givenMinMessageWidth := int16(666)
	givenMultiLineMessageAfterFields := true
	givenAllowMultiLineMessage := true
	givenPrintRootLogger := true
	givenValueFormatter := TextValueFunc(func(interface{}, log.Provider) ([]byte, error) {
		panic("should never be called")
	})
	givenKeySorter := fields.KeySorter(func(keys []string) {
		panic("should never be called")
	})

	actual := NewText(func(text *Text) {
		text.ColorMode = givenColorMode
		text.LevelColorizer = givenLevelColorzier
		text.TimeLayout = givenTimeLayout
		text.LevelWidth = &givenLevelWidth
		text.MinMessageWidth = &givenMinMessageWidth
		text.MultiLineMessageAfterFields = &givenMultiLineMessageAfterFields
		text.AllowMultiLineMessage = &givenAllowMultiLineMessage
		text.PrintRootLogger = &givenPrintRootLogger
		text.ValueFormatter = &givenValueFormatter
		text.KeySorter = givenKeySorter
	})

	assert.ToBeEqual(t, givenColorMode, actual.ColorMode)
	assert.ToBeEqual(t, givenLevelColorzier, actual.LevelColorizer)
	assert.ToBeEqual(t, givenTimeLayout, actual.TimeLayout)
	assert.ToBeSame(t, &givenLevelWidth, actual.LevelWidth)
	assert.ToBeSame(t, &givenMinMessageWidth, actual.MinMessageWidth)
	assert.ToBeSame(t, &givenMultiLineMessageAfterFields, actual.MultiLineMessageAfterFields)
	assert.ToBeSame(t, &givenAllowMultiLineMessage, actual.AllowMultiLineMessage)
	assert.ToBeSame(t, &givenPrintRootLogger, actual.PrintRootLogger)
	assert.ToBeSame(t, &givenValueFormatter, actual.ValueFormatter)
	assert.ToBeSame(t, givenKeySorter, actual.KeySorter)
}
