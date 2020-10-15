package level

import (
	"encoding"
	"flag"
	"fmt"

	"github.com/echocat/slf4g/level"
)

// AsNamed wraps the given level.Level into Named to it in a human readable
// format and provides the possibility to marshal and get used with
// flag (or compatible) packages.
func AsNamed(in *level.Level, aware NamesAware) Named {
	return &namedImpl{in, aware.GetLevelNames()}
}

// Named represents a level.Level in a human readable format and provides the
// possibility to marshal and get used with flag (or compatible) packages.
type Named interface {
	Unwrap() *level.Level

	encoding.TextMarshaler
	encoding.TextUnmarshaler
	flag.Getter
}

type namedImpl struct {
	level *level.Level
	names Names
}

func (instance *namedImpl) Get() interface{} {
	return instance.Unwrap()
}

func (instance *namedImpl) Unwrap() *level.Level {
	return instance.level
}

func (instance *namedImpl) MarshalText() (text []byte, err error) {
	name, err := instance.names.ToName(*instance.level)
	if err != nil {
		return nil, err
	}
	return []byte(name), nil
}

func (instance *namedImpl) UnmarshalText(text []byte) error {
	v, err := instance.names.ToLevel(string(text))
	if err != nil {
		return err
	}
	*instance.level = v
	return nil
}

func (instance *namedImpl) String() string {
	if text, err := instance.MarshalText(); err != nil {
		return fmt.Sprintf("illegal-level-%d", instance)
	} else {
		return string(text)
	}
}

func (instance *namedImpl) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}
