package color

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrIllegalSupport = errors.New("illegal color-support")
)

type Support uint8

const (
	SupportNone    Support = 0
	SupportNative  Support = 1
	SupportAssumed Support = 2
)

func DetectSupportForWriter(w io.Writer) (prepared io.Writer, support Support, err error) {
	supported, err := prepareForColors(w)
	if supported && err == nil {
		return w, SupportNative, nil
	}

	for _, d := range SupportAssumptionDetections {
		if d() {
			return w, SupportAssumed, nil
		}
	}

	return w, SupportNone, err
}

func (instance Support) IsSupported() bool {
	switch instance {
	case SupportNative, SupportAssumed:
		return true
	default:
		return false
	}
}

func (instance Support) MarshalText() (text []byte, err error) {
	switch instance {
	case SupportNone:
		return []byte("none"), nil
	case SupportNative:
		return []byte("native"), nil
	case SupportAssumed:
		return []byte("assumed"), nil
	default:
		return nil, fmt.Errorf("%w: %d", ErrIllegalSupport, instance)
	}
}

func (instance *Support) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "none", "no", "0", "never":
		*instance = SupportNone
		return nil
	case "native":
		*instance = SupportNative
		return nil
	case "assumed":
		*instance = SupportAssumed
		return nil
	default:
		return fmt.Errorf("%w: %v", ErrIllegalSupport, string(text))
	}
}

func (instance Support) String() string {
	if text, err := instance.MarshalText(); err != nil {
		return fmt.Sprintf("illegal-color-support-%d", instance)
	} else {
		return string(text)
	}
}

func (instance Support) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}
