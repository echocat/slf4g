package location

import (
	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/testing/recording"
)

func newEvent(l level.Level, values map[string]interface{}) log.Event {
	return recording.NewCoreLogger().NewEvent(l, values)
}
