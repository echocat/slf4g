package formatter

import (
	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/level"
	nlevel "github.com/echocat/slf4g/native/level"
)

// DefaultLevel is the default instance of Level which should cover the most of
// the cases.
var DefaultLevel Level = NewNamesBasedLevel(nlevel.NewNamesFacade(func() nlevel.Names {
	return nlevel.DefaultNames
}))

// Level is used to format a given level.Level.
type Level interface {
	// FormatLevel formats the given level.Level.
	FormatLevel(in level.Level, using log.Provider) (interface{}, error)
}

// NewNamesBasedLevel creates a new instance of Level which uses given nlevel.Names to
// resolve the name of a given log.Level and format it with it.
func NewNamesBasedLevel(names nlevel.Names) Level {
	return LevelFunc(func(in level.Level, using log.Provider) (interface{}, error) {
		result, err := names.ToName(in)
		return result, err
	})
}

// NewNamesBasedLevel creates a new instance of Level which formats the given
// level.Level by its ordinal.
func NewOrdinalBasedLevel() Level {
	return LevelFunc(func(in level.Level, using log.Provider) (interface{}, error) {
		return uint16(in), nil
	})
}

// LevelFunc is wrapping the given function into a Level.
type LevelFunc func(in level.Level, using log.Provider) (interface{}, error)

// FormatLevel implements Level.FormatLevel()
func (instance LevelFunc) FormatLevel(in level.Level, using log.Provider) (interface{}, error) {
	return instance(in, using)
}

// NewFacade creates a new facade instance of Formatter using the given
// provider.
func NewLevelFacade(provider func() Level) Level {
	return levelFacade(provider)
}

type levelFacade func() Level

func (instance levelFacade) FormatLevel(in level.Level, using log.Provider) (interface{}, error) {
	return instance().FormatLevel(in, using)
}

// NoopLevel provides a noop implementation of Level.
func NoopLevel() Level {
	return noopLevelV
}

var noopLevelV = LevelFunc(func(level.Level, log.Provider) (interface{}, error) {
	return []byte{}, nil
})
