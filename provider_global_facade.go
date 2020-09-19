package log

import "github.com/echocat/slf4g/fields"

var globalProviderFacadeV = &globalProviderFacade{}

type globalProviderFacade struct{}

func (instance *globalProviderFacade) GetName() string {
	return getProvider().GetName()
}

func (instance *globalProviderFacade) GetLogger(name string) Logger {
	return getProvider().GetLogger(name)
}

func (instance *globalProviderFacade) GetAllLevels() []Level {
	return getProvider().GetAllLevels()
}

func (instance *globalProviderFacade) GetFieldKeySpec() fields.KeysSpec {
	return getProvider().GetFieldKeySpec()
}

func (instance *globalProviderFacade) GetLevelNames() LevelNames {
	return getProvider().GetLevelNames()
}
