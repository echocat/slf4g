package sdk

import (
	"testing"

	"github.com/echocat/slf4g/testing/recording"

	"github.com/echocat/slf4g/internal/test/assert"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/level"
)

func Test_DefaultOnPanic(t *testing.T) {
	givenEvent := log.GetRootLogger().NewEvent(level.Warn, nil)

	defer func() {
		r := recover()
		if re, ok := r.(log.Event); ok {
			if givenEvent != re {
				assert.Failf(t, "Expected to panics with <%+v>; but got: <%+v>", givenEvent, re)
			}
		} else if r != nil {
			assert.Failf(t, "Expected to panics with <%+v>; but got: <%+v>", givenEvent, r)
		} else {
			assert.Failf(t, "Expected to panics with <%+v>; but it didn't", givenEvent)
		}
	}()

	DefaultOnPanic(givenEvent)
}

func Test_LoggerImpl_Print_withNoArgs(t *testing.T) {
	instance, logger, horror := prepareLoggerImpl()
	instance.PrintLevel = level.Warn

	instance.Print()

	horror.nothingShouldBeHappen(t)
	assert.ToBeEqual(t, 1, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(
		logger.NewEvent(level.Warn, nil),
	))
}

func Test_LoggerImpl_Print_configuredLevel(t *testing.T) {
	instance, logger, horror := prepareLoggerImpl()
	instance.PrintLevel = level.Warn

	instance.Print()

	horror.nothingShouldBeHappen(t)
	assert.ToBeEqual(t, 1, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(
		logger.NewEvent(level.Warn, nil),
	))
}

func Test_LoggerImpl_Print_noConfiguredLevel(t *testing.T) {
	instance, logger, horror := prepareLoggerImpl()

	instance.Print()

	horror.nothingShouldBeHappen(t)
	assert.ToBeEqual(t, 1, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(
		logger.NewEvent(level.Info, nil),
	))
}

func Test_LoggerImpl_Print_with1Arg(t *testing.T) {
	instance, logger, horror := prepareLoggerImpl()
	instance.PrintLevel = level.Warn

	instance.Print("a")

	horror.nothingShouldBeHappen(t)
	assert.ToBeEqual(t, 1, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(
		logger.NewEvent(level.Warn, nil).
			With("message", "a"),
	))
}

func Test_LoggerImpl_Print_with3Arg(t *testing.T) {
	instance, logger, horror := prepareLoggerImpl()
	instance.PrintLevel = level.Warn

	instance.Print("a", 1, "c")

	horror.nothingShouldBeHappen(t)
	assert.ToBeEqual(t, 1, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(
		logger.NewEvent(level.Warn, nil).
			With("message", []interface{}{"a", 1, "c"}),
	))
}

func Test_LoggerImpl_Printf(t *testing.T) {
	instance, logger, horror := prepareLoggerImpl()
	instance.PrintLevel = level.Error

	instance.Printf("fmt %d %s", 1, "c")

	horror.nothingShouldBeHappen(t)
	assert.ToBeEqual(t, 1, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(
		logger.NewEvent(level.Error, nil).
			Withf("message", "fmt %d %s", 1, "c"),
	))
}

func Test_LoggerImpl_Println(t *testing.T) {
	instance, logger, horror := prepareLoggerImpl()
	instance.PrintLevel = level.Info

	instance.Println("a", 1, "c")

	horror.nothingShouldBeHappen(t)
	assert.ToBeEqual(t, 1, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(
		logger.NewEvent(level.Info, nil).
			With("message", []interface{}{"a", 1, "c"}),
	))
}

func Test_LoggerImpl_Fatal(t *testing.T) {
	instance, logger, horror := prepareLoggerImpl()
	expected := logger.NewEvent(level.Fatal, nil).
		With("message", []interface{}{"a", 1, "c"})

	instance.Fatal("a", 1, "c")

	assert.ToBeEqual(t, 1, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(expected))
	assert.ToBeNil(t, horror.panicCalledWith)
	assert.ToBeEqualUsing(t, horror.fatalCalledWith, expected, log.AreEventsEqual)
}

func Test_LoggerImpl_Fatalf(t *testing.T) {
	instance, logger, horror := prepareLoggerImpl()
	expected := logger.NewEvent(level.Fatal, nil).
		Withf("message", "fmt %d %s", 1, "c")

	instance.Fatalf("fmt %d %s", 1, "c")

	assert.ToBeEqual(t, 1, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(expected))
	assert.ToBeNil(t, horror.panicCalledWith)
	assert.ToBeEqualUsing(t, horror.fatalCalledWith, expected, log.AreEventsEqual)
}

func Test_LoggerImpl_Fatalln(t *testing.T) {
	instance, logger, horror := prepareLoggerImpl()
	expected := logger.NewEvent(level.Fatal, nil).
		With("message", []interface{}{"a", 1, "c"})

	instance.Fatalln("a", 1, "c")

	assert.ToBeEqual(t, 1, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(expected))
	assert.ToBeNil(t, horror.panicCalledWith)
	assert.ToBeEqualUsing(t, horror.fatalCalledWith, expected, log.AreEventsEqual)
}

