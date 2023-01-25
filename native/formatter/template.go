package formatter

import (
	"bytes"
	"text/template"
	"time"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/native/color"
	"github.com/echocat/slf4g/native/formatter/functions"
	"github.com/echocat/slf4g/native/hints"
	nlevel "github.com/echocat/slf4g/native/level"
)

// Template is an implementation of Formatter which formats given log entries in a
// human readable format, formatted by a template it was created with.
type Template struct {
	// ColorMode defines when the output should be colorized. If not configured
	// color.ModeAuto will be used by default.
	ColorMode color.Mode

	// LevelColorizer is used to colorize output based on the level.Level of an
	// log.Event to be logged. If not set nlevel.DefaultColorizer will be used.
	LevelColorizer nlevel.Colorizer

	template *template.Template
}

// NewTemplate creates a new instance of Template which is ready to use for the
// given plain template.
func NewTemplate(plain string, customizer ...func(*Template)) (*Template, error) {
	return NewTemplateWithFuncMap(plain, nil, customizer...)
}

// MustNewTemplate is same as NewTemplate but will panic in case of errors.
func MustNewTemplate(plain string, customizer ...func(*Template)) *Template {
	result, err := NewTemplate(plain, customizer...)
	if err != nil {
		panic(err)
	}
	return result
}

// NewTemplateWithFuncMap creates a new instance of Template which is ready to use
// for the given plain template and funcMap.
func NewTemplateWithFuncMap(plain string, funcMap template.FuncMap, customizer ...func(*Template)) (*Template, error) {
	return NewTemplateByFactory(func(topFuncMap template.FuncMap) (*template.Template, error) {
		tmpl := template.New("logFormat")
		if topFuncMap != nil {
			tmpl = tmpl.Funcs(topFuncMap)
		}
		if funcMap != nil {
			tmpl = tmpl.Funcs(funcMap)
		}
		return tmpl.Parse(plain)
	}, customizer...)
}

// MustNewTemplateWithFuncMap is same as NewTemplateWithFuncMap but will panic in case of errors.
func MustNewTemplateWithFuncMap(plain string, funcMap template.FuncMap, customizer ...func(*Template)) *Template {
	result, err := NewTemplateWithFuncMap(plain, funcMap, customizer...)
	if err != nil {
		panic(err)
	}
	return result
}

// NewTemplateByFactory creates a new instance of Template which is ready to use
// using the given factory.
func NewTemplateByFactory(factory TemplateFactory, customizer ...func(*Template)) (*Template, error) {
	result := &Template{}
	tmpl, err := factory(result.toFuncMap())
	if err != nil {
		return nil, err
	}

	result.template = tmpl

	for _, c := range customizer {
		c(result)
	}

	return result, nil
}

// MustNewTemplateByFactory is same as NewTemplateByFactory but will panic in case of errors.
func MustNewTemplateByFactory(factory TemplateFactory, customizer ...func(*Template)) *Template {
	result, err := NewTemplateByFactory(factory, customizer...)
	if err != nil {
		panic(err)
	}
	return result
}

