package formatter

import (
	log "github.com/echocat/slf4g"
)

var (
	DefaultValueFormatter ValueFormatter = &SimpleValueFormatter{}
)

type ValueFormatter interface {
	FormatValue(interface{}, log.Provider) ([]byte, error)
}
