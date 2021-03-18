package formatter

import (
	"errors"
	"fmt"
	"testing"
	"time"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/native/color"
	"github.com/echocat/slf4g/native/execution"
	"github.com/echocat/slf4g/native/formatter/encoding"
	"github.com/echocat/slf4g/native/hints"
	nlevel "github.com/echocat/slf4g/native/level"
	"github.com/echocat/slf4g/testing/recording"
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

func Test_Text_Format(t *testing.T) {
	instance := NewText(func(text *Text) {
		text.ColorMode = color.ModeNever
	})
	givenProvider := recording.NewProvider()
	givenLogger := givenProvider.GetRootLogger()
	givenSpec := givenProvider.GetFieldKeysSpec()

	cases := []struct {
		event          map[string]interface{}
		allowMultiline bool
		expected       string
	}{{
		event: map[string]interface{}{
			givenSpec.GetMessage(): "hello, world",
		},
		expected: "[ INFO] hello, world                                      \n",
	}, {
		event: map[string]interface{}{
			givenSpec.GetMessage():   "hello, world",
			givenSpec.GetTimestamp(): mustParseTime("2021-01-02T13:14:15.1234"),
		},
		expected: "13:14:15.123[ INFO] hello, world                                      \n",
	}, {
		event: map[string]interface{}{
			givenSpec.GetMessage():   "hello, world",
			givenSpec.GetTimestamp(): mustParseTime("2021-01-02T13:14:15.1234"),
			"foo1":                   "bar1",
		},
		expected: "13:14:15.123[ INFO] hello, world                                       foo1=bar1\n",
	}, {
		event: map[string]interface{}{
			givenSpec.GetMessage():   "hello, world",
			givenSpec.GetTimestamp(): mustParseTime("2021-01-02T13:14:15.1234"),
			"foo1":                   "bar1",
			"foo2":                   2,
		},
		expected: "13:14:15.123[ INFO] hello, world                                       foo1=bar1 foo2=2\n",
	}, {
		event: map[string]interface{}{
			givenSpec.GetMessage():   "hello,\nworld",
			givenSpec.GetTimestamp(): mustParseTime("2021-01-02T13:14:15.1234"),
			"foo1":                   "bar1",
			"foo2":                   2,
		},
		expected: "13:14:15.123[ INFO] hello,⏎world                                       foo1=bar1 foo2=2\n",
	}, {
		event: map[string]interface{}{
			givenSpec.GetMessage():   "hello,\nworld",
			givenSpec.GetTimestamp(): mustParseTime("2021-01-02T13:14:15.1234"),
			"foo1":                   "bar1",
			"foo2":                   2,
		},
		allowMultiline: true,
		expected:       "13:14:15.123[ INFO] foo1=bar1 foo2=2\n\thello,\n\tworld\n",
	}}

	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			instance.AllowMultiLineMessage = &c.allowMultiline
			givenEvent := givenLogger.NewEvent(level.Info, c.event)

			actual, actualErr := instance.Format(givenEvent, givenProvider, nil)

			assert.ToBeNoError(t, actualErr)
			assert.ToBeEqual(t, c.expected, string(actual))
		})
	}
}

func Test_Text_getMessage(t *testing.T) {
	instance := NewText()
	givenProvider := recording.NewProvider()
	givenLogger := givenProvider.GetRootLogger()

	cases := []struct {
		given          interface{}
		allowMultiline bool
		expected       interface{}
	}{{
		given:          nil,
		allowMultiline: false,
		expected:       nil,
	}, {
		given:          "",
		allowMultiline: false,
		expected:       nil,
	}, {
		given:          "abc",
		allowMultiline: false,
		expected:       "abc",
	}, {
		given:          "  abc  ",
		allowMultiline: false,
		expected:       "  abc",
	}, {
		given:          "\u001Babc\u001B",
		allowMultiline: false,
		expected:       "abc",
	}, {
		given:          "\r\n  abc \r\n ",
		allowMultiline: false,
		expected:       "  abc",
	}, {
		given:          "abc\rdef",
		allowMultiline: false,
		expected:       "abcdef",
	}, {
		given:          "abc\ndef",
		allowMultiline: false,
		expected:       "abc⏎def",
	}, {
		given:          "abc\ndef",
		allowMultiline: true,
		expected:       "abc\ndef",
	}}

	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			instance.AllowMultiLineMessage = &c.allowMultiline
			givenEvent := givenLogger.NewEvent(level.Info, nil).
				With(givenProvider.GetFieldKeysSpec().GetMessage(), c.given)

			actual := instance.getMessage(givenEvent, givenProvider)

			if c.expected != nil {
				assert.ToBeNotNil(t, actual)
				assert.ToBeEqual(t, c.expected, *actual)
			} else {
				assert.ToBeNil(t, actual)
			}
		})
	}
}

