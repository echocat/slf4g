package formatter

import (
	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native/formatter/hints"
)

var (
	Default Formatter = DefaultConsole
)

type Formatter interface {
	Format(log.Event, log.Provider, hints.Hints) ([]byte, error)
}
