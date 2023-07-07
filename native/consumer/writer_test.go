package consumer

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/testing/recording"

	"github.com/echocat/slf4g/native/formatter"
	"github.com/echocat/slf4g/native/hints"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native/interceptor"

	"github.com/echocat/slf4g/native/color"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_NewWriter(t *testing.T) {
	givenOut := new(bytes.Buffer)
	instance := NewWriter(givenOut)

	assert.ToBeSame(t, givenOut, instance.out)
	assert.ToBeEqual(t, true, instance.Synchronized)
}

func Test_NewWriter_withCustomization(t *testing.T) {
	givenOut := new(bytes.Buffer)
	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.Synchronized = false
	})

	assert.ToBeSame(t, givenOut, instance.out)
	assert.ToBeEqual(t, false, instance.Synchronized)
}

func Test_Writer_Consume(t *testing.T) {
	givenOut := new(bytes.Buffer)
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Info, nil)

	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.Formatter = formatter.Func(func(actualEvent log.Event, actualProvider log.Provider, actualHints hints.Hints) ([]byte, error) {
			assert.ToBeSame(t, givenLogger.GetProvider(), actualProvider)
			assert.ToBeSame(t, givenEvent, actualEvent)
			return []byte("expectedResult"), nil
		})
	})

	instance.Consume(givenEvent, givenLogger)

	assert.ToBeEqual(t, "expectedResult", givenOut.String())
}

func Test_Writer_Consume_panicsOnFormatErrors(t *testing.T) {
	givenOut := new(bytes.Buffer)
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Info, nil)
	givenError := errors.New("expected")

	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.Formatter = formatter.Func(func(log.Event, log.Provider, hints.Hints) ([]byte, error) {
			return nil, givenError
		})
	})

	instance.Consume(givenEvent, givenLogger)

	assert.ToBeMatching(t, "^LOG_EVENT_FORMAT_ERROR \\(event: .+?, error: expected\\)$", givenOut.String())
}

func Test_Writer_Consume_callsHookOnFormatErrors(t *testing.T) {
	givenOut := new(bytes.Buffer)
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Info, nil)
	givenError := errors.New("expected")

	var instance *Writer
	instance = NewWriter(givenOut, func(writer *Writer) {
		writer.OnFormatError = func(actualInstance *Writer, actualOut io.Writer, actualErr error) {
			assert.ToBeSame(t, instance, actualInstance)
			assert.ToBeSame(t, givenOut, actualOut)
			assert.ToBeSame(t, givenError, actualErr)
		}
		writer.Formatter = formatter.Func(func(log.Event, log.Provider, hints.Hints) ([]byte, error) {
			return nil, givenError
		})
	})

	instance.Consume(givenEvent, givenLogger)

	assert.ToBeEqual(t, "", givenOut.String())
}

func Test_Writer_Consume_doNothingOnNilEvent(t *testing.T) {
	givenOut := new(bytes.Buffer)
	givenLogger := recording.NewLogger()
	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.Formatter = formatter.Func(func(log.Event, log.Provider, hints.Hints) ([]byte, error) {
			panic("should not be called")
		})
	})

	instance.Consume(nil, givenLogger)

	assert.ToBeEqual(t, "", givenOut.String())
}

func Test_Writer_Consume_doNothingOnNilOut(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Info, nil)
	instance := NewWriter(nil, func(writer *Writer) {
		writer.Interceptor = interceptor.OnBeforeLogFunc(func(log.Event, log.Provider) (intercepted log.Event) {
			panic("should not be called")
		})
		writer.Formatter = formatter.Func(func(log.Event, log.Provider, hints.Hints) ([]byte, error) {
			panic("should not be called")
		})
	})

	instance.Consume(givenEvent, givenLogger)
}

func Test_Writer_Consume_doNothingOnDisabledLevel(t *testing.T) {
	givenOut := new(bytes.Buffer)
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(0, nil)

	interceptorCalled := false
	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.Interceptor = interceptor.OnBeforeLogFunc(func(actualEvent log.Event, actualProvider log.Provider) (intercepted log.Event) {
			assert.ToBeSame(t, givenLogger.GetProvider(), actualProvider)
			assert.ToBeSame(t, givenEvent, actualEvent)
			interceptorCalled = true
			return givenEvent
		})
		writer.Formatter = formatter.Func(func(log.Event, log.Provider, hints.Hints) ([]byte, error) {
			panic("should not be called")
		})
	})

	instance.Consume(givenEvent, givenLogger)

	assert.ToBeEqual(t, true, interceptorCalled)
	assert.ToBeEqual(t, "", givenOut.String())
}

func Test_Writer_Consume_initIfRequired(t *testing.T) {
	givenOut := new(bytes.Buffer)
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Info, nil)
	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.colorSupported = nil
	})

	assert.ToBeNil(t, instance.colorSupported)

	instance.Consume(givenEvent, givenLogger)

	assert.ToBeNotNil(t, instance.colorSupported)
}