func Test_Text_printTimestampChecked(t *testing.T) {
	instance := NewText(func(text *Text) {
		text.ColorMode = color.ModeNever
	})
	givenProvider := recording.NewProvider()
	givenLogger := givenProvider.GetRootLogger()

	cases := []struct {
		givenTimestamp time.Time
		withColor      bool
		expected       string
	}{{
		givenTimestamp: mustParseTime("2021-01-02T13:14:15.1234"),
		withColor:      false,
		expected:       "13:14:15.123",
	}, {
		givenTimestamp: mustParseTime("2021-01-02T13:14:15.1234"),
		withColor:      true,
		expected: `[37m13:14:15.123[0m`,
	}, {
		givenTimestamp: time.Time{},
		expected:       ``,
	}}

	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			givenEncoder := encoding.NewBufferedTextEncoder()
			givenEvent := givenLogger.NewEvent(level.Info, nil).
				With(givenProvider.GetFieldKeysSpec().GetTimestamp(), c.givenTimestamp)

			var givenHints hints.Hints
			if c.withColor {
				givenHints = mockColorizingHints{}
			}

			actualExecution := instance.printTimestampChecked(givenEvent, givenProvider, givenHints, givenEncoder)
			if c.givenTimestamp.IsZero() {
				assert.ToBeNil(t, actualExecution)
				assert.ToBeEqual(t, c.expected, givenEncoder.String())
			} else {
				actualErr := actualExecution()

				assert.ToBeNoError(t, actualErr)
				assert.ToBeEqual(t, c.expected, givenEncoder.String())
			}
		})
	}
}

func Test_Text_printLevelChecked_nonColorized(t *testing.T) {
	instance := NewText(func(text *Text) {
		text.ColorMode = color.ModeNever
	})
	givenProvider := recording.NewProvider()

	cases := map[level.Level]string{
		level.Trace: "TRACE",
		level.Debug: "DEBUG",
		level.Info:  " INFO",
		level.Warn:  " WARN",
		level.Error: "ERROR",
		level.Fatal: "FATAL",
	}

	for l, n := range cases {
		t.Run(n, func(t *testing.T) {
			givenEncoder := encoding.NewBufferedTextEncoder()

			actualErr := instance.printLevelChecked(l, givenProvider, nil, givenEncoder)()

			assert.ToBeNoError(t, actualErr)
			assert.ToBeEqual(
				t,
				n,
				givenEncoder.String(),
			)
		})
	}
}

func Test_Text_printLevelChecked_colorized(t *testing.T) {
	instance := NewText()
	givenProvider := recording.NewProvider()
	givenHints := mockColorizingHints{}

	c := map[level.Level]string{
		level.Trace: "TRACE",
		level.Debug: "DEBUG",
		level.Info:  " INFO",
		level.Warn:  " WARN",
		level.Error: "ERROR",
		level.Fatal: "FATAL",
	}

	for l, n := range c {
		t.Run(n, func(t *testing.T) {
			givenEncoder := encoding.NewBufferedTextEncoder()

			actualErr := instance.printLevelChecked(l, givenProvider, givenHints, givenEncoder)()

			assert.ToBeNoError(t, actualErr)
			assert.ToBeEqual(
				t,
				fmt.Sprintf("%d(%s)", l/1000, n),
				givenEncoder.String(),
			)
		})
	}
}

func Test_Text_printFieldsChecked(t *testing.T) {
	instance := NewText(func(text *Text) {
		text.ColorMode = color.ModeNever
	})
	givenProvider := recording.NewProvider()
	givenLogger := givenProvider.GetRootLogger()
	givenHints := mockColorizingHints{}

	eventOf := func(fields map[string]interface{}) log.Event {
		return givenLogger.NewEvent(level.Info, fields)
	}

	cases := []struct {
		given    log.Event
		expected string
	}{{
		given:    eventOf(map[string]interface{}{}),
		expected: "",
	}, {
		given: eventOf(map[string]interface{}{
			"foo1": "bar1",
		}),
		expected: " 3(foo1)=bar1",
	}, {
		given: eventOf(map[string]interface{}{
			"foo1": "bar1",
			"foo2": 2,
		}),
		expected: " 3(foo1)=bar1 3(foo2)=2",
	}}

	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			givenEncoder := encoding.NewBufferedTextEncoder()
			var atLeastOneFieldPrinted bool
			actualErr := instance.printFieldsChecked(
				givenProvider,
				givenHints,
				c.given,
				givenEncoder,
				&atLeastOneFieldPrinted,
			)()

			assert.ToBeNoError(t, actualErr)
			assert.ToBeEqual(t, c.expected, givenEncoder.String())
			assert.ToBeEqual(t, c.expected != "", atLeastOneFieldPrinted)
		})
	}

}

