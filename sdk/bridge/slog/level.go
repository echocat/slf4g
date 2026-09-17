//go:build go1.21

package sdk

import (
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
	switch {
	case v < LevelDebug:
		return level.Trace, nil
	case v < LevelInfo:
		return level.Debug, nil
	case v < LevelWarn:
		return level.Info, nil
	case v < LevelError:
		return level.Warn, nil
	case v < LevelFatal:
		return level.Error, nil
	default:
		return level.Fatal, nil
	}
}

func (instance *defaultLevelMapper) ToSdk(v level.Level) (sdk.Level, error) {
	switch {
	case v < level.Debug:
		return LevelTrace, nil
	case v < level.Info:
		return LevelDebug, nil
	case v < level.Warn:
		return LevelInfo, nil
	case v < level.Error:
		return LevelWarn, nil
	case v < level.Fatal:
		return LevelError, nil
	default:
		return LevelFatal, nil
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
