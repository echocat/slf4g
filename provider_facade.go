package log

import (
	"github.com/echocat/slf4g/fields"
)

func NewProviderFacade(provider func() Provider) Provider {
	return &providerFacade{
		provider: provider,
	}
}

type providerFacade struct {
	provider func() Provider
}

func (instance *providerFacade) GetName() string {
	return instance.UnwrapProvider().GetName()
}

func (instance *providerFacade) GetLogger(name string) Logger {
	return NewLoggerFacade(func() CoreLogger {
		return instance.UnwrapProvider().GetLogger(name)
	})
}

func (instance *providerFacade) GetAllLevels() []Level {
	return instance.UnwrapProvider().GetAllLevels()
}

func (instance *providerFacade) GetFieldKeySpec() fields.KeysSpec {
	return instance.UnwrapProvider().GetFieldKeySpec()
}

func (instance *providerFacade) UnwrapProvider() Provider {
	return instance.provider()
}
