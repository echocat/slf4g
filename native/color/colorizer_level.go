package color

import (
	"github.com/echocat/slf4g/level"
)

var DefaultLevelBasedColorizer LevelBasedColorizer = LevelColorizerMap{
	level.Trace: `[30;1m`,
	level.Debug: `[36;1m`,
	level.Info: `[34;1m`,
	level.Warn: `[33;1m`,
	level.Error: `[31;1m`,
	level.Fatal: `[35;1m`,
}

type LevelBasedColorizer interface {
	Colorize(level.Level, string) string
}

type LevelColorizerMap map[level.Level]string

func (l LevelColorizerMap) Colorize(level level.Level, what string) string {
	prefix := l[level]
	if prefix == "" {
		prefix = `[37;1m`
	}
	return prefix + what + `[0m`
}