func Test_Text_printField(t *testing.T) {
	instance := NewText(func(text *Text) {
		text.ColorMode = color.ModeNever
	})
	provider := recording.NewProvider()

	cases := []struct {
		givenLevel                 level.Level
		givenKey                   string
		givenValue                 interface{}
		givenHints                 hints.Hints
		givenShouldPrintRootLogger bool
		expected                   string
	}{{
		givenLevel: level.Info,
		givenKey:   "foo",
		givenValue: "bar",
		givenHints: mockColorizingHints{},
		expected:   " 3(foo)=bar",
	}, {
		givenLevel: level.Info,
		givenKey:   "foo",
		givenValue: "bar",
		expected:   " foo=bar",
	}, {
		givenLevel: level.Info,
		givenKey:   "foo",
		givenValue: fields.LazyFormat("b%sr", "a"),
		expected:   " foo=bar",
	}, {
		givenLevel: level.Info,
		givenKey:   "logger",
		givenValue: "bar",
		expected:   " logger=bar",
	}, {
		givenLevel: level.Info,
		givenKey:   "logger",
		givenValue: "ROOT",
		expected:   "",
	}, {
		givenLevel:                 level.Info,
		givenKey:                   "logger",
		givenValue:                 "ROOT",
		givenShouldPrintRootLogger: true,
		expected:                   " logger=ROOT",
	}, {
		givenLevel: level.Info,
		givenKey:   "foo",
		givenValue: "",
		expected:   " foo=",
	}, {
		givenLevel: level.Info,
		givenKey:   "foo",
		givenValue: nil,
		expected:   "",
	}, {
		givenLevel: level.Info,
		givenKey:   provider.GetFieldKeysSpec().GetMessage(),
		givenValue: "bar",
		expected:   "",
	}, {
		givenLevel: level.Info,
		givenKey:   provider.GetFieldKeysSpec().GetTimestamp(),
		givenValue: "bar",
		expected:   "",
	}}

	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			instance.PrintRootLogger = &c.givenShouldPrintRootLogger
			givenEncoder := encoding.NewBufferedTextEncoder()
			actualPrinted, actualErr := instance.printField(
				c.givenLevel,
				c.givenKey,
				c.givenValue,
				c.givenHints,
				provider,
				givenEncoder,
			)

			assert.ToBeNoError(t, actualErr)
			assert.ToBeEqual(t, c.expected, givenEncoder.String())
			assert.ToBeEqual(t, c.expected != "", actualPrinted)
		})
	}
}
func Test_Text_printField_failsWithValueFormatter(t *testing.T) {
	expectedErr := errors.New("expected")
	instance := NewText(func(text *Text) {
		text.ColorMode = color.ModeNever
		text.ValueFormatter = TextValueFunc(func(i interface{}, provider log.Provider) ([]byte, error) {
			return nil, expectedErr
		})
	})
	provider := recording.NewProvider()

	givenEncoder := encoding.NewBufferedTextEncoder()
	_, actualErr := instance.printField(level.Info, "foo", "bar", nil, provider, givenEncoder)

	assert.ToBeSame(t, expectedErr, actualErr)
	assert.ToBeEqual(t, "", givenEncoder.String())
}

func Test_Text_printMessageAsSingleLineIfRequired(t *testing.T) {
	instance := NewText(func(text *Text) {
		v := int16(12)
		text.MinMessageWidth = &v
	})

	cases := []struct {
		given    string
		expected string
	}{{
		given:    "hello, world",
		expected: " hello, world",
	}, {
		given:    "hello",
		expected: " hello       ",
	}, {
		given:    "hello, world!!",
		expected: " hello, world!!",
	}}

	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			givenEncoder := encoding.NewBufferedTextEncoder()
			actualErr := instance.printMessageAsSingleLineIfRequiredChecked(&c.given, true, givenEncoder)()

			assert.ToBeNoError(t, actualErr)
			assert.ToBeEqual(t, c.expected, givenEncoder.String())
		})
	}
}

