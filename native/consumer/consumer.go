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
	formatter.Aware
	interceptor.Aware
}

type Aware interface {
	GetConsumer() Consumer
	SetConsumer(Consumer)
}
