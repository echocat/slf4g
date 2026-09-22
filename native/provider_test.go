package native

import (
	"context"
	"io"
	stdslog "log/slog"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	log "github.com/echocat/slf4g"

	"github.com/echocat/slf4g/native/location"

	nlevel "github.com/echocat/slf4g/native/level"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/native/consumer"
	slogbridge "github.com/echocat/slf4g/sdk/bridge/slog"
)

type strictCoreLogger struct {
	log.CoreLogger
	token    *int
	logCalls *atomic.Int32
}

type strictOwnershipCoreLogger struct {
	log.CoreLogger
}

type strictOwnershipEvent struct {
	log.Event
	owner *strictOwnershipCoreLogger
}

type fullLoggerWithoutUnwrap struct {
	log.Logger
	currentLevel level.Level
}

func (instance *fullLoggerWithoutUnwrap) GetLevel() level.Level {
	return instance.currentLevel
}

func (instance *fullLoggerWithoutUnwrap) SetLevel(value level.Level) {
	instance.currentLevel = value
}

func (instance *strictCoreLogger) NewEvent(v level.Level, values map[string]any) log.Event {
	return instance.CoreLogger.NewEvent(v, values).With("strict-owner", instance.token)
}

func (instance *strictCoreLogger) Accepts(event log.Event) bool {
	actual, exists := event.Get("strict-owner")
	return exists && actual == instance.token
}

func (instance *strictCoreLogger) Log(event log.Event, skipFrames uint16) {
	if instance.logCalls != nil {
		instance.logCalls.Add(1)
	}
	if !instance.Accepts(event) {
		panic("foreign event")
	}
	instance.CoreLogger.Log(event.Without("strict-owner"), skipFrames+1)
}

func (instance *strictOwnershipCoreLogger) NewEvent(v level.Level, values map[string]any) log.Event {
	return &strictOwnershipEvent{Event: instance.CoreLogger.NewEvent(v, values), owner: instance}
}

func (instance *strictOwnershipCoreLogger) NewEventWithFields(v level.Level, values fields.ForEachEnabled) log.Event {
	return &strictOwnershipEvent{Event: log.NewEventWithFields(instance.CoreLogger, v, values), owner: instance}
}

func (instance *strictOwnershipCoreLogger) Accepts(event log.Event) bool {
	actual, ok := event.(*strictOwnershipEvent)
	return ok && actual.owner == instance
}

func (instance *strictOwnershipCoreLogger) Log(event log.Event, skipFrames uint16) {
	actual, ok := event.(*strictOwnershipEvent)
	if !ok || actual.owner != instance {
		panic("foreign event")
	}
	instance.CoreLogger.Log(actual.Event, skipFrames+1)
}

func Test_Provider_GetName_specified(t *testing.T) {
	instance, _ := newProvider()
	instance.Name = "foo"

	assert.ToBeEqual(t, "foo", instance.GetName())
}

func Test_Provider_GetName_absent(t *testing.T) {
	instance, _ := newProvider()
	instance.Name = ""

	assert.ToBeEqual(t, "native", instance.GetName())
}

func Test_Provider_GetAllLevels_specified(t *testing.T) {
	instance, _ := newProvider()
	instance.LevelProvider = level.Levels{level.Warn, level.Fatal}.ToProvider("mock")

	assert.ToBeEqual(t, level.Levels{level.Warn, level.Fatal}, instance.GetAllLevels())
}

func Test_Provider_GetAllLevels_absent(t *testing.T) {
	instance, _ := newProvider()

	assert.ToBeEqual(t, level.GetProvider().GetLevels(), instance.GetAllLevels())
}

func Test_Provider_GetFieldKeysSpec_specified(t *testing.T) {
	givenSpec := &FieldKeysSpecImpl{}
	instance, _ := newProvider()
	instance.FieldKeysSpec = givenSpec

	assert.ToBeSame(t, givenSpec, instance.GetFieldKeysSpec())
}

func Test_Provider_GetFieldKeysSpec_globalDefault(t *testing.T) {
	before := DefaultFieldKeysSpec
	defer func() { DefaultFieldKeysSpec = before }()

	givenSpec := &FieldKeysSpecImpl{}
	DefaultFieldKeysSpec = givenSpec
	instance, _ := newProvider()

	assert.ToBeSame(t, givenSpec, instance.GetFieldKeysSpec())
}

func Test_Provider_GetFieldKeysSpec_fallback(t *testing.T) {
	before := DefaultFieldKeysSpec
	defer func() { DefaultFieldKeysSpec = before }()

	DefaultFieldKeysSpec = nil
	instance, _ := newProvider()

	assert.ToBeEqual(t, &FieldKeysSpecImpl{}, instance.GetFieldKeysSpec())
}

func Test_Provider_GetLevel_specified(t *testing.T) {
	instance, _ := newProvider()
	instance.Level = level.Warn

	assert.ToBeEqual(t, level.Warn, instance.GetLevel())
}

func Test_Provider_GetLevel_absent(t *testing.T) {
	instance, _ := newProvider()

	assert.ToBeEqual(t, level.Info, instance.GetLevel())
}

func Test_Provider_GetLevel(t *testing.T) {
	template, _ := newProvider()
	for _, l := range template.GetAllLevels() {
		instance, _ := newProvider()
		instance.Level = l
		assert.ToBeEqual(t, l, instance.GetLevel())
	}
}

