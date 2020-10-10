package interceptor

import (
	"os"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/native/internal/support"
)

var DefaultFatal = Fatal{
	ExitCode: support.PInt(13),
}

type Fatal struct {
	ExitCode *int
}

func (instance Fatal) OnBeforeLog(event log.Event, _ log.Provider) (intercepted log.Event) {
	return event
}

func (instance Fatal) OnAfterLog(event log.Event, _ log.Provider) (canContinue bool) {
	if code := instance.ExitCode; code != nil && level.Fatal.CompareTo(event.GetLevel()) <= 0 {
		os.Exit(*code)
		return false
	}
	return true
}
