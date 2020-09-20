package log

type CoreLogger interface {
	GetName() string
	Log(Event)
	IsLevelEnabled(Level) bool
	GetProvider() Provider
}
