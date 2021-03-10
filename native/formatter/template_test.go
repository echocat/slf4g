package formatter

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"text/template"

	"github.com/echocat/slf4g/native/hints"

	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/testing/recording"

	"github.com/echocat/slf4g/native/formatter/functions"

	nlevel "github.com/echocat/slf4g/native/level"

	"github.com/echocat/slf4g/native/color"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_Template_MustNewTemplate(t *testing.T) {
	actual := MustNewTemplate("foo{{ . | ensureWidth 5 true }}bar")

	assert.ToBeNotNil(t, actual)
	assert.ToBeNotNil(t, actual.template)
	assert.ToBeNil(t, actual.LevelColorizer)
	assert.ToBeEqual(t, color.ModeAuto, actual.ColorMode)

	buf := new(bytes.Buffer)
	actualErr := actual.template.Execute(buf, "hello, world")

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, "foohellobar", buf.String())
}

func Test_Template_MustNewTemplate_customize(t *testing.T) {
	actual := MustNewTemplate("foo{{ . | ensureWidth 5 true }}bar", func(template *Template) {
		template.LevelColorizer = nlevel.DefaultColorizer
		template.ColorMode = color.ModeAlways
	})

	assert.ToBeNotNil(t, actual)
	assert.ToBeNotNil(t, actual.template)
	assert.ToBeEqual(t, nlevel.DefaultColorizer, actual.LevelColorizer)
	assert.ToBeEqual(t, color.ModeAlways, actual.ColorMode)

	buf := new(bytes.Buffer)
	actualErr := actual.template.Execute(buf, "hello, world")

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, "foohellobar", buf.String())
}

func Test_Template_MustNewTemplate_failing(t *testing.T) {
	assert.Execution(t, func() {
		_ = MustNewTemplate("foo{{ . | foo }}bar")
	}).WillPanicWith("function \"foo\" not defined")
}

func Test_Template_MustNewTemplateWithFuncMap(t *testing.T) {
	actual := MustNewTemplateWithFuncMap("foo{{ myFunc }}bar", map[string]interface{}{
		"myFunc": func() string { return "some" },
	})

	assert.ToBeNotNil(t, actual)
	assert.ToBeNotNil(t, actual.template)
	assert.ToBeNil(t, actual.LevelColorizer)
	assert.ToBeEqual(t, color.ModeAuto, actual.ColorMode)

	buf := new(bytes.Buffer)
	actualErr := actual.template.Execute(buf, "")

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, "foosomebar", buf.String())
}

func Test_Template_MustNewTemplateWithFuncMap_customize(t *testing.T) {
	actual := MustNewTemplateWithFuncMap("foobar", nil, func(template *Template) {
		template.LevelColorizer = nlevel.DefaultColorizer
		template.ColorMode = color.ModeAlways
	})

	assert.ToBeNotNil(t, actual)
	assert.ToBeNotNil(t, actual.template)
	assert.ToBeEqual(t, nlevel.DefaultColorizer, actual.LevelColorizer)
	assert.ToBeEqual(t, color.ModeAlways, actual.ColorMode)
}

func Test_Template_MustNewTemplateWithFuncMap_failing(t *testing.T) {
	assert.Execution(t, func() {
		_ = MustNewTemplateWithFuncMap("foo{{ other }}bar", map[string]interface{}{
			"myFunc": func() string { return "some" },
		})
	}).WillPanicWith("function \"other\" not defined")
}

func Test_Template_MustNewTemplateByFactory(t *testing.T) {
	givenTemplate := template.New("expected")

	actual := MustNewTemplateByFactory(func(rootFuncMap template.FuncMap) (*template.Template, error) {
		assert.ToBeNotNil(t, rootFuncMap)
		assert.ToBeSame(t, functions.ColorizeByLevel, rootFuncMap["colorizeByLevel"])
		return givenTemplate, nil
	})

	assert.ToBeNotNil(t, actual)
	assert.ToBeSame(t, givenTemplate, actual.template)
	assert.ToBeNil(t, actual.LevelColorizer)
	assert.ToBeEqual(t, color.ModeAuto, actual.ColorMode)
}