func Test_Provider_SetLevel(t *testing.T) {
	instance, _ := newProvider()

	assert.ToBeEqual(t, level.Level(0), instance.Level)

	for _, l := range instance.GetAllLevels() {
		instance.SetLevel(l)
		assert.ToBeEqual(t, l, instance.Level)
		assert.ToBeEqual(t, l, instance.GetLevel())
	}
	instance.SetLevel(level.Level(^uint16(0)))
	assert.ToBeEqual(t, level.Level(^uint16(0)), instance.GetLevel())

	instance.SetLevel(0)
	assert.ToBeEqual(t, level.Level(0), instance.Level)
	assert.ToBeEqual(t, level.Info, instance.GetLevel())
}

func Test_Provider_SetLevel_concurrentlyWithGetLevel(t *testing.T) {
	instance, _ := newProvider()
	start := make(chan struct{})
	var wait sync.WaitGroup
	wait.Add(2)

	go func() {
		defer wait.Done()
		<-start
		for range 1000 {
			instance.SetLevel(level.Debug)
			instance.SetLevel(level.Info)
		}
	}()
	go func() {
		defer wait.Done()
		<-start
		for range 1000 {
			_ = instance.GetLevel()
		}
	}()

	close(start)
	wait.Wait()
	instance.SetLevel(level.Warn)
	assert.ToBeEqual(t, level.Warn, instance.GetLevel())
}

func Test_Provider_GetLevel_fromCopy(t *testing.T) {
	instance, _ := newProvider()
	instance.Level = level.Debug
	assert.ToBeEqual(t, level.Debug, instance.GetLevel())

	copied := *instance
	copied.Level = level.Warn

	assert.ToBeEqual(t, level.Warn, copied.GetLevel())
	assert.ToBeEqual(t, level.Debug, instance.GetLevel())
}

func Test_Provider_GetLevelNames_specified(t *testing.T) {
	givenNames := nlevel.NewNames()
	instance, _ := newProvider()
	instance.LevelNames = givenNames

	actual := instance.GetLevelNames()

	assert.ToBeSame(t, givenNames, actual)
}

func Test_Provider_GetLevelNames_globalDefault(t *testing.T) {
	beforeNames := nlevel.DefaultNames
	defer func() { nlevel.DefaultNames = beforeNames }()
	givenNames := nlevel.NewNames()
	nlevel.DefaultNames = givenNames

	instance, _ := newProvider()
	instance.LevelNames = nil

	actual := instance.GetLevelNames()

	assert.ToBeSame(t, givenNames, actual)
}

func Test_Provider_GetLevelNames_fallback(t *testing.T) {
	beforeNames := nlevel.DefaultNames
	defer func() { nlevel.DefaultNames = beforeNames }()
	nlevel.DefaultNames = nil

	instance, _ := newProvider()
	instance.LevelNames = nil

	actual := instance.GetLevelNames()

	assert.ToBeEqual(t, nlevel.NewNames(), actual)
}

func Test_Provider_getLocationDiscovery_specified(t *testing.T) {
	givenDiscovery := location.NewCallerDiscovery()
	instance, _ := newProvider()
	instance.LocationDiscovery = givenDiscovery

	actual := instance.getLocationDiscovery()

	assert.ToBeSame(t, givenDiscovery, actual)
}

func Test_Provider_getLocationDiscovery_globalDefault(t *testing.T) {
	before := location.DefaultDiscovery
	defer func() { location.DefaultDiscovery = before }()
	givenDiscovery := location.NewCallerDiscovery()
	location.DefaultDiscovery = givenDiscovery

	instance, _ := newProvider()
	instance.LocationDiscovery = nil

	actual := instance.getLocationDiscovery()

	assert.ToBeSame(t, givenDiscovery, actual)
}

func Test_Provider_getLocationDiscovery_fallback(t *testing.T) {
	before := location.DefaultDiscovery
	defer func() { location.DefaultDiscovery = before }()
	givenDiscovery := location.NoopDiscovery()
	location.DefaultDiscovery = nil

	instance, _ := newProvider()
	instance.LocationDiscovery = nil

	actual := instance.getLocationDiscovery()

	assert.ToBeEqual(t, givenDiscovery, actual)
}

func Test_Provider_SetConsumer_specified(t *testing.T) {
	givenConsumer := consumer.NewRecorder()
	instance, _ := newProvider()

	instance.SetConsumer(givenConsumer)

	assert.ToBeSame(t, givenConsumer, instance.Consumer)
	assert.ToBeSame(t, givenConsumer, instance.GetConsumer())
}

func Test_Provider_SetConsumer_concurrentlyWithLogging(t *testing.T) {
	first := consumer.NewRecorder()
	var secondCount atomic.Int32
	second := consumer.Func(func(log.Event, log.CoreLogger) {
		secondCount.Add(1)
	})
	instance, _ := newProvider()
	instance.Consumer = first
	logger := instance.GetLogger("test")
	start := make(chan struct{})
	var wait sync.WaitGroup
	wait.Add(2)

	go func() {
		defer wait.Done()
		<-start
		for range 1000 {
			instance.SetConsumer(first)
			instance.SetConsumer(second)
			runtime.Gosched()
		}
	}()
	go func() {
		defer wait.Done()
		<-start
		for range 1000 {
			logger.Info("message")
			runtime.Gosched()
		}
	}()

	close(start)
	wait.Wait()
	assert.ToBeEqual(t, 1000, first.Len()+int(secondCount.Load()))
}

