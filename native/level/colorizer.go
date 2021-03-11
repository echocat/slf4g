package level

import (
	"github.com/echocat/slf4g/level"
)

// DefaultColorizer is the default instance of Colorizer which should cover the
// most of the cases.
var DefaultColorizer Colorizer = ColorizerMap{
	level.Trace: `[30;1m`,
	level.Debug: `[36;1m`,
	level.Info: `[34;1m`,
	level.Warn: `[33;1m`,
	level.Error: `[31;1m`,
	level.Fatal: `[35;1m`,
}

// Colorizer is colorizing inputs for given levels by ANSI escape codes.
// See: https://en.wikipedia.org/wiki/ANSI_escape_code
type Colorizer interface {

	// ColorizeByLevel is colorizing the given input for the given level.Level.
	ColorizeByLevel(lvl level.Level, input string) string
}

// ColorizerMap is an implementation of Colorizer which simply holds for
// configured level.Level an ANSI escape code for colorizing. If there is no
// level.Level configured it defaults to a simple grey.
type ColorizerMap map[level.Level]string

// ColorizeByLevel implements Colorizer.ColorizeByLevel()
func (l ColorizerMap) ColorizeByLevel(lvl level.Level, what string) string {
	prefix := l[lvl]
	if prefix == "" {
		prefix = `[37;1m`
	}
	return prefix + what + `[0m`
}

// NewColorizerFacade creates a facade of Colorizer using the given provider.
func NewColorizerFacade(provider func() Colorizer) Colorizer {
	return colorizerFacade(provider)
}

type colorizerFacade func() Colorizer

func (instance colorizerFacade) ColorizeByLevel(lvl level.Level, input string) string {
	return instance.Unwrap().ColorizeByLevel(lvl, input)
}

func (instance colorizerFacade) Unwrap() Colorizer {
	return instance()
}

// NoopColorizer provides a noop implementation of Colorizer.
func NoopColorizer() Colorizer {
	return noopColorizerV
}

var noopColorizerV = &noopColorizer{}

type noopColorizer struct{}

func (instance noopColorizer) ColorizeByLevel(_ level.Level, what string) string {
	return what
}
