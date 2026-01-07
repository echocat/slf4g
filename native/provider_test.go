package native

import (
	"testing"

	log "github.com/echocat/slf4g"

	"github.com/echocat/slf4g/native/location"

	nlevel "github.com/echocat/slf4g/native/level"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/native/consumer"
)

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
	instance, _ := newProvider()

	assert.ToBeEqual(t, level.Info, instance.GetLevel())

	for _, l := range instance.GetAllLevels() {
		instance.Level = l
		assert.ToBeEqual(t, l, instance.GetLevel())
	}

	instance.Level = 0
	assert.ToBeEqual(t, level.Info, instance.GetLevel())
}

func Test_Provider_SetLevel(t *testing.T) {
	instance, _ := newProvider()

	assert.ToBeEqual(t, level.Level(0), instance.Level)

	for _, l := range instance.GetAllLevels() {
		instance.SetLevel(l)
		assert.ToBeEqual(t, l, instance.Level)
	}

	instance.SetLevel(0)
	assert.ToBeEqual(t, level.Level(0), instance.Level)
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
