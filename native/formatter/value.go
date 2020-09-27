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

type ValueFormatterFunc func(interface{}, log.Provider) ([]byte, error)

func (instance ValueFormatterFunc) FormatValue(value interface{}, provider log.Provider) ([]byte, error) {
	return instance(value, provider)
}
