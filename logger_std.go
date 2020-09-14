package log

type StdLogger interface {
	CoreLogger

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

func AsStdLogger(in CoreLogger) StdLogger {
	if v, ok := in.(StdLogger); ok {
		return v
	}
	return &loggerImpl{CoreLogger: in}
}
