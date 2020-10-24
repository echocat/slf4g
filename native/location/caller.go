package location

import (
	"path"
	"runtime"
	"strconv"
	"strings"

	log "github.com/echocat/slf4g"
)

// Caller describes the Location of the caller which leads to the initial
// creation of an log.Event.
type Caller interface {
	Location

	// GetFrame provides the runtime.Frame of the initial caller.
	GetFrame() runtime.Frame
}

// CallerReportingType defines how Caller should report the caller location, if
// requested.
type CallerReportingType uint8

const (
	// CallerReportingTypePrefersType will prefer to report the caller's type
	// (package/type/function). If this is not possible the file will be
	// reported instead.
	CallerReportingTypePrefersType CallerReportingType = 0

	// CallerReportingTypePrefersFile will prefer to report the caller's file.
	// If this is not possible the type (package/type/function) will be used
	// instead.
	CallerReportingTypePrefersFile CallerReportingType = 1
)

// CallerReportingDetail defines how detailed Caller will report the caller
// location, of requested.
type CallerReportingDetail uint8

const (
	// CallerReportingDetailDetailed will report the full available details of
	// the configured CallerReportingType.
	CallerReportingDetailDetailed CallerReportingDetail = 0

	// CallerReportingDetailSimplified will report a simplified version of all
	// available details of the configured CallerReportingType. It is focussed
	// to be still detailed enough to get a meaningful understanding about the
	// location but still keep it short.
	CallerReportingDetailSimplified CallerReportingDetail = 1
)

// CallerDiscovery implements Discovery that discovers the Caller's Location of
// a provided log.Event.
type CallerDiscovery struct {
	ReportingType   CallerReportingType
	ReportingDetail CallerReportingDetail
}

// NewCallerDiscovery create a new instance of CallerDiscovery which is ready to
// use.
func NewCallerDiscovery(customizer ...func(*CallerDiscovery)) *CallerDiscovery {
	result := &CallerDiscovery{
		ReportingType:   CallerReportingTypePrefersType,
		ReportingDetail: CallerReportingDetailSimplified,
	}
	for _, c := range customizer {
		c(result)
	}
	return result
}

// DiscoverLocation implements Discovery.DiscoverLocation().
func (instance *CallerDiscovery) DiscoverLocation(_ log.Event, skipFrames uint16) Location {
	pcs := make([]uintptr, 2)
	depth := runtime.Callers(int(skipFrames)+2, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	frame, _ := frames.Next()

	if frame.Function == "" {
		frame.Function = "???"
	}
	if frame.File == "" {
		frame.File = "???"
	}

	return &callerImpl{
		discovery: instance,
		frame:     &frame,
	}
}

type callerImpl struct {
	discovery *CallerDiscovery
	frame     *runtime.Frame
}

func (instance callerImpl) Get() interface{} {
	if instance.discovery.ReportingType == CallerReportingTypePrefersType && instance.frame.Function != "" {
		return instance.formatType()
	}
	if instance.discovery.ReportingType == CallerReportingTypePrefersFile && instance.frame.File != "" {
		return instance.formatFile()
	}
	if instance.frame.Function != "" {
		return instance.formatType()
	}
	if instance.frame.File != "" {
		return instance.formatFile()
	}
	return nil
}

func (instance callerImpl) formatFile() string {
	file := instance.frame.File
	if instance.discovery.ReportingDetail == CallerReportingDetailSimplified {
		file = path.Base(file)
	}
	result := file
	if instance.frame.Line > 0 {
		result += ":" + strconv.Itoa(instance.frame.Line)
	}
	return result
}

func (instance callerImpl) formatType() string {
	aPackage := strings.Split(instance.frame.Function, "/")
	lastPart := aPackage[len(aPackage)-1]
	lastSubParts := strings.SplitN(lastPart, ".", 3)

	var p string
	if instance.discovery.ReportingDetail == CallerReportingDetailSimplified {
		p = lastSubParts[0]
	} else {
		aPackage[len(aPackage)-1] = lastSubParts[0]
		p = strings.Join(aPackage, "/")
	}

	result := p + "." + strings.Join(lastSubParts[1:], ".")
	if instance.frame.Line > 0 {
		result += ":" + strconv.Itoa(instance.frame.Line)
	}
	return result
}

// GetFrame returns the contained frame
func (instance *callerImpl) GetFrame() runtime.Frame {
	return *instance.frame
}
