package location

import log "github.com/echocat/slf4g"

const Field = "location"

var DefaultFactory Factory = nil

type Location interface{}

type Factory func(event log.Event, callDepth int) Location

func NoopFactory(log.Event, int) Location {
	return nil
}

type FactoryAware interface {
	GetLocationFactory() Factory
	SetLocationFactory(Factory)
}
