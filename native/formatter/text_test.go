package formatter

import (
	"testing"
	"time"

	"github.com/echocat/slf4g/testing/recording"

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

func Test_Text_colorize_enabled(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Warn, nil)
	instance := NewText()

	actual := instance.colorize(givenEvent, "hello", mockColorsSupport(color.ModeAlways))

	expected := nlevel.DefaultColorizer.ColorizeByLevel(level.Warn, "hello")
	assert.ToBeEqual(t, expected, actual)
}

func Test_Text_colorize_disabled(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Warn, nil)
	instance := NewText()

	actual := instance.colorize(givenEvent, "hello", &struct{}{})

	assert.ToBeEqual(t, "hello", actual)
}

func Test_Text_shouldColorize_fromSupported(t *testing.T) {
	instance := NewText()

	actual := instance.shouldColorize(mockColorsSupport(color.ModeAlways))

	assert.ToBeEqual(t, true, actual)
}

func Test_Text_shouldColorize_fallback(t *testing.T) {
	instance := NewText()

	actual := instance.shouldColorize(&struct{}{})

	assert.ToBeEqual(t, false, actual)
}

func Test_Text_formatTime(t *testing.T) {
	givenTime := time.Date(2020, 01, 02, 03, 04, 05, 06, time.UTC)
	instance := NewText()

	actual := instance.formatTime(givenTime)

	assert.ToBeEqual(t, "03:04:05.000", actual)
}

func Test_Text_getLevelColorizer_explicit(t *testing.T) {
	givenLevelColorizer := nlevel.NewColorizerFacade(func() nlevel.Colorizer {
		panic("should never be called")
	})
	instance := NewText(func(text *Text) {
		text.LevelColorizer = givenLevelColorizer
	})

	actual := instance.getLevelColorizer()

	assert.ToBeSame(t, givenLevelColorizer, actual)
}

func Test_Text_getLevelColorizer_default(t *testing.T) {
	instance := NewText(func(text *Text) {
		text.LevelColorizer = nil
	})

	actual := instance.getLevelColorizer()

	assert.ToBeEqual(t, nlevel.DefaultColorizer, actual)
}

func Test_Text_getLevelColorizer_noop(t *testing.T) {
	old := nlevel.DefaultColorizer
	defer func() {
		nlevel.DefaultColorizer = old
	}()
	nlevel.DefaultColorizer = nil

	instance := NewText(func(text *Text) {
		text.LevelColorizer = nil
	})

	actual := instance.getLevelColorizer()

	assert.ToBeSame(t, nlevel.NoopColorizer(), actual)
}

func Test_Text_getTimeLayout_explicit(t *testing.T) {
	givenTimeLayout := "abc"
	instance := NewText(func(text *Text) {
		text.TimeLayout = givenTimeLayout
	})

	actual := instance.getTimeLayout()

	assert.ToBeEqual(t, givenTimeLayout, actual)
}

func Test_Text_TimeLayout_default(t *testing.T) {
	instance := NewText(func(text *Text) {
		text.TimeLayout = ""
	})

	actual := instance.getTimeLayout()

	assert.ToBeEqual(t, DefaultTimeLayout, actual)
}

func Test_Text_getLevelWidth_explicit(t *testing.T) {
	givenLevelWidth := int8(66)
	instance := NewText(func(text *Text) {
		text.LevelWidth = &givenLevelWidth
	})

	actual := instance.getLevelWidth()

	assert.ToBeEqual(t, givenLevelWidth, actual)
}

func Test_Text_LevelWidth_default(t *testing.T) {
	instance := NewText(func(text *Text) {
		text.LevelWidth = nil
	})

	actual := instance.getLevelWidth()

	assert.ToBeEqual(t, DefaultLevelWidth, actual)
}

func Test_Text_getMinMessageWidth_explicit(t *testing.T) {
	givenMinMessageWidth := int16(666)
	instance := NewText(func(text *Text) {
		text.MinMessageWidth = &givenMinMessageWidth
	})

	actual := instance.getMinMessageWidth()

	assert.ToBeEqual(t, givenMinMessageWidth, actual)
}

