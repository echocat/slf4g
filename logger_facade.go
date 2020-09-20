package log

import (
	"fmt"
	"github.com/echocat/slf4g/fields"
	"reflect"
)

func newLoggerFacade(provider func() Logger) Logger {
	return &loggerFacade{
		loggerImpl: loggerImpl{
			coreProvider: func() CoreLogger { return provider() },
			fields:       fields.Empty(),
		},
	}
}

type loggerFacade struct {
	loggerImpl
}

func (instance *loggerFacade) Unwrap() Logger {
	cl := instance.coreProvider()
	if l, ok := cl.(Logger); ok {
		return l
	}
	panic(fmt.Sprintf("expected a logger that implements Logger but got: %v", reflect.TypeOf(cl)))
}