func Test_Writer_Consume_beforeLog_continues(t *testing.T) {
	givenOut := new(bytes.Buffer)
	givenLogger := recording.NewLogger()
	givenEvent1 := givenLogger.NewEvent(level.Info, nil)
	givenEvent2 := givenLogger.NewEvent(level.Warn, nil)

	interceptorCalled := false
	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.Interceptor = interceptor.OnBeforeLogFunc(func(actualEvent log.Event, actualProvider log.Provider) (intercepted log.Event) {
			assert.ToBeSame(t, givenLogger.GetProvider(), actualProvider)
			assert.ToBeSame(t, givenEvent1, actualEvent)
			interceptorCalled = true
			return givenEvent2
		})
		writer.Formatter = formatter.Func(func(actualEvent log.Event, actualProvider log.Provider, actualHints hints.Hints) ([]byte, error) {
			assert.ToBeSame(t, givenLogger.GetProvider(), actualProvider)
			assert.ToBeSame(t, givenEvent2, actualEvent)
			return []byte("expectedResult"), nil
		})
	})

	instance.Consume(givenEvent1, givenLogger)

	assert.ToBeEqual(t, true, interceptorCalled)
	assert.ToBeEqual(t, "expectedResult", givenOut.String())
}

func Test_Writer_Consume_beforeLog_stops(t *testing.T) {
	givenOut := new(bytes.Buffer)
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Info, nil)

	interceptorCalled := false
	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.Formatter = formatter.Func(func(log.Event, log.Provider, hints.Hints) ([]byte, error) {
			panic("should not be called")
		})
		writer.Interceptor = interceptor.OnBeforeLogFunc(func(actualEvent log.Event, actualProvider log.Provider) (intercepted log.Event) {
			assert.ToBeSame(t, givenLogger.GetProvider(), actualProvider)
			assert.ToBeSame(t, givenEvent, actualEvent)
			interceptorCalled = true
			return nil
		})
	})

	instance.Consume(givenEvent, givenLogger)

	assert.ToBeEqual(t, true, interceptorCalled)
	assert.ToBeEqual(t, "", givenOut.String())
}

func Test_Writer_Consume_afterLog(t *testing.T) {
	givenOut := new(bytes.Buffer)
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Info, nil)

	interceptorCalled := false
	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.Interceptor = interceptor.OnAfterLogFunc(func(actualEvent log.Event, actualProvider log.Provider) bool {
			assert.ToBeSame(t, givenLogger.GetProvider(), actualProvider)
			assert.ToBeSame(t, givenEvent, actualEvent)
			interceptorCalled = true
			return true
		})
		writer.Formatter = formatter.Func(func(actualEvent log.Event, actualProvider log.Provider, actualHints hints.Hints) ([]byte, error) {
			assert.ToBeSame(t, givenLogger.GetProvider(), actualProvider)
			assert.ToBeSame(t, givenEvent, actualEvent)
			return []byte("expectedResult"), nil
		})
	})

	instance.Consume(givenEvent, givenLogger)

	assert.ToBeEqual(t, true, interceptorCalled)
	assert.ToBeEqual(t, "expectedResult", givenOut.String())
}

func Test_Writer_initIfRequired(t *testing.T) {
	old := color.SupportAssumptionDetections
	defer func() {
		color.SupportAssumptionDetections = old
	}()
	color.SupportAssumptionDetections = []color.SupportAssumptionDetection{func() (bool, error) {
		return true, nil
	}}

	givenOut := new(bytes.Buffer)
	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.colorSupported = nil
	})

	assert.ToBeNil(t, instance.colorSupported)
	instance.initIfRequired()

	assert.ToBeNotNil(t, instance.colorSupported)
	assert.ToBeEqual(t, color.SupportedAssumed, *instance.colorSupported)
}

func Test_Writer_initIfRequired_withErrors(t *testing.T) {
	old := color.SupportAssumptionDetections
	defer func() {
		color.SupportAssumptionDetections = old
	}()

	givenError := errors.New("foo")
	color.SupportAssumptionDetections = []color.SupportAssumptionDetection{func() (bool, error) {
		return false, givenError
	}}

	givenOut := new(bytes.Buffer)
	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.colorSupported = nil
	})

	assert.ToBeNil(t, instance.colorSupported)
	instance.initIfRequired()

	assert.ToBeNotNil(t, instance.colorSupported)
	assert.ToBeEqual(t, color.SupportedNone, *instance.colorSupported)
	assert.ToBeEqual(t, "", givenOut.String())
}

func Test_Writer_initIfRequired_onErrorsHandles(t *testing.T) {
	old := color.SupportAssumptionDetections
	defer func() {
		color.SupportAssumptionDetections = old
	}()

	givenError := errors.New("foo")
	color.SupportAssumptionDetections = []color.SupportAssumptionDetection{func() (bool, error) {
		return false, givenError
	}}

	givenOut := new(bytes.Buffer)
	hookCalled := false
	var instance *Writer
	instance = NewWriter(givenOut, func(writer *Writer) {
		writer.colorSupported = nil
		writer.OnColorInitializationError = func(actualInstance *Writer, actualOut io.Writer, actualErr error) {
			assert.ToBeSame(t, instance, actualInstance)
			assert.ToBeSame(t, givenOut, actualOut)
			assert.ToBeSame(t, givenError, actualErr)
			hookCalled = true
		}
	})

	assert.ToBeNil(t, instance.colorSupported)
	instance.initIfRequired()

	assert.ToBeNotNil(t, instance.colorSupported)
	assert.ToBeEqual(t, color.SupportedNone, *instance.colorSupported)
	assert.ToBeEqual(t, true, hookCalled)
	assert.ToBeEqual(t, "", givenOut.String())
}

