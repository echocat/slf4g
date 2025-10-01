package main

import (
	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	_ "github.com/echocat/slf4g/native"
	"github.com/echocat/slf4g/native/formatter"
	"github.com/echocat/slf4g/native/location"
)

func main() {
	formatter.Default = formatter.NewText()
	location.DefaultDiscovery = location.NewCallerDiscovery()

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
	log.WithAll(map[string]interface{}{
		"a": "1",
		"b": "2",
		"c": fields.LazyFunc(func() interface{} {
			return "3"
		}),
		"excluded": fields.Exclude,
		"excludedLazy": fields.LazyFunc(func() interface{} {
			return fields.Exclude
		}),
		"nil": nil,
		"nilLazy": fields.LazyFunc(func() interface{} {
			return nil
		}),
		"empty": "",
		"emptyLazy": fields.LazyFunc(func() interface{} {
			return ""
		}),
	}).Info("Some more variants in a map")
}
