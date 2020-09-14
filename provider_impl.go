package log

import (
	"github.com/echocat/slf4g/fields"
	"sync"
)

func NewProvider(name string, factory ProviderFactory, levels LevelProvider) Provider {
	return &providerImpl{
		name:    name,
		factory: factory,
		levels:  levels,
		global:  factory(GlobalLoggerName),
		loggers: make(map[string]Logger),
	}
}

type providerImpl struct {
	name    string
	factory ProviderFactory
	levels  LevelProvider

	global  Logger
	loggers map[string]Logger
	mutex   sync.RWMutex
}

type ProviderFactory func(name string) Logger

func (instance *providerImpl) GetName() string {
	return instance.name
}

func (instance *providerImpl) GetAllLevels() []Level {
	return instance.levels()
}

func (instance *providerImpl) GetLogger(name string) Logger {
	if name == GlobalLoggerName {
		return instance.global
	}

	instance.mutex.RLock()
	rLocked := true
	defer func() {
		if rLocked {
			instance.mutex.RUnlock()
		}
	}()

	if l, ok := instance.loggers[name]; ok {
		return l
	}

	instance.mutex.RUnlock()
	rLocked = false
	instance.mutex.Lock()
	defer instance.mutex.Unlock()

	if l, ok := instance.loggers[name]; ok {
		return l
	}

	l := instance.factory(name)
	if l == nil {
		l = instance.global
	}
	instance.loggers[name] = l

	return l
}

func (instance *providerImpl) GetFieldKeys() fields.Keys {
	return fields.DefaultKeys
}

func (instance *providerImpl) GetLevelNames() LevelNames {
	return DefaultLevelNames
}
