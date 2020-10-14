package interceptor

import (
	"math"
	"os"

	"github.com/echocat/slf4g/level"

	log "github.com/echocat/slf4g"
)

func NewFatal(customizer ...func(*Fatal)) *Fatal {
	result := &Fatal{
		ExitCode: 13,
	}
	for _, c := range customizer {
		c(result)
	}
	return result
}

type Fatal struct {
	ExitCode int
}

func (instance *Fatal) OnBeforeLog(event log.Event, _ log.Provider) (intercepted log.Event) {
	return event
}

func (instance *Fatal) OnAfterLog(event log.Event, _ log.Provider) (canContinue bool) {
	if level.Fatal.CompareTo(event.GetLevel()) <= 0 {
		os.Exit(instance.ExitCode)
		return false
	}
	return true
}

func (instance *Fatal) GetPriority() int16 {
	return math.MaxInt16
}
