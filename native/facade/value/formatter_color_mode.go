package value

import (
	"github.com/echocat/slf4g/native/color"
)

// FormatterColorMode enables the access the color.Mode of a FormatterTarget.
type FormatterColorMode struct {
	provider func() (color.Mode, bool)
	sink     func(color.Mode) error
	value    *color.Mode
}

// Set implements flag.Value.
func (instance *FormatterColorMode) Set(plain string) error {
	var buf color.Mode
	if err := buf.Set(plain); err != nil {
		return err
	}

	if err := instance.sink(buf); err != nil {
		return err
	}

	instance.value = &buf
	return nil
}

func (instance FormatterColorMode) get() color.Mode {
	if v, ok := instance.provider(); ok {
		return v
	}
	if v := instance.value; v != nil {
		return *v
	}
	return color.ModeAuto
}

// Get implements flag.Getter.
func (instance FormatterColorMode) Get() interface{} {
	return instance.get()
}

// String implements flag.Value.
func (instance FormatterColorMode) String() string {
	return instance.get().String()
}

// Type returns the type as a string.
func (instance FormatterColorMode) Type() string {
	return "logColorMode"
}
