package sdk

import (
	sdklog "log"

	"github.com/echocat/slf4g/level"

	log "github.com/echocat/slf4g"
)

// Configure configures the standard SDK logging framework to use slf4g with
// its root Logger and logs everything printed using Print(), Printf() and
// Println() on level.Info.
//
// # Limitations
//
// 1# Stuff logged using Fatal*() and Panic*() are on logged on the same Level
// as everything else.
//
// 2# ATTENTION! Fatal*() and Panic*() will still exit the whole application or
// panics afterwards. This cannot be prevented when you use the logger
// directly. If possible use NewLogger() to use an SDK compatible interface.
func Configure(customizer ...func(*log.LoggingWriter)) {
	ConfigureWith(log.GetRootLogger(), level.Info, customizer...)
}

// ConfigureWith configures the standard SDK logging framework to use slf4g with
// it's the given Logger and logs everything printed using Print(), Printf() and
// Println() on given level.Level.
//
// # Limitations
//
// 1# Stuff logged using Fatal*() and Panic*() are on logged on the same Level
// as everything else.
//
// 2# ATTENTION! Fatal*() and Panic*() will still exit the whole application or
// panics afterwards. This cannot be prevented when you use the logger
// directly. If possible use NewLogger() to use an SDK compatible interface.
func ConfigureWith(target log.CoreLogger, logAs level.Level, customizer ...func(*log.LoggingWriter)) {
	w := &log.LoggingWriter{
		Logger:         target,
		LevelExtractor: level.FixedLevelExtractor(logAs),
		SkipFrames:     2, // of the SDK based log
	}
	for _, c := range customizer {
		c(w)
	}
	sdklog.SetOutput(w)
	sdklog.SetPrefix("")
	sdklog.SetFlags(0)
}

// NewWrapper creates an SDK Logger which is use slf4g the provided Logger and
// logs everything printed using Print(), Printf() and Println() on given
// level.Level.
//
// # Limitations
//
// 1# Stuff logged using Fatal*() and Panic*() are on logged on the same Level
// as everything else.
//
// 2# ATTENTION! Fatal*() and Panic*() will still exit the whole application or
// panics afterwards. This cannot be prevented when you use the logger
// directly. If possible use NewLogger() to use an SDK compatible interface.
func NewWrapper(target log.CoreLogger, logAs level.Level) *sdklog.Logger {
	return sdklog.New(&log.LoggingWriter{
		Logger:         target,
		LevelExtractor: level.FixedLevelExtractor(logAs),
		SkipFrames:     2, // of the SDK based log
	}, "", 0)
}
