package level

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/echocat/slf4g/level"
)

// DefaultNames is the default instance of Names which should cover the most of
// the cases.
var DefaultNames = NewNames()

// NewNames creates a new default instance of a Names implementation.
func NewNames() level.Names {
	return &defaultNames{}
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
			return 0, fmt.Errorf("%w: %s", level.ErrIllegalLevel, name)
		} else {
			return level.Level(result), nil
		}
	}
}

// ErrIllegalLevel represents that an illegal level.Level value/name was
// provided.
//
// Deprecated: use level.ErrIllegalLevel instead.
var ErrIllegalLevel = level.ErrIllegalLevel

// Names is used to make readable names out of level.Level or the other way
// around.
//
// Deprecated: use level.Names instead.
type Names level.Names

// NamesAware represents an object that is aware of Names.
//
// Deprecated: use level.NamesAware instead.
type NamesAware level.NamesAware

// NewNamesFacade creates a facade of Names using the given provider.
//
// Deprecated: use level.NewNamesFacade instead.
func NewNamesFacade(provider func() Names) Names {
	return level.NewNamesFacade(func() level.Names {
		return provider()
	})
}
