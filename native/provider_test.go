package native

import (
	"testing"

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
	givenSpec := &FieldKeysSpecImpl{
		KeysSpecImpl: fields.KeysSpecImpl{
			Timestamp: "1",
			Message:   "2",
			Logger:    "3",
			Error:     "4",
		},
		Location: "5",
	}
	instance, _ := newProvider()
	instance.FieldKeysSpec = givenSpec

	assert.ToBeSame(t, givenSpec, instance.GetFieldKeysSpec())
}

func Test_Provider_GetFieldKeysSpec_absent(t *testing.T) {
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
