//go:build go1.21

package sdk

import (
	"reflect"
	"runtime"
	"strings"
)

// DetectSkipFrames defines a function handler to detect how many frames should be
// skipped while creating the log message.
type DetectSkipFrames func(skip uint16) uint16

// DefaultDetectSkipFrames is the default setting for DetectSkipFrames.
//
// By default, it ignores several relevant packages of the SDK and this package.
var DefaultDetectSkipFrames DetectSkipFrames = detectSkipFramesFromSdk

// Function values are not comparable. Comparing the entry point of this named,
// non-closure function preserves direct assignment compatibility of the public default.
var detectSkipFramesFromSdkPointer = reflect.ValueOf(detectSkipFramesFromSdk).Pointer()

func isBuiltInDetectSkipFrames(candidate DetectSkipFrames) bool {
	return candidate != nil && reflect.ValueOf(candidate).Pointer() == detectSkipFramesFromSdkPointer
}

func detectSkipFramesFromSdk(skip uint16) uint16 {
	pcs := make([]uintptr, 64)
	n := runtime.Callers(int(2+skip), pcs)
	if n == 0 {
		return 0
	}
	frames := runtime.CallersFrames(pcs[:n])

	skipped := skip
	for {
		f, more := frames.Next()
		pkg := packageOf(f.Function)
		if !isIgnoredPackage(pkg) {
			return skipped
		}
		skipped++
		if !more {
			return skipped
		}
	}
}

func packageOf(funcName string) string {
	if funcName == "" {
		return ""
	}
	lastSlash := strings.LastIndex(funcName, "/")
	start := 0
	if lastSlash >= 0 {
		start = lastSlash + 1
	}
	packagePart, _, _ := strings.Cut(funcName[start:], ".")
	return funcName[:start+len(packagePart)]
}

var ignoredPackages = map[string]struct{}{
	"log":       {},
	"log/slog":  {},
	"testing":   {},
	"io":        {},
	"os":        {},
	"io/ioutil": {},
	"bytes":     {},
	"bufio":     {},
	"strings":   {},
	"github.com/echocat/slf4g/sdk/bridge/slog": {},
}

func isIgnoredPackage(pkg string) bool {
	_, ok := ignoredPackages[pkg]
	return ok
}