func Test_Template_MustNewTemplateByFactory_customize(t *testing.T) {
	givenTemplate := template.New("expected")

	actual := MustNewTemplateByFactory(func(rootFuncMap template.FuncMap) (*template.Template, error) {
		assert.ToBeNotNil(t, rootFuncMap)
		assert.ToBeSame(t, functions.ColorizeByLevel, rootFuncMap["colorizeByLevel"])
		return givenTemplate, nil
	}, func(template *Template) {
		template.LevelColorizer = nlevel.DefaultColorizer
		template.ColorMode = color.ModeAlways
	})

	assert.ToBeNotNil(t, actual)
	assert.ToBeSame(t, givenTemplate, actual.template)
	assert.ToBeEqual(t, nlevel.DefaultColorizer, actual.LevelColorizer)
	assert.ToBeEqual(t, color.ModeAlways, actual.ColorMode)
}

func Test_Template_MustNewTemplateByFactory_failing(t *testing.T) {
	assert.Execution(t, func() {
		_ = MustNewTemplateByFactory(func(template.FuncMap) (*template.Template, error) {
			return nil, fmt.Errorf("expected")
		})
	}).WillPanicWith("^expected$")
}

func Test_Template_toFuncMap(t *testing.T) {
	actual := (&Template{}).toFuncMap()

	assert.ToBeNotNil(t, actual)
	assert.ToBeEqual(t, 6, len(actual))
	assert.ToBeSame(t, functions.ColorizeByLevel, actual["colorizeByLevel"])
	assert.ToBeSame(t, functions.Colorize, actual["colorize"])
	assert.ToBeSame(t, functions.ShouldColorize, actual["shouldColorize"])
	assert.ToBeSame(t, functions.LevelColorizer, actual["levelColorizer"])
	assert.ToBeSame(t, functions.IndentMultiline, actual["indentMultiline"])
	assert.ToBeSame(t, functions.EnsureWidth, actual["ensureWidth"])
}

func Test_Template_Format(t *testing.T) {
	givenProvider := recording.NewProvider()
	givenLogger := givenProvider.GetLogger("aLogger")
	givenEvent := givenLogger.NewEvent(level.Warn, nil).
		With("logger", givenLogger).
		With("message", "aMessage").
		With("timestamp", mustParseTime("2021-01-02T13:14:15.1234")).
		With("error", errors.New("anError")).
		With("foo", "bar")

	instance := MustNewTemplate("{{.Timestamp.Format `2006-01-02T15:04:05`}},{{.Logger}},{{.LevelName}},{{.Message | ensureWidth 5 true}},{{.Error}},{{.Fields}}")

	actual, actualErr := instance.Format(givenEvent, givenProvider, nil)

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, "2021-01-02T13:14:15,aLogger,WARN,aMess,anError,map[foo:bar]", string(actual))
}

func Test_Template_Format_failing(t *testing.T) {
	givenProvider := recording.NewProvider()
	givenLogger := givenProvider.GetLogger("aLogger")
	givenEvent := givenLogger.NewEvent(level.Warn, nil).
		With("logger", givenLogger).
		With("message", "aMessage").
		With("timestamp", mustParseTime("2021-01-02T13:14:15.1234")).
		With("error", errors.New("anError")).
		With("foo", "bar")

	instance := MustNewTemplate("{{.bla}}")

	actual, actualErr := instance.Format(givenEvent, givenProvider, nil)

	assert.ToBeNotNil(t, actualErr)
	assert.ToBeEqual(t, "", string(actual))
}

func Test_TemplateRenderingContext_Hints(t *testing.T) {
	givenHints := &mockColorizingHints{}
	instance := newTestTemplateRenderingContext(givenHints)

	actual := instance.Hints()

	assert.ToBeNotNil(t, actual)
	assert.ToBeSame(t, givenHints, actual.(templateHintsCombined).Hints)
	assert.ToBeSame(t, instance.Template, actual.(templateHintsCombined).Template)

	assert.ToBeSame(t, instance.Template, actual.(templateHintsCombined).Template)
}

func Test_TemplateRenderingContext_Hints_explicitValues(t *testing.T) {
	givenHints := &mockColorizingHints{}
	instance := newTestTemplateRenderingContext(givenHints)

	actual := instance.Hints()

	assert.ToBeNotNil(t, actual)
	assert.ToBeEqual(t, givenHints.ColorMode(), actual.(templateHintsCombined).ColorMode())
	assert.ToBeEqual(t, givenHints.IsColorSupported(), actual.(templateHintsCombined).IsColorSupported())
	assert.ToBeEqual(t, givenHints.LevelColorizer(), actual.(templateHintsCombined).LevelColorizer())
}