func Test_Provider_GetConsumer_fromCopy(t *testing.T) {
	first := consumer.NewRecorder()
	second := consumer.NewRecorder()
	instance, _ := newProvider()
	instance.Consumer = first
	assert.ToBeSame(t, first, instance.GetConsumer())

	copied := *instance
	copied.Consumer = second

	assert.ToBeSame(t, second, copied.GetConsumer())
	assert.ToBeSame(t, first, instance.GetConsumer())
}

func Test_Provider_GetConsumer_specified(t *testing.T) {
	givenConsumer := consumer.NewRecorder()
	instance, _ := newProvider()
	instance.Consumer = givenConsumer

	actual := instance.GetConsumer()

	assert.ToBeSame(t, givenConsumer, actual)
}

func Test_Provider_GetConsumer_globalDefault(t *testing.T) {
	before := consumer.Default
	defer func() { consumer.Default = before }()
	givenConsumer := consumer.NewRecorder()
	consumer.Default = givenConsumer

	instance, _ := newProvider()
	instance.Consumer = nil

	actual := instance.GetConsumer()

	assert.ToBeSame(t, givenConsumer, actual)
}

func Test_Provider_GetConsumer_fallback(t *testing.T) {
	before := consumer.Default
	defer func() { consumer.Default = before }()
	givenConsumer := consumer.Noop()
	consumer.Default = nil

	instance, _ := newProvider()
	instance.Consumer = nil

	actual := instance.GetConsumer()

	assert.ToBeEqual(t, givenConsumer, actual)
}

func Test_Provider_GetRootLogger(t *testing.T) {
	instance, _ := newProvider()

	actual1 := instance.GetRootLogger()
	assert.ToBeEqual(t, &CoreLogger{
		provider: instance,
		name:     rootLoggerName,
	}, log.UnwrapCoreLogger(actual1))

	actual2 := instance.GetRootLogger()
	assert.ToBeSame(t, actual1, actual2)

}

func Test_Provider_GetLogger(t *testing.T) {
	instance, _ := newProvider()

	actualA1 := instance.GetLogger("a")
	assert.ToBeEqual(t, &CoreLogger{
		provider: instance,
		name:     "a",
	}, log.UnwrapCoreLogger(actualA1))

	actualB1 := instance.GetLogger("b")
	assert.ToBeEqual(t, &CoreLogger{
		provider: instance,
		name:     "b",
	}, log.UnwrapCoreLogger(actualB1))

	actualA2 := instance.GetLogger("a")
	assert.ToBeSame(t, actualA1, actualA2)

	actualB2 := instance.GetLogger("b")
	assert.ToBeSame(t, actualB1, actualB2)
}

func Test_Provider_factory_usingCustomizer(t *testing.T) {
	givenCoreLogger := &CoreLogger{name: "bar"}

	instance, _ := newProvider()

	instance.CoreLoggerCustomizer = func(actualProvider *Provider, actualLogger *CoreLogger) log.CoreLogger {
		assert.ToBeSame(t, instance, actualProvider)
		assert.ToBeEqual(t, &CoreLogger{
			provider: instance,
			name:     "foo",
		}, actualLogger)

		return givenCoreLogger
	}

	actual := instance.factory("foo")

	assert.ToBeSame(t, givenCoreLogger, log.UnwrapCoreLogger(actual))
}

func Test_Provider_GetRootLogger_supportsReentrantCustomizer(t *testing.T) {
	instance, recorder := newProvider()
	var customizerCalls atomic.Int32
	instance.CoreLoggerCustomizer = func(actualProvider *Provider, actualLogger *CoreLogger) log.CoreLogger {
		customizerCalls.Add(1)
		actualProvider.GetRootLogger().Info("during customization")
		return actualLogger
	}

	done := make(chan log.Logger)
	go func() {
		done <- instance.GetRootLogger()
	}()

	select {
	case actual := <-done:
		assert.ToBeSame(t, instance.GetRootLogger(), actual)
	case <-time.After(time.Second):
		t.Fatal("reentrant root logger customizer did not complete")
	}
	assert.ToBeEqual(t, int32(1), customizerCalls.Load())
	assert.ToBeEqual(t, 1, recorder.Len())
}

func Test_Provider_GetRootLogger_supportsFallbackDetectionInReentrantCustomizer(t *testing.T) {
	instance, _ := newProvider()
	fallback := make(chan bool, 1)
	instance.CoreLoggerCustomizer = func(actualProvider *Provider, logger *CoreLogger) log.CoreLogger {
		fallback <- log.IsFallbackLogger(actualProvider.GetRootLogger())
		return logger
	}

	initialized := make(chan struct{})
	go func() {
		defer close(initialized)
		instance.GetRootLogger()
	}()

	select {
	case <-initialized:
		assert.ToBeEqual(t, false, <-fallback)
	case <-time.After(time.Second):
		t.Fatal("fallback detection blocked reentrant root logger customization")
	}
}