// Format implements Formatter.Format()
func (instance *Template) Format(event log.Event, using log.Provider, h hints.Hints) ([]byte, error) {
	ctx := instance.contextFor(event, using, h)

	buf := new(bytes.Buffer)
	if err := instance.template.Execute(buf, ctx); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (instance *Template) contextFor(event log.Event, using log.Provider, h hints.Hints) TemplateRenderingContext {
	return TemplateRenderingContext{
		Event:    event,
		Provider: using,
		Template: instance,
		hints:    h,
	}
}

func (instance *Template) toFuncMap() template.FuncMap {
	return map[string]interface{}{
		"colorizeByLevel": functions.ColorizeByLevel,
		"colorize":        functions.Colorize,
		"shouldColorize":  functions.ShouldColorize,
		"levelColorizer":  functions.LevelColorizer,

		"indentMultiline": functions.IndentMultiline,

		"ensureWidth": functions.EnsureWidth,
	}
}

// TemplateFactory will create for the given rootFuncMap a new instance
// of template.Template.
type TemplateFactory func(rootFuncMap template.FuncMap) (*template.Template, error)

// TemplateRenderingContext is used by Template.Format as a context object while rendering
// the log.Event.
type TemplateRenderingContext struct {
	// log.Event is the foundation of this context.
	log.Event

	// Provider contains the corresponding log.Provider of the current log.Event.
	Provider log.Provider

	// Template is the actual instance of Template which which is executing the
	// rendering of the log.Event.
	Template *Template

	hints hints.Hints
}

// Hints provides the hints.Hints the current log.Event should be rendered with.
func (instance TemplateRenderingContext) Hints() hints.Hints {
	return templateHintsCombined{
		Hints:    instance.hints,
		Template: instance.Template,
	}
}

// LevelNames is a convenience method to easy return the current nlevel.Names
// of the corresponding log.Provider.
func (instance TemplateRenderingContext) LevelNames() nlevel.Names {
	if v, ok := instance.Provider.(level.NamesAware); ok {
		if names := v.GetLevelNames(); names != nil {
			return names
		}
	}
	if v := nlevel.DefaultNames; v != nil {
		return v
	}
	return nlevel.NewNames()
}

// Level is a convenience method to easy return the current level.Level of the
// corresponding log.Event.
func (instance TemplateRenderingContext) Level() level.Level {
	return instance.GetLevel()
}

// LevelName is a convenience method to easy return the current level.Level of the
// corresponding log.Event.
func (instance TemplateRenderingContext) LevelName() (string, error) {
	return instance.LevelNames().ToName(instance.Level())
}

// Message is a convenience method to easy return the current message of the
// corresponding log.Event, can be nil.
func (instance TemplateRenderingContext) Message() *string {
	return log.GetMessageOf(instance.Event, instance.Provider)
}

// Error is a convenience method to easy return the current error of the
// corresponding log.Event, can be nil.
func (instance TemplateRenderingContext) Error() error {
	return log.GetErrorOf(instance.Event, instance.Provider)
}

// Timestamp is a convenience method to easy return the current timestamp of the
// corresponding log.Event, can be nil.
func (instance TemplateRenderingContext) Timestamp() *time.Time {
	return log.GetTimestampOf(instance.Event, instance.Provider)
}

// Logger is a convenience method to easy return the current logger name of the
// corresponding log.Event, can be nil.
func (instance TemplateRenderingContext) Logger() *string {
	return log.GetLoggerOf(instance.Event, instance.Provider)
}

// FieldKeysSpec is a convenience method to easy return the current fields.KeysSpec
// of the corresponding log.Event.
func (instance TemplateRenderingContext) FieldKeysSpec() fields.KeysSpec {
	return instance.Provider.GetFieldKeysSpec()
}

// Fields returns all fields of the log.Event (except the default ones).
func (instance TemplateRenderingContext) Fields() (result map[string]interface{}, err error) {
	result = make(map[string]interface{})
	err = instance.ForEach(func(key string, value interface{}) error {
		keysSpec := instance.FieldKeysSpec()
		switch key {
		case keysSpec.GetLogger(), keysSpec.GetMessage(), keysSpec.GetTimestamp(), keysSpec.GetError():
			return nil
		default:
			result[key] = value
			return nil
		}
	})
	return
}

type templateHintsCombined struct {
	hints.Hints
	*Template
}

func (instance templateHintsCombined) ColorMode() color.Mode {
	if v, ok := instance.Hints.(hints.ColorMode); ok {
		return v.ColorMode()
	}
	return instance.Template.ColorMode
}

func (instance templateHintsCombined) LevelColorizer() nlevel.Colorizer {
	if v, ok := instance.Hints.(hints.LevelColorizer); ok {
		return v.LevelColorizer()
	}
	return instance.Template.LevelColorizer
}

func (instance templateHintsCombined) IsColorSupported() color.Supported {
	if v, ok := instance.Hints.(hints.ColorsSupport); ok {
		return v.IsColorSupported()
	}
	return color.SupportedNone
}