func Test_TemplateRenderingContext_Hints_fallbackValues(t *testing.T) {
	instance := newTestTemplateRenderingContext(nil)
	instance.Template.ColorMode = 66
	instance.Template.LevelColorizer = nlevel.ColorizerMap{}

	actual := instance.Hints()

	assert.ToBeNotNil(t, actual)
	assert.ToBeEqual(t, instance.Template.ColorMode, actual.(templateHintsCombined).ColorMode())
	assert.ToBeEqual(t, color.SupportedNone, actual.(templateHintsCombined).IsColorSupported())
	assert.ToBeEqual(t, instance.Template.LevelColorizer, actual.(templateHintsCombined).LevelColorizer())
}

func Test_TemplateRenderingContext_LevelNames_ofProvider(t *testing.T) {
	givenNames := nlevel.NewNames()
	instance := newTestTemplateRenderingContext(nil)
	instance.Provider = &mockProviderWithLevelNames{
		Provider: recording.NewProvider(),
		Names:    givenNames,
	}

	actual := instance.LevelNames()

	assert.ToBeSame(t, givenNames, actual)
}

func Test_TemplateRenderingContext_LevelNames_ofProviderButNil(t *testing.T) {
	instance := newTestTemplateRenderingContext(nil)
	instance.Provider = &mockProviderWithLevelNames{
		Provider: recording.NewProvider(),
		Names:    nil,
	}

	actual := instance.LevelNames()

	assert.ToBeSame(t, nlevel.DefaultNames, actual)
}

func Test_TemplateRenderingContext_LevelNames_default(t *testing.T) {
	instance := newTestTemplateRenderingContext(nil)

	actual := instance.LevelNames()

	assert.ToBeSame(t, nlevel.DefaultNames, actual)
}

func Test_TemplateRenderingContext_LevelNames_fallback(t *testing.T) {
	beforeNames := nlevel.DefaultNames
	defer func() { nlevel.DefaultNames = beforeNames }()
	nlevel.DefaultNames = nil

	instance := newTestTemplateRenderingContext(nil)

	actual := instance.LevelNames()

	assert.ToBeEqual(t, nlevel.NewNames(), actual)
}

func Test_TemplateRenderingContext_Level(t *testing.T) {
	instance := newTestTemplateRenderingContext(nil)

	actual := instance.Level()

	assert.ToBeEqual(t, level.Warn, actual)
}

func Test_TemplateRenderingContext_LevelName(t *testing.T) {
	instance := newTestTemplateRenderingContext(nil)

	actual, actualErr := instance.LevelName()

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, "WARN", actual)
}

func Test_TemplateRenderingContext_Message(t *testing.T) {
	instance := newTestTemplateRenderingContext(nil)

	actual := instance.Message()

	assert.ToBeNotNil(t, actual)
	assert.ToBeEqual(t, "aMessage", *actual)
}

func Test_TemplateRenderingContext_Error(t *testing.T) {
	instance := newTestTemplateRenderingContext(nil)

	actual := instance.Error()

	assert.ToBeNotNil(t, actual)
	assert.ToBeEqual(t, "anError", actual.Error())
}

func Test_TemplateRenderingContext_Timestamp(t *testing.T) {
	instance := newTestTemplateRenderingContext(nil)

	actual := instance.Timestamp()

	assert.ToBeNotNil(t, actual)
	assert.ToBeEqual(t, mustParseTime("2021-01-02T13:14:15.1234"), *actual)
}

func Test_TemplateRenderingContext_Logger(t *testing.T) {
	instance := newTestTemplateRenderingContext(nil)

	actual := instance.Logger()

	assert.ToBeNotNil(t, actual)
	assert.ToBeEqual(t, "aLogger", *actual)
}

func Test_TemplateRenderingContext_Fields(t *testing.T) {
	instance := newTestTemplateRenderingContext(nil)

	actual, actualErr := instance.Fields()

	assert.ToBeNoError(t, actualErr)
	assert.ToBeNotNil(t, actual)
	assert.ToBeEqual(t, map[string]interface{}{"foo": "bar"}, actual)
}

func newTestTemplateRenderingContext(h hints.Hints) TemplateRenderingContext {
	givenProvider := recording.NewProvider()
	givenLogger := givenProvider.GetLogger("aLogger")
	givenEvent := givenLogger.NewEvent(level.Warn, nil).
		With("logger", givenLogger).
		With("message", "aMessage").
		With("timestamp", mustParseTime("2021-01-02T13:14:15.1234")).
		With("error", errors.New("anError")).
		With("foo", "bar")
	givenTemplate := MustNewTemplate("")

	return givenTemplate.contextFor(givenEvent, givenProvider, h)
}
