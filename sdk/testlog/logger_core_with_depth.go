//go:build !slf4gcompat && !go1.23

package testlog

import (
	"testing"
	"unsafe"
)

type testCommonT unsafe.Pointer

//go:linkname testingCommonLogDepth testing.(*common).logDepth
//goland:noinspection GoUnusedParameter
func testingCommonLogDepth(c testCommonT, s string, depth int)

func (instance *coreLogger) logLogDepth(str string, skipFrames uint16) {
	testingCommonLogDepth(testCommonT(instance.tb.(*testing.T)), str+"\n", int(skipFrames+2))
}
