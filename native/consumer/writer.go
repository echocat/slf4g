package consumer

import (
	"fmt"
	"io"
	"sync"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native/color"
	"github.com/echocat/slf4g/native/formatter"
	"github.com/echocat/slf4g/native/formatter/hints"
	"github.com/echocat/slf4g/native/interceptor"
)

type WritingConsumer struct {
	owner Owner
	out   io.Writer

	EventFormatter formatter.Formatter
	Interceptor    interceptor.Interceptor
	ColorMode      color.Mode
	HintsProvider  func(event log.Event, source log.CoreLogger) hints.Hints

	PrintErrorOnColorInitialization bool

	colorSupport color.Support
	initOnce     sync.Once
}

func NewWritingConsumer(owner Owner, out io.Writer) *WritingConsumer {
	return &WritingConsumer{
		owner: owner,
		out:   out,
	}
}

func (instance *WritingConsumer) Consume(event log.Event, source log.CoreLogger) {
	if event == nil {
		return
	}

	if event = instance.onBeforeLog(event, source); event == nil {
		return
	}

	if !source.IsLevelEnabled(event.GetLevel()) {
		return
	}

	out := instance.getOut()

	f := instance.getEventFormatter()
	h := instance.provideHints(event, source)
	content, err := f.Format(event, source.GetProvider(), h)
	if err != nil {
		content = []byte(fmt.Sprintf("ERR: Cannot format event %v: %v", event, err))
	}

	_, _ = out.Write(content)

	_ = instance.onAfterLog(event, source)
}

func (instance *WritingConsumer) initIfRequired() {
	instance.initOnce.Do(func() {
		var err error
		instance.out, instance.colorSupport, err = color.DetectSupportForWriter(instance.out)
		if err != nil && instance.PrintErrorOnColorInitialization {
			_, _ = fmt.Fprintf(instance.out, "ERR!!! Cannot intiate colors for target: %v\n", err)
		}
	})
}

func (instance *WritingConsumer) getOut() io.Writer {
	return instance.out
}

func (instance *WritingConsumer) GetInterceptor() interceptor.Interceptor {
	return instance.Interceptor
}

func (instance *WritingConsumer) SetInterceptor(v interceptor.Interceptor) {
	instance.Interceptor = v
}

func (instance *WritingConsumer) getInterceptor() interceptor.Interceptor {
	if i := instance.GetInterceptor(); i != nil {
		return i
	}
	if i := instance.owner.GetInterceptor(); i != nil {
		return i
	}
	return interceptor.Noop()
}

func (instance *WritingConsumer) onBeforeLog(event log.Event, source log.CoreLogger) log.Event {
	return instance.getInterceptor().OnBeforeLog(event, source.GetProvider())
}

func (instance *WritingConsumer) onAfterLog(event log.Event, source log.CoreLogger) (canContinue bool) {
	return instance.getInterceptor().OnAfterLog(event, source.GetProvider())
}

func (instance *WritingConsumer) GetFormatter() formatter.Formatter {
	return instance.EventFormatter
}

func (instance *WritingConsumer) SetFormatter(v formatter.Formatter) {
	instance.EventFormatter = v
}

func (instance *WritingConsumer) getEventFormatter() formatter.Formatter {
	if f := instance.GetFormatter(); f != nil {
		return f
	}
	if f := instance.owner.GetFormatter(); f != nil {
		return f
	}
	return formatter.DefaultConsole
}

func (instance *WritingConsumer) GetColorSupport() color.Support {
	instance.initIfRequired()
	return instance.colorSupport
}

func (instance *WritingConsumer) provideHints(event log.Event, source log.CoreLogger) hints.Hints {
	if v := instance.HintsProvider; v != nil {
		return v(event, source)
	}
	return instance
}
