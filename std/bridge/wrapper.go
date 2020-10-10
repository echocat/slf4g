package std

import (
	stdlog "log"

	"github.com/echocat/slf4g/level"

	log "github.com/echocat/slf4g"
)

func Configure() {
	ConfigureWith(log.GetRootLogger(), level.Info)
}

func ConfigureWith(target log.CoreLogger, logAs level.Level) {
	w := &log.LoggingWriter{
		CoreLogger: target,
		LogAs:      logAs,
	}
	stdlog.SetOutput(w)
	stdlog.SetPrefix("")
	stdlog.SetFlags(0)
}

func NewWrapper(target log.CoreLogger, logAs level.Level) *stdlog.Logger {
	return stdlog.New(&log.LoggingWriter{
		CoreLogger: target,
		LogAs:      logAs,
	}, "", 0)
}
