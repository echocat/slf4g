package hints

import "github.com/echocat/slf4g/native/color"

type ColorsSupport interface {
	Hints
	GetColorSupport() color.Supported
}
