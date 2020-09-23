package location

import (
	"fmt"
	log "github.com/echocat/slf4g"
	"os"
	"path"
	"runtime"
	"strconv"
)

type CallerAwareMode uint8
type CallerAwareDetail uint8

const (
	CallerAwareModePreferFile     CallerAwareMode = 0
	CallerAwareModePreferFunction CallerAwareMode = 1

	CallerAwareDetailSimplified CallerAwareDetail = 0
	CallerAwareDetailDetailed   CallerAwareDetail = 1
)

func NewCallerAwareFactory(mode CallerAwareMode, detail CallerAwareDetail) Factory {
	return func(event log.Event, callDepth int) Location {
		if event == nil {
			return nil
		}

		pcs := make([]uintptr, 2)
		depth := runtime.Callers(callDepth+2, pcs)
		frames := runtime.CallersFrames(pcs[:depth])

		frame, _ := frames.Next()

		doDebugCallerAwareFactory(callDepth)

		if frame.Line < 0 {
			frame.Line = 0
		}

		return &CallerAwareImpl{
			Mode:     mode,
			Detail:   detail,
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

	Mode   CallerAwareMode
	Detail CallerAwareDetail
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
	file := instance.File
	if instance.Detail == CallerAwareDetailSimplified {
		file = path.Base(file)
	}
	return file + ":" + strconv.Itoa(instance.Line)
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