func Test_Provider_GetRootLogger_customizesOnceDuringConcurrentInitialization(t *testing.T) {
	instance, _ := newProvider()
	customizerEntered := make(chan struct{})
	releaseCustomizer := make(chan struct{})
	var customizerCalls atomic.Int32
	instance.CoreLoggerCustomizer = func(_ *Provider, logger *CoreLogger) log.CoreLogger {
		if customizerCalls.Add(1) == 1 {
			close(customizerEntered)
		}
		<-releaseCustomizer
		return logger
	}

	firstResult := make(chan log.Logger, 1)
	go func() {
		firstResult <- instance.GetRootLogger()
	}()
	<-customizerEntered

	results := make(chan log.Logger, 31)
	var wait sync.WaitGroup
	for range cap(results) {
		wait.Add(1)
		go func() {
			defer wait.Done()
			results <- instance.GetRootLogger()
		}()
	}
	wait.Wait()
	close(results)
	close(releaseCustomizer)
	expected := <-firstResult

	for actual := range results {
		assert.ToBeSame(t, expected, actual)
	}
	assert.ToBeEqual(t, int32(1), customizerCalls.Load())
}

func Test_Provider_GetRootLogger_doesNotExposeProvisionalCoreDuringInitialization(t *testing.T) {
	for range 1000 {
		instance, _ := newProvider()
		customizerEntered := make(chan struct{})
		releaseCustomizer := make(chan struct{})
		instance.CoreLoggerCustomizer = func(_ *Provider, logger *CoreLogger) log.CoreLogger {
			close(customizerEntered)
			<-releaseCustomizer
			return &strictCoreLogger{CoreLogger: logger, token: new(int)}
		}

		initialized := make(chan struct{})
		go func() {
			defer close(initialized)
			instance.GetRootLogger()
		}()
		<-customizerEntered
		root := instance.GetRootLogger()

		stop := make(chan struct{})
		provisionalExposed := make(chan struct{}, 1)
		var readers sync.WaitGroup
		for range 8 {
			readers.Add(1)
			go func() {
				defer readers.Done()
				for {
					select {
					case <-stop:
						return
					default:
					}
					if _, ok := log.UnwrapCoreLogger(root).(*CoreLogger); ok {
						select {
						case provisionalExposed <- struct{}{}:
						default:
						}
						return
					}
				}
			}()
		}

		runtime.Gosched()
		close(releaseCustomizer)
		<-initialized
		close(stop)
		readers.Wait()
		select {
		case <-provisionalExposed:
			t.Fatal("UnwrapCoreLogger exposed the provisional root core")
		default:
		}
	}
}

func Test_Provider_GetRootLogger_boundsPendingLogs(t *testing.T) {
	stderr, err := os.CreateTemp(t.TempDir(), "stderr")
	assert.ToBeNoError(t, err)
	defer func() { _ = stderr.Close() }()
	oldStderr := os.Stderr
	os.Stderr = stderr
	defer func() { os.Stderr = oldStderr }()

	instance, recorder := newProvider()
	customizerEntered := make(chan struct{})
	releaseCustomizer := make(chan struct{})
	actualMaxPending := 0
	instance.CoreLoggerCustomizer = func(_ *Provider, logger *CoreLogger) log.CoreLogger {
		close(customizerEntered)
		<-releaseCustomizer
		return logger
	}

	initialized := make(chan struct{})
	go func() {
		defer close(initialized)
		instance.GetRootLogger()
	}()
	<-customizerEntered
	root := instance.GetRootLogger()
	for range maxPendingRootLoggerActions + 100 {
		root.Info("pending")
	}
	facade := root.(*rootLoggerFacade)
	facade.state.pendingMutex.Lock()
	actualMaxPending = len(facade.state.pending)
	facade.state.pendingMutex.Unlock()
	close(releaseCustomizer)
	<-initialized

	os.Stderr = oldStderr
	_, err = stderr.Seek(0, io.SeekStart)
	assert.ToBeNoError(t, err)
	actualStderr, err := io.ReadAll(stderr)
	assert.ToBeNoError(t, err)
	assert.ToBeEqual(t, maxPendingRootLoggerActions, actualMaxPending)
	assert.ToBeEqual(t, maxPendingRootLoggerActions, recorder.Len())
	assert.ToBeEqual(t, "{\"error\":\"ROOT_LOGGER_QUEUE_FULL\"}\n", string(actualStderr))
}

func Test_Provider_GetRootLogger_preservesCustomizedLogger(t *testing.T) {
	instance, recorder := newProvider()
	instance.CoreLoggerCustomizer = func(_ *Provider, logger *CoreLogger) log.CoreLogger {
		return log.NewLogger(logger).With("custom", "value")
	}

	instance.GetRootLogger().Info("message")

	assert.ToBeEqual(t, 1, recorder.Len())
	actual, exists := recorder.Get(0).Get("custom")
	assert.ToBeEqual(t, true, exists)
	assert.ToBeEqual(t, "value", actual)
}

func Test_Provider_GetRootLogger_doesNotLogBeforeCustomizationCompletes(t *testing.T) {
	instance, recorder := newProvider()
	customizerEntered := make(chan struct{})
	releaseCustomizer := make(chan struct{})
	instance.CoreLoggerCustomizer = func(_ *Provider, logger *CoreLogger) log.CoreLogger {
		close(customizerEntered)
		<-releaseCustomizer
		return log.NewLogger(logger).With("custom", "value")
	}

	initialized := make(chan struct{})
	go func() {
		defer close(initialized)
		instance.GetRootLogger()
	}()
	<-customizerEntered

	logged := make(chan struct{})
	go func() {
		defer close(logged)
		instance.GetRootLogger().Info("message")
	}()
	<-logged
	assert.ToBeEqual(t, 0, recorder.Len())

	close(releaseCustomizer)
	<-initialized

	assert.ToBeEqual(t, 1, recorder.Len())
	actual, exists := recorder.Get(0).Get("custom")
	assert.ToBeEqual(t, true, exists)
	assert.ToBeEqual(t, "value", actual)
}

