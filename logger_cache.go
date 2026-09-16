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
// rootFactory and factory. The factory may request loggers with other names
// from the created cache, but must not recursively request the name it is
// currently creating.
func NewLoggerCache(rootFactory func() Logger, factory func(name string) Logger) LoggerCache {
	root := rootFactory()
	if root == nil {
		panic("Root factory returned a nil root logger.")
	}
	return &loggerCache{
		factory:   factory,
		root:      root,
		loggers:   make(map[string]Logger),
		creations: make(map[string]*loggerCreation),
	}
}

type loggerCache struct {
	factory func(name string) Logger

	root      Logger
	loggers   map[string]Logger
	creations map[string]*loggerCreation
	mutex     sync.RWMutex
}

type loggerCreation struct {
	done   chan struct{}
	logger Logger
	failed bool
}

func (instance *loggerCache) GetRootLogger() Logger {
	return instance.root
}

func (instance *loggerCache) GetLogger(name string) Logger {
	instance.mutex.RLock()
	if l, ok := instance.loggers[name]; ok {
		instance.mutex.RUnlock()
		return l
	}
	instance.mutex.RUnlock()

	instance.mutex.Lock()
	if l, ok := instance.loggers[name]; ok {
		instance.mutex.Unlock()
		return l
	}
	if instance.creations == nil {
		instance.creations = make(map[string]*loggerCreation)
	}
	if creation := instance.creations[name]; creation != nil {
		instance.mutex.Unlock()
		<-creation.done
		if creation.failed {
			panic("Logger factory did not complete.")
		}
		return creation.logger
	}
	creation := &loggerCreation{done: make(chan struct{})}
	instance.creations[name] = creation
	instance.mutex.Unlock()

	completed := false
	defer func() {
		if !completed {
			instance.mutex.Lock()
			delete(instance.creations, name)
			creation.failed = true
			close(creation.done)
			instance.mutex.Unlock()
		}
	}()
	l := instance.factory(name)
	if l == nil {
		l = instance.root
	}

	instance.mutex.Lock()
	instance.loggers[name] = l
	delete(instance.creations, name)
	creation.logger = l
	close(creation.done)
	completed = true
	instance.mutex.Unlock()

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
