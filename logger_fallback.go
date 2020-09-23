package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/echocat/slf4g/fields"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	simpleTimeLayout = "0102 15:04:05.000000"
)

var (
	pid = os.Getpid()
)

type fallbackCoreLogger struct {
	*fallbackProvider
	name string
}

func (instance *fallbackCoreLogger) Log(event Event) {
	if !instance.IsLevelEnabled(event.GetLevel()) {
		return
	}

	if v := GetTimestampOf(event, instance); v == nil {
		event = event.With(instance.GetFieldKeySpec().GetTimestamp(), time.Now())
	}
	if v := GetLoggerOf(event, instance); v == nil {
		event = event.With(instance.GetFieldKeySpec().GetLogger(), instance.name)
	}

	s, err := instance.format(event)
	if err != nil {
		s = []byte(fmt.Sprintf("ERR!! Cannot format event %v: %v", event, err))
	}

	_, _ = instance.out.Write(s)
}

func (instance *fallbackCoreLogger) format(event Event) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := buf.WriteByte(instance.formatLevel(event.GetLevel())); err != nil {
		return nil, err
	}
	if ts := GetTimestampOf(event, instance); ts != nil {
		if _, err := buf.WriteString(*instance.formatTime(ts)); err != nil {
			return nil, err
		}
	}
	if err := buf.WriteByte(' '); err != nil {
		return nil, err
	}
	if _, err := buf.WriteString(*instance.formatPid()); err != nil {
		return nil, err
	}
	if err := buf.WriteByte(' '); err != nil {
		return nil, err
	}
	if _, err := buf.WriteString(*instance.formatLocation(event)); err != nil {
		return nil, err
	}
	if err := buf.WriteByte(']'); err != nil {
		return nil, err
	}

	if message := GetMessageOf(event, instance); message != nil {
		if err := buf.WriteByte(' '); err != nil {
			return nil, err
		}
		if _, err := buf.WriteString(*instance.formatMessage(message)); err != nil {
			return nil, err
		}
	}

	messageKey := instance.GetFieldKeySpec().GetMessage()
	loggerKey := instance.GetFieldKeySpec().GetLogger()
	timestampKey := instance.GetFieldKeySpec().GetTimestamp()
	if err := event.GetFields().ForEach(func(k string, vp interface{}) error {
		if vl, ok := vp.(fields.Lazy); ok {
			vp = vl.Get()
		}

		if k == loggerKey && vp == GlobalLoggerName {
			return nil
		}
		if k == messageKey || k == timestampKey {
			return nil
		}
		v, err := instance.formatValue(vp)
		if err != nil {
			return err
		}

		if err := buf.WriteByte(' '); err != nil {
			return err
		}
		if _, err := buf.WriteString(k); err != nil {
			return err
		}
		if err := buf.WriteByte('='); err != nil {
			return err
		}
		if _, err := buf.Write(v); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	buf.WriteByte('\n')

	return buf.Bytes(), nil
}

func (instance *fallbackCoreLogger) formatLevel(level Level) byte {
	switch level {
	case LevelTrace:
		return 'T'
	case LevelDebug:
		return 'D'
	case LevelInfo:
		return 'I'
	case LevelWarn:
		return 'W'
	case LevelError:
		return 'E'
	case LevelFatal:
		return 'F'
	default:
		return '?'
	}
}

func (instance *fallbackCoreLogger) formatPid() *string {
	result := strconv.Itoa(pid)
	return &result
}

func (instance *fallbackCoreLogger) formatLocation(event Event) *string {
	_, file, line, ok := runtime.Caller(3 + event.GetCallDepth())
	if !ok {
		file = "???"
		line = 1
	} else {
		file = path.Base(file)
	}
	if line < 0 {
		line = 0
	}
	result := file + ":" + strconv.Itoa(line)
	return &result
}

func (instance *fallbackCoreLogger) formatTime(time *time.Time) *string {
	result := time.Format(simpleTimeLayout)
	return &result
}

func (instance *fallbackCoreLogger) formatMessage(message *string) *string {
	*message = strings.TrimLeftFunc(*message, func(r rune) bool {
		return r == '\r' || r == '\n'
	})
	*message = strings.TrimRightFunc(*message, unicode.IsSpace)
	*message = strings.TrimFunc(*message, func(r rune) bool {
		return r == '\r' || !unicode.IsGraphic(r)
	})
	*message = strings.ReplaceAll(*message, "\n", "\u23CE")
	return message
}

func (instance *fallbackCoreLogger) formatValue(v interface{}) ([]byte, error) {
	if ve, ok := v.(error); ok {
		v = ve.Error()
	}
	return json.Marshal(v)
}

func (instance *fallbackCoreLogger) IsLevelEnabled(level Level) bool {
	return instance.level.CompareTo(level) <= 0
}

func (instance *fallbackCoreLogger) GetName() string {
	return instance.name
}

func (instance *fallbackCoreLogger) GetProvider() Provider {
	return instance
}