func Test_Provider_GetRootLogger_updatesDerivedLoggerAfterCustomization(t *testing.T) {
	instance, recorder := newProvider()
	customizerEntered := make(chan struct{})
	releaseCustomizer := make(chan struct{})
	instance.CoreLoggerCustomizer = func(_ *Provider, logger *CoreLogger) log.CoreLogger {
		close(customizerEntered)
		<-releaseCustomizer
		return log.NewLogger(logger).With("custom", "value")
	}

	initialized := make(chan struct{})
	go func() {
		defer close(initialized)
		instance.GetRootLogger()
	}()
	<-customizerEntered
	derived := instance.GetRootLogger().With("derived", "value")
	close(releaseCustomizer)
	<-initialized

	derived.Info("message")

	assert.ToBeEqual(t, 1, recorder.Len())
	for _, key := range []string{"custom", "derived"} {
		actual, exists := recorder.Get(0).Get(key)
		assert.ToBeEqual(t, true, exists)
		assert.ToBeEqual(t, "value", actual)
	}
}

func Test_Provider_GetRootLogger_rematerializesProvisionalEventAfterCustomization(t *testing.T) {
	instance, recorder := newProvider()
	customizerEntered := make(chan struct{})
	releaseCustomizer := make(chan struct{})
	var customizedLogCalls atomic.Int32
	instance.CoreLoggerCustomizer = func(_ *Provider, logger *CoreLogger) log.CoreLogger {
		close(customizerEntered)
		<-releaseCustomizer
		return &strictCoreLogger{CoreLogger: logger, token: new(int), logCalls: &customizedLogCalls}
	}

	initialized := make(chan struct{})
	go func() {
		defer close(initialized)
		instance.GetRootLogger()
	}()
	<-customizerEntered
	root := instance.GetRootLogger()
	event := root.NewEvent(level.Info, map[string]any{"initial": "value"}).
		With("added", "value").
		Without("initial")
	close(releaseCustomizer)
	<-initialized

	assert.ToBeEqual(t, true, root.Accepts(event))
	root.Log(event, 0)
	assert.ToBeEqual(t, int32(1), customizedLogCalls.Load())
	assert.ToBeEqual(t, 1, recorder.Len())
	actual, exists := recorder.Get(0).Get("added")
	assert.ToBeEqual(t, true, exists)
	assert.ToBeEqual(t, "value", actual)
	_, exists = recorder.Get(0).Get("initial")
	assert.ToBeEqual(t, false, exists)
}

func Test_Provider_GetRootLogger_rematerializesWrappedEventAfterAccepts(t *testing.T) {
	instance, recorder := newProvider()
	instance.LocationDiscovery = location.NewCallerDiscovery()
	customizerEntered := make(chan struct{})
	releaseCustomizer := make(chan struct{})
	instance.CoreLoggerCustomizer = func(_ *Provider, logger *CoreLogger) log.CoreLogger {
		close(customizerEntered)
		<-releaseCustomizer
		return &strictCoreLogger{CoreLogger: logger, token: new(int)}
	}

	initialized := make(chan struct{})
	go func() {
		defer close(initialized)
		instance.GetRootLogger()
	}()
	<-customizerEntered
	root := instance.GetRootLogger()
	event := rootLoggerEventWithProgramCounter{
		Event:          root.NewEvent(level.Info, map[string]any{"field": "value"}),
		programCounter: programCounterForSlogRecord(),
	}
	assert.ToBeEqual(t, true, root.Accepts(event))

	close(releaseCustomizer)
	<-initialized
	root.Log(event, 0)

	assert.ToBeEqual(t, 1, recorder.Len())
	actual, exists := recorder.Get(0).Get("field")
	assert.ToBeEqual(t, true, exists)
	assert.ToBeEqual(t, "value", actual)
	actual, exists = recorder.Get(0).Get(instance.getFieldKeysSpec().GetLocation())
	assert.ToBeEqual(t, true, exists)
	assert.ToBeEqual(t, "github.com/echocat/slf4g/native.programCounterForSlogRecord", actual.(location.Caller).GetFrame().Function)
}

func Test_Provider_GetRootLogger_rematerializesEventLoggedDuringCustomization(t *testing.T) {
	instance, recorder := newProvider()
	customizerEntered := make(chan struct{})
	releaseCustomizer := make(chan struct{})
	instance.CoreLoggerCustomizer = func(_ *Provider, logger *CoreLogger) log.CoreLogger {
		close(customizerEntered)
		<-releaseCustomizer
		return &strictCoreLogger{CoreLogger: logger, token: new(int)}
	}

	initialized := make(chan any)
	go func() {
		defer func() { initialized <- recover() }()
		instance.GetRootLogger()
	}()
	<-customizerEntered
	root := instance.GetRootLogger()
	root.Log(root.NewEvent(level.Info, nil), 0)
	close(releaseCustomizer)

	assert.ToBeEqual(t, nil, <-initialized)
	assert.ToBeEqual(t, 1, recorder.Len())
}

