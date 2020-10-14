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

const DefaultTimeLayout = "15:04:05.000"
const DefaultLevelWidth = int8(-5)
const DefaultMinMessageWidth = int16(50)
const DefaultMultiLineMessageAfterFields = true
const DefaultAllowMultiLineMessage = true
const DefaultPrintGlobalLogger = false

type Console struct {
	ColorMode      color.Mode
	LevelColorizer nlevel.Colorizer

	TimeLayout string

	LevelWidth *int8

	MinMessageWidth             *int16
	MultiLineMessageAfterFields *bool
	AllowMultiLineMessage       *bool

	PrintGlobalLogger *bool
	ValueFormatter    ValueFormatter
	FieldSorter       fields.KeySorter
}

func NewConsole(customizer ...func(*Console)) *Console {
	result := &Console{}
	for _, c := range customizer {
		c(result)
	}
	return result
}

func (instance *Console) Format(event log.Event, using log.Provider, h hints.Hints) (_ []byte, err error) {
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

func (instance *Console) formatMessage(message *string) *string {
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

func (instance *Console) printTimestamp(event log.Event, buf *bytes.Buffer, using log.Provider, h hints.Hints) (cn int, err error) {
	if ts := log.GetTimestampOf(event, using); ts != nil {
		if instance.ColorMode.ShouldColorizeByCheckingHints(h) {
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

	keysSpec := using.GetFieldKeysSpec()
	messageKey := keysSpec.GetMessage()
	loggerKey := keysSpec.GetLogger()
	timestampKey := keysSpec.GetTimestamp()

	printGlobalLogger := DefaultPrintGlobalLogger
	if v := instance.PrintGlobalLogger; v != nil {
		printGlobalLogger = *v
	}

	err = fields.SortedForEach(event, instance.getFieldSorter(), func(k string, v interface{}) error {
		if vl, ok := v.(fields.Lazy); ok {
			v = vl.Get()
		}
		if v == nil {
			return nil
		}
		if !printGlobalLogger && k == loggerKey && v == "ROOT" {
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
	if instance.ColorMode.ShouldColorizeByCheckingHints(h) {
		return instance.getLevelColorizer().ColorizeByLevel(event.GetLevel(), message)
	}
	return message
}

func (instance *Console) getLevelColorizer() nlevel.Colorizer {
	if v := instance.LevelColorizer; v != nil {
		return v
	}
	if v := nlevel.DefaultColorizer; v != nil {
		return v
	}
	return nlevel.NoopColorizer()
}

func (instance *Console) getFieldValueFormatter() ValueFormatter {
	if v := instance.ValueFormatter; v != nil {
		return v
	}
	return DefaultValueFormatter
}

func (instance *Console) getFieldSorter() fields.KeySorter {
	if v := instance.FieldSorter; v != nil {
		return v
	}
	if v := fields.DefaultKeySorter; v != nil {
		return v
	}
	return func([]string) {}
}

func (instance *Console) formatTime(time time.Time) string {
	lt := instance.TimeLayout
	if lt == "" {
		lt = DefaultTimeLayout
	}
	return time.Format(lt)
}

func (instance *Console) ensureMessageWidth(str string) string {
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
