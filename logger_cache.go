package log

import (
	"sync"
)

func NewLoggerCache(factory LoggerFactory) LoggerProvider {
	return &loggerCache{
		factory: factory,
		root:    factory(RootLoggerName),
		loggers: make(map[string]Logger),
	}
}

type loggerCache struct {
	factory LoggerFactory

	root    Logger
	loggers map[string]Logger
	mutex   sync.RWMutex
}

func (instance *loggerCache) GetLogger(name string) Logger {
	if name == RootLoggerName {
		return instance.root
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
		l = instance.root
	}
	instance.loggers[name] = l

	return l
}
