package log

import (
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

type providerFacade func() Provider

func (instance providerFacade) GetName() string {
	return instance.Unwrap().GetName()
}

func (instance providerFacade) GetRootLogger() Logger {
	return NewLoggerFacade(func() CoreLogger {
		return instance.Unwrap().GetRootLogger()
	})
}

func (instance providerFacade) GetLogger(name string) Logger {
	return NewLoggerFacade(func() CoreLogger {
		return instance.Unwrap().GetLogger(name)
	})
}

func (instance providerFacade) GetAllLevels() level.Levels {
	return instance.Unwrap().GetAllLevels()
}

func (instance providerFacade) GetFieldKeysSpec() fields.KeysSpec {
	return instance.Unwrap().GetFieldKeysSpec()
}

func (instance providerFacade) Unwrap() Provider {
	return instance()
}
