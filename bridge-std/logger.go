package std

import log "github.com/echocat/slf4g"

type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}

func NewLogger(target log.CoreLogger) Logger {
	return &LoggerImpl{
		CoreLogger: target,
		OnFatal:    DefaultOnFatal,
		OnPanic:    DefaultOnPanic,
	}
}
