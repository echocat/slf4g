package level

import (
	"github.com/echocat/slf4g/level"
)

var DefaultColorizer Colorizer = ColorizerMap{
	level.Trace: `[30;1m`,
	level.Debug: `[36;1m`,
	level.Info: `[34;1m`,
	level.Warn: `[33;1m`,
	level.Error: `[31;1m`,
	level.Fatal: `[35;1m`,
}

type Colorizer interface {
	ColorizeByLevel(level.Level, string) string
}

type ColorizerMap map[level.Level]string

func (l ColorizerMap) ColorizeByLevel(level level.Level, what string) string {
	prefix := l[level]
	if prefix == "" {
		prefix = `[37;1m`
	}
	return prefix + what + `[0m`
}

func NoopColorizer() Colorizer {
	return noopColorizerV
}

var noopColorizerV = &noopColorizer{}

type noopColorizer struct{}

func (instance noopColorizer) ColorizeByLevel(_ level.Level, what string) string {
	return what
}
