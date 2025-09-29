package testlog

import log "github.com/echocat/slf4g"

type coreLoggerRenamed struct {
	*coreLogger
	name string
}

// Log implements log.CoreLogger#Log(event).
func (instance *coreLoggerRenamed) Log(event log.Event, skipFrames uint16) {
	instance.tb.Helper()
	instance.log(instance.name, event, skipFrames+1)
}

// GetName implements log.CoreLogger#GetName().
func (instance *coreLoggerRenamed) GetName() string {
	return instance.name
}

// Helper wraps the helper of the testing framework into this logger.
// As this is called by the whole logging stack (if required) this will ensure
// the SDK logging framework respects the top entry as the log position.
func (instance *coreLoggerRenamed) Helper() func() {
	return instance.tb.Helper
}