func Test_Writer_getOut(t *testing.T) {
	givenOut := new(bytes.Buffer)
	instance := NewWriter(givenOut)

	actual := instance.GetOut()

	assert.ToBeSame(t, givenOut, actual)
}

func Test_Writer_getInterceptor_explicit(t *testing.T) {
	givenOut := new(bytes.Buffer)
	givenInterceptor := interceptor.OnBeforeLogFunc(func(event log.Event, provider log.Provider) (intercepted log.Event) {
		return event
	})
	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.Interceptor = givenInterceptor
	})

	actual := instance.getInterceptor()

	assert.ToBeSame(t, givenInterceptor, actual)
}

func Test_Writer_getInterceptor_default(t *testing.T) {
	old := interceptor.Default
	defer func() {
		interceptor.Default = old
	}()
	interceptor.Default = []interceptor.Interceptor{interceptor.OnBeforeLogFunc(func(event log.Event, provider log.Provider) (intercepted log.Event) {
		return event
	})}

	givenOut := new(bytes.Buffer)
	instance := NewWriter(givenOut)

	actual := instance.getInterceptor()

	assert.ToBeEqual(t, interceptor.Default, actual)
}

func Test_Writer_getInterceptor_noop(t *testing.T) {
	old := interceptor.Default
	defer func() {
		interceptor.Default = old
	}()
	interceptor.Default = nil

	givenOut := new(bytes.Buffer)
	instance := NewWriter(givenOut)

	actual := instance.getInterceptor()

	assert.ToBeEqual(t, interceptor.Noop(), actual)
}

func Test_Writer_SetFormatter(t *testing.T) {
	givenOut := new(bytes.Buffer)
	givenFormatter := formatter.Func(func(event log.Event, provider log.Provider, hints hints.Hints) ([]byte, error) {
		panic("should not be called")
	})
	instance := NewWriter(givenOut)

	instance.SetFormatter(givenFormatter)

	assert.ToBeSame(t, givenFormatter, instance.Formatter)
}

func Test_Writer_GetFormatter_explicit(t *testing.T) {
	givenOut := new(bytes.Buffer)
	givenFormatter := formatter.Func(func(event log.Event, provider log.Provider, hints hints.Hints) ([]byte, error) {
		panic("should not be called")
	})
	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.Formatter = givenFormatter
	})

	actual := instance.GetFormatter()

	assert.ToBeSame(t, givenFormatter, actual)
}

func Test_Writer_GetFormatter_default(t *testing.T) {
	old := formatter.Default
	defer func() {
		formatter.Default = old
	}()
	formatter.Default = formatter.Func(func(event log.Event, provider log.Provider, hints hints.Hints) ([]byte, error) {
		panic("should not be called")
	})

	givenOut := new(bytes.Buffer)
	instance := NewWriter(givenOut)

	actual := instance.GetFormatter()

	assert.ToBeEqual(t, formatter.Default, actual)
}

func Test_Writer_GetFormatter_noop(t *testing.T) {
	old := formatter.Default
	defer func() {
		formatter.Default = old
	}()
	formatter.Default = nil

	givenOut := new(bytes.Buffer)
	instance := NewWriter(givenOut)

	actual := instance.GetFormatter()

	assert.ToBeEqual(t, formatter.Noop(), actual)
}

func Test_Writer_provideHints_explicit(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Info, nil)
	givenOut := new(bytes.Buffer)
	givenHints := &struct{}{}
	givenProvider := func(actualEvent log.Event, actualLogger log.CoreLogger) hints.Hints {
		assert.ToBeSame(t, givenEvent, actualEvent)
		assert.ToBeSame(t, givenLogger, actualLogger)
		return givenHints
	}
	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.HintsProvider = givenProvider
	})

	actual := instance.provideHints(givenEvent, givenLogger)

	assert.ToBeSame(t, givenHints, actual)
}

func Test_Writer_provideHints_default(t *testing.T) {
	givenLogger := recording.NewLogger()
	givenEvent := givenLogger.NewEvent(level.Info, nil)
	givenOut := new(bytes.Buffer)
	givenSupport := color.Supported(66)
	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.colorSupported = &givenSupport
	})

	actual := instance.provideHints(givenEvent, givenLogger)

	assert.ToBeOfType(t, &writingConsumerHints{}, actual)
	assert.ToBeSame(t, instance, actual.(*writingConsumerHints).Writer)
	assert.ToBeEqual(t, givenSupport, actual.(*writingConsumerHints).IsColorSupported())
}
