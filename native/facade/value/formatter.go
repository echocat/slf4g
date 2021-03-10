package value

import (
	"fmt"

	"github.com/echocat/slf4g/native/formatter"
)

// FormatterTarget defines an object that receives the consumer.Consumer managed
// by the Formatter value facade.
type FormatterTarget interface {
	formatter.MutableAware
}

// Formatter is a value facade for transparent setting of consumer.Consumer
// for the slf4g/native implementation. This is quite handy for usage
// with flags package of the SDK or similar flag libraries. This might
// be usable, too in contexts where serialization might be required.
type Formatter struct {
	// Target is the instance of consumer.MutableAware which should be
	// configured by this facade.
	Target FormatterTarget

	// Codec is used to transform provided plain data. If this is not defined
	// DefaultFormatterCodec is used.
	Codec FormatterCodec
}

// NewFormatter creates a new instance of Formatter with the given target.
func NewFormatter(target FormatterTarget, customizer ...func(*Formatter)) Formatter {
	result := Formatter{
		Target: target,
	}

	for _, c := range customizer {
		c(&result)
	}

	return result
}

// Set implements flag.Value.
func (instance Formatter) Set(plain string) error {
	v, err := instance.getCodec().Parse(plain)
	if err != nil {
		return err
	}

	instance.Target.SetFormatter(v)
	return nil
}

// Get implements flag.Getter.
func (instance Formatter) Get() interface{} {
	return instance.Target.GetFormatter()
}

// String implements flag.Value.
func (instance Formatter) String() string {
	b, err := instance.MarshalText()
	if err != nil {
		return fmt.Sprintf("ERR-%v", err)
	}
	return string(b)
}

// MarshalText implements encoding.TextMarshaler
func (instance Formatter) MarshalText() (text []byte, err error) {
	if instance.Target == nil {
		return nil, nil
	}
	name, err := instance.getCodec().Format(instance.Target.GetFormatter())
	return []byte(name), err
}

// UnmarshalText implements encoding.TextUnmarshaler
func (instance Formatter) UnmarshalText(text []byte) error {
	return instance.Set(string(text))
}

func (instance Formatter) getCodec() FormatterCodec {
	if v := instance.Codec; v != nil {
		return v
	}
	if v := DefaultFormatterCodec; v != nil {
		return v
	}
	return NoopFormatterCodec()
}
