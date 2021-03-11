// Interceptors are used to intercept instances of log.Event that are requested
// to be logged. It can modify events (before they get logged) or even fully
// prevents to get logged.
package interceptor

import (
	"math"
	"sort"

	log "github.com/echocat/slf4g"
)

// Default is the default instance of Interceptors which should cover the most
// of the cases.
var Default Interceptors = Interceptors{}

// Interceptor is used to intercept instances of log.Event that are requested to
// be logged. It can modify events (before they get logged) or even fully
// prevents to get logged.
type Interceptor interface {
	// OnBeforeLog is called shortly before a log.Event should be logged.
	// Whatever is returned by this method will be logged (if not nil). If this
	// method returns nil logging of this event will be skipped and calling of
	// other interceptors will be skipped, too.
	OnBeforeLog(log.Event, log.Provider) (intercepted log.Event)

	// OnAfterLog is called shortly after a log.Event was logged. If canContinue
	// is false no other Interceptor will be called afterwards.
	OnAfterLog(log.Event, log.Provider) (canContinue bool)

	// GetPriority defines the order how all interceptors will be called. As
	// lower this value is, as earlier this Interceptor will be called.
	GetPriority() int16
}

// Interceptors is a collection of several interceptors.
//
// Add() should be used to modify this instance because it will ensure the
// correct order of all interceptors when this one is called.
type Interceptors []Interceptor

// Add appends a given Interceptor to this instance and ensures that everything
// inside is ordered according to Interceptor.GetPriority().
func (instance *Interceptors) Add(v Interceptor) *Interceptors {
	*instance = append(*instance, v)
	sort.Sort(*instance)
	return instance
}

// OnBeforeLog implements Interceptor.OnBeforeLog()
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

// OnAfterLog implements Interceptor.OnAfterLog()
func (instance Interceptors) OnAfterLog(event log.Event, provider log.Provider) (canContinue bool) {
	canContinue = true

	if event == nil {
		return false
	}

	for _, i := range instance {
		canContinue = i.OnAfterLog(event, provider)
		if !canContinue {
			return
		}
	}

	return
}

// Len provides the length of this instance.
func (instance Interceptors) Len() int {
	return len(instance)
}

// Swap swaps the elements with indexes i and j.
func (instance Interceptors) Swap(i, j int) {
	instance[i], instance[j] = instance[j], instance[i]
}

// Less reports whether the element with index i should sort before the element
// with index j.
func (instance Interceptors) Less(i, j int) bool {
	return instance[i].GetPriority() < instance[j].GetPriority()
}

// GetPriority implements Interceptor.GetPriority() and will return the lowest
// priority of all contained interceptors.
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

// OnBeforeLogFunc is a shortcut to Interceptor.OnBeforeLog().
type OnBeforeLogFunc func(event log.Event, provider log.Provider) (intercepted log.Event)

// OnBeforeLog implements Interceptor.OnBeforeLog().
func (instance OnBeforeLogFunc) OnBeforeLog(event log.Event, provider log.Provider) log.Event {
	return instance(event, provider)
}

// OnAfterLog implements Interceptor.OnAfterLog().
func (instance OnBeforeLogFunc) OnAfterLog(log.Event, log.Provider) bool {
	return true
}

// GetPriority implements Interceptor.GetPriority().
func (instance OnBeforeLogFunc) GetPriority() int16 {
	return 0
}

// OnAfterLogFunc is a shortcut to Interceptor.OnAfterLog().
type OnAfterLogFunc func(event log.Event, provider log.Provider) (canContinue bool)

// OnBeforeLog implements Interceptor.OnBeforeLog().
func (instance OnAfterLogFunc) OnBeforeLog(event log.Event, _ log.Provider) log.Event {
	return event
}

// OnAfterLog implements Interceptor.OnAfterLog().
func (instance OnAfterLogFunc) OnAfterLog(event log.Event, provider log.Provider) bool {
	return instance(event, provider)
}

// GetPriority implements Interceptor.GetPriority().
func (instance OnAfterLogFunc) GetPriority() int16 {
	return 0
}

// NewFacade creates a facade of Interceptor using the given provider.
func NewFacade(provider func() Interceptor) Interceptor {
	return facade(provider)
}

type facade func() Interceptor

func (instance facade) OnBeforeLog(event log.Event, provider log.Provider) (intercepted log.Event) {
	return instance.Unwrap().OnBeforeLog(event, provider)
}

func (instance facade) OnAfterLog(event log.Event, provider log.Provider) (canContinue bool) {
	return instance.Unwrap().OnAfterLog(event, provider)
}

func (instance facade) GetPriority() int16 {
	return instance.Unwrap().GetPriority()
}

func (instance facade) Unwrap() Interceptor {
	return instance()
}

// Noop provides a noop implementation of Interceptor.
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
