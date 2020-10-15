package color

import (
	"errors"
	"fmt"
	"strings"
)

// ErrIllegalMode will be returned in situations where illegal values or
// representations if a Mode are provided.
var ErrIllegalMode = errors.New("illegal color-mode")

// Mode defines how colors should be used.
type Mode uint8

const (
	// ModeAuto will result in that the application tries to detect
	// automatically if it is possible and meaningful to use colors in the
	// a given terminal.
	ModeAuto Mode = 0

	// ModeAlways tells the application to always use colorful output (if
	// supported).
	ModeAlways Mode = 1

	// ModeNever tells the application to never use colorful output.
	ModeNever Mode = 2
)

// AllModes returns all possible values of Mode.
func AllModes() Modes {
	return Modes{ModeAuto, ModeAlways, ModeNever}
}

// ShouldColorize will check for the given combination of this instance and the
// given Support value if the output should be colorized.
func (instance Mode) ShouldColorize(checking Supported) bool {
	switch instance {
	case ModeAuto:
		return checking.IsSupported()
	case ModeAlways:
		return true
	default:
		return false
	}
}

// MarshalText implements encoding.TextMarshaler
func (instance Mode) MarshalText() (text []byte, err error) {
	switch instance {
	case ModeAuto:
		return []byte("auto"), nil
	case ModeAlways:
		return []byte("always"), nil
	case ModeNever:
		return []byte("never"), nil
	default:
		return nil, fmt.Errorf("%w: %d", ErrIllegalMode, instance)
	}
}

// UnmarshalText implements encoding.TextMarshaler
func (instance *Mode) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "auto", "automatic", "detect":
		*instance = ModeAuto
		return nil
	case "always", "on", "true", "1":
		*instance = ModeAlways
		return nil
	case "never", "off", "false", "0":
		*instance = ModeNever
		return nil
	default:
		return fmt.Errorf("%w: %v", ErrIllegalMode, string(text))
	}
}

// String prints out a meaningful representation of this instance.
func (instance Mode) String() string {
	if text, err := instance.MarshalText(); err != nil {
		return fmt.Sprintf("illegal-color-mode-%d", instance)
	} else {
		return string(text)
	}
}

// Set will set this instance to the given plain value or errors.
func (instance *Mode) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

// Modes is a multiple version of Mode.
type Modes []Mode

// Strings returns a meaningful representation of all of it's values.
func (instance Modes) Strings() []string {
	result := make([]string, len(instance))
	for i, v := range instance {
		result[i] = v.String()
	}
	return result
}

// String returns a meaningful representation of this instance.
func (instance Modes) String() string {
	return strings.Join(instance.Strings(), ",")
}
