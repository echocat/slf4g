package recording

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/fields"

	"github.com/echocat/slf4g/level"

	log "github.com/echocat/slf4g"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_NewProvider(t *testing.T) {
	actual := NewProvider()

	assert.ToBeEqual(t, &Provider{}, actual)
}

func Test_Provider_HookGlobally(t *testing.T) {
	instance := NewProvider()

	assert.ToBeEqual(t, true, log.IsFallbackProvider(log.GetProvider()))

	finalizer := instance.HookGlobally()

	finalizerCalled := false
	defer func() {
		if !finalizerCalled {
			finalizer()
		}
	}()

	actual := log.UnwrapProvider(log.GetProvider())

	assert.ToBeEqual(t, false, log.IsFallbackProvider(log.GetProvider()))
	assert.ToBeSame(t, instance, actual)

	finalizer()
	finalizerCalled = true

	assert.ToBeEqual(t, true, log.IsFallbackProvider(log.GetProvider()))
}

func Test_Provider_MustContains(t *testing.T) {
	instance := NewProvider()
	instanceRootLogger := instance.GetRootLogger()
	instanceLogger := instance.GetLogger("foo")

	expectedEvent1 := instanceRootLogger.NewEvent(level.Info, nil).
		With("message", "a")
	expectedEvent2 := instanceRootLogger.NewEvent(level.Warn, nil).
		With("message", "b").
		With("foo", 1)
	unExpectedEvent3 := instanceRootLogger.NewEvent(level.Info, nil).
		With("message", "b").
		With("foo", 1)

	assert.ToBeEqual(t, 0, instance.Len())

	instanceRootLogger.Info("a")
	instanceLogger.With("foo", 1).Warn("b")

	actual1 := instance.MustContains(expectedEvent1)
	assert.ToBeEqual(t, true, actual1)

	actual2 := instance.MustContains(expectedEvent2)
	assert.ToBeEqual(t, true, actual2)

	actual3 := instance.MustContains(unExpectedEvent3)
	assert.ToBeEqual(t, false, actual3)
}

func Test_Provider_MustContains_panicsOnErrorsForRoot(t *testing.T) {
	previous := fields.DefaultValueEquality
	defer func() { fields.DefaultValueEquality = previous }()
	fields.DefaultValueEquality = fields.ValueEqualityFunc(func(name string, left, right interface{}) (bool, error) {
		if name == "message" && (left == "root" || right == "root") {
			return false, errors.New("expected")
		}
		return false, nil
	})

	instance := NewProvider()
	instance.GetRootLogger().Info("root")
	instance.GetLogger("bar").Info("fromOther")

	givenEvent := instance.GetRootLogger().NewEvent(level.Info, nil).
		With("message", "something")

	assert.Execution(t, func() {
		instance.MustContains(givenEvent)
	}).WillPanicWith("^expected$")
}

func Test_Provider_MustContains_panicsOnErrors(t *testing.T) {
	previous := fields.DefaultValueEquality
	defer func() { fields.DefaultValueEquality = previous }()
	fields.DefaultValueEquality = fields.ValueEqualityFunc(func(name string, left, right interface{}) (bool, error) {
		if name == "message" && (left == "fromOther" || right == "fromOther") {
			return false, errors.New("expected")
		}
		return false, nil
	})

	instance := NewProvider()
	instance.GetRootLogger().Info("root")
	instance.GetLogger("bar").Info("fromOther")

	givenEvent := instance.GetRootLogger().NewEvent(level.Info, nil).
		With("message", "something")

	assert.Execution(t, func() {
		instance.MustContains(givenEvent)
	}).WillPanicWith("^expected$")
}

