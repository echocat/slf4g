package level

import (
	"errors"
	"fmt"
	"github.com/echocat/slf4g"
	"strconv"
	"strings"
)

var (
	ErrIllegalLevel               = errors.New("illegal level")
	DefaultLevelNames       Names = &defaultNames{}
	DefaultLevelNamesFacade       = NewNamesFacade(func() Names {
		return DefaultLevelNames
	})
)

type Names interface {
	FromOrdinal(uint16) (string, error)
	ToOrdinal(string) (uint16, error)
}

func NewNamesFacade(provider func() Names) Names {
	return namesFacade(provider)
}

type defaultNames struct{}

func (instance *defaultNames) FromOrdinal(ordinal uint16) (string, error) {
	switch log.Level(ordinal) {
	case log.LevelTrace:
		return "TRACE", nil
	case log.LevelDebug:
		return "DEBUG", nil
	case log.LevelInfo:
		return "INFO", nil
	case log.LevelWarn:
		return "WARN", nil
	case log.LevelError:
		return "ERROR", nil
	case log.LevelFatal:
		return "FATAL", nil
	default:
		return fmt.Sprintf("%d", ordinal), nil
	}
}

func (instance *defaultNames) ToOrdinal(name string) (uint16, error) {
	switch strings.ToUpper(name) {
	case "TRACE":
		return uint16(log.LevelTrace), nil
	case "DEBUG", "VERBOSE":
		return uint16(log.LevelDebug), nil
	case "INFO", "INFORMATION":
		return uint16(log.LevelInfo), nil
	case "WARN", "WARNING":
		return uint16(log.LevelWarn), nil
	case "ERROR", "ERR":
		return uint16(log.LevelError), nil
	case "FATAL":
		return uint16(log.LevelFatal), nil
	default:
		if result, err := strconv.ParseUint(name, 10, 16); err != nil {
			return 0, fmt.Errorf("%w: %s", ErrIllegalLevel, name)
		} else {
			return uint16(result), nil
		}
	}
}

type NamesAware interface {
	GetLevelNames() Names
}

type namesFacade func() Names

func (instance namesFacade) FromOrdinal(ordinal uint16) (string, error) {
	return instance().FromOrdinal(ordinal)
}

func (instance namesFacade) ToOrdinal(name string) (uint16, error) {
	return instance().ToOrdinal(name)
}
