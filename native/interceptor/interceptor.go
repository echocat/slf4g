package interceptor

import (
	"math"
	"sort"

	log "github.com/echocat/slf4g"
)

var Default Interceptors

type Interceptor interface {
	OnBeforeLog(log.Event, log.Provider) (intercepted log.Event)
	OnAfterLog(log.Event, log.Provider) (canContinue bool)
	GetPriority() int16
}

type Interceptors []Interceptor

func (instance *Interceptors) Add(v Interceptor) *Interceptors {
	if *instance == nil {
		*instance = Interceptors{v}
	} else {
		*instance = append(*instance, v)
	}
	sort.Sort(*instance)
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

func (instance Interceptors) Len() int {
	return len(instance)
}

func (instance Interceptors) Swap(i, j int) {
	instance[i], instance[j] = instance[j], instance[i]
}

func (instance Interceptors) Less(i, j int) bool {
	return instance[i].GetPriority() < instance[j].GetPriority()
}

func (instance Interceptors) GetPriority() int16 {
	lowest := int16(math.MaxInt16)
	for _, i := range instance {
		c := i.GetPriority()
		if c < lowest {
			lowest = c
		}
	}
	return lowest
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

func (instance *noop) GetPriority() int16 {
	return math.MaxInt16
}

type Aware interface {
	GetInterceptor() Interceptor
}

type MutableAware interface {
	Aware
	SetInterceptor(Interceptor)
}
