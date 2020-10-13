package recording

import (
	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/level"
)

// Logger is a fully implemented version of CoreLogger and log.Logger.
type Logger struct {
	log.Logger
	*CoreLogger
}

// NewLogger creates a new instance of Logger which is ready to use.
func NewLogger() *Logger {
	coreLogger := NewCoreLogger()
	wrapped := log.NewLogger(coreLogger)
	return &Logger{
		Logger:     wrapped,
		CoreLogger: coreLogger,
	}
}

// GetName implements log.Logger#GetName()
func (instance *Logger) GetName() string {
	return instance.CoreLogger.GetName()
}

// GetProvider implements log.Logger#GetProvider()
func (instance *Logger) GetProvider() log.Provider {
	return instance.CoreLogger.GetProvider()
}

// IsLevelEnabled implements log.Logger#IsLevelEnabled()
func (instance *Logger) IsLevelEnabled(v level.Level) bool {
	return instance.CoreLogger.IsLevelEnabled(v)
}

// Log implements log.Logger#Log()
func (instance *Logger) Log(event log.Event) {
	instance.CoreLogger.Log(event.WithCallDepth(1))
}