func Test_Provider_MustContainsCustom(t *testing.T) {
	instance := NewProvider()
	instanceRootLogger := instance.GetRootLogger()
	instanceLogger := instance.GetLogger("foo")
	givenEquality := log.DefaultEventEquality.WithIgnoringKeys("foo", "timestamp", "logger")

	expectedEvent1 := instanceRootLogger.NewEvent(level.Info, nil).
		With("message", "a").
		With("foo", 666).
		With("bar", 2)
	expectedEvent2 := instanceRootLogger.NewEvent(level.Warn, nil).
		With("message", "b").
		With("foo", 666).
		With("bar", 2)
	unExpectedEvent3 := instanceRootLogger.NewEvent(level.Info, nil).
		With("message", "b").
		With("foo", 666).
		With("bar", 2)

	assert.ToBeEqual(t, 0, instance.Len())

	instanceRootLogger.
		With("foo", 1).
		With("bar", 2).
		Info("a")
	instanceLogger.
		With("foo", 1).
		With("bar", 2).
		Warn("b")

	assert.ToBeEqual(t, 2, instance.Len())

	actual1 := instance.MustContainsCustom(givenEquality, expectedEvent1)
	assert.ToBeEqual(t, true, actual1)

	actual2 := instance.MustContainsCustom(givenEquality, expectedEvent2)
	assert.ToBeEqual(t, true, actual2)

	actual3 := instance.MustContainsCustom(givenEquality, unExpectedEvent3)
	assert.ToBeEqual(t, false, actual3)
}

func Test_Provider_MustContainsCustom_panics(t *testing.T) {
	givenEquality := log.EventEqualityFunc(func(left, right log.Event) (bool, error) {
		return false, errors.New("expected")
	})

	instance := NewProvider()
	instance.GetRootLogger().Info("foo")

	givenEvent := instance.GetRootLogger().NewEvent(level.Info, nil).
		With("message", "foo")

	assert.Execution(t, func() {
		instance.MustContainsCustom(givenEquality, givenEvent)
	}).WillPanicWith("^expected$")
}

func Test_Provider_Len(t *testing.T) {
	instance := NewProvider()
	instanceRootLogger := instance.GetRootLogger()
	instanceLogger := instance.GetLogger("foo")

	assert.ToBeEqual(t, 0, instance.Len())

	instanceRootLogger.
		With("foo", 1).
		With("bar", 2).
		Info("a")
	instanceLogger.
		With("foo", 1).
		With("bar", 2).
		Warn("b")

	assert.ToBeEqual(t, 2, instance.Len())
}

func Test_Provider_GetAll(t *testing.T) {
	instance := NewProvider()
	instanceRootLogger := instance.GetRootLogger()
	instanceLogger := instance.GetLogger("foo")

	assert.ToBeEqual(t, 0, instance.Len())

	instanceRootLogger.
		With("foo", 1).
		With("bar", 2).
		Info("a")
	instanceLogger.
		With("foo", 1).
		With("bar", 2).
		Warn("b")

	actual := instance.GetAll()

	assert.ToBeEqual(t, 2, len(actual))
	assert.ToBeEqualUsing(t, instanceRootLogger.NewEvent(level.Info, nil).
		With("message", "a").
		With("foo", 1).
		With("bar", 2), actual[0], instance.defaultEventEquality().AreEventsEqual)
	assert.ToBeEqualUsing(t, instanceRootLogger.NewEvent(level.Warn, nil).
		With("message", "b").
		With("foo", 1).
		With("bar", 2), actual[1], instance.defaultEventEquality().AreEventsEqual)
}

func Test_Provider_GetAllRoot(t *testing.T) {
	instance := NewProvider()
	instanceRootLogger := instance.GetRootLogger()
	instanceLogger := instance.GetLogger("foo")

	assert.ToBeEqual(t, 0, instance.Len())

	instanceRootLogger.
		With("foo", 1).
		With("bar", 2).
		Info("a")
	instanceRootLogger.
		With("foo", 1).
		With("bar", 2).
		Warn("b")
	instanceLogger.
		With("foo", 1).
		With("bar", 2).
		Error("c")

	actual := instance.GetAllRoot()

	assert.ToBeEqual(t, 2, len(actual))
	assert.ToBeEqualUsing(t, instanceRootLogger.NewEvent(level.Info, nil).
		With("message", "a").
		With("foo", 1).
		With("bar", 2), actual[0], instance.defaultEventEquality().AreEventsEqual)
	assert.ToBeEqualUsing(t, instanceRootLogger.NewEvent(level.Warn, nil).
		With("message", "b").
		With("foo", 1).
		With("bar", 2), actual[1], instance.defaultEventEquality().AreEventsEqual)
}