func Test_Text_printMessageAsSingleLineIfRequired_ignoresNegativePredicate(t *testing.T) {
	instance := NewText()
	givenMessage := "hello, world"
	givenEncoder := encoding.NewBufferedTextEncoder()

	actual := instance.printMessageAsSingleLineIfRequiredChecked(&givenMessage, false, givenEncoder)

	assert.ToBeNil(t, actual)
}

func Test_Text_printMessageAsSingleLineIfRequired_ignoresNilMessage(t *testing.T) {
	instance := NewText()
	givenEncoder := encoding.NewBufferedTextEncoder()

	actual := instance.printMessageAsSingleLineIfRequiredChecked(nil, true, givenEncoder)

	assert.ToBeNil(t, actual)
}

func Test_Text_printMessageAsSingleLineIfRequired_ignoresEmptyMessage(t *testing.T) {
	instance := NewText()
	givenMessage := ""
	givenEncoder := encoding.NewBufferedTextEncoder()

	actual := instance.printMessageAsSingleLineIfRequiredChecked(&givenMessage, true, givenEncoder)

	assert.ToBeNil(t, actual)
}

func Test_Text_printMessageAsMultiLineIfRequired(t *testing.T) {
	instance := NewText()

	cases := []struct {
		given             string
		givenFieldPrinted bool
		expected          string
	}{{
		given:             "hello\nworld",
		givenFieldPrinted: false,
		expected:          " hello\n\tworld",
	}, {
		given:             "hello\nworld",
		givenFieldPrinted: true,
		expected:          "\n\thello\n\tworld",
	}, {
		given:             "hello",
		givenFieldPrinted: false,
		expected:          " hello",
	}, {
		given:             "hello",
		givenFieldPrinted: true,
		expected:          "\n\thello",
	}}

	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			givenEncoder := encoding.NewBufferedTextEncoder()
			actualErr := instance.printMessageAsMultiLineIfRequiredChecked(&c.given, true, &c.givenFieldPrinted, givenEncoder)()

			assert.ToBeNoError(t, actualErr)
			assert.ToBeEqual(t, c.expected, givenEncoder.String())
		})
	}
}

func Test_Text_printMessageAsMultiLineIfRequired_ignoresNegativePredicate(t *testing.T) {
	instance := NewText()
	givenMessage := "hello, world"
	givenAtLeastOneFieldPrinted := true
	givenEncoder := encoding.NewBufferedTextEncoder()

	actual := instance.printMessageAsMultiLineIfRequiredChecked(&givenMessage, false, &givenAtLeastOneFieldPrinted, givenEncoder)

	assert.ToBeNil(t, actual)
}

func Test_Text_printMessageAsMultiLineIfRequired_ignoresNilMessage(t *testing.T) {
	instance := NewText()
	givenAtLeastOneFieldPrinted := true
	givenEncoder := encoding.NewBufferedTextEncoder()

	actual := instance.printMessageAsMultiLineIfRequiredChecked(nil, true, &givenAtLeastOneFieldPrinted, givenEncoder)

	assert.ToBeNil(t, actual)
}

func Test_Text_printMessageAsMultiLineIfRequired_ignoresEmptyMessage(t *testing.T) {
	instance := NewText()
	givenMessage := ""
	givenAtLeastOneFieldPrinted := true
	givenEncoder := encoding.NewBufferedTextEncoder()

	actual := instance.printMessageAsMultiLineIfRequiredChecked(&givenMessage, true, &givenAtLeastOneFieldPrinted, givenEncoder)

	assert.ToBeNil(t, actual)
}

func Test_Text_textAsHints_ColorMode(t *testing.T) {
	instance := NewText()

	for _, expected := range []color.Mode{color.ModeAuto, color.ModeAlways, color.ModeNever} {
		t.Run(expected.String(), func(t *testing.T) {
			instance.ColorMode = expected
			actual := instance.textAsHints.ColorMode()

			assert.ToBeEqual(t, expected, actual)
		})
	}
}

func Test_Text_textAsHints_LevelColorizer(t *testing.T) {
	instance := NewText()
	instance.LevelColorizer = nlevel.ColorizerMap{level.Info: "abc"}

	actual := instance.textAsHints.LevelColorizer()

	assert.ToBeEqual(t, instance.LevelColorizer, actual)
}

