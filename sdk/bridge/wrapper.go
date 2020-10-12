package sdk

import (
	stdlog "log"

	"github.com/echocat/slf4g/level"

	log "github.com/echocat/slf4g"
)

// Configure configures the standard SDK logging framework to use slf4g with
// its root Logger and logs everything printed using Print(), Printf() and
// Println() on level.Info.
//
// Limitations
//
// 1# Stuff logged using Fatal*() and Panic*() are on logged on the same Level
// than everything else.
//
// 2# ATTENTION! Fatal*() and Panic*() will still exit the whole application or
// panics afterwards. This cannot be prevented when you use the logger
// directly. If possible use NewLogger() to use an SDK compatible interface.
func Configure() {
	ConfigureWith(log.GetRootLogger(), level.Info)
}

// ConfigureWith configures the standard SDK logging framework to use slf4g with
// its the given Logger and logs everything printed using Print(), Printf() and
// Println() on given level.Level.
//
// Limitations
//
// 1# Stuff logged using Fatal*() and Panic*() are on logged on the same Level
// than everything else.
//
// 2# ATTENTION! Fatal*() and Panic*() will still exit the whole application or
// panics afterwards. This cannot be prevented when you use the logger
// directly. If possible use NewLogger() to use an SDK compatible interface.
func ConfigureWith(target log.CoreLogger, logAs level.Level) {
	w := &log.LoggingWriter{
		Logger:         target,
		LevelExtractor: level.FixedLevelExtractor(logAs),
	}
	stdlog.SetOutput(w)
	stdlog.SetPrefix("")
	stdlog.SetFlags(0)
}

// NewWrapper creates a SDK Logger which is use slf4g the provided Logger and
// logs everything printed using Print(), Printf() and Println() on given
// level.Level.
//
// Limitations
//
// 1# Stuff logged using Fatal*() and Panic*() are on logged on the same Level
// than everything else.
//
// 2# ATTENTION! Fatal*() and Panic*() will still exit the whole application or
// panics afterwards. This cannot be prevented when you use the logger
// directly. If possible use NewLogger() to use an SDK compatible interface.
func NewWrapper(target log.CoreLogger, logAs level.Level) *stdlog.Logger {
	return stdlog.New(&log.LoggingWriter{
		Logger:         target,
		LevelExtractor: level.FixedLevelExtractor(logAs),
	}, "", 0)
}
