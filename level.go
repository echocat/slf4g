package log

const (
	LevelTrace = Level(1000)
	LevelDebug = Level(2000)
	LevelInfo  = Level(3000)
	LevelWarn  = Level(4000)
	LevelError = Level(5000)
	LevelFatal = Level(6000)
)

type Level uint16

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

type LevelProvider func() []Level

var DefaultLevelProvider LevelProvider = func() []Level {
	return []Level{LevelTrace, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal}
}
