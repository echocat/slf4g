package log

import (
	"sync"
)

type LoggerCache interface {
	GetLogger(name string) Logger
}

func NewLoggerCache(factory LoggerFactory) LoggerCache {
	return &loggerCache{
		factory: factory,
		global:  factory(GlobalLoggerName),
		loggers: make(map[string]Logger),
	}
}

type loggerCache struct {
	factory LoggerFactory

	global  Logger
	loggers map[string]Logger
	mutex   sync.RWMutex
}

func (instance *loggerCache) GetLogger(name string) Logger {
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
