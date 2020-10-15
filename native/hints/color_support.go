package hints

import "github.com/echocat/slf4g/native/color"

// ColorsSupport are Hints that provide the information if color is supported.
type ColorsSupport interface {
	Hints

	// IsColorSupported returns if color is supported and how.
	IsColorSupported() color.Supported
}
