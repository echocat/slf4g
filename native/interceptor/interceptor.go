package interceptor

import log "github.com/echocat/slf4g"

var Default Interceptor = Interceptors{
	DefaultFatal,
	DefaultPanics,
}

type Interceptor interface {
	OnBeforeLog(log.Event, log.Provider) (intercepted log.Event)
	OnAfterLog(log.Event, log.Provider) (canContinue bool)
}

type Aware interface {
	GetInterceptor() Interceptor
	SetInterceptor(Interceptor)
}

type Interceptors []Interceptor

func (instance *Interceptors) With(v Interceptor) *Interceptors {
	if *instance == nil {
		*instance = Interceptors{v}
	} else {
		*instance = append(*instance, v)
	}
	return instance
}

func (instance Interceptors) OnBeforeLog(event log.Event, provider log.Provider) (intercepted log.Event) {
	intercepted = event

	if intercepted == nil {
		return
	}

	for _, i := range instance {
		intercepted = i.OnBeforeLog(intercepted, provider)
		if intercepted == nil {
			return
		}
	}

	return
}

func (instance Interceptors) OnAfterLog(event log.Event, provider log.Provider) (canContinue bool) {
	canContinue = true

	if event == nil {
		return
	}

	for _, i := range instance {
		canContinue = i.OnAfterLog(event, provider)
		if !canContinue {
			return
		}
	}

	return
}

func Noop() Interceptor {
	return noopV
}

var noopV = &noop{}

type noop struct{}

func (instance *noop) OnBeforeLog(event log.Event, _ log.Provider) (intercepted log.Event) {
	return event
}

func (instance *noop) OnAfterLog(log.Event, log.Provider) (canContinue bool) {
	return true
}
