package color

import (
	"errors"
	"fmt"
	"strings"
)

type Mode uint8

var (
	ErrIllegalMode = errors.New("illegal color-mode")
)

const (
	ModeAuto   Mode = 0
	ModeAlways Mode = 1
	ModeNever  Mode = 2
)

func (instance Mode) IsEnabled(support Support) bool {
	switch instance {
	case ModeAuto:
		return support.IsSupported()
	case ModeAlways:
		return true
	default:
		return false
	}
}

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

func (instance Mode) String() string {
	if text, err := instance.MarshalText(); err != nil {
		return fmt.Sprintf("illegal-color-mode-%d", instance)
	} else {
		return string(text)
	}
}

func (instance Mode) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}
