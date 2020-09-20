package log

import "github.com/echocat/slf4g/fields"

func newProviderFacade(provider func() Provider) Provider {
	return &providerFacade{
		provider: provider,
	}
}

type providerFacade struct {
	provider func() Provider
}

func (instance *providerFacade) GetName() string {
	return instance.provider().GetName()
}

func (instance *providerFacade) GetLogger(name string) Logger {
	return newLoggerFacade(func() Logger {
		return getProvider().GetLogger(name)
	})
}

func (instance *providerFacade) GetAllLevels() []Level {
	return instance.provider().GetAllLevels()
}

func (instance *providerFacade) GetFieldKeySpec() fields.KeysSpec {
	return instance.provider().GetFieldKeySpec()
}

func (instance *providerFacade) GetLevelNames() LevelNames {
	return instance.provider().GetLevelNames()
}
