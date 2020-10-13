package log

import (
	"sync"
)

// LoggerCache could provide more than one time the same instance of a named
// Logger (by calling GetLogger(name string)) or the same root Logger (by
// calling GetRootLogger()).
type LoggerCache interface {
	// GetLogger returns a Logger for the given name.
	GetLogger(name string) Logger

	// GetRootLogger returns the root Logger.
	GetRootLogger() Logger

	// GetNames returns all names for all already known Logger which are
	// already received using GetLogger(name).
	GetNames() []string
}

// NewLoggerCache creates a new instance of LoggerCache by the given
// rootFactory and factory.
func NewLoggerCache(rootFactory func() Logger, factory func(name string) Logger) LoggerCache {
	root := rootFactory()
	if root == nil {
		panic("Root factory returned a nil root logger.")
	}
	return &loggerCache{
		factory: factory,
		root:    root,
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

func (instance *loggerCache) GetNames() (result []string) {
	instance.mutex.RLock()
	defer instance.mutex.RUnlock()

	result = make([]string, len(instance.loggers))

	i := 0
	for name := range instance.loggers {
		result[i] = name
		i++
	}

	return
}