func Test_Provider_GetRootLogger_preservesFinalCoreEventOwnership(t *testing.T) {
	instance, recorder := newProvider()
	instance.CoreLoggerCustomizer = func(_ *Provider, logger *CoreLogger) log.CoreLogger {
		return &strictOwnershipCoreLogger{CoreLogger: logger}
	}
	root := instance.GetRootLogger()
	event := root.NewEvent(level.Info, map[string]any{"field": "value"})

	assert.ToBeEqual(t, true, root.Accepts(event))
	root.Log(event, 0)

	assert.ToBeEqual(t, 1, recorder.Len())
	actual, exists := recorder.Get(0).Get("field")
	assert.ToBeEqual(t, true, exists)
	assert.ToBeEqual(t, "value", actual)
}

func Test_Provider_GetRootLogger_preservesOptionalInterfacesAndUnwrap(t *testing.T) {
	instance, _ := newProvider()
	expected := &CoreLogger{provider: instance, name: rootLoggerName}
	instance.CoreLoggerCustomizer = func(_ *Provider, _ *CoreLogger) log.CoreLogger {
		return log.NewLogger(expected).With("custom", "value")
	}

	actual := instance.GetRootLogger()

	assert.ToBeSame(t, expected, log.UnwrapCoreLogger(actual))
	_, isLoggerFacade := actual.(log.LoggerFacade)
	assert.ToBeEqual(t, true, isLoggerFacade)
	_, isEventFactoryWithFields := actual.(log.EventFactoryWithFields)
	assert.ToBeEqual(t, true, isEventFactoryWithFields)
}

func Test_Provider_GetRootLogger_unwrapsFullCustomLoggerWithoutUnwrap(t *testing.T) {
	instance, _ := newProvider()
	expected := &fullLoggerWithoutUnwrap{currentLevel: level.Debug}
	instance.CoreLoggerCustomizer = func(_ *Provider, logger *CoreLogger) log.CoreLogger {
		expected.Logger = log.NewLogger(logger)
		return expected
	}

	actual := instance.GetRootLogger()

	assert.ToBeSame(t, expected, log.UnwrapCoreLogger(actual))
	actualLevel, ok := level.Get(actual)
	assert.ToBeEqual(t, true, ok)
	assert.ToBeEqual(t, level.Debug, actualLevel)
	assert.ToBeEqual(t, true, level.Set(actual, level.Warn))
	assert.ToBeEqual(t, level.Warn, expected.currentLevel)
}

func Test_Provider_GetRootLogger_acceptsReentrantRootAsCustomization(t *testing.T) {
	instance, recorder := newProvider()
	instance.CoreLoggerCustomizer = func(actualProvider *Provider, _ *CoreLogger) log.CoreLogger {
		return actualProvider.GetRootLogger().With("custom", "value")
	}

	instance.GetRootLogger().Info("message")

	assert.ToBeEqual(t, 1, recorder.Len())
	actual, exists := recorder.Get(0).Get("custom")
	assert.ToBeEqual(t, true, exists)
	assert.ToBeEqual(t, "value", actual)
}

func Test_Provider_GetRootLogger_preservesCustomizedLoggerAfterPendingLogPanic(t *testing.T) {
	instance, recorder := newProvider()
	var customizerCalls atomic.Int32
	var consumerCalls atomic.Int32
	instance.Consumer = consumer.Func(func(event log.Event, source log.CoreLogger) {
		if consumerCalls.Add(1) == 1 {
			panic("expected")
		}
		recorder.Consume(event, source)
	})
	instance.CoreLoggerCustomizer = func(actualProvider *Provider, logger *CoreLogger) log.CoreLogger {
		customizerCalls.Add(1)
		actualProvider.GetRootLogger().Info("pending")
		return log.NewLogger(logger).With("custom", "value")
	}

	func() {
		defer func() { assert.ToBeEqual(t, "expected", recover()) }()
		instance.GetRootLogger()
	}()
	instance.GetRootLogger().Info("after")

	assert.ToBeEqual(t, int32(1), customizerCalls.Load())
	assert.ToBeEqual(t, int32(2), consumerCalls.Load())
	assert.ToBeEqual(t, 1, recorder.Len())
	actual, exists := recorder.Get(0).Get("custom")
	assert.ToBeEqual(t, true, exists)
	assert.ToBeEqual(t, "value", actual)
}

func Test_Provider_GetRootLogger_preservesCustomizedLoggerAfterPendingLogGoexit(t *testing.T) {
	instance, recorder := newProvider()
	var customizerCalls atomic.Int32
	var consumerCalls atomic.Int32
	instance.Consumer = consumer.Func(func(event log.Event, source log.CoreLogger) {
		if consumerCalls.Add(1) == 1 {
			runtime.Goexit()
		}
		recorder.Consume(event, source)
	})
	instance.CoreLoggerCustomizer = func(actualProvider *Provider, logger *CoreLogger) log.CoreLogger {
		customizerCalls.Add(1)
		actualProvider.GetRootLogger().Info("pending")
		return log.NewLogger(logger).With("custom", "value")
	}

	initialized := make(chan struct{})
	go func() {
		defer close(initialized)
		instance.GetRootLogger()
	}()
	<-initialized
	instance.GetRootLogger().Info("after")

	assert.ToBeEqual(t, int32(1), customizerCalls.Load())
	assert.ToBeEqual(t, int32(2), consumerCalls.Load())
	assert.ToBeEqual(t, 1, recorder.Len())
	actual, exists := recorder.Get(0).Get("custom")
	assert.ToBeEqual(t, true, exists)
	assert.ToBeEqual(t, "value", actual)
}

