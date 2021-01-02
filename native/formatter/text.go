package formatter

import (
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/native/color"
	"github.com/echocat/slf4g/native/hints"
	nlevel "github.com/echocat/slf4g/native/level"
)

const (
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
	DefaultAllowMultiLineMessage = true

	// DefaultPrintRootLogger is default setting if the root logger field should
	// be logged. See Text.PrintRootLogger / Json.PrintRootLogger for more
	// information.
	DefaultPrintRootLogger = false
)

// Text is an implementation of Formatter which formats given log entries in a
// human readable format. Additionally it can also colorize the formatted
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
	// in front of them. If not set set DefaultMultiLineMessageAfterFields will
	// be used.
	MultiLineMessageAfterFields *bool

	// AllowMultiLineMessage will allow (if set to true) multiline messages to
	// be printed as multiline to the output, too. If set to false linebreaks
	// will be replaced with ⏎. If not set set DefaultAllowMultiLineMessage will
	// be used.
	AllowMultiLineMessage *bool

	// PrintRootLogger will (if set to true) also print the field logger for the
	// root logger. If set to false the logger field will be only printed for
	// every logger but not for the root one. If not set set
	// DefaultPrintRootLogger will be used.
	PrintRootLogger *bool

	// ValueFormatter is used to format the field values (not the message). If
	// not set formatter.DefaultTextValue will be used.
	ValueFormatter TextValue

	// KeySorter is used to sort the field when they are printed. If not set
	// fields.DefaultKeySorter will be used.
	KeySorter fields.KeySorter
}

// NewText creates a new instance of Text which is ready to use.
func NewText(customizer ...func(*Text)) *Text {
	result := &Text{}
	for _, c := range customizer {
		c(result)
	}
	return result
}

// Format implements Formatter.Format()
func (instance *Text) Format(event log.Event, using log.Provider, h hints.Hints) ([]byte, error) {
	to := newBufferedTextEncoder()

	message := log.GetMessageOf(event, using)
	isMultilineMessage := false
	if message != nil {
		*message = instance.formatMessage(*message)
		if strings.ContainsRune(*message, '\n') {
			isMultilineMessage = instance.getMultiLineMessageAfterFields()
		} else {
			// Multiline message could be printed on a dedicated line
			*message = instance.ensureMessageWidth(*message)
		}
	}

	atLeastOneFieldPrinted := false
	if err := executeChecked(
		instance.printTimestamp(event, using, h, to),
		instance.printLevel(event, using, h, to),
		instance.printSingleLineMessageIfRequired(message, isMultilineMessage, to),
		instance.printFields(using, h, event, to, &atLeastOneFieldPrinted),
		instance.printMultiLineMessageIfRequired(message, isMultilineMessage, &atLeastOneFieldPrinted, to),
		to.WriteByteChecked('\n'),
	); err != nil {
		return nil, err
	}

	return to.Bytes(), nil
}

func (instance *Text) formatMessage(message string) string {
	message = strings.TrimLeftFunc(message, func(r rune) bool {
		return r == '\r' || r == '\n'
	})
	message = strings.TrimRightFunc(message, unicode.IsSpace)
	message = strings.TrimFunc(message, func(r rune) bool {
		return r == '\r' || !unicode.IsGraphic(r)
	})
	if instance.getAllowMultiLineMessage() {
		message = strings.ReplaceAll(message, "\n", "\u23CE")
	}
	return message
}

func (instance *Text) printTimestamp(event log.Event, using log.Provider, h hints.Hints, to textEncoder) checkedExecution {
	if v := log.GetTimestampOf(event, using); v != nil {
		if instance.shouldColorize(h) {
			return to.WriteStringChecked(`[30;1m` + instance.formatTime(*v) + `[0m `)
		}
		return to.WriteStringChecked(instance.formatTime(*v) + ` `)
	}
	return nil
}

func (instance *Text) printLevel(event log.Event, using log.Provider, h hints.Hints, to textEncoder) checkedExecution {
	v := instance.ensureLevelWidth(event.GetLevel(), using)

	return to.WriteStringChecked(`[` + instance.colorize(event, v, h) + `]`)
}

func (instance *Text) printFields(using log.Provider, h hints.Hints, event log.Event, to textEncoder, atLeastOneFieldPrinted *bool) checkedExecution {
	return func() error {
		formatter := instance.getValueFormatter()
		keysSpec := using.GetFieldKeysSpec()
		messageKey := keysSpec.GetMessage()
		loggerKey := keysSpec.GetLogger()
		timestampKey := keysSpec.GetTimestamp()
		printRootLogger := instance.getPrintRootLogger()
		fieldSorter := instance.getFieldSorter()

		return fields.SortedForEach(event, fieldSorter, func(k string, v interface{}) error {
			if vl, ok := v.(fields.Lazy); ok {
				v = vl.Get()
			}
			if v == nil {
				return nil
			}
			if !printRootLogger && k == loggerKey && v == "ROOT" {
				return nil
			}
			if k == messageKey || k == timestampKey {
				return nil
			}
			*atLeastOneFieldPrinted = true
			return instance.printField(event, k, v, using, formatter, h, to)
		})
	}
}

