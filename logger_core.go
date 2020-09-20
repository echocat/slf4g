package log

type CoreLogger interface {
	GetName() string
	LogEvent(Event)
	IsLevelEnabled(Level) bool
	GetProvider() Provider
}
