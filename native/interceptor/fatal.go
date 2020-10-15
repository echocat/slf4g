package interceptor

import (
	"math"
	"os"

	"github.com/echocat/slf4g/level"

	log "github.com/echocat/slf4g"
)

// Fatal will exit the application after logged events on level.Fatal.
type Fatal struct {
	// ExitCode to exit the application with.
	ExitCode int
}

// NewFatal creates a new instance of Fatal.
func NewFatal(customizer ...func(*Fatal)) *Fatal {
	result := &Fatal{
		ExitCode: 13,
	}
	for _, c := range customizer {
		c(result)
	}
	return result
}

// OnBeforeLog implements Interceptor.OnBeforeLog()
func (instance *Fatal) OnBeforeLog(event log.Event, _ log.Provider) (intercepted log.Event) {
	return event
}

// OnAfterLog implements Interceptor.OnAfterLog()
func (instance *Fatal) OnAfterLog(event log.Event, _ log.Provider) (canContinue bool) {
	if level.Fatal.CompareTo(event.GetLevel()) <= 0 {
		os.Exit(instance.ExitCode)
		return false
	}
	return true
}

// GetPriority implements Interceptor.GetPriority()
func (instance *Fatal) GetPriority() int16 {
	return math.MaxInt16
}
