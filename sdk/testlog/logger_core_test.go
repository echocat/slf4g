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
	instance := provider.coreRootLogger

	var actualMsg string
	var actualSkipFrames uint16
	instance.interceptLogDepth = func(msg string, skipFrames uint16) {
		actualMsg = msg
		actualSkipFrames = skipFrames
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
	assert.ToBeEqual(t, uint16(6), actualSkipFrames)
}

func Test_coreLogger_NewEvent(t *testing.T) {
	provider := NewProvider(t)
	provider.initIfRequired()
	instance := provider.coreRootLogger

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
	instance := provider.coreRootLogger

	givenAcceptable := instance.NewEvent(level.Level(666), map[string]interface{}{})

	assert.ToBeEqual(t, true, instance.Accepts(givenAcceptable))
}

func Test_coreLogger_Log_tooLowLevel(t *testing.T) {
	provider := NewProvider(t)
	provider.initIfRequired()
	instance := provider.coreRootLogger

	var actualMsg string
	var actualSkipFrames uint16
	instance.interceptLogDepth = func(msg string, skipFrames uint16) {
		actualMsg = msg
		actualSkipFrames = skipFrames
	}

	provider.GetRootLogger().Trace("foo")

	assert.ToBeEqual(t, "", actualMsg)
	assert.ToBeEqual(t, uint16(0), actualSkipFrames)
}

