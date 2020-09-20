package main

import (
	log "github.com/echocat/slf4g"
	_ "github.com/echocat/slf4g/native"
	_ "github.com/echocat/slf4g/std/bridge/hook"
	stdlog "log"
)

func main() {
	log.With("foo", "bar").Debug("hello, debug")
	log.With("a", "foo").
		With("c", "xyz").
		With("b", "bar").
		With("d2", "abc").
		With("d1", "zzz").
		With("d3", "abc").
		Info("hello, info")
	log.With("foo", 1).Warn("hello, warn")
	log.With("foo", "bar/s").Error("hello, error")
	log.With("bar", 234).Error()
	log.With("bar", 234).Info("hello\nworld")
	log.Info("\nhello\nworld")

	stdlog.Print("abc")
}
