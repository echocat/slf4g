package hints

import (
	"github.com/echocat/slf4g/native/color"
	"github.com/echocat/slf4g/native/level"
)

// ColorsSupport are Hints that provide the information if color is supported.
type ColorsSupport interface {
	Hints

	// IsColorSupported returns if color is supported and how.
	IsColorSupported() color.Supported
}

// ColorMode are Hints that provide the information how something should be colorized.
type ColorMode interface {
	Hints

	// ColorMode returns the desired color.Mode.
	ColorMode() color.Mode
}

// LevelColorizer are Hints that provides an assigned level.Colorizer.
type LevelColorizer interface {
	Hints

	// ColorMode returns the desired level.Colorizer.
	LevelColorizer() level.Colorizer
}
