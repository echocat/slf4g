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

// Writer is an implementation of Writer which formats the consumed log.Entry
// using a configured Formatter and logs it to the configured io.Writer.
//
// NewWriter() is used to create a new instance.
type Writer struct {
	// Formatter to format the consumed log.Event with. If nothing was provided
	// formatter.Default will be used.
	Formatter formatter.Formatter

	// Interceptor can be used to intercept the consumption of an event shortly
	// before the actual consumption or directly afterwards. If nothing was
	// provided interceptor.Default will be used.
	Interceptor interceptor.Interceptor

	// HintsProvider is used to determine an instance of hints.Hints for the
	// actual log.Event. These might be used by the actual Formatter to know how
	// to format the actual log.Event correctly. This could include of
	// colorization is supported and demanded or any other stuff. If nothing was
	// provided a default instance will be provided which provides:
	// 1. hints.ColorsSupport
	HintsProvider func(event log.Event, source log.CoreLogger) hints.Hints

	// Synchronized defines if this instance can be used in concurrent
	// environments; which is meaningful in the most context. It might have
	// additional performance costs.
	Synchronized bool

	// OnFormatError will be called if there as any kind of error while
	// formatting an log.Event using the configured Formatter. If nothing was
	// provided these errors will result in a panic.
	OnFormatError func(*Writer, io.Writer, error)

	// OnColorInitializationError will be called if there as any kind of error
	// while initialize the color support. If nothing was provided these errors
	// will be silently swallowed.
	OnColorInitializationError func(*Writer, io.Writer, error)

	out            io.Writer
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
	if event == nil {
		return
	}

	out := instance.GetOut()
	if out == nil {
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

	f := instance.GetFormatter()
	h := instance.provideHints(event, source)
	content, err := f.Format(event, source.GetProvider(), h)
	if err != nil {
		if v := instance.OnFormatError; v != nil {
			v(instance, out, err)
		} else {
			panic(fmt.Errorf("cannot format event %v: %w", event, err))
		}
	}

	_, _ = out.Write(content)

	_ = instance.onAfterLog(event, source)
}

func (instance *Writer) initIfRequired() {
	if instance.colorSupported == nil {
		out, supported, err := color.DetectSupportForWriter(instance.out)
		if err != nil {
			if v := instance.OnColorInitializationError; v != nil {
				v(instance, instance.GetOut(), err)
			}
		}
		instance.out = out
		instance.colorSupported = &supported
	}
}

// GetOut returns the actual io.Writer where the output will
// be written to.
func (instance *Writer) GetOut() io.Writer {
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

// GetFormatter implements formatter.Aware
func (instance *Writer) GetFormatter() formatter.Formatter {
	if v := instance.Formatter; v != nil {
		return v
	}
	if v := formatter.Default; v != nil {
		return v
	}
	return formatter.Noop()
}

// SetFormatter implements formatter.MutableAware
func (instance *Writer) SetFormatter(v formatter.Formatter) {
	instance.Formatter = v
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
