package formatter

import (
	"strings"
	"time"
	"unicode"

	"github.com/echocat/slf4g/native/formatter/encoding"

	"github.com/echocat/slf4g/native/execution"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/native/color"
	"github.com/echocat/slf4g/native/formatter/functions"
	"github.com/echocat/slf4g/native/hints"
	nlevel "github.com/echocat/slf4g/native/level"
)

var (
	// DefaultTimeLayout is the default format to format times of log entries
	// with. See Text.TimeLayout for more information.
	DefaultTimeLayout = "15:04:05.000"

	// DefaultLevelWidth is the default width of levels. See Text.LevelWidth for
	// more information.
	DefaultLevelWidth = int8(-5)

	// DefaultMinMessageWidth is the default width of messages. See
	// Text.LevelWidth for more information.
	DefaultMinMessageWidth = int16(50)

	// DefaultMultiLineMessageAfterFields is default setting if multiline
	// messages should be printed after fields. See
	// Text.MultiLineMessageAfterFields for more information.
	DefaultMultiLineMessageAfterFields = true

	// DefaultAllowMultiLineMessage is default setting if multiline should be
	// allowed to be multilines. See Text.AllowMultiLineMessage for more
	// information.
	DefaultAllowMultiLineMessage = false

	// DefaultPrintRootLogger is default setting if the root logger field should
	// be logged. See Text.PrintRootLogger / Json.PrintRootLogger for more
	// information.
	DefaultPrintRootLogger = false
)

// Text is an implementation of Formatter which formats given log entries in a
// human-readable format. Additionally, it can also colorize the formatted
// output.
type Text struct {
	// ColorMode defines when the output should be colorized. If not configured
	// color.ModeAuto will be used by default.
	ColorMode color.Mode

	// LevelColorizer is used to colorize output based on the level.Level of an
	// log.Event to be logged. If not set nlevel.DefaultColorizer will be used.
	LevelColorizer nlevel.Colorizer

	// TimeLayout defines how the time of log events should be formatted. Please
	// see time.Time#Format() for more details. If not set DefaultTimeLayout
	// will be used.
	TimeLayout string

	// LevelWidth defines the width of the string representation of the level
	// which will be printed to the output. If set to 0 the length will be kept
	// as it is. If bigger than 0 and has the exact same length of the printed
	// text; it will be printed as it is. If the representation is longer it
	// will be trimmed; if shorter whitespaces will be added. If smaller than 0
	// it behaves the same way but whitespaces will be added at the beginning
	// instead of the end. If not set DefaultLevelWidth will be used.
	LevelWidth *int8

	// MinMessageWidth defines the width of the message which will be printed to
	// the output. If set to 0 the length will be kept as it is. If bigger than
	// 0 and has the exact same length of the printed text or the printed text
	// will be longer ; it will be printed as it is. If the message is shorter
	// whitespaces will be added to the end. If smaller than 0 it behaves the
	// same way but whitespaces will be added at the beginning instead of the
	// end. If not set DefaultMinMessageWidth will be used.
	MinMessageWidth *int16

	// MultiLineMessageAfterFields will force multiline messages
	// (if set to true) to be printed behind the fields; instead (if default)
	// in front of them. If not set DefaultMultiLineMessageAfterFields will
	// be used.
	MultiLineMessageAfterFields *bool

	// AllowMultiLineMessage will allow (if set to true) multiline messages to
	// be printed as multiline to the output, too. If set to false linebreaks
	// will be replaced with ⏎. If not set DefaultAllowMultiLineMessage will
	// be used.
	AllowMultiLineMessage *bool

	// PrintRootLogger will (if set to true) also print the field logger for the
	// root logger. If set to false the logger field will be only printed for
	// every logger but not for the root one. If not set
	// DefaultPrintRootLogger will be used.
	PrintRootLogger *bool

	// ValueFormatter is used to format the field values (not the message). If
	// not set formatter.DefaultTextValue will be used.
	ValueFormatter TextValue

	// KeySorter is used to sort the field when they are printed. If not set
	// fields.DefaultKeySorter will be used.
	KeySorter fields.KeySorter

	textAsHints *textAsHints
}

// NewText creates a new instance of Text which is ready to use.
func NewText(customizer ...func(*Text)) *Text {
	result := &Text{}
	result.textAsHints = &textAsHints{
		Text: result,
	}
	for _, c := range customizer {
		c(result)
	}
	return result
}

