package std

import (
	"fmt"
	"os"

	"github.com/echocat/slf4g/level"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
)

var DefaultOnPanic = func(e log.Event) {
	if se, ok := e.(fmt.Stringer); ok {
		panic(se.String())
	} else {
		panic(fmt.Sprintf("%+v", e))
	}
}

var DefaultOnFatal = func(log.Event) {
	os.Exit(1)
}

type LoggerImpl struct {
	log.CoreLogger

	OnPanic func(log.Event)
	OnFatal func(log.Event)
}

func (instance *LoggerImpl) log(l level.Level, args ...interface{}) log.Event {
	provider := instance.GetProvider()

	e := log.NewEvent(provider, l, 3)
	if len(args) > 0 {
		e = e.With(provider.GetFieldKeysSpec().GetMessage(), fields.LazyFunc(func() interface{} {
			return fmt.Sprint(args...)
		}))
	}

	instance.Log(e)
	return e
}

func (instance *LoggerImpl) logf(l level.Level, format string, args ...interface{}) log.Event {
	provider := instance.GetProvider()

	e := log.NewEvent(provider, l, 3).
		Withf(provider.GetFieldKeysSpec().GetMessage(), format, args...)

	instance.Log(e)
	return e
}

func (instance *LoggerImpl) Print(args ...interface{}) {
	instance.log(level.Info, args...)
}

func (instance *LoggerImpl) Printf(s string, args ...interface{}) {
	instance.logf(level.Info, s, args...)
}

func (instance *LoggerImpl) Println(args ...interface{}) {
	instance.log(level.Info, args...)
}

func (instance *LoggerImpl) Fatal(args ...interface{}) {
	e := instance.log(level.Fatal, args...)
	instance.onFatal(e)
}

func (instance *LoggerImpl) Fatalf(s string, args ...interface{}) {
	e := instance.logf(level.Info, s, args...)
	instance.onFatal(e)
}

func (instance *LoggerImpl) Fatalln(args ...interface{}) {
	e := instance.log(level.Fatal, args...)
	instance.onFatal(e)
}

func (instance *LoggerImpl) Panic(args ...interface{}) {
	e := instance.log(level.Fatal, args...)
	instance.onPanic(e)
}

func (instance *LoggerImpl) Panicf(s string, args ...interface{}) {
	e := instance.logf(level.Info, s, args...)
	instance.onPanic(e)
}

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
