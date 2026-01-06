package testlog

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/echocat/slf4g/fields"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/level"
)

// RootLoggerName specifies the name of the root version of coreLogger
// instances which are managed by Provider.
const RootLoggerName = "ROOT"

type coreLogger struct {
	*Provider
	name  string
	level level.Level
}

// Log implements log.CoreLogger#Log(event).
func (instance *coreLogger) Log(event log.Event, skipFrames uint16) {
	instance.tb.Helper()
	instance.log(instance.GetName(), event, skipFrames+1)
}

func (instance *coreLogger) logDepth(msg string, skipFrames uint16) {
	instance.tb.Helper()
	i := instance.interceptLogDepth
	if i == nil {
		i = instance.logViaSdk
	}
	i(msg, skipFrames+1)
}

func (instance *coreLogger) logViaSdk(str string, _ uint16) {
	instance.tb.Helper()
	instance.tb.Log(str)
}

func (instance *coreLogger) fail() {
	i := instance.interceptFail
	if i == nil {
		i = instance.tb.Fail
	}
	i()
}

func (instance *coreLogger) failNow() {
	i := instance.interceptFailNow
	if i == nil {
		i = instance.tb.FailNow
	}
	i()
}

func (instance *coreLogger) log(loggerName string, event log.Event, skipFrames uint16) {
	instance.tb.Helper()
	l := event.GetLevel()
	if !instance.IsLevelEnabled(l) {
		return
	}

	if v := log.GetLoggerOf(event, instance); v == nil {
		event = event.With(instance.GetFieldKeysSpec().GetLogger(), loggerName)
	}

	instance.logDepth(instance.format(event), skipFrames+1)

	failNowAtLevel := instance.getFailNowAtLevel()
	if failNowAtLevel < NeverFailLevel && l >= failNowAtLevel {
		instance.failNow()
		return
	}

	failAtLevel := instance.getFailAtLevel()
	if failAtLevel < NeverFailLevel && l >= failAtLevel {
		instance.fail()
		return
	}
}

// GetLevel implements level.Aware#GetLevel(v). If there was no SetLevel called before,
// it will return the value of the holding Provider.
func (instance *coreLogger) GetLevel() level.Level {
	if v := instance.level; v != 0 {
		return v
	}
	return instance.Provider.GetLevel()
}

// SetLevel implements level.MutableAware#SetLevel(v).
//
// If set to 0 it will reset the handling back to the value of the holding Provider.
func (instance *coreLogger) SetLevel(v level.Level) {
	instance.level = v
}

// IsLevelEnabled implements log.CoreLogger#IsLevelEnabled()
func (instance *coreLogger) IsLevelEnabled(v level.Level) bool {
	return instance.GetLevel().CompareTo(v) <= 0
}

// GetName implements log.CoreLogger#GetName()
func (instance *coreLogger) GetName() string {
	return instance.name
}

// GetProvider implements log.CoreLogger#GetProvider()
func (instance *coreLogger) GetProvider() log.Provider {
	return instance.Provider
}

func (instance *coreLogger) NewEvent(l level.Level, values map[string]interface{}) log.Event {
	return instance.NewEventWithFields(l, fields.WithAll(values))
}

func (instance *coreLogger) NewEventWithFields(l level.Level, f fields.ForEachEnabled) log.Event {
	asFields, err := fields.AsFields(f)
	if err != nil {
		panic(err)
	}
	return &event{
		provider: instance.Provider,
		fields:   asFields,
		level:    l,
	}
}

func (instance *coreLogger) Accepts(e log.Event) bool {
	return e != nil
}

func (instance *coreLogger) format(event log.Event) string {
	buf := new(bytes.Buffer)

	_, _ = buf.WriteString(instance.formatTime(event))
	_, _ = buf.WriteString(instance.formatLevel(event.GetLevel()))
	_, _ = buf.WriteString(instance.formatMessage(event))
	messageKey := instance.GetFieldKeysSpec().GetMessage()
	loggerKey := instance.GetFieldKeysSpec().GetLogger()
	timestampKey := instance.GetFieldKeysSpec().GetTimestamp()
	if err := fields.SortedForEach(event, nil, func(k string, vp interface{}) error {
		if vl, ok := vp.(fields.Filtered); ok {
			fv, shouldBeRespected := vl.Filter(event)
			if !shouldBeRespected {
				return nil
			}
			vp = fv
		} else if vl, ok := vp.(fields.Lazy); ok {
			vp = vl.Get()
		}
		if vp == fields.Exclude {
			return nil
		}

		if k == loggerKey && vp == RootLoggerName {
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
		instance.tb.Fatalf("ERR!! Cannot format event %v: %v", event, err)
		return ""
	}

	return buf.String()
}

func (instance *coreLogger) formatLevel(l level.Level) string {
	return "[" + instance.getLevelFormatter().Format(l) + "]"
}

func (instance *coreLogger) formatTime(event log.Event) string {
	tf := instance.getTimeFormat()
	if tf == NoopTimeFormat {
		return ""
	}

	if tf == SinceTestStartedMcsTimeFormat {
		diff := runtimeNano() - instance.startedNs
		return strconv.FormatInt(diff/1000, 10) + " "
	}

	if v := log.GetTimestampOf(event, instance); v != nil {
		return v.Format(tf) + " "
	}
	return time.Now().Format(tf) + " "
}

func (instance *coreLogger) formatMessage(event log.Event) string {
	var message string
	if v := log.GetMessageOf(event, instance); v != nil {
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

func (instance *coreLogger) formatValue(v interface{}) ([]byte, error) {
	if ve, ok := v.(error); ok {
		v = ve.Error()
	}
	return json.Marshal(v)
}

// Helper wraps the helper of the testing framework into this logger.
// As this is called by the whole logging stack (if required) this will ensure
// the SDK logging framework respects the top entry as the log position.
func (instance *coreLogger) Helper() func() {
	return instance.tb.Helper
}