func Test_Text_MinMessageWidth_default(t *testing.T) {
	instance := NewText(func(text *Text) {
		text.MinMessageWidth = nil
	})

	actual := instance.getMinMessageWidth()

	assert.ToBeEqual(t, DefaultMinMessageWidth, actual)
}

func Test_Text_getMultiLineMessageAfterFields_explicit(t *testing.T) {
	givenMultiLineMessageAfterFields := true
	instance := NewText(func(text *Text) {
		text.MultiLineMessageAfterFields = &givenMultiLineMessageAfterFields
	})

	actual := instance.getMultiLineMessageAfterFields()

	assert.ToBeEqual(t, givenMultiLineMessageAfterFields, actual)
}

func Test_Text_getMultiLineMessageAfterFields_default(t *testing.T) {
	instance := NewText(func(text *Text) {
		text.MultiLineMessageAfterFields = nil
	})

	actual := instance.getMultiLineMessageAfterFields()

	assert.ToBeEqual(t, DefaultMultiLineMessageAfterFields, actual)
}

func Test_Text_getAllowMultiLineMessage_explicit(t *testing.T) {
	givenAllowMultiLineMessage := true
	instance := NewText(func(text *Text) {
		text.AllowMultiLineMessage = &givenAllowMultiLineMessage
	})

	actual := instance.getAllowMultiLineMessage()

	assert.ToBeEqual(t, givenAllowMultiLineMessage, actual)
}

func Test_Text_getAllowMultiLineMessage_default(t *testing.T) {
	instance := NewText(func(text *Text) {
		text.AllowMultiLineMessage = nil
	})

	actual := instance.getAllowMultiLineMessage()

	assert.ToBeEqual(t, DefaultAllowMultiLineMessage, actual)
}

func Test_Text_getPrintRootLogger_explicit(t *testing.T) {
	givenPrintRootLogger := true
	instance := NewText(func(text *Text) {
		text.PrintRootLogger = &givenPrintRootLogger
	})

	actual := instance.getPrintRootLogger()

	assert.ToBeEqual(t, givenPrintRootLogger, actual)
}

func Test_Text_getPrintRootLogger_default(t *testing.T) {
	instance := NewText(func(text *Text) {
		text.PrintRootLogger = nil
	})

	actual := instance.getPrintRootLogger()

	assert.ToBeEqual(t, DefaultPrintRootLogger, actual)
}

func Test_Text_getValueFormatter_explicit(t *testing.T) {
	givenValueFormatter := TextValueFunc(func(interface{}, log.Provider) ([]byte, error) {
		panic("should never be called")
	})
	instance := NewText(func(text *Text) {
		text.ValueFormatter = givenValueFormatter
	})

	actual := instance.getValueFormatter()

	assert.ToBeSame(t, givenValueFormatter, actual)
}

func Test_Text_getValueFormatter_default(t *testing.T) {
	instance := NewText(func(text *Text) {
		text.ValueFormatter = nil
	})

	actual := instance.getValueFormatter()

	assert.ToBeSame(t, DefaultTextValue, actual)
}

func Test_Text_getValueFormatter_noop(t *testing.T) {
	old := DefaultTextValue
	defer func() {
		DefaultTextValue = old
	}()
	DefaultTextValue = nil

	instance := NewText(func(text *Text) {
		text.ValueFormatter = nil
	})

	actual := instance.getValueFormatter()

	assert.ToBeSame(t, NoopTextValue(), actual)
}

func Test_Text_getFieldSorter_explicit(t *testing.T) {
	givenKeySorter := fields.KeySorter(func(keys []string) {
		panic("should never be called")
	})
	instance := NewText(func(text *Text) {
		text.KeySorter = givenKeySorter
	})

	actual := instance.getFieldSorter()

	assert.ToBeSame(t, givenKeySorter, actual)
}

func Test_Text_getFieldSorter_default(t *testing.T) {
	instance := NewText(func(text *Text) {
		text.KeySorter = nil
	})

	actual := instance.getFieldSorter()

	assert.ToBeSame(t, fields.DefaultKeySorter, actual)
}

type mockColorsSupport color.Supported

func (instance mockColorsSupport) IsColorSupported() color.Supported {
	return color.Supported(instance)
}
