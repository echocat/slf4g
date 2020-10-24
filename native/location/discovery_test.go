package location

import (
	"testing"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
)

func Test_DiscoveryFunc_DiscoverLocation(t *testing.T) {
	givenEvent := newEvent(level.Warn, nil)
	givenLocation := &callerImpl{}

	instance := DiscoveryFunc(func(actualEvent log.Event, actualSkipFrames uint16) Location {
		assert.ToBeSame(t, givenEvent, actualEvent)
		assert.ToBeEqual(t, uint16(667), actualSkipFrames)

		return givenLocation
	})

	actual := instance.DiscoverLocation(givenEvent, 666)

	assert.ToBeSame(t, givenLocation, actual)
}

func Test_NewDiscoveryFacade(t *testing.T) {
	delegate := DiscoveryFunc(func(actualEvent log.Event, actualSkipFrames uint16) Location {
		panic("should not be called")
	})

	actual := NewDiscoveryFacade(func() Discovery {
		return delegate
	})

	assert.ToBeSame(t, delegate, actual.(discoveryFacade)())
}

func Test_discoveryFacade_DiscoverLocation(t *testing.T) {
	givenEvent := newEvent(level.Warn, nil)
	givenLocation := &callerImpl{}

	delegate := DiscoveryFunc(func(actualEvent log.Event, actualSkipFrames uint16) Location {
		assert.ToBeSame(t, givenEvent, actualEvent)
		assert.ToBeEqual(t, uint16(668), actualSkipFrames)

		return givenLocation
	})
	instance := discoveryFacade(func() Discovery {
		return delegate
	})

	actual := instance.DiscoverLocation(givenEvent, 666)

	assert.ToBeSame(t, givenLocation, actual)
}

func Test_NoopDiscovery(t *testing.T) {
	actual := NoopDiscovery()

	assert.ToBeSame(t, noopV, actual)
}

func Test_noop_DiscoverLocation(t *testing.T) {
	givenEvent := newEvent(level.Warn, nil)

	actual := noopV.DiscoverLocation(givenEvent, 666)

	assert.ToBeNil(t, actual)
}
