package log

import "github.com/echocat/slf4g/fields"

type Provider interface {
	GetName() string

	GetLogger(name string) Logger

	GetAllLevels() []Level

	GetFieldKeySpec() fields.KeysSpec

	GetLevelNames() LevelNames
}
