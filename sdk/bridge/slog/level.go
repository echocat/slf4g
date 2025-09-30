package sdk

import (
	"fmt"
	sdk "log/slog"

	"github.com/echocat/slf4g/level"
)

const (
	LevelTrace = sdk.Level(-8)
	LevelDebug = sdk.LevelDebug
	LevelInfo  = sdk.LevelInfo
	LevelWarn  = sdk.LevelWarn
	LevelError = sdk.LevelError
	LevelFatal = sdk.Level(12)
)

type LevelMapper interface {
	FromSdk(sdk.Level) (level.Level, error)
	ToSdk(level.Level) (sdk.Level, error)
}

var DefaultLevelMapper LevelMapper = NewLevelMapper()

func NewLevelMapper() LevelMapper {
	return &defaultLevelMapper{}
}

type defaultLevelMapper struct{}

func (instance *defaultLevelMapper) FromSdk(v sdk.Level) (level.Level, error) {
	switch v {
	case LevelTrace:
		return level.Trace, nil
	case LevelDebug:
		return level.Debug, nil
	case LevelInfo:
		return level.Info, nil
	case LevelWarn:
		return level.Warn, nil
	case LevelError:
		return level.Error, nil
	case LevelFatal:
		return level.Fatal, nil
	default:
		return 0, fmt.Errorf("unknown slog level: %d", v)
	}
}

func (instance *defaultLevelMapper) ToSdk(v level.Level) (sdk.Level, error) {
	switch v {
	case level.Trace:
		return LevelTrace, nil
	case level.Debug:
		return LevelDebug, nil
	case level.Info:
		return LevelInfo, nil
	case level.Warn:
		return LevelWarn, nil
	case level.Error:
		return LevelError, nil
	case level.Fatal:
		return LevelFatal, nil
	default:
		return 0, fmt.Errorf("unknown log level: %d", v)
	}
}

// NewLevelMapperFacade creates a facade of LevelMapper using the given provider.
func NewLevelMapperFacade(provider func() LevelMapper) LevelMapper {
	return levelMapperFacade(provider)
}

type levelMapperFacade func() LevelMapper

func (instance levelMapperFacade) FromSdk(v sdk.Level) (level.Level, error) {
	return instance.Unwrap().FromSdk(v)
}

func (instance levelMapperFacade) ToSdk(v level.Level) (sdk.Level, error) {
	return instance.Unwrap().ToSdk(v)
}

func (instance levelMapperFacade) Unwrap() LevelMapper {
	return instance()
}