func Test_LoggerImpl_Panic(t *testing.T) {
	instance, logger, horror := prepareLoggerImpl()
	expected := logger.NewEvent(level.Fatal, nil).
		With("message", []interface{}{"a", 1, "c"})

	instance.Panic("a", 1, "c")

	assert.ToBeEqual(t, 1, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(expected))
	assert.ToBeNil(t, horror.fatalCalledWith)
	assert.ToBeEqualUsing(t, horror.panicCalledWith, expected, log.AreEventsEqual)
}

func Test_LoggerImpl_Panicf(t *testing.T) {
	instance, logger, horror := prepareLoggerImpl()
	expected := logger.NewEvent(level.Fatal, nil).
		Withf("message", "fmt %d %s", 1, "c")

	instance.Panicf("fmt %d %s", 1, "c")

	assert.ToBeEqual(t, 1, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(expected))
	assert.ToBeNil(t, horror.fatalCalledWith)
	assert.ToBeEqualUsing(t, horror.panicCalledWith, expected, log.AreEventsEqual)
}

func Test_LoggerImpl_Panicln(t *testing.T) {
	instance, logger, horror := prepareLoggerImpl()
	expected := logger.NewEvent(level.Fatal, nil).
		With("message", []interface{}{"a", 1, "c"})

	instance.Panicln("a", 1, "c")

	assert.ToBeEqual(t, 1, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(expected))
	assert.ToBeNil(t, horror.fatalCalledWith)
	assert.ToBeEqualUsing(t, horror.panicCalledWith, expected, log.AreEventsEqual)
}

func Test_LoggerImpl_onFatal_configured(t *testing.T) {
	givenEvent := log.GetRootLogger().NewEvent(level.Fatal, nil).
		Withf("message", "fmt %d %s", 1, "c")

	called := false

	old := DefaultOnFatal
	defer func() { DefaultOnFatal = old }()
	DefaultOnFatal = func(log.Event) { panic("should not be called ") }

	instance := LoggerImpl{
		OnFatal: func(event log.Event) {
			assert.ToBeSame(t, givenEvent, event)
			called = true
		},
	}

	instance.onFatal(givenEvent)

	assert.ToBeEqual(t, true, called)
}

func Test_LoggerImpl_onFatal_notConfigured(t *testing.T) {
	givenEvent := log.GetRootLogger().NewEvent(level.Fatal, nil).
		Withf("message", "fmt %d %s", 1, "c")

	called := false

	old := DefaultOnFatal
	defer func() { DefaultOnFatal = old }()
	DefaultOnFatal = func(event log.Event) {
		assert.ToBeSame(t, givenEvent, event)
		called = true
	}

	instance := LoggerImpl{}

	instance.onFatal(givenEvent)

	assert.ToBeEqual(t, true, called)
}

func Test_LoggerImpl_onPanic_configured(t *testing.T) {
	givenEvent := log.GetRootLogger().NewEvent(level.Fatal, nil).
		Withf("message", "fmt %d %s", 1, "c")

	called := false

	old := DefaultOnPanic
	defer func() { DefaultOnPanic = old }()
	DefaultOnPanic = func(log.Event) { panic("should not be called ") }

	instance := LoggerImpl{
		OnPanic: func(event log.Event) {
			assert.ToBeSame(t, givenEvent, event)
			called = true
		},
	}

	instance.onPanic(givenEvent)

	assert.ToBeEqual(t, true, called)
}

func Test_LoggerImpl_onPanic_notConfigured(t *testing.T) {
	givenEvent := log.GetRootLogger().NewEvent(level.Fatal, nil).
		Withf("message", "fmt %d %s", 1, "c")

	called := false

	old := DefaultOnPanic
	defer func() { DefaultOnPanic = old }()
	DefaultOnPanic = func(event log.Event) {
		assert.ToBeSame(t, givenEvent, event)
		called = true
	}

	instance := LoggerImpl{}

	instance.onPanic(givenEvent)

	assert.ToBeEqual(t, true, called)
}

func prepareLoggerImpl() (*LoggerImpl, *recording.CoreLogger, *horrorEventsHook) {
	recorder := recording.NewCoreLogger()
	result := &LoggerImpl{
		Delegate: recorder,
	}
	horror := configureWithHorrorEvents(result)
	return result, recorder, horror
}

func configureWithHorrorEvents(l *LoggerImpl) *horrorEventsHook {
	result := &horrorEventsHook{}
	l.OnFatal = result.onFatal
	l.OnPanic = result.onPanic
	return result
}

type horrorEventsHook struct {
	panicCalledWith log.Event
	fatalCalledWith log.Event
}

func (instance *horrorEventsHook) onPanic(event log.Event) {
	if event == nil {
		event = log.GetRootLogger().NewEvent(level.Fatal, nil).
			With("message", "NIL EVENT")
	}
	instance.panicCalledWith = event
}

func (instance *horrorEventsHook) onFatal(event log.Event) {
	if event == nil {
		event = log.GetRootLogger().NewEvent(level.Fatal, nil).
			With("message", "NIL EVENT")
	}
	instance.fatalCalledWith = event
}

func (instance *horrorEventsHook) nothingShouldBeHappen(t *testing.T) {
	t.Helper()
	assert.ToBeNil(t, instance.panicCalledWith)
	assert.ToBeNil(t, instance.fatalCalledWith)
}
