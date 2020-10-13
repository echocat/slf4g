package names

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// FullLoggerNameGenerator creates a meaningful name for loggers with the full
// name out of given objects.
func FullLoggerNameGenerator(something interface{}) string {
	if f := FullLoggerNameCustomizer; f != nil {
		return f(something)
	}
	switch v := something.(type) {
	case string:
		return v
	case *string:
		return *v
	}

	t := reflect.TypeOf(something)
	for t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	var result string
	if t != nil {
		result = t.PkgPath() + "." + t.Name()
	}
	if result == "" || result[0] == '.' {
		panic(fmt.Sprintf("invalid value to receive a package from: %+v", something))
	}
	return result
}

// CurrentPackageLoggerNameGenerator creates a meaningful name for loggers with
// the package of the caller who calls this method (respecting the frameToSkip).
func CurrentPackageLoggerNameGenerator(framesToSkip int) string {
	if f := CurrentPackageLoggerNameCustomizer; f != nil {
		return f(framesToSkip + 1)
	}
	pcs := make([]uintptr, 3)
	depth := runtime.Callers(framesToSkip+2, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	frame, _ := frames.Next()

	if frame.Function == "" {
		panic("cannot capture a valid position for the package name.")
	}

	allParts := strings.Split(frame.Function, "/")
	lastPart := allParts[len(allParts)-1]

	lastSubParts := strings.SplitN(lastPart, ".", 2)

	allParts[len(allParts)-1] = lastSubParts[0]

	result := strings.Join(allParts, "/")

	return result
}

// FullLoggerNameCustomizer will override the default behavior of
// FullLoggerNameGenerator if set.
var FullLoggerNameCustomizer func(something interface{}) string

// CurrentPackageLoggerNameCustomizer will override the default behavior of
// CurrentPackageLoggerNameGenerator if set.
var CurrentPackageLoggerNameCustomizer func(framesToSkip int) string