func Test_coreLogger_Log_fail(t *testing.T) {
	provider := NewProvider(t)
	provider.initIfRequired()
	instance := provider.coreRootLogger

	var actualMsg string
	var actualSkipFrames uint16
	instance.interceptLogDepth = func(msg string, skipFrames uint16) {
		actualMsg = msg
		actualSkipFrames = skipFrames
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
	assert.ToBeEqual(t, uint16(6), actualSkipFrames)
	assert.ToBeEqual(t, true, actualFail)
	assert.ToBeEqual(t, false, actualFailNow)
}

func Test_coreLogger_Log_failNow(t *testing.T) {
	provider := NewProvider(t)
	provider.initIfRequired()
	instance := provider.coreRootLogger

	var actualMsg string
	var actualSkipFrames uint16
	instance.interceptLogDepth = func(msg string, skipFrames uint16) {
		actualMsg = msg
		actualSkipFrames = skipFrames
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
	assert.ToBeEqual(t, uint16(6), actualSkipFrames)
	assert.ToBeEqual(t, false, actualFail)
	assert.ToBeEqual(t, true, actualFailNow)
}

func Test_coreLogger_formatTime_sinceTestStartedMcs(t *testing.T) {
	provider := NewProvider(t, TimeFormat(SinceTestStartedMcsTimeFormat))
	provider.initIfRequired()
	instance := provider.coreRootLogger

	givenTs, _ := time.Parse(dateTimeFormat, "2024-07-25 18:56:13")
	givenEvent := instance.NewEvent(level.Info, map[string]interface{}{
		"timestamp": givenTs,
	})

	assert.ToBeMatching(t, `^\d+ $`, instance.formatTime(givenEvent))
}

func Test_coreLogger_formatTime_noop(t *testing.T) {
	provider := NewProvider(t, TimeFormat(NoopTimeFormat))
	provider.initIfRequired()
	instance := provider.coreRootLogger

	givenTs, _ := time.Parse(dateTimeFormat, "2024-07-25 18:56:13")
	givenEvent := instance.NewEvent(level.Info, map[string]interface{}{
		"timestamp": givenTs,
	})

	assert.ToBeEqual(t, "", instance.formatTime(givenEvent))
}

func Test_coreLogger_formatTime_ts(t *testing.T) {
	provider := NewProvider(t, TimeFormat(dateTimeFormat))
	provider.initIfRequired()
	instance := provider.coreRootLogger

	givenTs, _ := time.Parse(dateTimeFormat, "2024-07-25 18:56:13")
	givenEvent := instance.NewEvent(level.Info, map[string]interface{}{
		"timestamp": givenTs,
	})

	assert.ToBeEqual(t, "2024-07-25 18:56:13 ", instance.formatTime(givenEvent))
}

func Test_coreLogger_formatTime_ts_defaultNow(t *testing.T) {
	provider := NewProvider(t, TimeFormat(time.RFC3339))
	provider.initIfRequired()
	instance := provider.coreRootLogger

	givenEvent := instance.NewEvent(level.Info, map[string]interface{}{})

	now := time.Now()
	actualTs, actualErr := time.Parse(time.RFC3339+" ", instance.formatTime(givenEvent))
	assert.ToBeNoError(t, actualErr)

	diff := now.Sub(actualTs)
	if diff > time.Second*10 || diff < -time.Second*10 {
		t.Fatalf("the difference between %v(now) and %v(ts) should not be greater than 10s; bug was: %v", now, actualTs, diff)
	}
}

func Test_coreLogger_GetLevel_default(t *testing.T) {
	provider := NewProvider(t, TimeFormat(time.RFC3339))
	provider.initIfRequired()

	coreInstance := provider.getLogger(RootLoggerName)
	assert.ToBeEqual(t, level.Level(0), coreInstance.loggerNameToLevel[RootLoggerName])
	assert.ToBeEqual(t, level.Debug, coreInstance.GetLevel())

	otherInstance := provider.getLogger("other")
	assert.ToBeEqual(t, level.Level(0), otherInstance.loggerNameToLevel["other"])
	assert.ToBeEqual(t, level.Debug, otherInstance.GetLevel())
}

func Test_coreLogger_GetLevel_setDirect(t *testing.T) {
	provider := NewProvider(t, TimeFormat(time.RFC3339))
	provider.initIfRequired()

	coreInstance := provider.getLogger(RootLoggerName)
	coreInstance.loggerNameToLevel[RootLoggerName] = level.Warn
	assert.ToBeEqual(t, level.Warn, coreInstance.GetLevel())

	otherInstance := provider.getLogger("other")
	otherInstance.loggerNameToLevel["other"] = level.Error
	assert.ToBeEqual(t, level.Error, otherInstance.GetLevel())
}

func Test_coreLogger_SetLevel(t *testing.T) {
	provider := NewProvider(t, TimeFormat(time.RFC3339))
	provider.initIfRequired()

	coreInstance := provider.getLogger(RootLoggerName)
	assert.ToBeEqual(t, level.Level(0), coreInstance.loggerNameToLevel[RootLoggerName])
	coreInstance.SetLevel(level.Debug)
	assert.ToBeEqual(t, level.Debug, coreInstance.loggerNameToLevel[RootLoggerName])
	coreInstance.SetLevel(level.Info)
	assert.ToBeEqual(t, level.Info, coreInstance.loggerNameToLevel[RootLoggerName])
	coreInstance.SetLevel(0)
	assert.ToBeEqual(t, level.Level(0), coreInstance.loggerNameToLevel[RootLoggerName])

	otherInstance := provider.getLogger("other")
	otherInstance2 := provider.getLogger("other")
	assert.ToBeEqual(t, level.Level(0), otherInstance.loggerNameToLevel["other"])
	otherInstance.SetLevel(level.Debug)
	assert.ToBeEqual(t, level.Debug, otherInstance.loggerNameToLevel["other"])
	assert.ToBeEqual(t, level.Debug, otherInstance2.GetLevel())
	otherInstance.SetLevel(level.Info)
	assert.ToBeEqual(t, level.Info, otherInstance.loggerNameToLevel["other"])
	assert.ToBeEqual(t, level.Info, otherInstance2.GetLevel())
	otherInstance.SetLevel(0)
	assert.ToBeEqual(t, level.Level(0), otherInstance.loggerNameToLevel["other"])
	assert.ToBeEqual(t, provider.GetLevel(), otherInstance2.GetLevel())

	otherInstance.SetLevel(level.Error)

	otherInstance3 := provider.getLogger("other")
	assert.ToBeEqual(t, level.Error, otherInstance3.GetLevel())
}
