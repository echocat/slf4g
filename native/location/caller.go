package location

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	log "github.com/echocat/slf4g"
)

type CallerAware interface {
	Location

	GetFrame() runtime.Frame
}

type CallerAwareMode uint8

const (
	CallerAwareModePreferFunction CallerAwareMode = 0
	CallerAwareModePreferFile     CallerAwareMode = 1
)

type CallerAwareDetail uint8

const (
	CallerAwareDetailDetailed   CallerAwareDetail = 0
	CallerAwareDetailSimplified CallerAwareDetail = 1
)

func NewCallerAwareDiscovery(customizer ...func(*CallerAwareDiscovery)) *CallerAwareDiscovery {
	result := &CallerAwareDiscovery{
		Mode:   CallerAwareModePreferFunction,
		Detail: CallerAwareDetailSimplified,
	}
	for _, c := range customizer {
		c(result)
	}
	return result
}

type CallerAwareDiscovery struct {
	Mode   CallerAwareMode
	Detail CallerAwareDetail
}

func (instance *CallerAwareDiscovery) DiscoveryLocation(event log.Event, callDepth int) Location {
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

	return &callerAwareImpl{instance, &frame}
}

type callerAwareImpl struct {
	*CallerAwareDiscovery

	frame *runtime.Frame
}

func (instance callerAwareImpl) Get() interface{} {
	if instance.Mode == CallerAwareModePreferFunction && instance.frame.Function != "" {
		return instance.formatFunction()
	}
	if instance.Mode == CallerAwareModePreferFile && instance.frame.File != "" {
		return instance.formatFile()
	}
	if instance.frame.Function != "" {
		return instance.formatFunction()
	}
	if instance.frame.File != "" {
		return instance.formatFile()
	}
	return nil
}

func (instance callerAwareImpl) formatFile() string {
	file := instance.frame.File
	if instance.Detail == CallerAwareDetailSimplified {
		file = path.Base(file)
	}
	return file
}

func (instance callerAwareImpl) formatFunction() string {
	aPackage := strings.Split(instance.frame.Function, "/")
	lastPart := aPackage[len(aPackage)-1]
	lastSubParts := strings.SplitN(lastPart, ".", 3)

	var p string
	if instance.Detail == CallerAwareDetailSimplified {
		p = lastSubParts[0]
	} else {
		aPackage[len(aPackage)-1] = lastSubParts[0]
		p = strings.Join(aPackage, "/")
	}

	return p + "." + strings.Join(lastSubParts[1:], ".") + ":" + strconv.Itoa(instance.frame.Line)
}

func (instance *callerAwareImpl) GetFrame() runtime.Frame {
	return *instance.frame
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
