package std

import (
	"fmt"
	"os"

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

func (instance *LoggerImpl) log(level log.Level, args ...interface{}) log.Event {
	f := fields.Empty()
	if len(args) > 0 {
		f = fields.With(instance.GetProvider().GetFieldKeysSpec().GetMessage(), fields.LazyFunc(func() interface{} {
			return fmt.Sprint(args...)
		}))
	}
	e := log.NewEvent(level, f, 3)
	instance.Log(e)
	return e
}

func (instance *LoggerImpl) logf(level log.Level, format string, args ...interface{}) log.Event {
	f := fields.Withf(instance.GetProvider().GetFieldKeysSpec().GetMessage(), format, args...)
	e := log.NewEvent(level, f, 3)
	instance.Log(e)
	return e
}

func (instance *LoggerImpl) Print(args ...interface{}) {
	instance.log(log.LevelInfo, args...)
}

func (instance *LoggerImpl) Printf(s string, args ...interface{}) {
	instance.logf(log.LevelInfo, s, args...)
}

func (instance *LoggerImpl) Println(args ...interface{}) {
	instance.log(log.LevelInfo, args...)
}

func (instance *LoggerImpl) Fatal(args ...interface{}) {
	e := instance.log(log.LevelFatal, args...)
	instance.onFatal(e)
}

func (instance *LoggerImpl) Fatalf(s string, args ...interface{}) {
	e := instance.logf(log.LevelInfo, s, args...)
	instance.onFatal(e)
}

func (instance *LoggerImpl) Fatalln(args ...interface{}) {
	e := instance.log(log.LevelFatal, args...)
	instance.onFatal(e)
}

func (instance *LoggerImpl) Panic(args ...interface{}) {
	e := instance.log(log.LevelFatal, args...)
	instance.onPanic(e)
}

func (instance *LoggerImpl) Panicf(s string, args ...interface{}) {
	e := instance.logf(log.LevelInfo, s, args...)
	instance.onPanic(e)
}

func (instance *LoggerImpl) Panicln(args ...interface{}) {
	e := instance.log(log.LevelFatal, args...)
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
