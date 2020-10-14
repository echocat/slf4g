package level

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/echocat/slf4g/level"
)

var (
	ErrIllegalLevel = errors.New("illegal level")
	DefaultNames    = NewNames()
)

type Names interface {
	FromOrdinal(uint16) (string, error)
	ToOrdinal(string) (uint16, error)
}

func NewLevelNamesFacade(provider func() Names) Names {
	return namesFacade(provider)
}

func NewNames() Names {
	return &defaultNames{}
}

type defaultNames struct{}

func (instance *defaultNames) FromOrdinal(ordinal uint16) (string, error) {
	switch level.Level(ordinal) {
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
		return fmt.Sprintf("%d", ordinal), nil
	}
}

func (instance *defaultNames) ToOrdinal(name string) (uint16, error) {
	switch strings.ToUpper(name) {
	case "TRACE":
		return uint16(level.Trace), nil
	case "DEBUG", "VERBOSE":
		return uint16(level.Debug), nil
	case "INFO", "INFORMATION":
		return uint16(level.Info), nil
	case "WARN", "WARNING":
		return uint16(level.Warn), nil
	case "ERROR", "ERR":
		return uint16(level.Error), nil
	case "FATAL":
		return uint16(level.Fatal), nil
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
