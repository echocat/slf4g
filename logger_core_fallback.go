package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

const (
	simpleTimeLayout       = "0102 15:04:05.000000"
	fallbackRootLoggerName = "ROOT"
)

var (
	pid = os.Getpid()
)

// IsFallbackLogger will return true of the given (Core)Logger is the fallback
// logger. This usually indicates that currently there is no other
// implementation of slf4g registered.
func IsFallbackLogger(candidate CoreLogger) bool {
	for candidate != nil {
		if _, ok := candidate.(*fallbackCoreLogger); ok {
			return true
		}
		candidate = UnwrapCoreLogger(candidate)
	}
	return false
}

type fallbackCoreLogger struct {
	*fallbackProvider
	name  string
	level level.Level
}

func (instance *fallbackCoreLogger) Log(event Event, skipFrames uint16) {
	if !instance.IsLevelEnabled(event.GetLevel()) {
		return
	}

	if v := GetLoggerOf(event, instance); v == nil {
		event = event.With(instance.GetFieldKeysSpec().GetLogger(), instance.name)
	}

	b := instance.format(event, skipFrames+1)
	_, _ = instance.out.Write(b)
}

func (instance *fallbackCoreLogger) format(event Event, skipFrames uint16) []byte {
	buf := new(bytes.Buffer)

	_ = buf.WriteByte(instance.formatLevel(event.GetLevel()))
	_, _ = buf.WriteString(*instance.formatTime(event))
	_ = buf.WriteByte(' ')
	_, _ = buf.WriteString(instance.formatPid())
	_ = buf.WriteByte(' ')
	_, _ = buf.WriteString(instance.formatLocation(skipFrames + 1))
	_ = buf.WriteByte(']')
	_, _ = buf.WriteString(instance.formatMessage(event))
	messageKey := instance.GetFieldKeysSpec().GetMessage()
	loggerKey := instance.GetFieldKeysSpec().GetLogger()
	timestampKey := instance.GetFieldKeysSpec().GetTimestamp()
	if err := fields.SortedForEach(event, nil, func(k string, vp interface{}) error {
		if vl, ok := vp.(fields.Lazy); ok {
			vp = vl.Get()
		}
		if vp == fields.Exclude {
			return nil
		}

		if k == loggerKey && vp == fallbackRootLoggerName {
			return nil
		}
		if k == messageKey || k == timestampKey {
			return nil
		}
		v, err := instance.formatValue(vp)
		if err != nil {
			return err
		}

		_ = buf.WriteByte(' ')
		_, _ = buf.WriteString(k)
		_ = buf.WriteByte('=')
		_, _ = buf.Write(v)
		return nil
	}); err != nil {
		return []byte(fmt.Sprintf("ERR!! Cannot format event %v: %v", event, err))
	}

	buf.WriteByte('\n')

	return buf.Bytes()
}

func (instance *fallbackCoreLogger) formatLevel(l level.Level) byte {
	switch l {
	case level.Trace:
		return 'T'
	case level.Debug:
		return 'D'
	case level.Info:
		return 'I'
	case level.Warn:
		return 'W'
	case level.Error:
		return 'E'
	case level.Fatal:
		return 'F'
	default:
		return '?'
	}
}

func (instance *fallbackCoreLogger) formatPid() string {
	return strconv.Itoa(pid)
}

func (instance *fallbackCoreLogger) formatLocation(skipFrames uint16) string {
	_, file, line, ok := runtime.Caller(int(skipFrames + 1))
	if !ok {
		file = "???"
	} else {
		file = path.Base(file)
	}
	if line > 0 {
		return file + ":" + strconv.Itoa(line)
	}
	return file + ":?"
}

func (instance *fallbackCoreLogger) formatTime(event Event) *string {
	var result string
	if v := GetTimestampOf(event, instance); v != nil {
		result = v.Format(simpleTimeLayout)
	} else {
		result = time.Now().Format(simpleTimeLayout)
	}
	return &result
}

func (instance *fallbackCoreLogger) formatMessage(event Event) string {
	var message string
	if v := GetMessageOf(event, instance); v != nil {
		message = *v

		message = strings.TrimLeftFunc(message, func(r rune) bool {
			return r == '\r' || r == '\n'
		})
		message = strings.TrimRightFunc(message, unicode.IsSpace)
		message = strings.TrimFunc(message, func(r rune) bool {
			return r == '\r' || !unicode.IsGraphic(r)
		})
		message = strings.ReplaceAll(message, "\n", "\u23CE")
		if message != "" {
			message = " " + message
		}
	}
	return message
}

func (instance *fallbackCoreLogger) formatValue(v interface{}) ([]byte, error) {
	if ve, ok := v.(error); ok {
		v = ve.Error()
	}
	return json.Marshal(v)
}

func (instance *fallbackCoreLogger) IsLevelEnabled(v level.Level) bool {
	return instance.GetLevel().CompareTo(v) <= 0
}

func (instance *fallbackCoreLogger) GetName() string {
	return instance.name
}

func (instance *fallbackCoreLogger) GetProvider() Provider {
	return instance.fallbackProvider
}

func (instance *fallbackCoreLogger) GetLevel() level.Level {
	if v := instance.level; v != 0 {
		return v
	}
	return instance.fallbackProvider.GetLevel()
}

func (instance *fallbackCoreLogger) SetLevel(in level.Level) {
	instance.level = in
}

func (instance *fallbackCoreLogger) NewEvent(l level.Level, values map[string]interface{}) Event {
	return instance.NewEventWithFields(l, fields.WithAll(values))
}

func (instance *fallbackCoreLogger) NewEventWithFields(l level.Level, f fields.ForEachEnabled) Event {
	asFields, err := fields.AsFields(f)
	if err != nil {
		panic(err)
	}
	return &fallbackEvent{
		provider: instance,
		fields:   asFields,
		level:    l,
	}
}

func (instance *fallbackCoreLogger) Accepts(e Event) bool {
	return e != nil
}
