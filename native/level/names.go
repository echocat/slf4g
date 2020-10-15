package level

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/echocat/slf4g/level"
)

// ErrIllegalLevel represents that an illegal level.Level value/name was
// provided.
var ErrIllegalLevel = errors.New("illegal level")

// DefaultNames is the default instance of Names which should cover the most of
// the cases.
var DefaultNames = NewNames()

// Names is used to make readable names out of level.Level or the other way
// around.
type Names interface {
	// ToName converts a given level.Level to a human readable name. If this
	// level is unknown by this instance an error is returned. Most likely
	// ErrIllegalLevel.
	ToName(level.Level) (string, error)

	// ToLevel converts a given human readable name to a level.Level. If this
	// name is unknown by this instance an error is returned. Most likely
	// ErrIllegalLevel.
	ToLevel(string) (level.Level, error)
}

// NewNames creates a new default instance of a Names implementation.
func NewNames() Names {
	return &defaultNames{}
}

// NamesAware represents an object that is aware of Names.
type NamesAware interface {
	// GetLevelNames returns an instance of level.Names that support by
	// formatting levels in a human readable format.
	GetLevelNames() Names
}

// NewNamesFacade creates a facade of Names using the given provider.
func NewNamesFacade(provider func() Names) Names {
	return namesFacade(provider)
}

type namesFacade func() Names

func (instance namesFacade) ToName(lvl level.Level) (string, error) {
	return instance().ToName(lvl)
}

func (instance namesFacade) ToLevel(name string) (level.Level, error) {
	return instance().ToLevel(name)
}

type defaultNames struct{}

func (instance *defaultNames) ToName(lvl level.Level) (string, error) {
	switch lvl {
	case level.Trace:
		return "TRACE", nil
	case level.Debug:
		return "DEBUG", nil
	case level.Info:
		return "INFO", nil
	case level.Warn:
		return "WARN", nil
	case level.Error:
		return "ERROR", nil
	case level.Fatal:
		return "FATAL", nil
	default:
		return fmt.Sprintf("%d", lvl), nil
	}
}

func (instance *defaultNames) ToLevel(name string) (level.Level, error) {
	switch strings.ToUpper(name) {
	case "TRACE":
		return level.Trace, nil
	case "DEBUG", "VERBOSE":
		return level.Debug, nil
	case "INFO", "INFORMATION":
		return level.Info, nil
	case "WARN", "WARNING":
		return level.Warn, nil
	case "ERROR", "ERR":
		return level.Error, nil
	case "FATAL":
		return level.Fatal, nil
	default:
		if result, err := strconv.ParseUint(name, 10, 16); err != nil {
			return 0, fmt.Errorf("%w: %s", ErrIllegalLevel, name)
		} else {
			return level.Level(result), nil
		}
	}
}
