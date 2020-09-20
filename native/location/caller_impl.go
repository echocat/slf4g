package location

import (
	"fmt"
	log "github.com/echocat/slf4g"
	"os"
	"runtime"
)

type CallerAwareMode uint8

const (
	CallerAwareModePreferFile     CallerAwareMode = 0
	CallerAwareModePreferFunction CallerAwareMode = 1
)

func CallerAwareFactory(mode CallerAwareMode) Factory {
	return func(event log.Event, callDepth int) Location {
		if event == nil {
			return nil
		}

		pcs := make([]uintptr, 2)
		depth := runtime.Callers(callDepth+2, pcs)
		frames := runtime.CallersFrames(pcs[:depth])

		frame, _ := frames.Next()

		doDebugCallerAwareFactory(callDepth)

		return &CallerAwareImpl{
			Mode:     mode,
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		}
	}
}

var debugCallerAwareFactory = false

func doDebugCallerAwareFactory(callDepth int) {
	//goland:noinspection GoBoolExpressions
	if debugCallerAwareFactory {
		pcs := make([]uintptr, 20)
		depth := runtime.Callers(0, pcs)
		frames := runtime.CallersFrames(pcs[:depth])

		_, _ = fmt.Fprintf(os.Stderr, "callerDepth: %d\n", callDepth+2)
		i := 1
		for f, again := frames.Next(); again; f, again = frames.Next() {
			_, _ = fmt.Fprintf(os.Stderr, "\t[%d] %s:%d\n", i, f.File, f.Line)
			i++
		}
	}
}

type CallerAwareImpl struct {
	Function string
	File     string
	Line     int

	Mode CallerAwareMode
}

func (instance CallerAwareImpl) Get() interface{} {
	if instance.Mode == CallerAwareModePreferFile && instance.File != "" {
		return instance.formatFile()
	}
	if instance.Mode == CallerAwareModePreferFunction && instance.Function != "" {
		return instance.formatFunction()
	}
	if instance.File != "" {
		return instance.formatFile()
	}
	if instance.Function != "" {
		return instance.formatFunction()
	}
	return nil
}

func (instance CallerAwareImpl) formatFile() string {
	return fmt.Sprintf("%s:%d", instance.File, instance.Line)
}

func (instance CallerAwareImpl) formatFunction() string {
	return instance.Function
}

func (instance *CallerAwareImpl) GetFunction() string {
	return instance.Function
}

func (instance *CallerAwareImpl) GetFile() string {
	return instance.File
}

func (instance *CallerAwareImpl) GetLine() int {
	return instance.Line
}