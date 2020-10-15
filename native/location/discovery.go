package location

import log "github.com/echocat/slf4g"

// DefaultDiscovery is the default instance of Discovery which should cover the
// majority of the cases.
var DefaultDiscovery = NewCallerDiscovery()

// Discovery is used to discover the Location where an log.Event happened.
type Discovery interface {
	// DiscoverLocation discovers the Location for the given log.Event.
	DiscoverLocation(event log.Event, extraCallDepth int) Location
}

// DiscoveryFunc is wrapping a given function into a Discovery.
type DiscoveryFunc func(event log.Event, extraCallDepth int) Location

// DiscoverLocation implements Discovery.DiscoverLocation()
func (instance DiscoveryFunc) DiscoverLocation(event log.Event, extraCallDepth int) Location {
	return instance(event, extraCallDepth+1)
}

// NewDiscoveryFacade creates a facade of KeysSpec using the given provider.
func NewDiscoveryFacade(provider func() Discovery) Discovery {
	return discoveryFacade(provider)
}

type discoveryFacade func() Discovery

func (instance discoveryFacade) DiscoverLocation(event log.Event, extraCallDepth int) Location {
	return instance().DiscoverLocation(event, extraCallDepth+1)
}

var noopV = DiscoveryFunc(func(log.Event, int) Location {
	return nil
})

// NoopDiscovery provides a noop implementation of Discovery.
func NoopDiscovery() Discovery {
	return noopV
}
