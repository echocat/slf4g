package log

import (
	"sync"
)

type LoggerCache interface {
	GetLogger(name string) Logger
	GetRootLogger() Logger
}

func NewLoggerCache(rootFactory func() Logger, factory func(name string) Logger) LoggerCache {
	return &loggerCache{
		factory: factory,
		root:    rootFactory(),
		loggers: make(map[string]Logger),
	}
}

type loggerCache struct {
	factory func(name string) Logger

	root    Logger
	loggers map[string]Logger
	mutex   sync.RWMutex
}

func (instance *loggerCache) GetRootLogger() Logger {
	return instance.root
}

func (instance *loggerCache) GetLogger(name string) Logger {
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
		l = instance.root
	}
	instance.loggers[name] = l

	return l
}
