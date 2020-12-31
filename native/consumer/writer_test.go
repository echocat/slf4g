package consumer

import (
	"bytes"
	"errors"
	"testing"

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
