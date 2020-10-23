package native

import (
	"errors"
	"testing"
	"time"

	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/native/consumer"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/native/location"

	log "github.com/echocat/slf4g"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_CoreLogger_Log_withoutLoggerAndMessageButTimestamp(t *testing.T) {
	givenProvider, recorder := newProvider()

	for _, configuredLevel := range givenProvider.GetAllLevels() {
		t.Run(givenProvider.getLevelName(configuredLevel), func(t *testing.T) {

			givenProvider.SetLevel(configuredLevel)
			for _, eventLevel := range givenProvider.GetAllLevels() {
				t.Run(givenProvider.getLevelName(eventLevel), func(t *testing.T) {

					defer recorder.Reset()
					instance := newCoreLoggerWith(givenProvider)
					givenTimestamp := time.Now()
					givenEvent := newEvent(givenProvider, eventLevel).
						With("timestamp", givenTimestamp)
					expectedEvent := newEvent(givenProvider, eventLevel).
						With("timestamp", givenTimestamp).
						With("logger", "test")

					instance.Log(givenEvent, 123)

					if configuredLevel.CompareTo(eventLevel) <= 0 {
						assert.ToBeEqual(t, 1, recorder.Len())
						assert.ToBeEqual(t, expectedEvent, recorder.Get(0))
					} else {
						assert.ToBeEqual(t, 0, recorder.Len())
					}
				})
			}
		})
	}
}

func Test_CoreLogger_Log_withoutMessageButTimestampAndLogger(t *testing.T) {
	givenProvider, recorder := newProvider()

	for _, configuredLevel := range givenProvider.GetAllLevels() {
		t.Run(givenProvider.getLevelName(configuredLevel), func(t *testing.T) {

			givenProvider.SetLevel(configuredLevel)
			for _, eventLevel := range givenProvider.GetAllLevels() {
				t.Run(givenProvider.getLevelName(eventLevel), func(t *testing.T) {

					defer recorder.Reset()
					instance := newCoreLoggerWith(givenProvider)
					givenTimestamp := time.Now()
					givenEvent := newEvent(givenProvider, eventLevel).
						With("timestamp", givenTimestamp).
						With("logger", "test")
					expectedEvent := newEvent(givenProvider, eventLevel).
						With("timestamp", givenTimestamp).
						With("logger", "test")

					instance.Log(givenEvent, 123)

					if configuredLevel.CompareTo(eventLevel) <= 0 {
						assert.ToBeEqual(t, 1, recorder.Len())
						assert.ToBeEqual(t, expectedEvent, recorder.Get(0))
					} else {
						assert.ToBeEqual(t, 0, recorder.Len())
					}
				})
			}
		})
	}
}

func Test_CoreLogger_Log_withoutMessageButTimestampAndDifferentLogger(t *testing.T) {
	givenProvider, recorder := newProvider()

	for _, configuredLevel := range givenProvider.GetAllLevels() {
		t.Run(givenProvider.getLevelName(configuredLevel), func(t *testing.T) {

			givenProvider.SetLevel(configuredLevel)
			for _, eventLevel := range givenProvider.GetAllLevels() {
				t.Run(givenProvider.getLevelName(eventLevel), func(t *testing.T) {

					defer recorder.Reset()
					instance := newCoreLoggerWith(givenProvider)
					givenTimestamp := time.Now()
					givenEvent := newEvent(givenProvider, eventLevel).
						With("timestamp", givenTimestamp).
						With("logger", "somethingElse")
					expectedEvent := newEvent(givenProvider, eventLevel).
						With("timestamp", givenTimestamp).
						With("logger", "somethingElse").
						With("logger", "test")

					instance.Log(givenEvent, 123)

					if configuredLevel.CompareTo(eventLevel) <= 0 {
						assert.ToBeEqual(t, 1, recorder.Len())
						assert.ToBeEqual(t, expectedEvent, recorder.Get(0))
					} else {
						assert.ToBeEqual(t, 0, recorder.Len())
					}
				})
			}
		})
	}
}

func Test_CoreLogger_Log_withoutLoggerButTimestampAndMessage(t *testing.T) {
	givenProvider, recorder := newProvider()

	for _, configuredLevel := range givenProvider.GetAllLevels() {
		t.Run(givenProvider.getLevelName(configuredLevel), func(t *testing.T) {

			givenProvider.SetLevel(configuredLevel)
			for _, eventLevel := range givenProvider.GetAllLevels() {
				t.Run(givenProvider.getLevelName(eventLevel), func(t *testing.T) {

					defer recorder.Reset()
					instance := newCoreLoggerWith(givenProvider)
					givenTimestamp := time.Now()
					givenEvent := newEvent(givenProvider, eventLevel).
						With("timestamp", givenTimestamp).
						With("message", "foo")
					expectedEvent := newEvent(givenProvider, eventLevel).
						With("timestamp", givenTimestamp).
						With("message", "foo").
						With("logger", "test")

					instance.Log(givenEvent, 123)

					if configuredLevel.CompareTo(eventLevel) <= 0 {
						assert.ToBeEqual(t, 1, recorder.Len())
						assert.ToBeEqual(t, expectedEvent, recorder.Get(0))
					} else {
						assert.ToBeEqual(t, 0, recorder.Len())
					}
				})
			}
		})
	}
}

func Test_CoreLogger_Log_withTimestamp(t *testing.T) {
	givenProvider, recorder := newProvider()

	for _, configuredLevel := range givenProvider.GetAllLevels() {
		t.Run(givenProvider.getLevelName(configuredLevel), func(t *testing.T) {

			givenProvider.SetLevel(configuredLevel)
			for _, eventLevel := range givenProvider.GetAllLevels() {
				t.Run(givenProvider.getLevelName(eventLevel), func(t *testing.T) {

					defer recorder.Reset()
					instance := newCoreLoggerWith(givenProvider)
					givenEvent := newEvent(givenProvider, eventLevel).
						With("message", "foo")
					expectedEvent := newEvent(givenProvider, eventLevel).
						With("timestamp", time.Now()).
						With("message", "foo").
						With("logger", "test")

					instance.Log(givenEvent, 123)

					if configuredLevel.CompareTo(eventLevel) <= 0 {
						assert.ToBeEqual(t, 1, recorder.Len())
						assert.ToBeEqualUsing(t, expectedEvent, recorder.Get(0), log.DefaultEventEquality.WithIgnoringKeys("timestamp").AreEventsEqual)
						assert.ToBeNotNil(t, log.GetTimestampOf(recorder.Get(0), givenProvider))
					} else {
						assert.ToBeEqual(t, 0, recorder.Len())
					}
				})
			}
		})
	}
}

func Test_CoreLogger_Log_withLocation(t *testing.T) {
	givenProvider, recorder := newProvider(func(provider *Provider) {
		provider.LocationDiscovery = location.NewDepthOnlyDiscovery()
	})

	for _, configuredLevel := range givenProvider.GetAllLevels() {
		t.Run(givenProvider.getLevelName(configuredLevel), func(t *testing.T) {

			givenProvider.SetLevel(configuredLevel)
			for _, eventLevel := range givenProvider.GetAllLevels() {
				t.Run(givenProvider.getLevelName(eventLevel), func(t *testing.T) {

					defer recorder.Reset()
					instance := newCoreLoggerWith(givenProvider)
					givenTimestamp := time.Now()
					givenEvent := newEvent(givenProvider, eventLevel).
						With("timestamp", givenTimestamp).
						With("message", "foo")
					expectedEvent := newEvent(givenProvider, eventLevel).
						With("timestamp", givenTimestamp).
						With("message", "foo").
						With("logger", "test").
						With("location", location.DepthOnly(124))

					instance.Log(givenEvent, 123)

					if configuredLevel.CompareTo(eventLevel) <= 0 {
						assert.ToBeEqual(t, 1, recorder.Len())
						assert.ToBeEqual(t, expectedEvent, recorder.Get(0))
					} else {
						assert.ToBeEqual(t, 0, recorder.Len())
					}
				})
			}
		})
	}
}

func Test_CoreLogger_Log_nil(t *testing.T) {
	givenProvider, recorder := newProvider(func(provider *Provider) {
		provider.LocationDiscovery = location.NewDepthOnlyDiscovery()
	})

	for _, configuredLevel := range givenProvider.GetAllLevels() {
		t.Run(givenProvider.getLevelName(configuredLevel), func(t *testing.T) {
			defer recorder.Reset()
			givenProvider.SetLevel(configuredLevel)
			instance := newCoreLoggerWith(givenProvider)

			instance.Log(nil, 123)

			assert.ToBeEqual(t, 0, recorder.Len())
		})
	}
}

func Test_CoreLogger_IsLevelEnabled(t *testing.T) {
	givenProvider, _ := newProvider()

	for _, configuredLevel := range givenProvider.GetAllLevels() {
		t.Run(givenProvider.getLevelName(configuredLevel), func(t *testing.T) {

			givenProvider.SetLevel(configuredLevel)
			for _, givenLevel := range givenProvider.GetAllLevels() {
				t.Run(givenProvider.getLevelName(givenLevel), func(t *testing.T) {
					instance := newCoreLoggerWith(givenProvider)
					expected := configuredLevel.CompareTo(givenLevel) <= 0

					actual := instance.IsLevelEnabled(givenLevel)

					assert.ToBeEqual(t, expected, actual)
				})
			}
		})
	}
}

func Test_CoreLogger_GetName_specified(t *testing.T) {
	instance, _ := newCoreLogger()
	instance.name = "foo"

	assert.ToBeEqual(t, "foo", instance.GetName())
}

func Test_CoreLogger_GetName_panicsIfAbsent(t *testing.T) {
	instance, _ := newCoreLogger()
	instance.name = ""

	assert.Execution(t, func() {
		instance.GetName()
	}).WillPanicWith("^This .+ was not initiated by a .+\\.$")
}

func Test_CoreLogger_GetLevel_specified(t *testing.T) {
	instance, _ := newCoreLogger()
	instance.Level = level.Warn

	assert.ToBeEqual(t, level.Warn, instance.GetLevel())
}

func Test_CoreLogger_GetLevel_absent(t *testing.T) {
	instance, _ := newCoreLogger()
	instance.provider.SetLevel(level.Error)

	assert.ToBeEqual(t, level.Error, instance.GetLevel())
}

func Test_CoreLogger_SetLevel(t *testing.T) {
	instance, _ := newCoreLogger()

	assert.ToBeEqual(t, level.Level(0), instance.Level)

	for _, l := range level.GetProvider().GetLevels() {
		instance.SetLevel(l)
		assert.ToBeEqual(t, l, instance.Level)
	}

	instance.SetLevel(0)
	assert.ToBeEqual(t, level.Level(0), instance.Level)
}

func Test_CoreLogger_GetProvider(t *testing.T) {
	givenProvider, _ := newProvider()
	instance := newCoreLoggerWith(givenProvider)

	actual := instance.GetProvider()

	assert.ToBeSame(t, givenProvider, actual)
}

func Test_CoreLogger_GetProvider_panicsOnNil(t *testing.T) {
	instance := newCoreLoggerWith(nil)

	assert.Execution(t, func() {
		instance.GetProvider()
	}).WillPanicWith("^This .+ was not initiated by a .+\\.$")
}

func Test_CoreLogger_NewEvent(t *testing.T) {
	instance, _ := newCoreLogger()

	assert.ToBeEqual(t, &event{
		provider: instance.provider,
		fields:   fields.Empty(),
		level:    level.Fatal,
	}, instance.NewEvent(level.Fatal, nil))

	assert.ToBeEqual(t, &event{
		provider: instance.provider,
		fields:   fields.WithAll(map[string]interface{}{"foo": "bar"}),
		level:    level.Fatal,
	}, instance.NewEvent(level.Fatal, map[string]interface{}{"foo": "bar"}))
}

func Test_CoreLogger_NewEventWithFields(t *testing.T) {
	instance, _ := newCoreLogger()

	assert.ToBeEqual(t, &event{
		provider: instance.provider,
		fields:   fields.Empty(),
		level:    level.Fatal,
	}, instance.NewEventWithFields(level.Fatal, nil))

	assert.ToBeEqual(t, &event{
		provider: instance.provider,
		fields:   fields.With("foo", "bar"),
		level:    level.Fatal,
	}, instance.NewEventWithFields(level.Fatal, fields.With("foo", "bar")))
}

func Test_CoreLogger_NewEventWithFields_panicsOnError(t *testing.T) {
	instance, _ := newCoreLogger()

	assert.Execution(t, func() {
		instance.NewEventWithFields(level.Fatal, fields.ForEachFunc(func(func(string, interface{}) error) error {
			return errors.New("expected")
		}))
	}).WillPanicWith("^expected$")
}

func Test_CoreLogger_Accepts(t *testing.T) {
	instance, _ := newCoreLogger()

	assert.ToBeEqual(t, false, instance.Accepts(nil))
	assert.ToBeEqual(t, true, instance.Accepts(&event{}))
}

func Test_CoreLogger_getConsumer_specified(t *testing.T) {
	givenConsumer := consumer.NewRecorder()
	instance, _ := newCoreLogger()
	instance.Consumer = givenConsumer

	actual := instance.getConsumer()

	assert.ToBeSame(t, givenConsumer, actual)
}

func Test_CoreLogger_getConsumer_fromProvider(t *testing.T) {
	givenConsumer := consumer.NewRecorder()
	instance, _ := newCoreLogger()
	instance.provider.Consumer = givenConsumer

	actual := instance.getConsumer()

	assert.ToBeSame(t, givenConsumer, actual)
}

func Test_CoreLogger_getLocationDiscovery_specified(t *testing.T) {
	givenLocationDiscovery := location.NewCallerDiscovery()
	instance, _ := newCoreLogger()
	instance.LocationDiscovery = givenLocationDiscovery

	actual := instance.getLocationDiscovery()

	assert.ToBeSame(t, givenLocationDiscovery, actual)
}

func Test_CoreLogger_getLocationDiscovery_fromProvider(t *testing.T) {
	givenLocationDiscovery := location.NewCallerDiscovery()
	instance, _ := newCoreLogger()
	instance.provider.LocationDiscovery = givenLocationDiscovery

	actual := instance.getLocationDiscovery()

	assert.ToBeSame(t, givenLocationDiscovery, actual)
}

func newCoreLogger(customizer ...func(*CoreLogger)) (*CoreLogger, *consumer.Recorder) {
	provider, recorder := newProvider()
	return newCoreLoggerWith(provider, customizer...), recorder
}

func newCoreLoggerWith(provider *Provider, customizer ...func(*CoreLogger)) *CoreLogger {
	result := &CoreLogger{
		name:     "test",
		provider: provider,
	}
	for _, c := range customizer {
		c(result)
	}
	return result
}
