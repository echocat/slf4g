package sdk

import (
	"os"

	"github.com/echocat/slf4g/fields"

	"github.com/echocat/slf4g/level"

	log "github.com/echocat/slf4g"
)

// DefaultOnPanic defines what happens by default when someone calls one of the
// Logger.Panic(), Logger.Panicf() or Logger.Panicln() methods. The initial
// behavior will be that it panics after logging the event.
var DefaultOnPanic = func(e log.Event) {
	panic(e)
}

// DefaultOnFatal defines what happens by default when someone calls one of the
// Logger.Fatal(), Logger.Fatalf() or Logger.Fatalln() methods. The initial
// behavior will be that it exit with error code 1 after logging the event.
var DefaultOnFatal = func(log.Event) {
	os.Exit(1)
}

// LoggerImpl is a default implementation of the Logger interface. It forwards
// all logged events to the configured Delegate.
type LoggerImpl struct {
	// Delegate is the Logger of the slf4g framework where to forward all logged
	// events of this implementation to.
	Delegate log.CoreLogger

	// PrintLevel defines the regular level.Level to log everything one if
	// methods Print(), Printf() or Println() are used. If this is not set
	// level.Info is used.
	PrintLevel level.Level

	// OnPanic defines what to do if someone calls one of the Logger.Panic(),
	// Logger.Panicf() or Logger.Panicln() methods. Be default DefaultOnPanic
	// is used.
	OnPanic func(log.Event)

	// OnFatal defines what to do if someone calls one of the Logger.Fatal(),
	// Logger.Fatalf() or Logger.Fatalln() methods. Be default DefaultOnFatal
	// is used.
	OnFatal func(log.Event)
}

func (instance *LoggerImpl) log(l level.Level, args ...interface{}) (doLog func() log.Event, helper func()) {
	logger := instance.Delegate
	helper = helperOf(logger)
	helper()

	var values map[string]interface{}
	if len(args) > 0 {
		values = make(map[string]interface{}, 1)

		if len(args) == 1 {
			values[logger.GetProvider().GetFieldKeysSpec().GetMessage()] = args[0]
		} else if len(args) > 1 {
			values[logger.GetProvider().GetFieldKeysSpec().GetMessage()] = args
		}

	}
	e := logger.NewEvent(l, values)

	return func() log.Event {
		helper()
		logger.Log(e, 2)
		return e
	}, helper
}

func (instance *LoggerImpl) logf(l level.Level, format string, args ...interface{}) (doLog func() log.Event, helper func()) {
	logger := instance.Delegate
	helper = helperOf(logger)
	helper()

	e := logger.NewEvent(l, map[string]interface{}{
		logger.GetProvider().GetFieldKeysSpec().GetMessage(): fields.LazyFormat(format, args...),
	})

	return func() log.Event {
		helper()
		logger.Log(e, 2)
		return e
	}, helper
}

func (instance *LoggerImpl) printLevel() level.Level {
	if v := instance.PrintLevel; v > 0 {
		return v
	}
	return level.Info
}

// Print implements Logger.Print
func (instance *LoggerImpl) Print(args ...interface{}) {
	l, helper := instance.log(instance.printLevel(), args...)
	helper()
	l()
}

// Printf implements Logger.Printf
func (instance *LoggerImpl) Printf(s string, args ...interface{}) {
	l, helper := instance.logf(instance.printLevel(), s, args...)
	helper()
	l()
}

// Println implements Logger.Println
func (instance *LoggerImpl) Println(args ...interface{}) {
	l, helper := instance.log(instance.printLevel(), args...)
	helper()
	l()
}

// Fatal implements Logger.Fatal
func (instance *LoggerImpl) Fatal(args ...interface{}) {
	l, helper := instance.log(level.Fatal, args...)
	helper()
	instance.onFatal(l())
}

// Fatalf implements Logger.Fatalf
func (instance *LoggerImpl) Fatalf(s string, args ...interface{}) {
	l, helper := instance.logf(level.Fatal, s, args...)
	helper()
	instance.onFatal(l())
}

// Fatalln implements Logger.Fatalln
func (instance *LoggerImpl) Fatalln(args ...interface{}) {
	l, helper := instance.log(level.Fatal, args...)
	helper()
	instance.onFatal(l())
}

// Panic implements Logger.Panic
func (instance *LoggerImpl) Panic(args ...interface{}) {
	l, helper := instance.log(level.Fatal, args...)
	helper()
	instance.onPanic(l())
}

// Panicf implements Logger.Panicf
func (instance *LoggerImpl) Panicf(s string, args ...interface{}) {
	l, helper := instance.logf(level.Fatal, s, args...)
	helper()
	instance.onPanic(l())
}

// Panicln implements Logger.Panicln
func (instance *LoggerImpl) Panicln(args ...interface{}) {
	l, helper := instance.log(level.Fatal, args...)
	helper()
	instance.onPanic(l())
}

func (instance *LoggerImpl) onFatal(e log.Event) {
	if f := instance.OnFatal; f != nil {
		f(e)
	} else if f := DefaultOnFatal; f != nil {
		f(e)
	}
}

func (instance *LoggerImpl) onPanic(e log.Event) {
	if f := instance.OnPanic; f != nil {
		f(e)
	} else if f := DefaultOnPanic; f != nil {
		f(e)
	}
}

func helperOf(instance log.CoreLogger) func() {
	if wh, ok := instance.(interface {
		Helper() func()
	}); ok {
		return wh.Helper()
	}
	return func() {}
}
