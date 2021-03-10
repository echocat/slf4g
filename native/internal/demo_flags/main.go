package main

import (
	"flag"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native"
	_ "github.com/echocat/slf4g/native"
	"github.com/echocat/slf4g/native/facade/value"
	_ "github.com/echocat/slf4g/sdk/bridge/hook"
)

func main() {
	pv := value.NewProvider(native.DefaultProvider)
	flag.Var(pv.Consumer.Formatter, "log.format", "configures the log format.")
	flag.Var(pv.Level, "log.level", "configures the log level.")
	flag.Parse()

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
}
