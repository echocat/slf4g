package main

import (
	log "github.com/echocat/slf4g"
	_ "github.com/echocat/slf4g/native"
)

func main() {
	log.With("foo", "bar").Debug("hello, debug")
	log.With("foo", "bar").Info("hello, info")
	log.With("foo", 1).Warn("hello, warn")
	log.With("foo", "bar/s").Error("hello, error")
	log.With("bar", 234).Error()
}
