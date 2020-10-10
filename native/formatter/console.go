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
	"github.com/echocat/slf4g/native/formatter/hints"
	nlevel "github.com/echocat/slf4g/native/level"
)

var (
	DefaultConsole = NewConsole()
)

type Console struct {
	ColorMode           color.Mode
	LevelBasedColorizer color.LevelBasedColorizer

	TimeLayout string

	LevelWidth int8

	MinMessageWidth             int16
	MultiLineMessageAfterFields bool
	AllowMultiLineMessage       bool

	PrintGlobalLogger   bool
	FieldValueFormatter ValueFormatter
	FieldSorter         fields.KeySorter
}

func NewConsole() *Console {
	return &Console{
		LevelBasedColorizer: color.DefaultLevelBasedColorizer,
		ColorMode:           color.ModeAuto,

		TimeLayout: "15:04:05.000",

		LevelWidth: -5,

		MinMessageWidth:             50,
		MultiLineMessageAfterFields: true,
		AllowMultiLineMessage:       true,

		FieldValueFormatter: DefaultValueFormatter,
		PrintGlobalLogger:   false,
		FieldSorter:         fields.DefaultKeySorter,
	}
}

func (instance *Console) Format(event log.Event, using log.Provider, h hints.Hints) (_ []byte, err error) {
	buf := new(bytes.Buffer)

	message := log.GetMessageOf(event, using)
	multiLineMessage := false
	if message != nil {
		message = instance.formatMessage(message)
		if strings.ContainsRune(*message, '\n') {
			multiLineMessage = instance.MultiLineMessageAfterFields
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

func (instance *Console) formatMessage(message *string) *string {
	*message = strings.TrimLeftFunc(*message, func(r rune) bool {
		return r == '\r' || r == '\n'
	})
	*message = strings.TrimRightFunc(*message, unicode.IsSpace)
	*message = strings.TrimFunc(*message, func(r rune) bool {
		return r == '\r' || !unicode.IsGraphic(r)
	})
	if !instance.AllowMultiLineMessage {
		*message = strings.ReplaceAll(*message, "\n", "\u23CE")
	}
	return message
}

func (instance *Console) shouldColorize(h hints.Hints) bool {
	switch instance.ColorMode {
	case color.ModeAlways:
		return true
	case color.ModeAuto:
		if ca, ok := h.(hints.ColorsSupport); ok {
			return ca.GetColorSupport().IsSupported()
		}
		return false
	default:
		return false
	}
}

func (instance *Console) printTimestamp(event log.Event, buf *bytes.Buffer, using log.Provider, h hints.Hints) (cn int, err error) {
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

func (instance *Console) printLevel(event log.Event, buf *bytes.Buffer, using log.Provider, h hints.Hints) (cn int, err error) {
	v := instance.ensureLevelWidth(event.GetLevel(), using)

	_, err = buf.WriteString(`[` + instance.colorize(event, v, h) + `]`)
	cn = 1 + len(v) + 1

	return
}

func (instance *Console) printFields(event log.Event, buf *bytes.Buffer, using log.Provider, h hints.Hints) (printed bool, err error) {
	formatter := instance.getFieldValueFormatter()

	messageKey := using.GetFieldKeysSpec().GetMessage()
	loggerKey := using.GetFieldKeysSpec().GetLogger()
	timestampKey := using.GetFieldKeysSpec().GetTimestamp()

	err = fields.SortedForEach(event, instance.FieldSorter, func(k string, v interface{}) error {
		if vl, ok := v.(fields.Lazy); ok {
			v = vl.Get()
		}
		if v == nil {
			return nil
		}
		if !instance.PrintGlobalLogger && k == loggerKey && v == "ROOT" {
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

func (instance *Console) printField(event log.Event, key string, value interface{}, buf *bytes.Buffer, using log.Provider, formatter ValueFormatter, h hints.Hints) error {
	v, err := formatter.FormatValue(value, using)
	if err != nil {
		return err
	}
	_, err = buf.WriteString(` ` + instance.colorize(event, key, h) + `=` + string(v))
	return err
}

func (instance *Console) colorize(event log.Event, message string, h hints.Hints) string {
	if instance.shouldColorize(h) {
		return instance.getLevelBasedColorizer().Colorize(event.GetLevel(), message)
	}
	return message
}

func (instance *Console) getLevelBasedColorizer() color.LevelBasedColorizer {
	if v := instance.LevelBasedColorizer; v != nil {
		return v
	}
	return color.DefaultLevelBasedColorizer
}

func (instance *Console) getFieldValueFormatter() ValueFormatter {
	if v := instance.FieldValueFormatter; v != nil {
		return v
	}
	return DefaultValueFormatter
}

func (instance *Console) formatTime(time time.Time) string {
	return time.Format(instance.TimeLayout)
}

func (instance *Console) ensureMessageWidth(str string) string {
	width := instance.MinMessageWidth
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

func (instance *Console) printWithIdent(str string, firstLine, otherLines string, buf *bytes.Buffer) error {
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

func (instance *Console) ensureLevelWidth(l level.Level, using log.Provider) string {
	str := nlevel.AsSerializable(&l, using.(nlevel.NamesAware)).String()
	width := instance.LevelWidth
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
