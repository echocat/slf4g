//go:build mock
// +build mock

package color

import (
	"io"
	"os"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

// Currently we do not have perfect automated testing for this feature because
// the output will be always wrapped in tests; so currently we only do minimal
// testing here.
func Test_DetectSupportForWriter_native(t *testing.T) {
	oldAssumptions := SupportAssumptionDetections
	defer func() {
		SupportAssumptionDetections = oldAssumptions
	}()
	oldPrepare := prepareCallback
	defer func() {
		prepareCallback = oldPrepare
	}()

	SupportAssumptionDetections = nil // Do use assumptions
	prepareCallback = nil             // Do not support any callback here

	actualWriter1, actualSupported1, actualErr1 := DetectSupportForWriter(os.Stdout)
	assert.ToBeNil(t, actualErr1)
	assert.ToBeSame(t, os.Stdout, actualWriter1)
	assert.ToBeEqual(t, SupportedNone, actualSupported1)

	prepareCallback = func(w io.Writer) (bool, error) {
		return true, nil
	}

	actualWriter2, actualSupported2, actualErr2 := DetectSupportForWriter(os.Stdout)
	assert.ToBeNil(t, actualErr2)
	assert.ToBeSame(t, os.Stdout, actualWriter2)
	assert.ToBeEqual(t, SupportedNative, actualSupported2)

}
