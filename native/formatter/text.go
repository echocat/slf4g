package formatter

import (
	"bytes"
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
	// not set formatter.DefaultTextValueFormatter will be used.
	ValueFormatter TextValueFormatter

	// FieldSorter is used to sort the field when they are printed. If not set
	// fields.DefaultKeySorter will be used.
	FieldSorter fields.KeySorter
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
func (instance *Text) Format(event log.Event, using log.Provider, h hints.Hints) (_ []byte, err error) {
	buf := new(bytes.Buffer)

	message := log.GetMessageOf(event, using)
	multiLineMessage := false
	if message != nil {
		message = instance.formatMessage(message)
		if strings.ContainsRune(*message, '\n') {
			if v := instance.MultiLineMessageAfterFields; v != nil {
				multiLineMessage = *instance.MultiLineMessageAfterFields
			} else {
				multiLineMessage = DefaultMultiLineMessageAfterFields
			}
		} else {
			// Multiline message could be printed on a dedicated line
			*message = instance.ensureMessageWidth(*message)
		}
	}

	if _, err = instance.printTimestamp(event, buf, using, h); err != nil {
		return nil, err
	}
	if _, err = instance.printLevel(event, buf, using, h); err != nil {
		return nil, err
	}

	if !multiLineMessage && message != nil {
		if _, err := buf.WriteString(` ` + *message); err != nil {
			return nil, err
		}
	}

	fieldsPrinted := false
	if fieldsPrinted, err = instance.printFields(event, buf, using, h); err != nil {
		return nil, err
	}

	if multiLineMessage && message != nil {
		otherIdent := "\t"
		firstIdent := otherIdent
		if fieldsPrinted {
			if err := buf.WriteByte('\n'); err != nil {
				return nil, err
			}
		} else {
			firstIdent = " "
		}
		if err := instance.printWithIdent(*message, firstIdent, otherIdent, buf); err != nil {
			return nil, err
		}
	}

	buf.WriteByte('\n')

	return buf.Bytes(), nil
}

func (instance *Text) formatMessage(message *string) *string {
	*message = strings.TrimLeftFunc(*message, func(r rune) bool {
		return r == '\r' || r == '\n'
	})
	*message = strings.TrimRightFunc(*message, unicode.IsSpace)
	*message = strings.TrimFunc(*message, func(r rune) bool {
		return r == '\r' || !unicode.IsGraphic(r)
	})
	if (instance.AllowMultiLineMessage != nil && *instance.AllowMultiLineMessage) ||
		(instance.AllowMultiLineMessage == nil && DefaultAllowMultiLineMessage) {
		*message = strings.ReplaceAll(*message, "\n", "\u23CE")
	}
	return message
}

func (instance *Text) printTimestamp(event log.Event, buf *bytes.Buffer, using log.Provider, h hints.Hints) (cn int, err error) {
	if ts := log.GetTimestampOf(event, using); ts != nil {
		if instance.shouldColorize(h) {
			_, err = buf.WriteString(`[30;1m` + instance.formatTime(*ts) + `[0m `)
			cn = len(instance.formatTime(*ts)) + 1
		} else {
			_, err = buf.WriteString(instance.formatTime(*ts) + ` `)
			cn = len(instance.formatTime(*ts)) + 1
		}
	}
	return
}

func (instance *Text) shouldColorize(h hints.Hints) bool {
	supported := color.SupportedNone
	if aware, ok := h.(hints.ColorsSupport); ok {
		supported = aware.IsColorSupported()
	}
	return instance.ColorMode.ShouldColorize(supported)
}

func (instance *Text) printLevel(event log.Event, buf *bytes.Buffer, using log.Provider, h hints.Hints) (cn int, err error) {
	v := instance.ensureLevelWidth(event.GetLevel(), using)

	_, err = buf.WriteString(`[` + instance.colorize(event, v, h) + `]`)
	cn = 1 + len(v) + 1

	return
}

func (instance *Text) printFields(event log.Event, buf *bytes.Buffer, using log.Provider, h hints.Hints) (printed bool, err error) {
	formatter := instance.getFieldValueFormatter()

	keysSpec := using.GetFieldKeysSpec()
	messageKey := keysSpec.GetMessage()
	loggerKey := keysSpec.GetLogger()
	timestampKey := keysSpec.GetTimestamp()

	printRootLogger := DefaultPrintRootLogger
	if v := instance.PrintRootLogger; v != nil {
		printRootLogger = *v
	}

	err = fields.SortedForEach(event, instance.getFieldSorter(), func(k string, v interface{}) error {
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
		printed = true
		return instance.printField(event, k, v, buf, using, formatter, h)
	})

	return
}

func (instance *Text) printField(event log.Event, key string, value interface{}, buf *bytes.Buffer, using log.Provider, formatter TextValueFormatter, h hints.Hints) error {
	v, err := formatter.FormatValue(value, using)
	if err != nil {
		return err
	}
	_, err = buf.WriteString(` ` + instance.colorize(event, key, h) + `=` + string(v))
	return err
}

func (instance *Text) colorize(event log.Event, message string, h hints.Hints) string {
	if instance.shouldColorize(h) {
		return instance.getLevelColorizer().ColorizeByLevel(event.GetLevel(), message)
	}
	return message
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

func (instance *Text) getFieldValueFormatter() TextValueFormatter {
	if v := instance.ValueFormatter; v != nil {
		return v
	}
	if v := DefaultTextValueFormatter; v != nil {
		return v
	}
	return NoopTextValueFormatter()
}

func (instance *Text) getFieldSorter() fields.KeySorter {
	if v := instance.FieldSorter; v != nil {
		return v
	}
	if v := fields.DefaultKeySorter; v != nil {
		return v
	}
	return fields.NoopKeySorter()
}

func (instance *Text) formatTime(time time.Time) string {
	lt := instance.TimeLayout
	if lt == "" {
		lt = DefaultTimeLayout
	}
	return time.Format(lt)
}

func (instance *Text) ensureMessageWidth(str string) string {
	width := DefaultMinMessageWidth
	if v := instance.MinMessageWidth; v != nil {
		width = *v
	}
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

func (instance *Text) printWithIdent(str string, firstLine, otherLines string, buf *bytes.Buffer) error {
	for i, line := range strings.Split(str, "\n") {
		ident := firstLine
		if i > 0 {
			ident = otherLines
			if _, err := buf.WriteRune('\n'); err != nil {
				return err
			}
		}

		if _, err := buf.WriteString(ident); err != nil {
			return err
		}

		if _, err := buf.WriteString(line); err != nil {
			return err
		}
	}

	return nil
}

func (instance *Text) ensureLevelWidth(l level.Level, using log.Provider) string {
	var names nlevel.Names
	if v, ok := using.(nlevel.NamesAware); ok {
		names = v.GetLevelNames()
	} else {
		names = nlevel.DefaultNames
	}
	str := nlevel.AsNamed(&l, names).String()
	width := DefaultLevelWidth
	if v := instance.LevelWidth; v != nil {
		width = *v
	}

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
