package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

const (
	simpleTimeLayout = "15:04:05.000"
)

type fallbackCoreLogger struct {
	*fallbackProvider
	name  string
	level *Level
}

func (instance *fallbackCoreLogger) Log(event Event) {
	if instance.level.CompareTo(event.GetLevel()) > 0 {
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

	if ts := GetTimestampOf(event, instance); ts != nil {
		if _, err := buf.WriteString(instance.formatTime(*ts) + " "); err != nil {
			return nil, err
		}
	}

	if _, err := fmt.Fprintf(buf, "[%5v]", event.GetLevel()); err != nil {
		return nil, err
	}

	if message := GetMessageOf(event, instance); message != nil {
		if _, err := buf.WriteString(" " + *message); err != nil {
			return nil, err
		}
	}

	messageKey := instance.GetFieldKeySpec().GetMessage()
	loggerKey := instance.GetFieldKeySpec().GetLogger()
	timestampKey := instance.GetFieldKeySpec().GetTimestamp()
	if err := event.GetFields().ForEach(func(key string, value interface{}) error {
		if key == loggerKey && value == GlobalLoggerName {
			return nil
		}
		if key == messageKey || key == timestampKey {
			return nil
		}
		v, err := instance.formatValue(value)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(buf, " %s=%s", key, string(v))
		return err
	}); err != nil {
		return nil, err
	}

	buf.WriteByte('\n')

	return buf.Bytes(), nil
}

func (instance *fallbackCoreLogger) formatTime(time time.Time) string {
	return time.Format(simpleTimeLayout)
}

func (instance *fallbackCoreLogger) formatValue(value interface{}) ([]byte, error) {
	return json.Marshal(value)
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
