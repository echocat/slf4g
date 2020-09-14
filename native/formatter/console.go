package formatter

import (
	"bytes"
	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native/color"
	"github.com/echocat/slf4g/native/formatter/hints"
	"strings"
	"time"
	"unicode"
)

const (
	DefaultConsoleTimeLayout      string = "15:04:05.000"
	DefaultConsoleMinMessageWidth int16  = 50
	DefaultConsoleLevelWidth      int8   = -5
)

var (
	DefaultConsole = NewConsole()
)

type Console struct {
	PrintGlobalLogger bool
	ColorMode         color.Mode

	MinMessageWidth *int16
	LevelWidth      *int8

	LevelBasedColorizer color.LevelBasedColorizer

	TimeLayout          string
	FieldValueFormatter ValueFormatter
}

func NewConsole() *Console {
	return &Console{}
}

func (instance *Console) Format(event log.Event, using log.Provider, h hints.Hints) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := instance.printTimestamp(event, buf, using, h); err != nil {
		return nil, err
	} else if err := instance.printLevel(event, buf, using, h); err != nil {
		return nil, err
	} else if err := instance.printMessage(event, buf, using, h); err != nil {
		return nil, err
	} else if err := instance.printFields(event, buf, using, h); err != nil {
		return nil, err
	}

	buf.WriteByte('\n')

	return buf.Bytes(), nil
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

func (instance *Console) printTimestamp(event log.Event, buf *bytes.Buffer, using log.Provider, h hints.Hints) (err error) {
	if ts := log.GetTimestampOf(event, using); ts != nil {
		if instance.shouldColorize(h) {
			_, err = buf.WriteString(`[30;1m` + instance.formatTime(*ts) + `[0m `)
		} else {
			_, err = buf.WriteString(instance.formatTime(*ts) + ` `)
		}
	}
	return
}

func (instance *Console) printLevel(event log.Event, buf *bytes.Buffer, _ log.Provider, h hints.Hints) (err error) {
	v := instance.ensureLevelWidth(event.GetLevel())

	_, err = buf.WriteString(`[` + instance.colorize(event, v, h) + `]`)

	return
}

func (instance *Console) printMessage(event log.Event, buf *bytes.Buffer, using log.Provider, _ hints.Hints) (err error) {
	if message := log.GetMessageOf(event, using); message != nil {
		target := strings.TrimRightFunc(*message, unicode.IsSpace)
		target = instance.ensureMessageWidth(target)
		_, err = buf.WriteString(` ` + target)
	}
	return
}

func (instance *Console) printFields(event log.Event, buf *bytes.Buffer, using log.Provider, h hints.Hints) (err error) {
	formatter := instance.getFieldValueFormatter()

	messageKey := using.GetFieldKeys().GetMessage()
	loggerKey := using.GetFieldKeys().GetLogger()
	timestampKey := using.GetFieldKeys().GetTimestamp()

	return event.ForEach(func(key string, value interface{}) error {
		if (!instance.PrintGlobalLogger && key == loggerKey && value == log.GlobalLoggerName) ||
			key == messageKey ||
			key == timestampKey {
			return nil
		}
		return instance.printField(event, key, value, buf, using, formatter, h)
	})
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
	return time.Format(instance.timeLayout())
}

func (instance *Console) timeLayout() string {
	if v := instance.TimeLayout; v != "" {
		return v
	}
	return DefaultConsoleTimeLayout
}

func (instance *Console) ensureMessageWidth(str string) string {
	width := instance.minMessageWidth()
	l2r := true
	if width < 0 {
		width *= -1
		l2r = false
	}
	if width == 0 {
		return str
	}
	if len(str) >= int(width) {
		return str
	}
	if l2r {
		return str + strings.Repeat(" ", int(width)-len(str))
	} else {
		return strings.Repeat(" ", int(width)-len(str)) + str
	}
}

func (instance *Console) ensureLevelWidth(level log.Level) string {
	str := level.String()
	width := instance.levelWidth()
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

func (instance *Console) minMessageWidth() int16 {
	if v := instance.MinMessageWidth; v != nil {
		return *v
	}
	return DefaultConsoleMinMessageWidth
}

func (instance *Console) levelWidth() int8 {
	if v := instance.LevelWidth; v != nil {
		return *v
	}
	return DefaultConsoleLevelWidth
}