func (instance *Text) printField(event log.Event, key string, value interface{}, using log.Provider, formatter TextValue, h hints.Hints, to textEncoder) error {
	v, err := formatter.FormatTextValue(value, using)
	if err != nil {
		return err
	}
	return to.WriteString(` ` + instance.colorize(event, key, h) + `=` + string(v))
}

func (instance *Text) printSingleLineMessageIfRequired(message *string, handleAsMultiline bool, to textEncoder) checkedExecution {
	if !handleAsMultiline && message != nil {
		return to.WriteStringChecked(` ` + *message)
	}
	return nil
}

func (instance *Text) printMultiLineMessageIfRequired(message *string, handleAsMultiline bool, atLeastOneFieldPrinted *bool, to textEncoder) checkedExecution {
	if handleAsMultiline && message != nil {
		return func() error {
			otherIdent := "\t"
			firstIdent := otherIdent

			var prefixExecution checkedExecution
			if *atLeastOneFieldPrinted {
				prefixExecution = to.WriteByteChecked('\n')
			} else {
				firstIdent = " "
			}

			return executeChecked(
				prefixExecution,
				instance.printMultilineWithIdent(*message, firstIdent, otherIdent, to),
			)
		}
	}
	return nil
}

func (instance *Text) printMultilineWithIdent(str string, firstLine, otherLines string, to textEncoder) checkedExecution {
	return func() error {
		for i, line := range strings.Split(str, "\n") {
			ident := &firstLine
			var prefixExecution checkedExecution
			if i > 0 {
				ident = &otherLines
				prefixExecution = to.WriteByteChecked('\n')
			}

			if err := executeChecked(
				prefixExecution,
				to.WriteStringPChecked(ident),
				to.WriteStringChecked(line),
			); err != nil {
				return err
			}
		}
		return nil
	}
}

func (instance *Text) colorize(event log.Event, message string, h hints.Hints) string {
	if instance.shouldColorize(h) {
		return instance.getLevelColorizer().ColorizeByLevel(event.GetLevel(), message)
	}
	return message
}

func (instance *Text) shouldColorize(h hints.Hints) bool {
	supported := color.SupportedNone
	if v, ok := h.(hints.ColorsSupport); ok {
		supported = v.IsColorSupported()
	}
	return instance.ColorMode.ShouldColorize(supported)
}

func (instance *Text) formatTime(time time.Time) string {
	return time.Format(instance.getTimeLayout())
}

func (instance *Text) ensureMessageWidth(str string) string {
	width := instance.getMinMessageWidth()
	l2r := true
	if width < 0 {
		width *= -1
		l2r = false
	}
	if width == 0 {
		return str
	}
	l := utf8.RuneCountInString(str)
	if l >= int(width) {
		return str
	}
	if l2r {
		return str + strings.Repeat(" ", int(width)-l)
	} else {
		return strings.Repeat(" ", int(width)-l) + str
	}
}

func (instance *Text) ensureLevelWidth(l level.Level, using log.Provider) string {
	var names nlevel.Names
	if v, ok := using.(nlevel.NamesAware); ok {
		names = v.GetLevelNames()
	} else {
		names = nlevel.DefaultNames
	}
	str := nlevel.AsNamed(&l, names).String()
	width := instance.getLevelWidth()

	l2r := true
	if width < 0 {
		width *= -1
		l2r = false
	}
	if width == 0 {
		return str
	}
	if len(str) >= int(width) {
		return str[:width]
	}
	if l2r {
		return str + strings.Repeat(" ", int(width)-len(str))
	} else {
		return strings.Repeat(" ", int(width)-len(str)) + str
	}
}

func (instance *Text) getLevelColorizer() nlevel.Colorizer {
	if v := instance.LevelColorizer; v != nil {
		return v
	}
	if v := nlevel.DefaultColorizer; v != nil {
		return v
	}
	return nlevel.NoopColorizer()
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
	return DefaultMultiLineMessageAfterFields
}

func (instance *Text) getAllowMultiLineMessage() bool {
	if v := instance.AllowMultiLineMessage; v != nil {
		return *v
	}
	return DefaultAllowMultiLineMessage
}

func (instance *Text) getPrintRootLogger() bool {
	if v := instance.PrintRootLogger; v != nil {
		return *v
	}
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