func Test_Provider_GetOf(t *testing.T) {
	instance := NewProvider()
	instanceRootLogger := instance.GetRootLogger()
	instanceLogger := instance.GetLogger("foo")

	assert.ToBeEqual(t, 0, instance.Len())

	instanceRootLogger.
		With("foo", 1).
		With("bar", 2).
		Info("a")
	instanceLogger.
		With("foo", 1).
		With("bar", 2).
		Warn("b")
	instanceLogger.
		With("foo", 1).
		With("bar", 2).
		Error("c")

	actual1 := instance.GetAllOf("foo")

	assert.ToBeEqual(t, 2, len(actual1))
	assert.ToBeEqualUsing(t, instanceRootLogger.NewEvent(level.Warn, nil).
		With("message", "b").
		With("foo", 1).
		With("bar", 2), actual1[0], instance.defaultEventEquality().AreEventsEqual)
	assert.ToBeEqualUsing(t, instanceRootLogger.NewEvent(level.Error, nil).
		With("message", "c").
		With("foo", 1).
		With("bar", 2), actual1[1], instance.defaultEventEquality().AreEventsEqual)

	actual2 := instance.GetAllOf("bar")
	assert.ToBeEqual(t, 0, len(actual2))
}

func Test_Provider_ResetAll(t *testing.T) {
	instance := NewProvider()
	instanceRootLogger := instance.getRootLogger()
	instanceLogger := instance.getLogger("foo")

	assert.ToBeEqual(t, 0, instance.Len())
	assert.ToBeEqual(t, 0, instanceRootLogger.Len())
	assert.ToBeEqual(t, 0, instanceLogger.Len())

	instanceRootLogger.
		With("foo", 1).
		With("bar", 2).
		Info("a")
	instanceLogger.
		With("foo", 1).
		With("bar", 2).
		Warn("b")

	assert.ToBeEqual(t, 2, instance.Len())
	assert.ToBeEqual(t, 1, instanceRootLogger.Len())
	assert.ToBeEqual(t, 1, instanceLogger.Len())

	instance.ResetAll()

	assert.ToBeEqual(t, 0, instance.Len())
	assert.ToBeEqual(t, 0, instanceRootLogger.Len())
	assert.ToBeEqual(t, 0, instanceLogger.Len())
}

func Test_Provider_ResetRoot(t *testing.T) {
	instance := NewProvider()
	instanceRootLogger := instance.getRootLogger()
	instanceLogger := instance.getLogger("foo")

	assert.ToBeEqual(t, 0, instance.Len())
	assert.ToBeEqual(t, 0, instanceRootLogger.Len())
	assert.ToBeEqual(t, 0, instanceLogger.Len())

	instanceRootLogger.
		With("foo", 1).
		With("bar", 2).
		Info("a")
	instanceLogger.
		With("foo", 1).
		With("bar", 2).
		Warn("b")

	assert.ToBeEqual(t, 2, instance.Len())
	assert.ToBeEqual(t, 1, instanceRootLogger.Len())
	assert.ToBeEqual(t, 1, instanceLogger.Len())

	instance.ResetRoot()

	assert.ToBeEqual(t, 1, instance.Len())
	assert.ToBeEqual(t, 0, instanceRootLogger.Len())
	assert.ToBeEqual(t, 1, instanceLogger.Len())
}

