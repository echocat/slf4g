package native

import (
	"io"
	"maps"
	"os"
	"runtime"
	"sync"
	"sync/atomic"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

const (
	// Absorb transient initialization contention without retaining logs indefinitely.
	maxPendingRootLoggerActions = 1024
	rootLoggerQueueFullFallback = "{\"error\":\"ROOT_LOGGER_QUEUE_FULL\"}\n"
)

type rootLoggerState struct {
	current         atomic.Pointer[rootLoggerSnapshot]
	locationCore    *CoreLogger
	pendingMutex    sync.Mutex
	pending         []func()
	pendingOverflow bool
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
	owner    *rootLoggerFacade
	logger   log.Logger
	version  uint64
	delegate log.Event
}

type rootLoggerEventWithProgramCounter struct {
	log.Event
	programCounter uintptr
}

type rootLoggerFacade struct {
	state     *rootLoggerState
	transform func(log.Logger) log.Logger

	cacheMutex sync.Mutex
	cache      atomic.Pointer[rootLoggerCache]
}

func newRootLoggerFacade(delegate log.CoreLogger, locationCore *CoreLogger) *rootLoggerFacade {
	state := &rootLoggerState{locationCore: locationCore}
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

	instance.state.pendingMutex.Lock()
	pending := instance.state.pending
	pendingOverflow := instance.state.pendingOverflow
	instance.state.pending = nil
	instance.state.pendingOverflow = false
	instance.state.pendingMutex.Unlock()
	if pendingOverflow {
		_, _ = io.WriteString(os.Stderr, rootLoggerQueueFullFallback)
	}
	for _, action := range pending {
		action()
	}
}

func (instance *rootLoggerFacade) failInitialization() {
	previous := instance.state.current.Load()
	if previous.version > 0 {
		return
	}
	instance.state.current.Store(&rootLoggerSnapshot{
		version: previous.version + 1,
		failed:  true,
	})
	instance.state.pendingMutex.Lock()
	instance.state.pending = nil
	instance.state.pendingOverflow = false
	instance.state.pendingMutex.Unlock()
}

func (instance *rootLoggerFacade) hasPendingCapacityLocked() bool {
	if len(instance.state.pending) >= maxPendingRootLoggerActions {
		instance.state.pendingOverflow = true
		return false
	}
	return true
}

func (instance *rootLoggerFacade) appendPendingLocked(action func()) {
	if !instance.hasPendingCapacityLocked() {
		return
	}
	instance.state.pending = append(instance.state.pending, action)
}

func (instance *rootLoggerFacade) isInitializing() bool {
	return instance.state.current.Load().version == 0
}

func (instance *rootLoggerFacade) deferUntilReady(action func(log.Logger)) bool {
	if instance.state.current.Load().version > 0 {
		return false
	}

	// Queue writes instead of blocking, so calls from inside the customizer remain reentrant.
	instance.state.pendingMutex.Lock()
	defer instance.state.pendingMutex.Unlock()
	if instance.state.current.Load().version > 0 {
		return false
	}
	instance.appendPendingLocked(func() { action(instance.current()) })
	return true
}

func (instance *rootLoggerFacade) deferLogUntilReady(
	v level.Level,
	skipFrames uint16,
	args []any,
	action func(log.Logger, []any),
) bool {
	if instance.state.current.Load().version > 0 {
		return false
	}

	instance.state.pendingMutex.Lock()
	defer instance.state.pendingMutex.Unlock()
	if instance.state.current.Load().version > 0 {
		return false
	}
	if !instance.hasPendingCapacityLocked() {
		return true
	}

	current := instance.state.current.Load()
	_, ok := current.core.(*CoreLogger)
	argsCopy := append([]any(nil), args...)
	if !ok {
		instance.state.pending = append(instance.state.pending, func() { action(instance.current(), argsCopy) })
		return true
	}
	pcs := make([]uintptr, 1)
	if runtime.Callers(int(skipFrames)+3, pcs) == 0 {
		instance.state.pending = append(instance.state.pending, func() { action(instance.current(), argsCopy) })
		return true
	}
	instance.state.pending = append(instance.state.pending, func() {
		current := instance.current()
		currentCore := instance.state.current.Load().core
		locationCore, ok := currentCore.(*CoreLogger)
		if !ok {
			locationCore = instance.state.locationCore
		}
		event := rootLoggerEventWithProgramCounter{
			Event:          locationCore.NewEvent(v, nil),
			programCounter: pcs[0],
		}
		location := locationCore.getLocationDiscovery().DiscoverLocation(event, 0)
		if location != nil {
			current = current.With(locationCore.getProvider().getFieldKeysSpec().GetLocation(), location)
		}
		action(current, argsCopy)
	})
	return true
}

func (instance *rootLoggerFacade) snapshot() (log.CoreLogger, log.Logger, uint64) {
	current := instance.state.current.Load()
	if current.failed {
		panic("Root logger customizer did not complete.")
	}
	return current.core, current.logger, current.version
}

func (instance *rootLoggerFacade) current() log.Logger {
	current, _ := instance.currentWithVersion()
	return current
}

func (instance *rootLoggerFacade) currentWithVersion() (log.Logger, uint64) {
	_, current, version := instance.snapshot()
	if instance.transform == nil {
		return current, version
	}

	if cached := instance.cache.Load(); cached != nil && cached.version == version {
		return cached.logger, version
	}

	instance.cacheMutex.Lock()
	defer instance.cacheMutex.Unlock()
	if cached := instance.cache.Load(); cached != nil && cached.version == version {
		return cached.logger, version
	}
	result := &rootLoggerCache{version: version, logger: instance.transform(current)}
	instance.cache.Store(result)
	return result.logger, version
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
	if instance.state.current.Load().version == 0 {
		return nil
	}
	if instance.transform != nil {
		return instance.current()
	}
	return current
}

func (instance *rootLoggerFacade) Log(event log.Event, skipFrames uint16) {
	if event == nil {
		return
	}
	if instance.isInitializing() && instance.deferUntilReady(func(current log.Logger) {
		if own, ok := rootLoggerEventOf(event); ok && own.owner.state == instance.state {
			current, rematerialized := own.resolve(event)
			current.Log(rematerialized, skipFrames+1)
			return
		}
		rematerialized := rematerializeRootLoggerEvent(current, event)
		current.Log(rematerialized, skipFrames+1)
	}) {
		return
	}
	if own, ok := rootLoggerEventOf(event); ok && own.owner.state == instance.state {
		current, resolved := own.resolve(event)
		current.Log(resolved, skipFrames+1)
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
	current, version := instance.currentWithVersion()
	return &rootLoggerEvent{owner: instance, logger: current, version: version, delegate: current.NewEvent(v, values)}
}

func (instance *rootLoggerFacade) NewEventWithFields(v level.Level, values fields.ForEachEnabled) log.Event {
	current, version := instance.currentWithVersion()
	return &rootLoggerEvent{owner: instance, logger: current, version: version, delegate: log.NewEventWithFields(current, v, values)}
}

func (instance *rootLoggerFacade) Accepts(event log.Event) bool {
	if own, ok := rootLoggerEventOf(event); ok && own.owner.state == instance.state {
		current, resolved := own.resolve(event)
		return current.Accepts(resolved)
	}
	return instance.current().Accepts(event)
}

func (instance *rootLoggerFacade) GetProvider() log.Provider {
	return instance.current().GetProvider()
}

func (instance *rootLoggerFacade) DoLog(v level.Level, skipFrames uint16, args ...any) {
	if instance.isInitializing() && instance.deferLogUntilReady(v, skipFrames, args, func(delegate log.Logger, args []any) {
		if delegate, ok := delegate.(log.LoggerFacade); ok {
			delegate.DoLog(v, skipFrames+1, args...)
			return
		}
		log.NewLoggerFacade(func() log.CoreLogger { return delegate }).DoLog(v, skipFrames+1, args...)
	}) {
		return
	}
	delegate := instance.current()
	if delegate, ok := delegate.(log.LoggerFacade); ok {
		delegate.DoLog(v, skipFrames+1, args...)
		return
	}
	log.NewLoggerFacade(func() log.CoreLogger { return delegate }).DoLog(v, skipFrames+1, args...)
}

func (instance *rootLoggerFacade) DoLogf(v level.Level, skipFrames uint16, format string, args ...any) {
	if instance.isInitializing() && instance.deferLogUntilReady(v, skipFrames, args, func(delegate log.Logger, args []any) {
		if delegate, ok := delegate.(log.LoggerFacade); ok {
			delegate.DoLogf(v, skipFrames+1, format, args...)
			return
		}
		log.NewLoggerFacade(func() log.CoreLogger { return delegate }).DoLogf(v, skipFrames+1, format, args...)
	}) {
		return
	}
	delegate := instance.current()
	if delegate, ok := delegate.(log.LoggerFacade); ok {
		delegate.DoLogf(v, skipFrames+1, format, args...)
		return
	}
	log.NewLoggerFacade(func() log.CoreLogger { return delegate }).DoLogf(v, skipFrames+1, format, args...)
}

func (instance *rootLoggerFacade) Trace(args ...any) {
	if instance.isInitializing() && instance.deferLogUntilReady(level.Trace, 0, args, func(current log.Logger, args []any) { current.Trace(args...) }) {
		return
	}
	instance.current().Trace(args...)
}

func (instance *rootLoggerFacade) Tracef(format string, args ...any) {
	if instance.isInitializing() && instance.deferLogUntilReady(level.Trace, 0, args, func(current log.Logger, args []any) { current.Tracef(format, args...) }) {
		return
	}
	instance.current().Tracef(format, args...)
}

func (instance *rootLoggerFacade) IsTraceEnabled() bool {
	return instance.current().IsTraceEnabled()
}

func (instance *rootLoggerFacade) Debug(args ...any) {
	if instance.isInitializing() && instance.deferLogUntilReady(level.Debug, 0, args, func(current log.Logger, args []any) { current.Debug(args...) }) {
		return
	}
	instance.current().Debug(args...)
}

func (instance *rootLoggerFacade) Debugf(format string, args ...any) {
	if instance.isInitializing() && instance.deferLogUntilReady(level.Debug, 0, args, func(current log.Logger, args []any) { current.Debugf(format, args...) }) {
		return
	}
	instance.current().Debugf(format, args...)
}

func (instance *rootLoggerFacade) IsDebugEnabled() bool {
	return instance.current().IsDebugEnabled()
}

func (instance *rootLoggerFacade) Info(args ...any) {
	if instance.isInitializing() && instance.deferLogUntilReady(level.Info, 0, args, func(current log.Logger, args []any) { current.Info(args...) }) {
		return
	}
	instance.current().Info(args...)
}

func (instance *rootLoggerFacade) Infof(format string, args ...any) {
	if instance.isInitializing() && instance.deferLogUntilReady(level.Info, 0, args, func(current log.Logger, args []any) { current.Infof(format, args...) }) {
		return
	}
	instance.current().Infof(format, args...)
}

func (instance *rootLoggerFacade) IsInfoEnabled() bool {
	return instance.current().IsInfoEnabled()
}

func (instance *rootLoggerFacade) Warn(args ...any) {
	if instance.isInitializing() && instance.deferLogUntilReady(level.Warn, 0, args, func(current log.Logger, args []any) { current.Warn(args...) }) {
		return
	}
	instance.current().Warn(args...)
}

func (instance *rootLoggerFacade) Warnf(format string, args ...any) {
	if instance.isInitializing() && instance.deferLogUntilReady(level.Warn, 0, args, func(current log.Logger, args []any) { current.Warnf(format, args...) }) {
		return
	}
	instance.current().Warnf(format, args...)
}

func (instance *rootLoggerFacade) IsWarnEnabled() bool {
	return instance.current().IsWarnEnabled()
}

func (instance *rootLoggerFacade) Error(args ...any) {
	if instance.isInitializing() && instance.deferLogUntilReady(level.Error, 0, args, func(current log.Logger, args []any) { current.Error(args...) }) {
		return
	}
	instance.current().Error(args...)
}

func (instance *rootLoggerFacade) Errorf(format string, args ...any) {
	if instance.isInitializing() && instance.deferLogUntilReady(level.Error, 0, args, func(current log.Logger, args []any) { current.Errorf(format, args...) }) {
		return
	}
	instance.current().Errorf(format, args...)
}

func (instance *rootLoggerFacade) IsErrorEnabled() bool {
	return instance.current().IsErrorEnabled()
}

func (instance *rootLoggerFacade) Fatal(args ...any) {
	if instance.isInitializing() && instance.deferLogUntilReady(level.Fatal, 0, args, func(current log.Logger, args []any) { current.Fatal(args...) }) {
		return
	}
	instance.current().Fatal(args...)
}

func (instance *rootLoggerFacade) Fatalf(format string, args ...any) {
	if instance.isInitializing() && instance.deferLogUntilReady(level.Fatal, 0, args, func(current log.Logger, args []any) { current.Fatalf(format, args...) }) {
		return
	}
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
	return &rootLoggerEvent{
		owner:    instance.owner,
		logger:   instance.logger,
		version:  instance.version,
		delegate: delegate,
	}
}

func (instance *rootLoggerEvent) resolve(source log.Event) (log.Logger, log.Event) {
	current, version := instance.owner.currentWithVersion()
	if version == instance.version {
		if source == instance {
			return instance.logger, instance.delegate
		}
		return instance.logger, source
	}
	return current, rematerializeRootLoggerEvent(current, source)
}

func rootLoggerEventOf(event log.Event) (*rootLoggerEvent, bool) {
	for depth := 0; event != nil && depth < 16; depth++ {
		if own, ok := event.(*rootLoggerEvent); ok {
			return own, true
		}
		unwrapper, ok := event.(interface{ UnwrapEvent() log.Event })
		if !ok {
			return nil, false
		}
		event = unwrapper.UnwrapEvent()
	}
	return nil, false
}

func rematerializeRootLoggerEvent(target log.Logger, source log.Event) log.Event {
	values, err := fields.AsMap(source)
	if err != nil {
		panic(err)
	}
	result := target.NewEvent(source.GetLevel(), values)
	if programCounter, ok := rootLoggerEventProgramCounter(source); ok {
		candidate := rootLoggerEventWithProgramCounter{Event: result, programCounter: programCounter}
		if target.Accepts(candidate) {
			result = candidate
		}
	}
	return result
}

func rootLoggerEventProgramCounter(event log.Event) (uintptr, bool) {
	if source, ok := event.(interface{ GetProgramCounter() uintptr }); ok {
		return source.GetProgramCounter(), true
	}
	if own, ok := event.(*rootLoggerEvent); ok {
		return rootLoggerEventProgramCounter(own.delegate)
	}
	return 0, false
}

func (instance rootLoggerEventWithProgramCounter) GetProgramCounter() uintptr {
	return instance.programCounter
}

func (instance rootLoggerEventWithProgramCounter) UnwrapEvent() log.Event {
	return instance.Event
}

func (instance rootLoggerEventWithProgramCounter) With(key string, value any) log.Event {
	instance.Event = instance.Event.With(key, value)
	return instance
}

func (instance rootLoggerEventWithProgramCounter) Withf(key string, format string, args ...any) log.Event {
	instance.Event = instance.Event.Withf(key, format, args...)
	return instance
}

func (instance rootLoggerEventWithProgramCounter) WithError(err error) log.Event {
	instance.Event = instance.Event.WithError(err)
	return instance
}

func (instance rootLoggerEventWithProgramCounter) WithAll(values map[string]any) log.Event {
	instance.Event = instance.Event.WithAll(values)
	return instance
}

func (instance rootLoggerEventWithProgramCounter) Without(keys ...string) log.Event {
	instance.Event = instance.Event.Without(keys...)
	return instance
}
