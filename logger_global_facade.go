package log

import "github.com/echocat/slf4g/fields"

var globalLoggerFacadeV = &loggerImpl{
	getCoreLogger: func() CoreLogger { return GetProvider().GetLogger(GlobalLoggerName) },
	fields:        fields.Empty(),
}
