package recording

import (
	"errors"
	"sync"
	"testing"
	"time"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_NewCoreLogger(t *testing.T) {
	actual := NewCoreLogger()

	assert.ToBeEqual(t, &CoreLogger{}, actual)
}

func Test_CoreLogger_MustContains(t *testing.T) {
	instance := NewLogger()

	expectedEvent1 := instance.NewEvent(level.Info, nil).
		With("message", "a")
	expectedEvent2 := instance.NewEvent(level.Warn, nil).
		With("message", "b").
		With("foo", 1)
	unExpectedEvent3 := instance.NewEvent(level.Info, nil).
		With("message", "b").
		With("foo", 1)

	assert.ToBeEqual(t, 0, instance.Len())

	instance.Info("a")
	instance.With("foo", 1).Warn("b")

	actual1 := instance.MustContains(expectedEvent1)
	assert.ToBeEqual(t, true, actual1)

	actual2 := instance.MustContains(expectedEvent2)
	assert.ToBeEqual(t, true, actual2)

	actual3 := instance.MustContains(unExpectedEvent3)
	assert.ToBeEqual(t, false, actual3)
}

func Test_CoreLogger_MustContains_panicsOnErrorsForRoot(t *testing.T) {
	previous := fields.DefaultValueEquality
	defer func() { fields.DefaultValueEquality = previous }()
	fields.DefaultValueEquality = fields.ValueEqualityFunc(func(name string, left, right interface{}) (bool, error) {
		return false, errors.New("expected")
	})

	instance := NewLogger()
	instance.Info("something")

	givenEvent := instance.NewEvent(level.Info, nil).
		With("message", "something")

	assert.Execution(t, func() {
		instance.MustContains(givenEvent)
	}).WillPanicWith("^expected$")
}

func Test_CoreLogger_MustContainsCustom(t *testing.T) {
	instance := NewLogger()
	givenEquality := log.DefaultEventEquality.WithIgnoringKeys("foo", "timestamp", "logger")

	expectedEvent1 := instance.NewEvent(level.Info, nil).
		With("message", "a").
		With("foo", 666).
		With("bar", 2)
	expectedEvent2 := instance.NewEvent(level.Warn, nil).
		With("message", "b").
		With("foo", 666).
		With("bar", 2)
	unExpectedEvent3 := instance.NewEvent(level.Info, nil).
		With("message", "b").
		With("foo", 666).
		With("bar", 2)

	assert.ToBeEqual(t, 0, instance.Len())

	instance.
		With("foo", 1).
		With("bar", 2).
		Info("a")
	instance.
		With("foo", 1).
		With("bar", 2).
		Warn("b")

	assert.ToBeEqual(t, 2, instance.Len())

	actual1 := instance.MustContainsCustom(givenEquality, expectedEvent1)
	assert.ToBeEqual(t, true, actual1)

	actual2 := instance.MustContainsCustom(givenEquality, expectedEvent2)
	assert.ToBeEqual(t, true, actual2)

	actual3 := instance.MustContainsCustom(givenEquality, unExpectedEvent3)
	assert.ToBeEqual(t, false, actual3)
}

func Test_CoreLogger_MustContainsCustom_panics(t *testing.T) {
	givenEquality := log.EventEqualityFunc(func(left, right log.Event) (bool, error) {
		return false, errors.New("expected")
	})

	instance := NewLogger()
	instance.Info("foo")

	givenEvent := instance.NewEvent(level.Info, nil).
		With("message", "foo")

	assert.Execution(t, func() {
		instance.MustContainsCustom(givenEquality, givenEvent)
	}).WillPanicWith("^expected$")
}

func Test_CoreLogger_ContainsCustom_doesNotHoldLockDuringEquality(t *testing.T) {
	instance := NewLogger()
	instance.Info("a")
	expected := instance.NewEvent(level.Info, nil).With("message", "a")
	entered := make(chan struct{})
	release := make(chan struct{})
	result := make(chan struct {
		matches bool
		err     error
	}, 1)
	updatesDone := make(chan struct{})
	var releaseOnce sync.Once
	releaseEquality := func() {
		releaseOnce.Do(func() { close(release) })
	}
	defer releaseEquality()

	go func() {
		actual, err := instance.ContainsCustom(log.EventEqualityFunc(func(left, right log.Event) (bool, error) {
			close(entered)
			<-release
			return true, nil
		}), expected)
		result <- struct {
			matches bool
			err     error
		}{actual, err}
	}()
	waitForSignal(t, entered, "event equality was not called")
	go func() {
		instance.Reset()
		instance.Info("b")
		close(updatesDone)
	}()

	blocked := false
	select {
	case <-updatesDone:
	case <-time.After(time.Second):
		blocked = true
	}
	releaseEquality()
	var actual struct {
		matches bool
		err     error
	}
	select {
	case actual = <-result:
	case <-time.After(time.Second):
		t.Fatal("event equality did not complete after being released")
	}
	if blocked {
		waitForSignal(t, updatesDone, "recorder updates did not complete after event equality was released")
		t.Fatal("recorder updates blocked while event equality was running")
	}

	assert.ToBeNoError(t, actual.err)
	assert.ToBeEqual(t, true, actual.matches)
	assert.ToBeEqual(t, 1, instance.Len())
	assert.ToBeEqual(t, true, instance.MustContains(instance.NewEvent(level.Info, nil).With("message", "b")))
}

func Test_CoreLogger_Len(t *testing.T) {
	instance := NewLogger()

	assert.ToBeEqual(t, 0, instance.Len())

	instance.
		With("foo", 1).
		With("bar", 2).
		Info("a")
	instance.
		With("foo", 1).
		With("bar", 2).
		Warn("b")

	assert.ToBeEqual(t, 2, instance.Len())
}

func Test_CoreLogger_Log(t *testing.T) {
	instance := NewLogger()

	expectedEvent1 := instance.NewEvent(level.Info, nil).
		With("message", "a")
	expectedEvent2 := instance.NewEvent(level.Warn, nil).
		With("message", "b").
		With("foo", 1)
	unExpectedEvent3 := instance.NewEvent(level.Debug, nil).
		With("message", "c").
		With("foo", 1)

	instance.Log(expectedEvent1, 2)
	instance.Log(expectedEvent2, 2)
	instance.Log(unExpectedEvent3, 2)

	assert.ToBeEqual(t, 2, len(instance.recorded))
	assert.ToBeEqualUsing(t, expectedEvent1, instance.recorded[0], instance.defaultEventEquality().AreEventsEqual)
	assert.ToBeEqualUsing(t, expectedEvent2, instance.recorded[1], instance.defaultEventEquality().AreEventsEqual)
}

func Test_CoreLogger_Log_doesNotHoldLockDuringLazyEvaluation(t *testing.T) {
	instance := NewLogger()
	instance.Info("old")
	entered := make(chan struct{})
	release := make(chan struct{})
	logDone := make(chan struct{})
	resetDone := make(chan struct{})
	var releaseOnce sync.Once
	releaseLazy := func() {
		releaseOnce.Do(func() { close(release) })
	}
	defer releaseLazy()
	timestampKey := instance.GetProvider().GetFieldKeysSpec().GetTimestamp()
	event := instance.NewEvent(level.Info, map[string]interface{}{
		timestampKey: fields.LazyFunc(func() interface{} {
			close(entered)
			<-release
			return time.Now()
		}),
	})

	go func() {
		instance.Log(event, 0)
		close(logDone)
	}()
	waitForSignal(t, entered, "lazy event value was not evaluated")
	go func() {
		instance.Reset()
		close(resetDone)
	}()

	blocked := false
	select {
	case <-resetDone:
	case <-time.After(time.Second):
		blocked = true
	}
	releaseLazy()
	waitForSignal(t, logDone, "logging did not complete after lazy evaluation was released")
	if blocked {
		waitForSignal(t, resetDone, "recorder reset did not complete after lazy evaluation was released")
		t.Fatal("recorder reset blocked while a lazy event value was evaluated")
	}

	assert.ToBeEqual(t, 1, instance.Len())
}

func Test_CoreLogger_GetAll(t *testing.T) {
	instance := NewLogger()

	assert.ToBeEqual(t, 0, instance.Len())

	instance.
		With("foo", 1).
		With("bar", 2).
		Info("a")
	instance.
		With("foo", 1).
		With("bar", 2).
		Warn("b")

	actual := instance.GetAll()

	assert.ToBeEqual(t, 2, len(actual))
	assert.ToBeEqualUsing(t, instance.NewEvent(level.Info, nil).
		With("message", "a").
		With("foo", 1).
		With("bar", 2), actual[0], instance.defaultEventEquality().AreEventsEqual)
	assert.ToBeEqualUsing(t, instance.NewEvent(level.Warn, nil).
		With("message", "b").
		With("foo", 1).
		With("bar", 2), actual[1], instance.defaultEventEquality().AreEventsEqual)
}

func Test_CoreLogger_Get(t *testing.T) {
	instance := NewLogger()

	assert.ToBeEqual(t, 0, instance.Len())

	instance.
		With("foo", 1).
		With("bar", 2).
		Info("a")
	instance.
		With("foo", 1).
		With("bar", 2).
		Warn("b")

	assert.ToBeEqual(t, 2, instance.Len())
	actual0 := instance.Get(0)
	actual1 := instance.Get(1)

	assert.ToBeEqualUsing(t, instance.NewEvent(level.Info, nil).
		With("message", "a").
		With("foo", 1).
		With("bar", 2), actual0, instance.defaultEventEquality().AreEventsEqual)
	assert.ToBeEqualUsing(t, instance.NewEvent(level.Warn, nil).
		With("message", "b").
		With("foo", 1).
		With("bar", 2), actual1, instance.defaultEventEquality().AreEventsEqual)
}

func Test_CoreLogger_Get_panicsOnOutOfRange(t *testing.T) {
	instance := NewLogger()

	assert.ToBeEqual(t, 0, instance.Len())

	instance.
		With("foo", 1).
		With("bar", 2).
		Info("a")
	instance.
		With("foo", 1).
		With("bar", 2).
		Warn("b")

	assert.ToBeEqual(t, 2, instance.Len())

	assert.Execution(t, func() {
		instance.Get(2)
	}).WillPanicWith("^Index 2 requested but the amount of recorded events is only 2$")
}

func Test_CoreLogger_Reset(t *testing.T) {
	instance := NewLogger()

	assert.ToBeEqual(t, 0, instance.Len())

	instance.
		With("foo", 1).
		With("bar", 2).
		Info("a")
	instance.
		With("foo", 1).
		With("bar", 2).
		Warn("b")

	assert.ToBeEqual(t, 2, instance.Len())

	instance.Reset()

	assert.ToBeEqual(t, 0, instance.Len())
}

func Test_CoreLogger_GetName_specified(t *testing.T) {
	instance := NewCoreLogger()
	instance.Name = "foo"

	assert.ToBeEqual(t, "foo", instance.GetName())
}

func Test_CoreLogger_GetName_absent(t *testing.T) {
	instance := NewCoreLogger()

	assert.ToBeEqual(t, RootLoggerName, instance.GetName())
}

func Test_CoreLogger_GetLevel_specified(t *testing.T) {
	instance := NewCoreLogger()
	instance.Level = level.Warn

	assert.ToBeEqual(t, level.Warn, instance.GetLevel())
}

func Test_CoreLogger_GetLevel_absent(t *testing.T) {
	instance := NewCoreLogger()

	assert.ToBeEqual(t, level.Info, instance.GetLevel())
}

func Test_CoreLogger_GetLevel(t *testing.T) {
	instance := NewCoreLogger()

	assert.ToBeEqual(t, level.Info, instance.GetLevel())

	for _, l := range level.GetProvider().GetLevels() {
		instance.Level = l
		assert.ToBeEqual(t, l, instance.GetLevel())
	}

	instance.Level = 0
	assert.ToBeEqual(t, level.Info, instance.GetLevel())
}

func Test_CoreLogger_SetLevel(t *testing.T) {
	instance := NewCoreLogger()

	assert.ToBeEqual(t, level.Level(0), instance.Level)

	for _, l := range level.GetProvider().GetLevels() {
		instance.SetLevel(l)
		assert.ToBeEqual(t, l, instance.Level)
	}

	instance.SetLevel(0)
	assert.ToBeEqual(t, level.Level(0), instance.Level)
}

func Test_CoreLogger_NewEvent(t *testing.T) {
	instance := NewCoreLogger()

	assert.ToBeEqual(t, &event{
		provider: instance.Provider,
		fields:   fields.Empty(),
		level:    level.Fatal,
	}, instance.NewEvent(level.Fatal, nil))

	assert.ToBeEqual(t, &event{
		provider: instance.Provider,
		fields:   fields.WithAll(map[string]interface{}{"foo": "bar"}),
		level:    level.Fatal,
	}, instance.NewEvent(level.Fatal, map[string]interface{}{"foo": "bar"}))
}

func Test_CoreLogger_NewEventWithFields(t *testing.T) {
	instance := NewCoreLogger()

	assert.ToBeEqual(t, &event{
		provider: instance.Provider,
		fields:   fields.Empty(),
		level:    level.Fatal,
	}, instance.NewEventWithFields(level.Fatal, nil))

	assert.ToBeEqual(t, &event{
		provider: instance.Provider,
		fields:   fields.With("foo", "bar"),
		level:    level.Fatal,
	}, instance.NewEventWithFields(level.Fatal, fields.With("foo", "bar")))
}

func Test_CoreLogger_NewEventWithFields_panicsOnError(t *testing.T) {
	instance := NewCoreLogger()

	assert.Execution(t, func() {
		instance.NewEventWithFields(level.Fatal, fields.ForEachFunc(func(func(string, interface{}) error) error {
			return errors.New("expected")
		}))
	}).WillPanicWith("^expected$")
}

func waitForSignal(t *testing.T, signal <-chan struct{}, failure string) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(time.Second):
		t.Fatal(failure)
	}
}

func Test_CoreLogger_Accepts(t *testing.T) {
	instance := NewCoreLogger()

	assert.ToBeEqual(t, false, instance.Accepts(nil))
	assert.ToBeEqual(t, true, instance.Accepts(&event{}))
}
