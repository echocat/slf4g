package location

import log "github.com/echocat/slf4g"

var DefaultFactory Factory = nil

type Location interface{}

type Factory func(event log.Event, callDepth int) Location

func NoopFactory(log.Event, int) Location {
	return nil
}