func Test_Provider_GetRootLogger_exposesCustomizerFailure(t *testing.T) {
	instance, _ := newProvider()
	instance.CoreLoggerCustomizer = func(_ *Provider, _ *CoreLogger) log.CoreLogger {
		panic("expected")
	}

	func() {
		defer func() { assert.ToBeEqual(t, "expected", recover()) }()
		instance.GetRootLogger()
	}()
	func() {
		defer func() { assert.ToBeEqual(t, "Root logger customizer did not complete.", recover()) }()
		instance.GetRootLogger().Info("message")
	}()
}

func Test_Provider_GetRootLogger_rejectsProvisionalEventAfterCustomizerFailure(t *testing.T) {
	instance, _ := newProvider()
	customizerEntered := make(chan struct{})
	releaseCustomizer := make(chan struct{})
	instance.CoreLoggerCustomizer = func(_ *Provider, _ *CoreLogger) log.CoreLogger {
		close(customizerEntered)
		<-releaseCustomizer
		panic("expected")
	}

	done := make(chan struct{})
	go func() {
		defer func() {
			assert.ToBeEqual(t, "expected", recover())
			close(done)
		}()
		instance.GetRootLogger()
	}()
	<-customizerEntered
	root := instance.GetRootLogger()
	event := root.NewEvent(level.Info, nil)
	close(releaseCustomizer)
	<-done

	defer func() { assert.ToBeEqual(t, "Root logger customizer did not complete.", recover()) }()
	root.Log(event, 0)
}

func Test_Provider_GetRootLogger_exposesCustomizerGoexit(t *testing.T) {
	instance, _ := newProvider()
	instance.CoreLoggerCustomizer = func(_ *Provider, _ *CoreLogger) log.CoreLogger {
		runtime.Goexit()
		return nil
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		instance.GetRootLogger()
	}()
	<-done

	defer func() { assert.ToBeEqual(t, "Root logger customizer did not complete.", recover()) }()
	instance.GetRootLogger().Info("message")
}

func logThroughCustomizedRoot(instance *Provider) {
	instance.GetRootLogger().Info("message")
}

func logThroughPendingCustomizedRoot(instance *Provider) {
	instance.GetRootLogger().Info("message")
}

func Test_Provider_GetRootLogger_preservesCallerLocation(t *testing.T) {
	instance, recorder := newProvider()
	instance.LocationDiscovery = location.NewCallerDiscovery()
	instance.CoreLoggerCustomizer = func(_ *Provider, logger *CoreLogger) log.CoreLogger {
		return logger
	}

	logThroughCustomizedRoot(instance)

	actual, exists := recorder.Get(0).Get(instance.getFieldKeysSpec().GetLocation())
	assert.ToBeEqual(t, true, exists)
	assert.ToBeEqual(t, "github.com/echocat/slf4g/native.logThroughCustomizedRoot", actual.(location.Caller).GetFrame().Function)
}

func Test_Provider_GetRootLogger_preservesPendingCallerLocation(t *testing.T) {
	instance, recorder := newProvider()
	instance.LocationDiscovery = location.NewCallerDiscovery()
	customizerEntered := make(chan struct{})
	releaseCustomizer := make(chan struct{})
	instance.CoreLoggerCustomizer = func(_ *Provider, logger *CoreLogger) log.CoreLogger {
		close(customizerEntered)
		<-releaseCustomizer
		return logger
	}

	initialized := make(chan struct{})
	go func() {
		defer close(initialized)
		instance.GetRootLogger()
	}()
	<-customizerEntered
	logThroughPendingCustomizedRoot(instance)
	close(releaseCustomizer)
	<-initialized

	actual, exists := recorder.Get(0).Get(instance.getFieldKeysSpec().GetLocation())
	assert.ToBeEqual(t, true, exists)
	assert.ToBeEqual(t, "github.com/echocat/slf4g/native.logThroughPendingCustomizedRoot", actual.(location.Caller).GetFrame().Function)
}

func Test_Provider_GetRootLogger_usesCustomizedLocationDiscoveryForPendingLog(t *testing.T) {
	instance, recorder := newProvider()
	instance.LocationDiscovery = location.DiscoveryFunc(func(log.Event, uint16) location.Location {
		return "provisional"
	})
	instance.CoreLoggerCustomizer = func(actualProvider *Provider, logger *CoreLogger) log.CoreLogger {
		logger.LocationDiscovery = location.DiscoveryFunc(func(log.Event, uint16) location.Location {
			return "customized"
		})
		actualProvider.GetRootLogger().Info("pending")
		return &strictCoreLogger{CoreLogger: logger, token: new(int)}
	}

	instance.GetRootLogger()

	actual, exists := recorder.Get(0).Get(instance.getFieldKeysSpec().GetLocation())
	assert.ToBeEqual(t, true, exists)
	assert.ToBeEqual(t, "customized", actual)
}

func programCounterForSlogRecord() uintptr {
	pcs := make([]uintptr, 1)
	runtime.Callers(1, pcs)
	return pcs[0]
}

