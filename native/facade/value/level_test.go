package value

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
	nlevel "github.com/echocat/slf4g/native/level"
)

func Test_NewLevel(t *testing.T) {
	givenTarget := &mockLevelTarget{}
	instance := NewLevel(givenTarget)

	assert.ToBeNotNil(t, instance)
	assert.ToBeSame(t, givenTarget, instance.Target)
	assert.ToBeNil(t, instance.Names)
}

func Test_NewLevel_customized(t *testing.T) {
	givenTarget := &mockLevelTarget{}
	givenNames := &failingLevelNames{}
	instance := NewLevel(givenTarget, func(level *Level) {
		level.Names = givenNames
	})

	assert.ToBeNotNil(t, instance)
	assert.ToBeSame(t, givenTarget, instance.Target)
	assert.ToBeSame(t, givenNames, instance.Names)
}

func Test_Level_UnmarshalText(t *testing.T) {
	givenTarget := &mockLevelTarget{}
	instance := NewLevel(givenTarget)

	actualErr := instance.UnmarshalText([]byte("warn"))

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, level.Warn, givenTarget.level)
}

func Test_Level_UnmarshalText_number(t *testing.T) {
	givenTarget := &mockLevelTarget{}
	instance := NewLevel(givenTarget)

	actualErr := instance.UnmarshalText([]byte("666"))

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, level.Level(666), givenTarget.level)
}

func Test_Level_UnmarshalText_failing(t *testing.T) {
	givenTarget := &mockLevelTarget{level: level.Info}
	instance := NewLevel(givenTarget)

	actualErr := instance.UnmarshalText([]byte("foo"))

	assert.ToBeMatching(t, "^illegal level: foo$", actualErr)
	assert.ToBeEqual(t, level.Info, givenTarget.level)
}

func Test_Level_String(t *testing.T) {
	givenTarget := &mockLevelTarget{level: level.Warn}
	instance := NewLevel(givenTarget)

	actual := instance.String()

	assert.ToBeEqual(t, "WARN", actual)
}

func Test_Level_String_withEmptyTarget(t *testing.T) {
	instance := NewLevel(nil)

	actual := instance.String()

	assert.ToBeEqual(t, "", actual)
}

func Test_Level_String_failing(t *testing.T) {
	givenTarget := &mockLevelTarget{level: level.Level(666)}
	instance := NewLevel(givenTarget, func(l *Level) {
		l.Names = failingLevelNames{nlevel.NewNames()}
	})

	actual := instance.String()

	assert.ToBeEqual(t, "ERR-expectedError", actual)
}

func Test_Level_Get(t *testing.T) {
	givenTarget := &mockLevelTarget{level: level.Warn}
	instance := NewLevel(givenTarget)

	actual := instance.Get()

	assert.ToBeEqual(t, level.Warn, actual)
}

func Test_Level_getNames_explicit(t *testing.T) {
	instance := &Level{
		Names: &failingLevelNames{},
	}

	actual := instance.getNames()

	assert.ToBeSame(t, instance.Names, actual)
}

func Test_Level_getNames_namesAware(t *testing.T) {
	instance := &Level{
		Target: &mockLevelTarget{},
	}

	actual := instance.getNames()

	assert.ToBeSame(t, mockLevelNames, actual)
}

func Test_Level_getNames_default(t *testing.T) {
	instance := &Level{
		Names: nil,
	}

	actual := instance.getNames()

	assert.ToBeEqual(t, nlevel.DefaultNames, actual)
}

func Test_Level_getNames_noop(t *testing.T) {
	old := nlevel.DefaultNames
	defer func() {
		nlevel.DefaultNames = old
	}()
	nlevel.DefaultNames = nil

	instance := &Level{
		Names: nil,
	}

	actual := instance.getNames()

	assert.ToBeEqual(t, nlevel.NewNames(), actual)
}

var mockLevelNames = nlevel.NewNames()

type mockLevelTarget struct {
	level level.Level
}

func (instance *mockLevelTarget) GetLevel() level.Level {
	return instance.level
}

func (instance *mockLevelTarget) SetLevel(v level.Level) {
	instance.level = v
}

func (instance *mockLevelTarget) GetLevelNames() nlevel.Names {
	return mockLevelNames
}

type failingLevelNames struct {
	nlevel.Names
}

func (instance failingLevelNames) ToName(level.Level) (string, error) {
	return "", errors.New("expectedError")
}
