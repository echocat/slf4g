package main

import (
	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

func main() {
	log.With("foo", "bar").
		Debug("hello, debug")

	log.With("foo", "bar").
		Info("hello, info")

	log.With("foo", "bar").
		With("filtered1", fields.RequireMaximalLevel(level.Info, "visibleUntilLevelInfo")).
		With("filtered2", fields.RequireMaximalLevel(level.Debug, "visibleUntilLevelDebug")).
		Info()

	log.With("foo", 1).
		Warn("hello, warn")

	log.Error("hello, error")
}