// Format implements Formatter.Format()
func (instance *Text) Format(event log.Event, using log.Provider, h hints.Hints) ([]byte, error) {
	to := encoding.NewBufferedTextEncoder()

	message := instance.getMessage(event, using)
	printMessageAsMultiLine := message != nil &&
		strings.ContainsRune(*message, '\n') &&
		instance.getMultiLineMessageAfterFields()

	atLeastOneFieldPrinted := false
	if err := execution.Execute(
		instance.printTimestampChecked(event, using, h, to),

		to.WriteByteChecked('['),
		instance.printLevelChecked(event.GetLevel(), using, h, to),
		to.WriteByteChecked(']'),

		instance.printMessageAsSingleLineIfRequiredChecked(message, !printMessageAsMultiLine, to),

		instance.printFieldsChecked(using, h, event, to, &atLeastOneFieldPrinted),

		instance.printMessageAsMultiLineIfRequiredChecked(message, printMessageAsMultiLine, &atLeastOneFieldPrinted, to),

		to.WriteByteChecked('\n'),
	); err != nil {
		return nil, err
	}

	return to.Bytes(), nil
}

func (instance *Text) getMessage(of log.Event, using log.Provider) *string {
	message := log.GetMessageOf(of, using)

	if message != nil {
		*message = instance.sanitizeMessage(*message)
	}

	if message != nil && *message == "" {
		return nil
	}

	return message
}

func (instance *Text) sanitizeMessage(message string) string {
	message = strings.TrimLeftFunc(message, func(r rune) bool {
		return r == '\r' || r == '\n' || !unicode.IsGraphic(r)
	})
	message = strings.TrimRightFunc(message, func(r rune) bool {
		return unicode.IsSpace(r) || !unicode.IsGraphic(r)
	})
	message = strings.ReplaceAll(message, "\r", "")
	if !instance.getAllowMultiLineMessage() {
		message = strings.ReplaceAll(message, "\n", "\u23CE")
	}
	return message
}

func (instance *Text) printTimestampChecked(event log.Event, using log.Provider, h hints.Hints, to encoding.TextEncoder) execution.Execution {
	if v := log.GetTimestampOf(event, using); v != nil {
		formatted := instance.formatTime(*v)
		colorized := functions.Colorize("37", instance.wrapHints(h), formatted)
		return to.WriteStringChecked(colorized)
	}
	return nil
}

func (instance *Text) printLevelChecked(l level.Level, using log.Provider, h hints.Hints, to encoding.TextEncoder) execution.Execution {
	var v string
	vp := &v
	return execution.Join(
		instance.formatLevelChecked(l, using, vp),
		instance.ensureWidthChecked(vp, int32(instance.getLevelWidth()), true),
		instance.colorizeChecked(l, vp, h),
		to.WriteStringPChecked(vp),
	)
}

func (instance *Text) printFieldsChecked(using log.Provider, h hints.Hints, event log.Event, to encoding.TextEncoder, atLeastOneFieldPrinted *bool) execution.Execution {
	return func() error {
		return fields.SortedForEach(event, instance.getFieldSorter(), func(k string, v interface{}) error {
			printed, err := instance.printField(event, k, v, h, using, to)
			if printed {
				*atLeastOneFieldPrinted = printed
			}
			return err
		})
	}
}

func (instance *Text) printField(ctx fields.FilterContext, k string, v interface{}, h hints.Hints, using log.Provider, to encoding.TextEncoder) (bool, error) {
	if vl, ok := v.(fields.Filtered); ok {
		fv, shouldBeRespected := vl.Filter(ctx)
		if !shouldBeRespected {
			return false, nil
		}
		v = fv
	} else if vl, ok := v.(fields.Lazy); ok {
		v = vl.Get()
	}
	if v == fields.Exclude {
		return false, nil
	}

	keysSpec := using.GetFieldKeysSpec()

	if v == "ROOT" && k == keysSpec.GetLogger() && !instance.getPrintRootLogger() {
		return false, nil
	}
	if k == keysSpec.GetMessage() || k == keysSpec.GetTimestamp() {
		return false, nil
	}
	b, err := instance.getValueFormatter().FormatTextValue(v, using)
	if err != nil {
		return false, err
	}
	return true, to.WriteString(` ` + instance.colorize(ctx.GetLevel(), k, h) + `=` + string(b))
}

func (instance *Text) printMessageAsSingleLineIfRequiredChecked(message *string, predicate bool, to encoding.TextEncoder) execution.Execution {
	if predicate && message != nil && *message != "" {
		v := functions.EnsureWidth(int32(instance.getMinMessageWidth()), false, *message)
		return to.WriteStringChecked(` ` + v)
	}
	return nil
}

