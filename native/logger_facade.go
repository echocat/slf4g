package native

import (
	"maps"
	"sync"
	"sync/atomic"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

type rootLoggerState struct {
	current atomic.Pointer[rootLoggerSnapshot]
}

type rootLoggerSnapshot struct {
	core    log.CoreLogger
	logger  log.Logger
	version uint64
	failed  bool
}

type rootLoggerCache struct {
	version uint64
	logger  log.Logger
}

type rootLoggerEvent struct {
	state    *rootLoggerState
	logger   log.Logger
	delegate log.Event
}

type rootLoggerFacade struct {
	state     *rootLoggerState
	transform func(log.Logger) log.Logger

	cacheMutex sync.Mutex
	cache      atomic.Pointer[rootLoggerCache]
}

func newRootLoggerFacade(delegate log.CoreLogger) *rootLoggerFacade {
	state := &rootLoggerState{}
	logger := log.NewLogger(delegate)
	state.current.Store(&rootLoggerSnapshot{core: log.UnwrapCoreLogger(logger), logger: logger})
	return &rootLoggerFacade{state: state}
}

func (instance *rootLoggerFacade) setDelegate(delegate log.CoreLogger) {
	if own, ok := delegate.(*rootLoggerFacade); ok && own.state == instance.state {
		delegate = own.current()
	}
	logger := log.NewLogger(delegate)
	previous := instance.state.current.Load()
	instance.state.current.Store(&rootLoggerSnapshot{
		core:    log.UnwrapCoreLogger(logger),
		logger:  logger,
		version: previous.version + 1,
	})
}

func (instance *rootLoggerFacade) setFailed() {
	previous := instance.state.current.Load()
	instance.state.current.Store(&rootLoggerSnapshot{
		version: previous.version + 1,
		failed:  true,
	})
}

func (instance *rootLoggerFacade) snapshot() (log.CoreLogger, log.Logger, uint64) {
	current := instance.state.current.Load()
	if current.failed {
		panic("Root logger customizer did not complete.")
	}
	return current.core, current.logger, current.version
}

func (instance *rootLoggerFacade) current() log.Logger {
	_, current, version := instance.snapshot()
	if instance.transform == nil {
		return current
	}

	if cached := instance.cache.Load(); cached != nil && cached.version == version {
		return cached.logger
	}

	instance.cacheMutex.Lock()
	defer instance.cacheMutex.Unlock()
	if cached := instance.cache.Load(); cached != nil && cached.version == version {
		return cached.logger
	}
	result := &rootLoggerCache{version: version, logger: instance.transform(current)}
	instance.cache.Store(result)
	return result.logger
}

func (instance *rootLoggerFacade) derive(transform func(log.Logger) log.Logger) log.Logger {
	parent := instance.transform
	result := &rootLoggerFacade{
		state: instance.state,
		transform: func(current log.Logger) log.Logger {
			if parent != nil {
				current = parent(current)
			}
			return transform(current)
		},
	}
	result.current()
	return result
}

func (instance *rootLoggerFacade) Unwrap() log.CoreLogger {
	current, _, _ := instance.snapshot()
	if instance.transform != nil {
		return instance.current()
	}
	return current
}

func (instance *rootLoggerFacade) Log(event log.Event, skipFrames uint16) {
	if own, ok := event.(*rootLoggerEvent); ok && own.state == instance.state {
		instance.snapshot()
		own.logger.Log(own.delegate, skipFrames+1)
		return
	}
	instance.current().Log(event, skipFrames+1)
}

func (instance *rootLoggerFacade) IsLevelEnabled(v level.Level) bool {
	return instance.current().IsLevelEnabled(v)
}

func (instance *rootLoggerFacade) GetName() string {
	return instance.current().GetName()
}

func (instance *rootLoggerFacade) NewEvent(v level.Level, values map[string]any) log.Event {
	current := instance.current()
	return &rootLoggerEvent{state: instance.state, logger: current, delegate: current.NewEvent(v, values)}
}

func (instance *rootLoggerFacade) NewEventWithFields(v level.Level, values fields.ForEachEnabled) log.Event {
	current := instance.current()
	return &rootLoggerEvent{state: instance.state, logger: current, delegate: log.NewEventWithFields(current, v, values)}
}

func (instance *rootLoggerFacade) Accepts(event log.Event) bool {
	if own, ok := event.(*rootLoggerEvent); ok && own.state == instance.state {
		instance.snapshot()
		return own.logger.Accepts(own.delegate)
	}
	return instance.current().Accepts(event)
}

func (instance *rootLoggerFacade) GetProvider() log.Provider {
	return instance.current().GetProvider()
}

func (instance *rootLoggerFacade) DoLog(v level.Level, skipFrames uint16, args ...any) {
	delegate := instance.current()
	if delegate, ok := delegate.(log.LoggerFacade); ok {
		delegate.DoLog(v, skipFrames+1, args...)
		return
	}
	log.NewLoggerFacade(func() log.CoreLogger { return delegate }).DoLog(v, skipFrames+1, args...)
}

func (instance *rootLoggerFacade) DoLogf(v level.Level, skipFrames uint16, format string, args ...any) {
	delegate := instance.current()
	if delegate, ok := delegate.(log.LoggerFacade); ok {
		delegate.DoLogf(v, skipFrames+1, format, args...)
		return
	}
	log.NewLoggerFacade(func() log.CoreLogger { return delegate }).DoLogf(v, skipFrames+1, format, args...)
}

func (instance *rootLoggerFacade) Trace(args ...any) {
	instance.current().Trace(args...)
}

func (instance *rootLoggerFacade) Tracef(format string, args ...any) {
	instance.current().Tracef(format, args...)
}

func (instance *rootLoggerFacade) IsTraceEnabled() bool {
	return instance.current().IsTraceEnabled()
}

func (instance *rootLoggerFacade) Debug(args ...any) {
	instance.current().Debug(args...)
}

func (instance *rootLoggerFacade) Debugf(format string, args ...any) {
	instance.current().Debugf(format, args...)
}

func (instance *rootLoggerFacade) IsDebugEnabled() bool {
	return instance.current().IsDebugEnabled()
}

func (instance *rootLoggerFacade) Info(args ...any) {
	instance.current().Info(args...)
}

func (instance *rootLoggerFacade) Infof(format string, args ...any) {
	instance.current().Infof(format, args...)
}

func (instance *rootLoggerFacade) IsInfoEnabled() bool {
	return instance.current().IsInfoEnabled()
}

func (instance *rootLoggerFacade) Warn(args ...any) {
	instance.current().Warn(args...)
}

func (instance *rootLoggerFacade) Warnf(format string, args ...any) {
	instance.current().Warnf(format, args...)
}

func (instance *rootLoggerFacade) IsWarnEnabled() bool {
	return instance.current().IsWarnEnabled()
}

func (instance *rootLoggerFacade) Error(args ...any) {
	instance.current().Error(args...)
}

func (instance *rootLoggerFacade) Errorf(format string, args ...any) {
	instance.current().Errorf(format, args...)
}

func (instance *rootLoggerFacade) IsErrorEnabled() bool {
	return instance.current().IsErrorEnabled()
}

func (instance *rootLoggerFacade) Fatal(args ...any) {
	instance.current().Fatal(args...)
}

func (instance *rootLoggerFacade) Fatalf(format string, args ...any) {
	instance.current().Fatalf(format, args...)
}

func (instance *rootLoggerFacade) IsFatalEnabled() bool {
	return instance.current().IsFatalEnabled()
}

func (instance *rootLoggerFacade) With(name string, value any) log.Logger {
	return instance.derive(func(current log.Logger) log.Logger { return current.With(name, value) })
}

func (instance *rootLoggerFacade) Withf(name string, format string, args ...any) log.Logger {
	argsCopy := append([]any(nil), args...)
	return instance.derive(func(current log.Logger) log.Logger { return current.Withf(name, format, argsCopy...) })
}

func (instance *rootLoggerFacade) WithError(err error) log.Logger {
	return instance.derive(func(current log.Logger) log.Logger { return current.WithError(err) })
}

func (instance *rootLoggerFacade) WithAll(values map[string]any) log.Logger {
	valuesCopy := make(map[string]any, len(values))
	maps.Copy(valuesCopy, values)
	return instance.derive(func(current log.Logger) log.Logger { return current.WithAll(valuesCopy) })
}

func (instance *rootLoggerFacade) Without(keys ...string) log.Logger {
	keysCopy := append([]string(nil), keys...)
	return instance.derive(func(current log.Logger) log.Logger { return current.Without(keysCopy...) })
}

func (instance *rootLoggerFacade) Helper() func() {
	if delegate, ok := instance.current().(interface{ Helper() func() }); ok {
		return delegate.Helper()
	}
	return func() {}
}

func (instance *rootLoggerEvent) GetLevel() level.Level {
	return instance.delegate.GetLevel()
}

func (instance *rootLoggerEvent) ForEach(consumer func(key string, value any) error) error {
	return instance.delegate.ForEach(consumer)
}

func (instance *rootLoggerEvent) Get(key string) (any, bool) {
	return instance.delegate.Get(key)
}

func (instance *rootLoggerEvent) Len() int {
	return instance.delegate.Len()
}

func (instance *rootLoggerEvent) With(key string, value any) log.Event {
	return instance.wrap(instance.delegate.With(key, value))
}

func (instance *rootLoggerEvent) Withf(key string, format string, args ...any) log.Event {
	return instance.wrap(instance.delegate.Withf(key, format, args...))
}

func (instance *rootLoggerEvent) WithError(err error) log.Event {
	return instance.wrap(instance.delegate.WithError(err))
}

func (instance *rootLoggerEvent) WithAll(values map[string]any) log.Event {
	return instance.wrap(instance.delegate.WithAll(values))
}

func (instance *rootLoggerEvent) Without(keys ...string) log.Event {
	return instance.wrap(instance.delegate.Without(keys...))
}

func (instance *rootLoggerEvent) wrap(delegate log.Event) log.Event {
	return &rootLoggerEvent{state: instance.state, logger: instance.logger, delegate: delegate}
}