func Test_Provider_slogHandlerUsesRecordProgramCounter(t *testing.T) {
	instance, recorder := newProvider()
	instance.LocationDiscovery = location.NewCallerDiscovery()
	handler := slogbridge.NewHandler(instance.GetRootLogger())
	record := stdslog.NewRecord(time.Now(), stdslog.LevelInfo, "message", programCounterForSlogRecord())
	record.AddAttrs(stdslog.String("location", "untrusted"))

	actualErr := handler.Handle(context.Background(), record)

	assert.ToBeNoError(t, actualErr)
	actual, exists := recorder.Get(0).Get(instance.getFieldKeysSpec().GetLocation())
	assert.ToBeEqual(t, true, exists)
	assert.ToBeEqual(t, "github.com/echocat/slf4g/native.programCounterForSlogRecord", actual.(location.Caller).GetFrame().Function)
}

func Test_Provider_slogHandlerUsesRecordProgramCounterDuringCustomization(t *testing.T) {
	instance, recorder := newProvider()
	instance.LocationDiscovery = location.NewCallerDiscovery()
	customizerEntered := make(chan struct{})
	releaseCustomizer := make(chan struct{})
	instance.CoreLoggerCustomizer = func(_ *Provider, logger *CoreLogger) log.CoreLogger {
		close(customizerEntered)
		<-releaseCustomizer
		return logger
	}

	initialized := make(chan struct{})
	go func() {
		defer close(initialized)
		instance.GetRootLogger()
	}()
	<-customizerEntered

	handler := slogbridge.NewHandler(instance.GetRootLogger())
	record := stdslog.NewRecord(time.Now(), stdslog.LevelInfo, "message", programCounterForSlogRecord())
	assert.ToBeNoError(t, handler.Handle(context.Background(), record))
	close(releaseCustomizer)
	<-initialized

	actual, exists := recorder.Get(0).Get(instance.getFieldKeysSpec().GetLocation())
	assert.ToBeEqual(t, true, exists)
	assert.ToBeEqual(t, "github.com/echocat/slf4g/native.programCounterForSlogRecord", actual.(location.Caller).GetFrame().Function)
}

func Test_Provider_slogHandlerPreservesFinalCoreEventOwnershipDuringCustomization(t *testing.T) {
	instance, recorder := newProvider()
	customizerEntered := make(chan struct{})
	releaseCustomizer := make(chan struct{})
	instance.CoreLoggerCustomizer = func(_ *Provider, logger *CoreLogger) log.CoreLogger {
		close(customizerEntered)
		<-releaseCustomizer
		return &strictOwnershipCoreLogger{CoreLogger: logger}
	}

	initialized := make(chan any)
	go func() {
		defer func() { initialized <- recover() }()
		instance.GetRootLogger()
	}()
	<-customizerEntered

	handler := slogbridge.NewHandler(instance.GetRootLogger())
	record := stdslog.NewRecord(time.Now(), stdslog.LevelInfo, "message", programCounterForSlogRecord())
	assert.ToBeNoError(t, handler.Handle(context.Background(), record))
	close(releaseCustomizer)

	assert.ToBeEqual(t, nil, <-initialized)
	assert.ToBeEqual(t, 1, recorder.Len())
}

func Test_Provider_levelAware(t *testing.T) {
	instance := &Provider{}

	actual, actualOk := level.Get(instance)
	assert.ToBeEqual(t, level.Info, actual)
	assert.ToBeEqual(t, true, actualOk)
	assert.ToBeEqual(t, true, level.Set(instance, level.Level(666)))

	actual2, actualOk2 := level.Get(instance)
	assert.ToBeEqual(t, level.Level(666), actual2)
	assert.ToBeEqual(t, true, actualOk2)
}

func Test_Provider_levelAwareLogger(t *testing.T) {
	instance := &Provider{}

	fooLogger := instance.GetLogger("foo")
	actual, actualOk := level.Get(fooLogger)
	assert.ToBeEqual(t, level.Info, actual)
	assert.ToBeEqual(t, true, actualOk)
	assert.ToBeEqual(t, true, level.Set(fooLogger, level.Level(666)))

	fooLogger2 := instance.GetLogger("foo")
	actual2, actualOk2 := level.Get(fooLogger2)
	assert.ToBeEqual(t, level.Level(666), actual2)
	assert.ToBeEqual(t, true, actualOk2)
}

func Test_init_providerWasRegistered(t *testing.T) {
	for _, candidate := range log.GetAllProviders() {
		if candidate == DefaultProvider {
			return
		}
	}
	assert.Fail(t, "Expected all providers to contain contain <%+v>; but got: <%+v>", DefaultProvider, log.GetAllProviders())
}

func newProvider(customizer ...func(*Provider)) (*Provider, *consumer.Recorder) {
	recorder := consumer.NewRecorder()
	result := &Provider{
		Name:     "mock",
		Consumer: recorder,
	}
	for _, c := range customizer {
		c(result)
	}
	return result, recorder
}

func newEvent(provider *Provider, l level.Level, customizer ...func(*event)) *event {
	result := &event{
		provider: provider,
		fields:   fields.Empty(),
		level:    l,
	}
	for _, c := range customizer {
		c(result)
	}
	return result
}

func (instance *Provider) getLevelName(l level.Level) string {
	result, err := instance.GetLevelNames().ToName(l)
	if err != nil {
		panic(err)
	}
	return result
}
