package level

import (
	"fmt"

	"github.com/echocat/slf4g/level"
)

type Formatter interface {
	Format(level.Level) string
}

type FormatterFunc func(level.Level) string

func (instance FormatterFunc) Format(v level.Level) string {
	return instance(v)
}

var DefaultFormatter Formatter = FormatterFunc(func(l level.Level) string {
	switch l {
	case level.Trace:
		return "TRACE"
	case level.Debug:
		return "DEBUG"
	case level.Info:
		return " INFO"
	case level.Warn:
		return " WARN"
	case level.Error:
		return "ERROR"
	case level.Fatal:
		return "FATAL"
	default:
		return fmt.Sprintf("%5d", l)
	}
})
