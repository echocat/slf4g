package std

import (
	stdlog "log"

	log "github.com/echocat/slf4g"
)

func Configure() {
	ConfigureWith(log.GetRootLogger(), log.LevelInfo)
}

func ConfigureWith(target log.CoreLogger, logAs log.Level) {
	w := &log.LoggingWriter{
		CoreLogger: target,
		LogAs:      logAs,
	}
	stdlog.SetOutput(w)
	stdlog.SetPrefix("")
	stdlog.SetFlags(0)
}

func NewWrapper(target log.CoreLogger, logAs log.Level) *stdlog.Logger {
	return stdlog.New(&log.LoggingWriter{
		CoreLogger: target,
		LogAs:      logAs,
	}, "", 0)
}
