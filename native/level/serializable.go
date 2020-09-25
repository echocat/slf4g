package level

import (
	"encoding"
	"flag"
	"fmt"
	"github.com/echocat/slf4g"
)

func AsSerializable(level *log.Level, aware NamesAware) Serializable {
	return &serializable{level, aware.GetLevelNames()}
}

type Serializable interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	flag.Value

	AsLevel() *log.Level
}

type serializable struct {
	*log.Level
	levelNames Names
}

func (instance serializable) AsLevel() *log.Level {
	return instance.Level
}

func (instance serializable) MarshalText() (text []byte, err error) {
	name, err := instance.levelNames.FromOrdinal(uint16(*instance.Level))
	if err != nil {
		return nil, err
	}
	return []byte(name), nil
}

func (instance serializable) UnmarshalText(text []byte) error {
	ordinal, err := instance.levelNames.ToOrdinal(string(text))
	if err != nil {
		return err
	}
	l := log.Level(ordinal)
	instance.Level = &l
	return nil
}

func (instance serializable) String() string {
	if text, err := instance.MarshalText(); err != nil {
		return fmt.Sprintf("illegal-level-%d", instance)
	} else {
		return string(text)
	}
}

func (instance *serializable) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}