func Test_Text_wrapHints_ColorMode(t *testing.T) {
	instance := NewText()
	h := instance.wrapHints(mockColorsSupport(color.SupportedAssumed))

	for _, expected := range []color.Mode{color.ModeAuto, color.ModeAlways, color.ModeNever} {
		t.Run(expected.String(), func(t *testing.T) {
			instance.ColorMode = expected
			actual := h.(hints.ColorMode).ColorMode()

			assert.ToBeEqual(t, expected, actual)
		})
	}
}

func Test_Text_wrapHints_LevelColorizer(t *testing.T) {
	instance := NewText()
	instance.LevelColorizer = nlevel.ColorizerMap{level.Info: "abc"}
	h := instance.wrapHints(mockColorsSupport(color.SupportedAssumed))

	actual := h.(hints.LevelColorizer).LevelColorizer()

	assert.ToBeEqual(t, instance.LevelColorizer, actual)
}

func Test_Text_wrapHints_ColorsSupport(t *testing.T) {
	instance := NewText()
	instance.LevelColorizer = nlevel.ColorizerMap{level.Info: "abc"}
	h := instance.wrapHints(mockColorsSupport(color.SupportedAssumed))

	actual, actualOk := h.(hints.ColorsSupport)

	assert.ToBeEqual(t, true, actualOk)
	if actualOk {
		assert.ToBeEqual(t, color.SupportedAssumed, actual.IsColorSupported())
	}
}

func Test_Text_wrapHints_ColorsSupport_fallback(t *testing.T) {
	instance := NewText()
	instance.LevelColorizer = nlevel.ColorizerMap{level.Info: "abc"}
	h := instance.wrapHints(nil)

	actual, actualOk := h.(hints.ColorsSupport)

	assert.ToBeEqual(t, true, actualOk)
	if actualOk {
		assert.ToBeEqual(t, color.SupportedNone, actual.IsColorSupported())
	}
}

func Test_Text_formatTime(t *testing.T) {
	givenTime := time.Date(2020, 01, 02, 03, 04, 05, 06, time.UTC)
	instance := NewText()

	actual := instance.formatTime(givenTime)

	assert.ToBeEqual(t, "03:04:05.000", actual)
}

func Test_Text_ensureWidthChecked(t *testing.T) {
	instance := NewText()

	actualString := "hello, world"

	actualExecution := instance.ensureWidthChecked(&actualString, 11, true)

	assert.ToBeOfType(t, (execution.Execution)(nil), actualExecution)
	assert.ToBeEqual(t, "hello, world", actualString)

	actualErr := actualExecution()
	assert.ToBeNoError(t, actualErr)

	assert.ToBeEqual(t, "hello, worl", actualString)
}

func Test_Text_formatLevelChecked(t *testing.T) {
	givenProvider := &someProvider{}
	instance := NewText()

	for _, l := range level.GetProvider().GetLevels() {
		t.Run(fmt.Sprint(l), func(t *testing.T) {
			expectedString := givenProvider.toName(l)

			var actualString string
			actualExecution := instance.formatLevelChecked(l, givenProvider, &actualString)

			assert.ToBeOfType(t, (execution.Execution)(nil), actualExecution)
			assert.ToBeEqual(t, "", actualString)

			actualErr := actualExecution()
			assert.ToBeNoError(t, actualErr)

			assert.ToBeEqual(t, expectedString, actualString)
		})
	}
}

func Test_Text_getLevelNames_byProvider(t *testing.T) {
	givenProvider := &someProvider{}
	instance := NewText()

	actual := instance.getLevelNames(givenProvider)

	assert.ToBeSame(t, givenProvider, actual)
}

func Test_Text_getLevelNames_default(t *testing.T) {
	instance := NewText()

	actual := instance.getLevelNames(nil)

	assert.ToBeEqual(t, nlevel.DefaultNames, actual)
}

func Test_Text_getLevelNames_noop(t *testing.T) {
	old := nlevel.DefaultNames
	defer func() {
		nlevel.DefaultNames = old
	}()
	nlevel.DefaultNames = nil

	instance := NewText()

	actual := instance.getLevelNames(nil)

	assert.ToBeEqual(t, nlevel.NewNames(), actual)
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

	//goland:noinspection GoBoolExpressions
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

	//goland:noinspection GoBoolExpressions
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

	//goland:noinspection GoBoolExpressions
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
