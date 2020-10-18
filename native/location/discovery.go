package location

import log "github.com/echocat/slf4g"

// DefaultDiscovery is the default instance of Discovery which should cover the
// majority of the cases.
var DefaultDiscovery = NoopDiscovery()

// Discovery is used to discover the Location where an log.Event happened.
type Discovery interface {
	// DiscoverLocation discovers the Location for the given log.Event.
	DiscoverLocation(event log.Event, skipFrames uint16) Location
}

// DiscoveryFunc is wrapping a given function into a Discovery.
type DiscoveryFunc func(event log.Event, skipFrames uint16) Location

// DiscoverLocation implements Discovery.DiscoverLocation()
func (instance DiscoveryFunc) DiscoverLocation(event log.Event, skipFrames uint16) Location {
	return instance(event, skipFrames+1)
}

// NewDiscoveryFacade creates a facade of KeysSpec using the given provider.
func NewDiscoveryFacade(provider func() Discovery) Discovery {
	return discoveryFacade(provider)
}

type discoveryFacade func() Discovery

func (instance discoveryFacade) DiscoverLocation(event log.Event, skipFrames uint16) Location {
	return instance().DiscoverLocation(event, skipFrames+1)
}

var noopV = DiscoveryFunc(func(log.Event, uint16) Location {
	return nil
})

// NoopDiscovery provides a noop implementation of Discovery.
func NoopDiscovery() Discovery {
	return noopV
}