func Test_Provider_Reset(t *testing.T) {
	instance := NewProvider()
	instanceRootLogger := instance.getRootLogger()
	instanceLogger1 := instance.getLogger("foo")
	instanceLogger2 := instance.getLogger("bar")

	assert.ToBeEqual(t, 0, instance.Len())
	assert.ToBeEqual(t, 0, instanceRootLogger.Len())
	assert.ToBeEqual(t, 0, instanceLogger1.Len())
	assert.ToBeEqual(t, 0, instanceLogger2.Len())

	instanceRootLogger.
		With("foo", 1).
		With("bar", 2).
		Info("a")
	instanceLogger1.
		With("foo", 1).
		With("bar", 2).
		Warn("b")
	instanceLogger2.
		With("foo", 1).
		With("bar", 2).
		Error("c")

	assert.ToBeEqual(t, 3, instance.Len())
	assert.ToBeEqual(t, 1, instanceRootLogger.Len())
	assert.ToBeEqual(t, 1, instanceLogger1.Len())
	assert.ToBeEqual(t, 1, instanceLogger2.Len())

	instance.Reset("foo")

	assert.ToBeEqual(t, 2, instance.Len())
	assert.ToBeEqual(t, 1, instanceRootLogger.Len())
	assert.ToBeEqual(t, 0, instanceLogger1.Len())
	assert.ToBeEqual(t, 1, instanceLogger2.Len())
}

func Test_Provider_GetName_specified(t *testing.T) {
	instance := NewProvider()
	instance.Name = "foo"

	assert.ToBeEqual(t, "foo", instance.GetName())
}

func Test_Provider_GetName_absent(t *testing.T) {
	instance := NewProvider()

	assert.ToBeEqual(t, DefaultProviderName, instance.GetName())
}

func Test_Provider_GetAllLevels_specified(t *testing.T) {
	instance := NewProvider()
	instance.AllLevels = level.Levels{level.Warn, level.Fatal}

	assert.ToBeEqual(t, level.Levels{level.Warn, level.Fatal}, instance.GetAllLevels())
}

func Test_Provider_GetAllLevels_absent(t *testing.T) {
	instance := NewProvider()

	assert.ToBeEqual(t, level.GetProvider().GetLevels(), instance.GetAllLevels())
}

func Test_Provider_GetFieldKeysSpec_specified(t *testing.T) {
	givenSpec := &fields.KeysSpecImpl{
		Timestamp: "1",
		Message:   "2",
		Logger:    "3",
		Error:     "4",
	}
	instance := NewProvider()
	instance.FieldKeysSpec = givenSpec

	assert.ToBeSame(t, givenSpec, instance.GetFieldKeysSpec())
}

func Test_Provider_GetFieldKeysSpec_absent(t *testing.T) {
	instance := NewProvider()

	assert.ToBeEqual(t, &fields.KeysSpecImpl{}, instance.GetFieldKeysSpec())
}

func Test_Provider_GetLevel_specified(t *testing.T) {
	instance := NewProvider()
	instance.Level = level.Warn

	assert.ToBeEqual(t, level.Warn, instance.GetLevel())
}

func Test_Provider_GetLevel_absent(t *testing.T) {
	instance := NewProvider()

	assert.ToBeEqual(t, level.Info, instance.GetLevel())
}

func Test_Provider_GetLevel(t *testing.T) {
	instance := NewProvider()

	assert.ToBeEqual(t, level.Info, instance.GetLevel())

	for _, l := range instance.GetAllLevels() {
		instance.Level = l
		assert.ToBeEqual(t, l, instance.GetLevel())
	}

	instance.Level = 0
	assert.ToBeEqual(t, level.Info, instance.GetLevel())
}

func Test_Provider_SetLevel(t *testing.T) {
	instance := NewProvider()

	assert.ToBeEqual(t, level.Level(0), instance.Level)

	for _, l := range instance.GetAllLevels() {
		instance.SetLevel(l)
		assert.ToBeEqual(t, l, instance.Level)
	}

	instance.SetLevel(0)
	assert.ToBeEqual(t, level.Level(0), instance.Level)
}
