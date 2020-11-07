package consumer

import (
	"fmt"
	"io"
	"sync"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native/color"
	"github.com/echocat/slf4g/native/formatter"
	"github.com/echocat/slf4g/native/hints"
	"github.com/echocat/slf4g/native/interceptor"
)

type Writer struct {
	out io.Writer

	Formatter     formatter.Formatter
	Interceptor   interceptor.Interceptor
	HintsProvider func(event log.Event, source log.CoreLogger) hints.Hints

	// Synchronized defines if this instance can be used in concurrent
	// environments; which is meaningful in the most context. It might have
	// additional performance costs.
	Synchronized                    bool
	PrintErrorOnColorInitialization bool

	colorSupported *color.Supported
	mutex          sync.Mutex
}

// NewWriter creates a new instance of Writer which can be customized using
// customizer and is ready to use. The created instance is synchronized by
// default (See Writer.Synchronized).
func NewWriter(out io.Writer, customizer ...func(*Writer)) *Writer {
	result := &Writer{
		out:          out,
		Synchronized: true,
	}
	for _, c := range customizer {
		c(result)
	}
	return result
}

// Consume implements Consumer.Consume()
func (instance *Writer) Consume(event log.Event, source log.CoreLogger) {
	if event == nil || instance.out == nil {
		return
	}

	if instance.Synchronized {
		instance.mutex.Lock()
		defer instance.mutex.Unlock()
	}
	instance.initIfRequired()

	if event = instance.onBeforeLog(event, source); event == nil {
		return
	}

	if !source.IsLevelEnabled(event.GetLevel()) {
		return
	}

	out := instance.getOut()

	f := instance.getFormatter()
	h := instance.provideHints(event, source)
	content, err := f.Format(event, source.GetProvider(), h)
	if err != nil {
		content = []byte(fmt.Sprintf("ERR: Cannot format event %v: %v", event, err))
	}

	_, _ = out.Write(content)

	_ = instance.onAfterLog(event, source)
}

func (instance *Writer) initIfRequired() {
	if instance.colorSupported == nil {
		out, supported, err := color.DetectSupportForWriter(instance.out)
		if err != nil && instance.PrintErrorOnColorInitialization {
			_, _ = fmt.Fprintf(instance.out, "ERR!!! Cannot intiate colors for target: %v\n", err)
		}
		instance.out = out
		instance.colorSupported = &supported
	}
}

func (instance *Writer) getOut() io.Writer {
	return instance.out
}

func (instance *Writer) onBeforeLog(event log.Event, source log.CoreLogger) log.Event {
	return instance.getInterceptor().OnBeforeLog(event, source.GetProvider())
}

func (instance *Writer) onAfterLog(event log.Event, source log.CoreLogger) (canContinue bool) {
	return instance.getInterceptor().OnAfterLog(event, source.GetProvider())
}

func (instance *Writer) getInterceptor() interceptor.Interceptor {
	if v := instance.Interceptor; v != nil {
		return v
	}
	if v := interceptor.Default; v != nil {
		return v
	}
	return interceptor.Noop()
}

func (instance *Writer) getFormatter() formatter.Formatter {
	if v := instance.Formatter; v != nil {
		return v
	}
	if v := formatter.Default; v != nil {
		return v
	}
	return formatter.Noop()
}

func (instance *Writer) provideHints(event log.Event, source log.CoreLogger) hints.Hints {
	if v := instance.HintsProvider; v != nil {
		return v(event, source)
	}
	return &writingConsumerHints{instance}
}

type writingConsumerHints struct {
	*Writer
}

func (instance *Writer) IsColorSupported() color.Supported {
	return *instance.colorSupported
}
