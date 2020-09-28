package log

import (
	"github.com/echocat/slf4g/fields"
)

type Provider interface {
	GetLogger(name string) Logger
	GetName() string
	GetAllLevels() Levels
	GetFieldKeysSpec() fields.KeysSpec
}
