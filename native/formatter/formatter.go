package formatter

import (
	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native/formatter/hints"
)

var (
	Default Formatter = NewConsole()
)

type Formatter interface {
	Format(log.Event, log.Provider, hints.Hints) ([]byte, error)
}

type Func func(log.Event, log.Provider, hints.Hints) ([]byte, error)

func (instance Func) Format(event log.Event, provider log.Provider, h hints.Hints) ([]byte, error) {
	return instance(event, provider, h)
}

func NewFacade(provider func() Formatter) Formatter {
	return facade(provider)
}

type facade func() Formatter

func (instance facade) Format(event log.Event, provider log.Provider, h hints.Hints) ([]byte, error) {
	return instance().Format(event, provider, h)
}

type Aware interface {
	GetFormatter() Formatter
}

type MutableAware interface {
	Aware
	SetFormatter(Formatter)
}
