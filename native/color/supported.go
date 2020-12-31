package color

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrIllegalSupport will be returned in situations where illegal values or
// representations if a Supported are provided.
var ErrIllegalSupport = errors.New("illegal color-support")

// Supported expresses if color is supported or not in the current context.
type Supported uint8

const (
	// SupportedNone clearly says that color is not supported and cannot be
	// used.
	SupportedNone Supported = 0

	// SupportedNative expresses that color is natively supported and most
	// likely guaranteed to be working.
	SupportedNative Supported = 1

	// SupportedAssumed expresses that color is assumed to work, but there is no
	// guarantee that this will really work. Most likely this is discovered
	// based on environment variables that expresses the existence of some
	// assumed environment, ...
	SupportedAssumed Supported = 2
)

// AllSupports returns all possible values of Supported.
func AllSupports() Supports {
	return Supports{SupportedNone, SupportedNative, SupportedAssumed}
}

// DetectSupportForWriter detects for the given io.Writer if color is supported
// or not. Additionally it prepares the given logger with the color mode and
// will return a modified instance of it. If color is not supported the original
// io.Writer is still returned. Errors are only returned in cases where
// something bad happens while detecting or preparing for color.
//
// See SupportAssumptionDetections for assumed detections.
func DetectSupportForWriter(w io.Writer) (prepared io.Writer, supported Supported, err error) {
	actual, err := prepareForColors(w)
	if actual && err == nil {
		return w, SupportedNative, nil
	}

	for _, d := range SupportAssumptionDetections {
		if v, err := d(); err != nil {
			return w, SupportedNone, err
		} else if v {
			return w, SupportedAssumed, nil
		}
	}

	return w, SupportedNone, err
}

// IsSupported returns true when color should most likely work.
func (instance Supported) IsSupported() bool {
	switch instance {
	case SupportedNative, SupportedAssumed:
		return true
	default:
		return false
	}
}

// MarshalText implements encoding.TextMarshaler
func (instance Supported) MarshalText() (text []byte, err error) {
	switch instance {
	case SupportedNone:
		return []byte("none"), nil
	case SupportedNative:
		return []byte("native"), nil
	case SupportedAssumed:
		return []byte("assumed"), nil
	default:
		return nil, fmt.Errorf("%w: %d", ErrIllegalSupport, instance)
	}
}

// UnmarshalText implements encoding.TextMarshaler
func (instance *Supported) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "none", "no", "0", "never", "off", "false":
		*instance = SupportedNone
		return nil
	case "native":
		*instance = SupportedNative
		return nil
	case "assumed":
		*instance = SupportedAssumed
		return nil
	default:
		return fmt.Errorf("%w: %v", ErrIllegalSupport, string(text))
	}
}

// String prints out a meaningful representation of this instance.
func (instance Supported) String() string {
	if text, err := instance.MarshalText(); err != nil {
		return fmt.Sprintf("illegal-color-support-%d", instance)
	} else {
		return string(text)
	}
}

// Set will set this instance to the given plain value or errors.
func (instance *Supported) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

// Supports is a multiple version of Supported.
type Supports []Supported

// Strings returns a meaningful representation of all of it's values.
func (instance Supports) Strings() []string {
	result := make([]string, len(instance))
	for i, v := range instance {
		result[i] = v.String()
	}
	return result
}

// String returns a meaningful representation of this instance.
func (instance Supports) String() string {
	return strings.Join(instance.Strings(), ",")
}
