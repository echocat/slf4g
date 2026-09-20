package consumer

import (
	"io"
	"os"
	"sync"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native/color"
	"github.com/echocat/slf4g/native/formatter"
	"github.com/echocat/slf4g/native/hints"
	"github.com/echocat/slf4g/native/interceptor"
)

const (
	// Absorb transient contention without retaining events indefinitely behind a stalled sink.
	maxPendingWriterRequests = 1024
	formatErrorFallback      = "{\"error\":\"LOG_EVENT_FORMAT_ERROR\"}\n"
	writeErrorFallback       = "{\"error\":\"LOG_EVENT_WRITE_ERROR\"}\n"
	queueFullFallback        = "{\"error\":\"LOG_EVENT_QUEUE_FULL\"}\n"
)

// Writer is an implementation of Writer which formats the consumed log.Entry
// using a configured Formatter and logs it to the configured io.Writer.
//
// NewWriter() is used to create a new instance.
type Writer struct {
	// Formatter to format the consumed log.Event with. If nothing was provided
	// formatter.Default will be used. Set this field only before the Writer is
	// first used; use SetFormatter for runtime changes.
	Formatter formatter.Formatter

	// Interceptor can be used to intercept the consumption of an event shortly
	// before the actual consumption or directly afterward. If nothing was
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
	// additional performance costs. If an event is already being consumed,
	// concurrent or reentrant calls are queued and can return before their event
	// was written. If the bounded queue is full, new events are dropped and the
	// active caller writes one safe diagnostic to stderr per drain cycle. If
	// consumption panics or exits its goroutine, accepted requests remain queued
	// until the next call resumes draining them.
	Synchronized bool

	// OnFormatError will be called if their as any kind of error while
	// formatting an log.Event using the configured Formatter. If nothing was
	// provided these errors will result in a generic fallback message without
	// details from the event or error. A configured callback is responsible for
	// safely encoding error details before writing them to out.
	OnFormatError func(*Writer, io.Writer, error)

	// OnColorInitializationError will be called if their as any kind of error
	// while initialize the color support. If nothing was provided these errors
	// will be silently swallowed.
	OnColorInitializationError func(*Writer, io.Writer, error)

	out             io.Writer
	colorSupported  *color.Supported
	mutex           sync.Mutex
	consuming       bool
	pending         []writerRequest
	pendingOverflow bool
}

type writerRequest struct {
	event  log.Event
	source log.CoreLogger
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

	if !instance.Synchronized {
		instance.consume(event, source)
		return
	}

	// A single caller drains the queue without holding the mutex while invoking
	// callbacks. Reentrant calls can therefore enqueue safely, while formatting
	// and output remain serialized.
	instance.mutex.Lock()
	if instance.consuming {
		if len(instance.pending) >= maxPendingWriterRequests {
			instance.pendingOverflow = true
			instance.mutex.Unlock()
			return
		}
		instance.pending = append(instance.pending, writerRequest{event, source})
		instance.mutex.Unlock()
		return
	}
	request := writerRequest{event, source}
	instance.consuming = true
	if len(instance.pending) > 0 {
		request = instance.pending[0]
		instance.pending[0] = writerRequest{}
		instance.pending = append(instance.pending[1:], writerRequest{event, source})
	}
	instance.mutex.Unlock()

	instance.consumePending(request)
}

func (instance *Writer) consumePending(request writerRequest) {
	completed := false
	overflowReported := false
	defer func() {
		if !completed {
			// Do not leave the writer blocked if user-provided code panics or exits
			// its goroutine. A later Consume call resumes the accepted queue.
			instance.mutex.Lock()
			instance.consuming = false
			instance.mutex.Unlock()
		}
	}()

	for {
		instance.consume(request.event, request.source)

		instance.mutex.Lock()
		if instance.pendingOverflow && !overflowReported {
			instance.pendingOverflow = false
			instance.mutex.Unlock()
			_, _ = io.WriteString(os.Stderr, queueFullFallback)
			overflowReported = true
			instance.mutex.Lock()
		}
		if len(instance.pending) == 0 {
			instance.pending = nil
			instance.pendingOverflow = false
			instance.consuming = false
			completed = true
			instance.mutex.Unlock()
			return
		}
		request = instance.pending[0]
		instance.pending[0] = writerRequest{}
		instance.pending = instance.pending[1:]
		instance.mutex.Unlock()
	}
}

func (instance *Writer) consume(event log.Event, source log.CoreLogger) {
	out := instance.GetOut()
	if out == nil {
		return
	}
	instance.initIfRequired()
	out = instance.GetOut()

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
		content = nil
		if v := instance.OnFormatError; v != nil {
			v(instance, out, err)
		} else {
			content = []byte(formatErrorFallback)
		}
	}

	written, writeErr := out.Write(content)
	if writeErr != nil || written != len(content) {
		_, _ = io.WriteString(os.Stderr, writeErrorFallback)
	}

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
	instance.mutex.Lock()
	v := instance.Formatter
	instance.mutex.Unlock()
	if v != nil {
		return v
	}
	if v := formatter.Default; v != nil {
		return v
	}
	return formatter.Noop()
}

// SetFormatter implements formatter.MutableAware
func (instance *Writer) SetFormatter(v formatter.Formatter) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
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

func (instance *writingConsumerHints) IsColorSupported() color.Supported {
	return *instance.colorSupported
}
