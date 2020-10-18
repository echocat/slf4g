package sdk

import (
	"os"

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

func (instance *LoggerImpl) log(l level.Level, args ...interface{}) log.Event {
	provider := instance.Delegate.GetProvider()

	e := log.NewEvent(provider, l)
	if len(args) == 1 {
		e = e.With(provider.GetFieldKeysSpec().GetMessage(), args[0])
	} else if len(args) > 1 {
		e = e.With(provider.GetFieldKeysSpec().GetMessage(), args)
	}

	instance.Delegate.Log(e, 2)
	return e
}

func (instance *LoggerImpl) logf(l level.Level, format string, args ...interface{}) log.Event {
	provider := instance.Delegate.GetProvider()

	e := log.NewEvent(provider, l).
		Withf(provider.GetFieldKeysSpec().GetMessage(), format, args...)

	instance.Delegate.Log(e, 2)
	return e
}

func (instance *LoggerImpl) printLevel() level.Level {
	if v := instance.PrintLevel; v > 0 {
		return v
	}
	return level.Info
}

// Print implements Logger.Print
func (instance *LoggerImpl) Print(args ...interface{}) {
	instance.log(instance.printLevel(), args...)
}

// Printf implements Logger.Printf
func (instance *LoggerImpl) Printf(s string, args ...interface{}) {
	instance.logf(instance.printLevel(), s, args...)
}

// Println implements Logger.Println
func (instance *LoggerImpl) Println(args ...interface{}) {
	instance.log(instance.printLevel(), args...)
}

// Fatal implements Logger.Fatal
func (instance *LoggerImpl) Fatal(args ...interface{}) {
	e := instance.log(level.Fatal, args...)
	instance.onFatal(e)
}

// Fatalf implements Logger.Fatalf
func (instance *LoggerImpl) Fatalf(s string, args ...interface{}) {
	e := instance.logf(level.Fatal, s, args...)
	instance.onFatal(e)
}

// Fatalln implements Logger.Fatalln
func (instance *LoggerImpl) Fatalln(args ...interface{}) {
	e := instance.log(level.Fatal, args...)
	instance.onFatal(e)
}

// Panic implements Logger.Panic
func (instance *LoggerImpl) Panic(args ...interface{}) {
	e := instance.log(level.Fatal, args...)
	instance.onPanic(e)
}

// Panicf implements Logger.Panicf
func (instance *LoggerImpl) Panicf(s string, args ...interface{}) {
	e := instance.logf(level.Fatal, s, args...)
	instance.onPanic(e)
}

// Panicln implements Logger.Panicln
func (instance *LoggerImpl) Panicln(args ...interface{}) {
	e := instance.log(level.Fatal, args...)
	instance.onPanic(e)
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