func (instance *Text) printMessageAsMultiLineIfRequiredChecked(message *string, predicate bool, atLeastOneFieldPrinted *bool, to encoding.TextEncoder) execution.Execution {
	if predicate && message != nil && *message != "" {
		return func() error {
			if *atLeastOneFieldPrinted {
				if err := to.WriteByte('\n'); err != nil {
					return err
				}
				return functions.EncodeMultilineWithIndent("\t", "\t", to, *message)
			}

			return functions.EncodeMultilineWithIndent(" ", "\t", to, *message)
		}
	}
	return nil
}

func (instance *Text) colorize(l level.Level, message string, h hints.Hints) string {
	return functions.ColorizeByLevel(
		l,
		instance.wrapHints(h),
		message,
	)
}

func (instance *Text) colorizeChecked(l level.Level, message *string, h hints.Hints) execution.Execution {
	return func() error {
		*message = instance.colorize(l, *message, h)
		return nil
	}
}

func (instance *Text) GetColorMode() color.Mode {
	return instance.ColorMode
}

func (instance *Text) SetColorMode(v color.Mode) {
	instance.ColorMode = v
}

func (instance *Text) wrapHints(h hints.Hints) hints.Hints {
	return textHintsCombined{h, instance.textAsHints}
}

func (instance *Text) formatTime(time time.Time) string {
	return time.Format(instance.getTimeLayout())
}

func (instance *Text) ensureWidthChecked(of *string, width int32, cutOffToLong bool) execution.Execution {
	return func() error {
		*of = functions.EnsureWidth(width, cutOffToLong, *of)
		return nil
	}
}

func (instance *Text) formatLevelChecked(l level.Level, using log.Provider, to *string) execution.Execution {
	return func() error {
		v, err := instance.getLevelNames(using).ToName(l)
		*to = v
		return err
	}
}

func (instance *Text) getLevelNames(using log.Provider) level.Names {
	if v, ok := using.(level.NamesAware); ok {
		return v.GetLevelNames()
	}
	if v := nlevel.DefaultNames; v != nil {
		return v
	}
	return nlevel.NewNames()
}

func (instance *Text) getTimeLayout() string {
	if v := instance.TimeLayout; v != "" {
		return v
	}
	return DefaultTimeLayout
}

func (instance *Text) getLevelWidth() int8 {
	if v := instance.LevelWidth; v != nil {
		return *v
	}
	return DefaultLevelWidth
}

func (instance *Text) getMinMessageWidth() int16 {
	if v := instance.MinMessageWidth; v != nil {
		return *v
	}
	return DefaultMinMessageWidth
}

func (instance *Text) getMultiLineMessageAfterFields() bool {
	if v := instance.MultiLineMessageAfterFields; v != nil {
		return *v
	}
	//goland:noinspection GoBoolExpressions
	return DefaultMultiLineMessageAfterFields
}

func (instance *Text) getAllowMultiLineMessage() bool {
	if v := instance.AllowMultiLineMessage; v != nil {
		return *v
	}
	//goland:noinspection GoBoolExpressions
	return DefaultAllowMultiLineMessage
}

func (instance *Text) getPrintRootLogger() bool {
	if v := instance.PrintRootLogger; v != nil {
		return *v
	}
	//goland:noinspection GoBoolExpressions
	return DefaultPrintRootLogger
}

func (instance *Text) getValueFormatter() TextValue {
	if v := instance.ValueFormatter; v != nil {
		return v
	}
	if v := DefaultTextValue; v != nil {
		return v
	}
	return NoopTextValue()
}

func (instance *Text) getFieldSorter() fields.KeySorter {
	if v := instance.KeySorter; v != nil {
		return v
	}
	return fields.DefaultKeySorter
}

type textAsHints struct {
	*Text
}

func (instance *textAsHints) ColorMode() color.Mode {
	return instance.Text.ColorMode
}

func (instance *textAsHints) LevelColorizer() nlevel.Colorizer {
	return instance.Text.LevelColorizer
}

type textHintsCombined struct {
	hints.Hints
	*textAsHints
}

func (instance textHintsCombined) ColorMode() color.Mode {
	if v, ok := instance.Hints.(hints.ColorMode); ok {
		return v.ColorMode()
	}
	return instance.textAsHints.ColorMode()
}

func (instance textHintsCombined) LevelColorizer() nlevel.Colorizer {
	if v, ok := instance.Hints.(hints.LevelColorizer); ok {
		return v.LevelColorizer()
	}
	return instance.textAsHints.LevelColorizer()
}

func (instance textHintsCombined) IsColorSupported() color.Supported {
	if v, ok := instance.Hints.(hints.ColorsSupport); ok {
		return v.IsColorSupported()
	}
	return color.SupportedNone
}
