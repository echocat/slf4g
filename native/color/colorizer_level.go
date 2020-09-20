package color

import (
	log "github.com/echocat/slf4g"
)

var DefaultLevelBasedColorizer LevelBasedColorizer = LevelColorizerMap{
	log.LevelTrace: `[30;1m`,
	log.LevelDebug: `[36;1m`,
	log.LevelInfo: `[34;1m`,
	log.LevelWarn: `[33;1m`,
	log.LevelError: `[31;1m`,
	log.LevelFatal: `[35;1m`,
}

type LevelBasedColorizer interface {
	Colorize(log.Level, string) string
}

type LevelColorizerMap map[log.Level]string

func (l LevelColorizerMap) Colorize(level log.Level, what string) string {
	prefix := l[level]
	if prefix == "" {
		prefix = `[37;1m`
	}
	return prefix + what + `[0m`
}
