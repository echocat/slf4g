package testlog

import (
	"errors"
	"testing"
	"time"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
)

const (
	dateTimeFormat = "2006-01-02 15:04:05"
)

func Test_coreLogger_Log_regular(t *testing.T) {
	provider := NewProvider(t)
	provider.initIfRequired()
	instance := provider.coreLogger

	var actualMsg string
	var actualskipFrames uint16
	instance.interceptLogDepth = func(msg string, skipFrames uint16) {
		actualMsg = msg
		actualskipFrames = skipFrames
	}

	provider.GetRootLogger().
		WithError(errors.New("testError")).
		With("stringField", "bar").
		With("intField", 123).
		With("lazyField", fields.LazyFunc(func() interface{} { return "lazy" })).
		With("nilField", nil).
		With("excludedField", fields.Exclude).
		With("ignoredByLevelField", fields.IgnoreLevels(level.Debug, level.Info+1, "ignored")).
		With("respectedByLevelField", fields.IgnoreLevels(level.Warn, level.Error, "respected")).
		Info("foo")

	assert.ToBeMatching(t, `^\d+ \[ INFO] foo error="testError" intField=123 lazyField="lazy" nilField=null respectedByLevelField="respected" stringField="bar"$`, actualMsg)
	assert.ToBeEqual(t, uint16(6), actualskipFrames)
}

func Test_coreLogger_NewEvent(t *testing.T) {
	provider := NewProvider(t)
	provider.initIfRequired()
	instance := provider.coreLogger

	actual := instance.NewEvent(level.Level(666), map[string]interface{}{
		"foo": 123,
		"bar": "str",
	})

	assert.ToBeEqual(t, &event{
		provider,
		fields.WithAll(map[string]interface{}{
			"foo": 123,
			"bar": "str",
		}),
		666,
	}, actual)
}

func Test_coreLogger_Accepts(t *testing.T) {
	provider := NewProvider(t, Level(level.Level(700)))
	provider.initIfRequired()
	instance := provider.coreLogger

	givenAcceptable := instance.NewEvent(level.Level(666), map[string]interface{}{})

	assert.ToBeEqual(t, true, instance.Accepts(givenAcceptable))
}

func Test_coreLogger_Log_tooLowLevel(t *testing.T) {
	provider := NewProvider(t)
	provider.initIfRequired()
	instance := provider.coreLogger

	var actualMsg string
	var actualskipFrames uint16
	instance.interceptLogDepth = func(msg string, skipFrames uint16) {
		actualMsg = msg
		actualskipFrames = skipFrames
	}

	provider.GetRootLogger().Trace("foo")

	assert.ToBeEqual(t, "", actualMsg)
	assert.ToBeEqual(t, uint16(0), actualskipFrames)
}

func Test_coreLogger_Log_fail(t *testing.T) {
	provider := NewProvider(t)
	provider.initIfRequired()
	instance := provider.coreLogger

	var actualMsg string
	var actualskipFrames uint16
	instance.interceptLogDepth = func(msg string, skipFrames uint16) {
		actualMsg = msg
		actualskipFrames = skipFrames
	}
	var actualFail bool
	instance.interceptFail = func() {
		actualFail = true
	}
	var actualFailNow bool
	instance.interceptFailNow = func() {
		actualFailNow = true
	}

	provider.GetRootLogger().Error("foo")

	assert.ToBeMatching(t, `^\d+ \[ERROR] foo$`, actualMsg)
	assert.ToBeEqual(t, uint16(6), actualskipFrames)
	assert.ToBeEqual(t, true, actualFail)
	assert.ToBeEqual(t, false, actualFailNow)
}

func Test_coreLogger_Log_failNow(t *testing.T) {
	provider := NewProvider(t)
	provider.initIfRequired()
	instance := provider.coreLogger

	var actualMsg string
	var actualskipFrames uint16
	instance.interceptLogDepth = func(msg string, skipFrames uint16) {
		actualMsg = msg
		actualskipFrames = skipFrames
	}
	var actualFail bool
	instance.interceptFailNow = func() {
		actualFail = true
	}
	var actualFailNow bool
	instance.interceptFailNow = func() {
		actualFailNow = true
	}

	provider.GetRootLogger().Fatal("foo")

	assert.ToBeMatching(t, `^\d+ \[FATAL] foo$`, actualMsg)
	assert.ToBeEqual(t, uint16(6), actualskipFrames)
	assert.ToBeEqual(t, false, actualFail)
	assert.ToBeEqual(t, true, actualFailNow)
}

func Test_coreLogger_formatTime_sinceTestStartedMcs(t *testing.T) {
	provider := NewProvider(t, TimeFormat(SinceTestStartedMcsTimeFormat))
	provider.initIfRequired()
	instance := provider.coreLogger

	givenTs, _ := time.Parse(dateTimeFormat, "2024-07-25 18:56:13")
	givenEvent := instance.NewEvent(level.Info, map[string]interface{}{
		"timestamp": givenTs,
	})

	assert.ToBeMatching(t, `^\d+ $`, instance.formatTime(givenEvent))
}

func Test_coreLogger_formatTime_noop(t *testing.T) {
	provider := NewProvider(t, TimeFormat(NoopTimeFormat))
	provider.initIfRequired()
	instance := provider.coreLogger

	givenTs, _ := time.Parse(dateTimeFormat, "2024-07-25 18:56:13")
	givenEvent := instance.NewEvent(level.Info, map[string]interface{}{
		"timestamp": givenTs,
	})

	assert.ToBeEqual(t, "", instance.formatTime(givenEvent))
}

func Test_coreLogger_formatTime_ts(t *testing.T) {
	provider := NewProvider(t, TimeFormat(dateTimeFormat))
	provider.initIfRequired()
	instance := provider.coreLogger

	givenTs, _ := time.Parse(dateTimeFormat, "2024-07-25 18:56:13")
	givenEvent := instance.NewEvent(level.Info, map[string]interface{}{
		"timestamp": givenTs,
	})

	assert.ToBeEqual(t, "2024-07-25 18:56:13 ", instance.formatTime(givenEvent))
}

func Test_coreLogger_formatTime_ts_defaultNow(t *testing.T) {
	provider := NewProvider(t, TimeFormat(time.RFC3339))
	provider.initIfRequired()
	instance := provider.coreLogger

	givenEvent := instance.NewEvent(level.Info, map[string]interface{}{})

	now := time.Now()
	actualTs, actualErr := time.Parse(time.RFC3339+" ", instance.formatTime(givenEvent))
	assert.ToBeNoError(t, actualErr)

	diff := now.Sub(actualTs)
	if diff > time.Second*10 || diff < -time.Second*10 {
		t.Fatalf("the difference between %v(now) and %v(ts) should not be greather than 10s; bug was: %v", now, actualTs, diff)
	}
}
