package main

import (
	log "github.com/echocat/slf4g"
)

func main() {
	log.With("foo", "bar").Debug("hello, debug")
	log.With("foo", "bar").Info("hello, info")
	log.With("foo", 1).Warn("hello, warn")
	log.Error("hello, error")
}
