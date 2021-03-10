package formatter

import (
	"fmt"
	"time"

	"github.com/echocat/slf4g/testing/recording"

	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/native/color"
	nlevel "github.com/echocat/slf4g/native/level"
)

type mockColorsSupport color.Supported

func (instance mockColorsSupport) IsColorSupported() color.Supported {
	return color.Supported(instance)
}

type mockColorizingHints struct{}

func (instance mockColorizingHints) IsColorSupported() color.Supported {
	return color.SupportedAssumed
}

func (instance mockColorizingHints) LevelColorizer() nlevel.Colorizer {
	return instance
}

func (instance mockColorizingHints) ColorMode() color.Mode {
	return color.ModeAlways
}

func (instance mockColorizingHints) ColorizeByLevel(l level.Level, input string) string {
	return fmt.Sprintf("%d(%s)", l/1000, input)
}

func mustParseTime(in string) time.Time {
	v, err := time.Parse("2006-01-02T15:04:05.9999", in)
	if err != nil {
		panic(err)
	}
	return v
}

type mockProviderWithLevelNames struct {
	*recording.Provider
	Names nlevel.Names
}

func (instance mockProviderWithLevelNames) GetLevelNames() nlevel.Names {
	return instance.Names
}
