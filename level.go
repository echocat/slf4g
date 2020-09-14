package log

import (
	"errors"
	"fmt"
)

var (
	ErrIllegalLevel = errors.New("illegal level")

	defaultLevels = []Level{LevelTrace, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal, LevelPanic}
)

const (
	LevelTrace = Level(1000)
	LevelDebug = Level(2000)
	LevelInfo  = Level(3000)
	LevelWarn  = Level(4000)
	LevelError = Level(5000)
	LevelFatal = Level(6000)
	LevelPanic = Level(7000)
)

type Level uint16

func (instance Level) MarshalText() (text []byte, err error) {
	name, err := GetProvider().GetLevelNames().FromOrdinal(uint16(instance))
	if err != nil {
		return nil, err
	}
	return []byte(name), nil
}

func (instance *Level) UnmarshalText(text []byte) error {
	ordinal, err := GetProvider().GetLevelNames().ToOrdinal(string(text))
	if err != nil {
		return err
	}
	*instance = Level(ordinal)
	return nil
}

func (instance Level) String() string {
	if text, err := instance.MarshalText(); err != nil {
		return fmt.Sprintf("illegal-level-%d", instance)
	} else {
		return string(text)
	}
}

func (instance Level) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance Level) CompareTo(o Level) int {
	return int(instance) - int(o)
}

type Levels []Level

func (instance Levels) Len() int {
	return len(instance)
}

func (instance Levels) Swap(i, j int) {
	instance[i], instance[j] = instance[j], instance[i]
}

func (instance Levels) Less(i, j int) bool {
	return instance[i].CompareTo(instance[j]) < 0
}

type LevelAware interface {
	GetLevel() Level
	SetLevel(Level)
}

type LevelProvider func() []Level

var DefaultLevelProvider LevelProvider = func() []Level {
	return defaultLevels
}
