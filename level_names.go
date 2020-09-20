package log

import (
	"fmt"
	"strconv"
	"strings"
)

type LevelNames interface {
	FromOrdinal(uint16) (string, error)
	ToOrdinal(string) (uint16, error)
}

var DefaultLevelNames LevelNames = &defaultLevelNames{}

type defaultLevelNames struct{}

func (instance *defaultLevelNames) FromOrdinal(ordinal uint16) (string, error) {
	switch Level(ordinal) {
	case LevelTrace:
		return "TRACE", nil
	case LevelDebug:
		return "DEBUG", nil
	case LevelInfo:
		return "INFO", nil
	case LevelWarn:
		return "WARN", nil
	case LevelError:
		return "ERROR", nil
	case LevelFatal:
		return "FATAL", nil
	default:
		return fmt.Sprintf("%d", ordinal), nil
	}
}

func (instance *defaultLevelNames) ToOrdinal(name string) (uint16, error) {
	switch strings.ToUpper(name) {
	case "TRACE":
		return uint16(LevelTrace), nil
	case "DEBUG", "VERBOSE":
		return uint16(LevelDebug), nil
	case "INFO", "INFORMATION":
		return uint16(LevelInfo), nil
	case "WARN", "WARNING":
		return uint16(LevelWarn), nil
	case "ERROR", "ERR":
		return uint16(LevelError), nil
	case "FATAL":
		return uint16(LevelFatal), nil
	default:
		if result, err := strconv.ParseUint(name, 10, 16); err != nil {
			return 0, fmt.Errorf("%w: %s", ErrIllegalLevel, name)
		} else {
			return uint16(result), nil
		}
	}
}
