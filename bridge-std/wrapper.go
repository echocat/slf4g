package std

import (
	log "github.com/echocat/slf4g"
	stdlog "log"
)

func NewWrapper(target log.CoreLogger, logAs log.Level) *stdlog.Logger {
	return stdlog.New(&log.Writer{
		CoreLogger: target,
		LogAs:      logAs,
	}, "", 0)
}

func ConfigureStd() {
	ConfigureStdWith(log.GetGlobalLogger(), log.LevelInfo)
}

func ConfigureStdWith(target log.CoreLogger, logAs log.Level) {
	w := &log.Writer{
		CoreLogger: target,
		LogAs:      logAs,
	}
	stdlog.SetOutput(w)
	stdlog.SetPrefix("")
	stdlog.SetFlags(0)
}
