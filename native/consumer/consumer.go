package consumer

import (
	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native/formatter"
	"github.com/echocat/slf4g/native/interceptor"
)

type Consumer interface {
	Consume(event log.Event, source log.CoreLogger)
}

type Owner interface {
	GetFormatter() formatter.Formatter
	GetInterceptor() interceptor.Interceptor
}
