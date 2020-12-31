package consumer

import (
	"bytes"
	"errors"
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

func Test_Writer_initIfRequired_withErrorsAndPrintThem(t *testing.T) {
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
		writer.PrintErrorOnColorInitialization = true
	})

	assert.ToBeNil(t, instance.colorSupported)
	instance.initIfRequired()

	assert.ToBeNotNil(t, instance.colorSupported)
	assert.ToBeEqual(t, color.SupportedNone, *instance.colorSupported)
	assert.ToBeEqual(t, "WARNING!!! Cannot initiate colors for target: foo; falling back to no color support.\n", givenOut.String())
}

func Test_Writer_getOut(t *testing.T) {
	givenOut := new(bytes.Buffer)
	instance := NewWriter(givenOut)

	actual := instance.getOut()

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

func Test_Writer_getFormatter_explicit(t *testing.T) {
	givenOut := new(bytes.Buffer)
	givenFormatter := formatter.Func(func(event log.Event, provider log.Provider, hints hints.Hints) ([]byte, error) {
		panic("should not be called")
	})
	instance := NewWriter(givenOut, func(writer *Writer) {
		writer.Formatter = givenFormatter
	})

	actual := instance.getFormatter()

	assert.ToBeSame(t, givenFormatter, actual)
}

func Test_Writer_getFormatter_default(t *testing.T) {
	old := formatter.Default
	defer func() {
		formatter.Default = old
	}()
	formatter.Default = formatter.Func(func(event log.Event, provider log.Provider, hints hints.Hints) ([]byte, error) {
		panic("should not be called")
	})

	givenOut := new(bytes.Buffer)
	instance := NewWriter(givenOut)

	actual := instance.getFormatter()

	assert.ToBeEqual(t, formatter.Default, actual)
}

func Test_Writer_getFormatter_noop(t *testing.T) {
	old := formatter.Default
	defer func() {
		formatter.Default = old
	}()
	formatter.Default = nil

	givenOut := new(bytes.Buffer)
	instance := NewWriter(givenOut)

	actual := instance.getFormatter()

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
