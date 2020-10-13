package names

import (
	"testing"

	"github.com/echocat/slf4g/fields"

	"github.com/echocat/slf4g/internal/test/assert"
)

var forPackageSomethingFromInit = CurrentPackageLoggerNameGenerator(0)

func Test_FullLoggerNameGenerator_usesCustomizer(t *testing.T) {
	before := FullLoggerNameCustomizer
	defer func() { FullLoggerNameCustomizer = before }()

	givenSomething := struct{ foo string }{"bar"}
	called := false
	FullLoggerNameCustomizer = func(something interface{}) string {
		assert.ToBeEqual(t, givenSomething, something)
		called = true
		return "xyz"
	}

	actual := FullLoggerNameGenerator(givenSomething)
	assert.ToBeEqual(t, true, called)
	assert.ToBeEqual(t, "xyz", actual)
}

func Test_FullLoggerNameGenerator_panics_withNil(t *testing.T) {
	assert.Execution(t, func() {
		FullLoggerNameGenerator(nil)
	}).WillPanicWith("^invalid value to receive a package from: <nil>$")
}

func Test_FullLoggerNameGenerator_panics_primitive(t *testing.T) {
	assert.Execution(t, func() {
		FullLoggerNameGenerator(123)
	}).WillPanicWith("^invalid value to receive a package from: 123$")
}

func Test_FullLoggerNameGenerator_withString(t *testing.T) {
	givenString := "123"
	assert.ToBeEqual(t, givenString, FullLoggerNameGenerator(givenString))
	assert.ToBeEqual(t, givenString, FullLoggerNameGenerator(&givenString))
}

func Test_FullLoggerNameGenerator_regularCases(t *testing.T) {
	assert.ToBeEqual(t, "testing.T", FullLoggerNameGenerator(t))
	assert.ToBeEqual(t, "github.com/echocat/slf4g/names.someStruct", FullLoggerNameGenerator(&someStruct{}))
	assert.ToBeEqual(t, "github.com/echocat/slf4g/names.someStruct", FullLoggerNameGenerator(someStruct{}))
	assert.ToBeEqual(t, "github.com/echocat/slf4g/names.someStruct", FullLoggerNameGenerator((*someStruct)(nil)))
	assert.ToBeEqual(t, "github.com/echocat/slf4g/fields.empty", FullLoggerNameGenerator(fields.Empty()))
}

func Test_CurrentPackageLoggerNameGenerator_usesCustomizer(t *testing.T) {
	before := CurrentPackageLoggerNameCustomizer
	defer func() { CurrentPackageLoggerNameCustomizer = before }()

	called := false
	CurrentPackageLoggerNameCustomizer = func(framesToSkip int) string {
		assert.ToBeEqual(t, 124, framesToSkip)
		called = true
		return "xyz"
	}

	actual := CurrentPackageLoggerNameGenerator(123)
	assert.ToBeEqual(t, true, called)
	assert.ToBeEqual(t, "xyz", actual)
}

func Test_CurrentPackageLoggerNameGenerator(t *testing.T) {
	assert.ToBeEqual(t, "github.com/echocat/slf4g/names", forPackageSomethingFromInit)
	assert.ToBeEqual(t, "github.com/echocat/slf4g/names", (&someStruct{}).somethingFromAMethodInAStruct())
	assert.ToBeEqual(t, "github.com/echocat/slf4g/names", someStruct{}.somethingFromAMethodInAStruct())
	assert.ToBeEqual(t, "github.com/echocat/slf4g/names", CurrentPackageLoggerNameGenerator(0))
}

type someStruct struct{}

func (instance someStruct) somethingFromAMethodInAStruct() string {
	return CurrentPackageLoggerNameGenerator(0)
}
