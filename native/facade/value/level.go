package value

import (
	"fmt"
	"strconv"

	"github.com/echocat/slf4g/level"
	nlevel "github.com/echocat/slf4g/native/level"
)

// LevelTarget defines an object that receives the level.Level managed
// by the Level value facade.
type LevelTarget interface {
	level.MutableAware
}

// Level is a value facade for transparent setting of level.Level for
// the slf4g/native implementation. This is quite handy for usage with
// flags package of the SDK or similar flag libraries. This might be
// usable, too in contexts where serialization might be required.
type Level struct {
	// Target is the instance of LevelTarget which should be configured
	// by this facade.
	Target LevelTarget

	// Names is used to transform provided plain data. If this is not defined
	// the Target is assumed as nlevel.NamesAware or if this even does not work
	// nlevel.DefaultNames is used.
	Names nlevel.Names
}

// NewLevel creates a new instance of Level with the given target.
func NewLevel(target LevelTarget, customizer ...func(*Level)) Level {
	result := Level{
		Target: target,
	}

	for _, c := range customizer {
		c(&result)
	}

	return result
}

// Set implements flag.Value.
func (instance Level) Set(plain string) error {
	if l, err := strconv.ParseUint(plain, 10, 16); err == nil {
		instance.Target.SetLevel(level.Level(l))
		return nil
	}

	l, err := instance.getNames().ToLevel(plain)
	if err != nil {
		return err
	}

	instance.Target.SetLevel(l)
	return nil
}

// Get implements flag.Getter.
func (instance Level) Get() interface{} {
	return instance.Target.GetLevel()
}

// String implements flag.Value.
func (instance Level) String() string {
	b, err := instance.MarshalText()
	if err != nil {
		return fmt.Sprintf("ERR-%v", err)
	}
	return string(b)
}

// MarshalText implements encoding.TextMarshaler
func (instance Level) MarshalText() (text []byte, err error) {
	if instance.Target == nil {
		return nil, nil
	}
	name, err := instance.getNames().ToName(instance.Target.GetLevel())
	return []byte(name), err
}

// UnmarshalText implements encoding.TextUnmarshaler
func (instance Level) UnmarshalText(text []byte) error {
	return instance.Set(string(text))
}

func (instance Level) getNames() nlevel.Names {
	if v := instance.Names; v != nil {
		return v
	}
	if va, ok := instance.Target.(nlevel.NamesAware); ok {
		if v := va.GetLevelNames(); v != nil {
			return v
		}
	}
	if v := nlevel.DefaultNames; v != nil {
		return v
	}
	return nlevel.NewNames()
}
