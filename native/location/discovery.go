package location

import log "github.com/echocat/slf4g"

var DefaultDiscovery = NewCallerAwareDiscovery()

type Discovery interface {
	DiscoveryLocation(event log.Event, callDepth int) Location
}

type DiscoveryFunc func(event log.Event, callDepth int) Location

func (instance DiscoveryFunc) DiscoveryLocation(event log.Event, callDepth int) Location {
	return instance(event, callDepth+1)
}

func NewDiscoveryFacade(provider func() Discovery) Discovery {
	return discoveryFacade(provider)
}

type discoveryFacade func() Discovery

func (instance discoveryFacade) DiscoveryLocation(event log.Event, callDepth int) Location {
	return instance().DiscoveryLocation(event, callDepth+1)
}

var noopV = DiscoveryFunc(func(log.Event, int) Location {
	return nil
})

func NoopDiscovery() Discovery {
	return noopV
}
