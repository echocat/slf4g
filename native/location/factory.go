package location

import log "github.com/echocat/slf4g"

var (
	DefaultFactory       = NewCallerAwareFactory(CallerAwareModePreferFile, CallerAwareDetailSimplified)
	DefaultFactoryFacade = NewFacade(func() Factory {
		return DefaultFactory
	})
)

type Factory func(event log.Event, callDepth int) Location

func NoopFactory(log.Event, int) Location {
	return nil
}

func NewFacade(provider func() Factory) Factory {
	return func(event log.Event, callDepth int) Location {
		if f := provider(); f != nil {
			return f(event, callDepth+1)
		}
		return NoopFactory(event, callDepth+1)
	}
}
