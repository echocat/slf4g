package functions

import (
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/native/color"
	"github.com/echocat/slf4g/native/hints"
	nlevel "github.com/echocat/slf4g/native/level"
)

func ColorizeByLevel(l level.Level, h hints.Hints, what string) string {
	if !ShouldColorize(h) {
		return what
	}
	colorizer := LevelColorizer(h)
	return colorizer.ColorizeByLevel(l, what)
}

func Colorize(colorCode string, h hints.Hints, what string) string {
	if !ShouldColorize(h) {
		return what
	}
	return `[` + colorCode + `m` + what + `[0m`
}

func ShouldColorize(h hints.Hints) bool {
	supported := color.SupportedNone
	mode := color.ModeAuto
	if v, ok := h.(hints.ColorsSupport); ok {
		supported = v.IsColorSupported()
	}
	if v, ok := h.(hints.ColorMode); ok {
		mode = v.ColorMode()
	}
	return mode.ShouldColorize(supported)
}

func LevelColorizer(h hints.Hints) nlevel.Colorizer {
	if vh, ok := h.(hints.LevelColorizer); ok {
		if v := vh.LevelColorizer(); v != nil {
			return v
		}
	}
	if v := nlevel.DefaultColorizer; v != nil {
		return v
	}
	return nlevel.NoopColorizer()
}
