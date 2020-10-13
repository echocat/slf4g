package names

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// FullLoggerNameGenerator creates a meaningful name for loggers with the full
// name out of given objects.
var FullLoggerNameGenerator = defaultFullLoggerNameGenerator

// CurrentPackageLoggerNameGenerator creates a meaningful name for loggers with
// the package of the caller who calls this method (respecting the frameToSkip).
var CurrentPackageLoggerNameGenerator = defaultCurrentPackageLoggerNameGenerator

func defaultFullLoggerNameGenerator(something interface{}) (result string) {
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
	if t != nil {
		result = t.PkgPath() + "." + t.Name()
	}
	if result == "" || result[0] == '.' {
		panic(fmt.Sprintf("invalid value to receive a package from: %+v", something))
	}
	return result
}

func defaultCurrentPackageLoggerNameGenerator(framesToSkip int) string {
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
