package interceptor

import (
	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native/formatter"
)

var DefaultPanics = Panics{
	Panics: true,
}

type Panics struct {
	Panics bool
}

func (instance Panics) OnBeforeLog(event log.Event, _ log.Provider) (intercepted log.Event) {
	return event
}

func (instance Panics) OnAfterLog(event log.Event, provider log.Provider) (canContinue bool) {
	if instance.Panics && log.LevelPanic.CompareTo(event.GetLevel()) <= 0 {
		formatted, err := formatter.Default.Format(event, provider, nil)
		if err != nil {
			panic("Panic log event triggered; see logs above.")
		}
		panic(string(formatted))
	}
	return true
}
