package fields

import (
	"fmt"
	"github.com/echocat/slf4g/level"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

var veryComplexValue = struct{}{}
var filterContextWithLeveInfo = filterContext{
	level: level.Info,
}
var filterContextWithLeveDebug = filterContext{
	level: level.Debug,
}

func ExampleRequireMaximalLevel() {
	filteredValue := RequireMaximalLevel(level.Debug, veryComplexValue)

	// Will be <nil>, <false>
	fmt.Println(filteredValue.Filter(filterContextWithLeveInfo))

	// Will be <veryComplexValue>, <true>
	fmt.Println(filteredValue.Filter(filterContextWithLeveDebug))
}

func Test_RequireMaximalLevelLazy_Get(t *testing.T) {
	expected := struct{ foo string }{foo: "bar"}
	givenLazy := LazyFunc(func() interface{} { return expected })

	actualInstance := RequireMaximalLevelLazy(level.Info, givenLazy)
	actual := actualInstance.Get()

	assert.ToBeEqual(t, expected, actual)
}

func Test_RequireMaximalLevelLazy_Filter_respected(t *testing.T) {
	expected := struct{ foo string }{foo: "bar"}
	givenLazy := LazyFunc(func() interface{} { return expected })

	actualInstance := RequireMaximalLevelLazy(level.Debug, givenLazy)
	actual, actualRespected := actualInstance.Filter(filterContextWithLeveDebug)

	assert.ToBeEqual(t, expected, actual)
	assert.ToBeEqual(t, true, actualRespected)
}

func Test_RequireMaximalLevelLazy_Filter_ignored(t *testing.T) {
	givenLazy := LazyFunc(func() interface{} { return struct{ foo string }{foo: "bar"} })

	actualInstance := RequireMaximalLevelLazy(level.Debug, givenLazy)
	actual, actualRespected := actualInstance.Filter(filterContextWithLeveInfo)

	assert.ToBeNil(t, actual)
	assert.ToBeEqual(t, false, actualRespected)
}

func Test_RequireMaximalLevel_Filter_respected(t *testing.T) {
	expected := struct{ foo string }{foo: "bar"}

	actualInstance := RequireMaximalLevel(level.Debug, expected)
	actual, actualRespected := actualInstance.Filter(filterContextWithLeveDebug)

	assert.ToBeEqual(t, expected, actual)
	assert.ToBeEqual(t, true, actualRespected)
}

func Test_RequireMaximalLevel_Filter_ignored(t *testing.T) {
	actualInstance := RequireMaximalLevel(level.Debug, struct{ foo string }{foo: "bar"})
	actual, actualRespected := actualInstance.Filter(filterContextWithLeveInfo)

	assert.ToBeNil(t, actual)
	assert.ToBeEqual(t, false, actualRespected)
}

func Test_IgnoreLevelsLazy_Get(t *testing.T) {
	expected := struct{ foo string }{foo: "bar"}
	givenLazy := LazyFunc(func() interface{} { return expected })

	actualInstance := IgnoreLevelsLazy(level.Info, level.Warn, givenLazy)
	actual := actualInstance.Get()

	assert.ToBeEqual(t, expected, actual)
}

func Test_IgnoreLevelsLazy_Filter_respectedBelow(t *testing.T) {
	expected := struct{ foo string }{foo: "bar"}
	givenLazy := LazyFunc(func() interface{} { return expected })

	actualInstance := IgnoreLevelsLazy(level.Info, level.Warn, givenLazy)
	actual, actualRespected := actualInstance.Filter(filterContext{level: level.Info - 1})

	assert.ToBeEqual(t, expected, actual)
	assert.ToBeEqual(t, true, actualRespected)
}

func Test_IgnoreLevelsLazy_Filter_respectedAbove(t *testing.T) {
	expected := struct{ foo string }{foo: "bar"}
	givenLazy := LazyFunc(func() interface{} { return expected })

	actualInstance := IgnoreLevelsLazy(level.Info, level.Warn, givenLazy)
	actual, actualRespected := actualInstance.Filter(filterContext{level: level.Warn})

	assert.ToBeEqual(t, expected, actual)
	assert.ToBeEqual(t, true, actualRespected)
}

func Test_IgnoreLevelsLazy_Filter_ignored(t *testing.T) {
	givenLazy := LazyFunc(func() interface{} { return struct{ foo string }{foo: "bar"} })

	actualInstance := IgnoreLevelsLazy(level.Info, level.Warn, givenLazy)
	actual, actualRespected := actualInstance.Filter(filterContextWithLeveInfo)

	assert.ToBeNil(t, actual)
	assert.ToBeEqual(t, false, actualRespected)
}

func Test_IgnoreLevels_Filter_respectedBelow(t *testing.T) {
	expected := struct{ foo string }{foo: "bar"}

	actualInstance := IgnoreLevels(level.Info, level.Warn, expected)
	actual, actualRespected := actualInstance.Filter(filterContext{level: level.Info - 1})

	assert.ToBeEqual(t, expected, actual)
	assert.ToBeEqual(t, true, actualRespected)
}

func Test_IgnoreLevels_Filter_respectedAbove(t *testing.T) {
	expected := struct{ foo string }{foo: "bar"}

	actualInstance := IgnoreLevels(level.Info, level.Warn, expected)
	actual, actualRespected := actualInstance.Filter(filterContext{level: level.Warn})

	assert.ToBeEqual(t, expected, actual)
	assert.ToBeEqual(t, true, actualRespected)
}

func Test_IgnoreLevels_Filter_ignored(t *testing.T) {
	actualInstance := IgnoreLevels(level.Info, level.Warn, struct{ foo string }{foo: "bar"})
	actual, actualRespected := actualInstance.Filter(filterContextWithLeveInfo)

	assert.ToBeNil(t, actual)
	assert.ToBeEqual(t, false, actualRespected)
}

type filterContext struct {
	level  level.Level
	fields map[string]interface{}
}

func (instance filterContext) GetLevel() level.Level {
	return instance.level
}

func (instance filterContext) Get(key string) (value interface{}, exists bool) {
	if instance.fields == nil {
		return nil, false
	}
	v, ok := instance.fields[key]
	return v, ok
}
