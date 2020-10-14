package level

import (
	"encoding"
	"flag"
	"fmt"

	"github.com/echocat/slf4g/level"
)

func AsSerializable(in *level.Level, aware NamesAware) Serializable {
	return &serializableImpl{in, aware.GetLevelNames()}
}

type Serializable interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	flag.Value

	AsLevel() *level.Level
}

type serializableImpl struct {
	*level.Level
	names Names
}

func (instance serializableImpl) AsLevel() *level.Level {
	return instance.Level
}

func (instance serializableImpl) MarshalText() (text []byte, err error) {
	name, err := instance.names.FromOrdinal(uint16(*instance.Level))
	if err != nil {
		return nil, err
	}
	return []byte(name), nil
}

func (instance serializableImpl) UnmarshalText(text []byte) error {
	ordinal, err := instance.names.ToOrdinal(string(text))
	if err != nil {
		return err
	}
	l := level.Level(ordinal)
	instance.Level = &l
	return nil
}

func (instance serializableImpl) String() string {
	if text, err := instance.MarshalText(); err != nil {
		return fmt.Sprintf("illegal-level-%d", instance)
	} else {
		return string(text)
	}
}

func (instance *serializableImpl) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}
