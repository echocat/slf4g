package location

import (
	"runtime"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func someFuncForCallerDiscoveryDiscoverTest(instance *CallerDiscovery, skipFrames uint16) Location {
	// Leave it here because the line where this method was called is important for the test!
	return instance.DiscoverLocation(nil, skipFrames)
}

func Test_CallerDiscovery_Discover(t *testing.T) {
	// Leave it here because the line where this method was called is important for the test!
	instance := NewCallerDiscovery()

	actual1 := someFuncForCallerDiscoveryDiscoverTest(instance, 0)
	assert.ToBeOfType(t, &callerImpl{}, actual1)
	assert.ToBeSame(t, instance, actual1.(*callerImpl).discovery)
	assert.ToBeEqual(t, "github.com/echocat/slf4g/native/location.someFuncForCallerDiscoveryDiscoverTest", actual1.(*callerImpl).frame.Function)
	assert.ToBeEqual(t, 12, actual1.(*callerImpl).frame.Line)

	actual2 := someFuncForCallerDiscoveryDiscoverTest(instance, 1)
	assert.ToBeEqual(t, "github.com/echocat/slf4g/native/location.Test_CallerDiscovery_Discover", actual2.(*callerImpl).frame.Function)
	assert.ToBeEqual(t, 25, actual2.(*callerImpl).frame.Line)

	actual3 := someFuncForCallerDiscoveryDiscoverTest(instance, 255)
	assert.ToBeEqual(t, "???", actual3.(*callerImpl).frame.Function)
	assert.ToBeEqual(t, "???", actual3.(*callerImpl).frame.File)
	assert.ToBeEqual(t, 0, actual3.(*callerImpl).frame.Line)
}

func Test_NewCallerDiscovery(t *testing.T) {
	actual := NewCallerDiscovery(func(actualDiscovery *CallerDiscovery) {
		assert.ToBeEqual(t, CallerReportingTypePrefersType, actualDiscovery.ReportingType)
		assert.ToBeEqual(t, CallerReportingDetailSimplified, actualDiscovery.ReportingDetail)

		// Customize!
		actualDiscovery.ReportingType = 123
	})

	assert.ToBeEqual(t, CallerReportingType(123), actual.ReportingType)
	assert.ToBeEqual(t, CallerReportingDetailSimplified, actual.ReportingDetail)
}

func Test_callerImpl_formatFile_simplified(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingDetail = CallerReportingDetailSimplified
		}),
		frame: &runtime.Frame{
			File: "github.com/foo/bar/aFile.go",
			Line: 123,
		},
	}

	actual1 := instance.formatFile()
	assert.ToBeEqual(t, "aFile.go:123", actual1)
}

func Test_callerImpl_formatFile_simplified_withoutLine(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingDetail = CallerReportingDetailSimplified
		}),
		frame: &runtime.Frame{
			File: "github.com/foo/bar/aFile.go",
			Line: 0,
		},
	}

	actual1 := instance.formatFile()
	assert.ToBeEqual(t, "aFile.go", actual1)
}

func Test_callerImpl_formatFile_detailed(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingDetail = CallerReportingDetailDetailed
		}),
		frame: &runtime.Frame{
			File: "github.com/foo/bar/aFile.go",
			Line: 123,
		},
	}

	actual1 := instance.formatFile()
	assert.ToBeEqual(t, "github.com/foo/bar/aFile.go:123", actual1)
}

func Test_callerImpl_formatFile_detailed_withoutLine(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingDetail = CallerReportingDetailDetailed
		}),
		frame: &runtime.Frame{
			File: "github.com/foo/bar/aFile.go",
			Line: 0,
		},
	}

	actual1 := instance.formatFile()
	assert.ToBeEqual(t, "github.com/foo/bar/aFile.go", actual1)
}

func Test_callerImpl_formatType_withType_detailed(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingDetail = CallerReportingDetailDetailed
		}),
		frame: &runtime.Frame{
			Function: "github.com/foo/bar/aPackage.aType.aFunc",
			Line:     123,
		},
	}

	actual1 := instance.formatType()
	assert.ToBeEqual(t, "github.com/foo/bar/aPackage.aType.aFunc:123", actual1)
}

func Test_callerImpl_formatType_withType_detailed_withoutLine(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingDetail = CallerReportingDetailDetailed
		}),
		frame: &runtime.Frame{
			Function: "github.com/foo/bar/aPackage.aType.aFunc",
			Line:     0,
		},
	}

	actual1 := instance.formatType()
	assert.ToBeEqual(t, "github.com/foo/bar/aPackage.aType.aFunc", actual1)
}

func Test_callerImpl_formatType_withoutType_detailed(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingDetail = CallerReportingDetailDetailed
		}),
		frame: &runtime.Frame{
			Function: "github.com/foo/bar/aPackage.aFunc",
			Line:     123,
		},
	}

	actual1 := instance.formatType()
	assert.ToBeEqual(t, "github.com/foo/bar/aPackage.aFunc:123", actual1)
}

func Test_callerImpl_formatType_withoutType_detailed_withoutLine(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingDetail = CallerReportingDetailDetailed
		}),
		frame: &runtime.Frame{
			Function: "github.com/foo/bar/aPackage.aFunc",
			Line:     0,
		},
	}

	actual1 := instance.formatType()
	assert.ToBeEqual(t, "github.com/foo/bar/aPackage.aFunc", actual1)
}

func Test_callerImpl_formatType_withType_simplified(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingDetail = CallerReportingDetailSimplified
		}),
		frame: &runtime.Frame{
			Function: "github.com/foo/bar/aPackage.aType.aFunc",
			Line:     123,
		},
	}

	actual1 := instance.formatType()
	assert.ToBeEqual(t, "aPackage.aType.aFunc:123", actual1)
}

func Test_callerImpl_formatType_withType_simplified_withoutLine(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingDetail = CallerReportingDetailSimplified
		}),
		frame: &runtime.Frame{
			Function: "github.com/foo/bar/aPackage.aType.aFunc",
			Line:     0,
		},
	}

	actual1 := instance.formatType()
	assert.ToBeEqual(t, "aPackage.aType.aFunc", actual1)
}

func Test_callerImpl_formatType_withoutType_simplified(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingDetail = CallerReportingDetailSimplified
		}),
		frame: &runtime.Frame{
			Function: "github.com/foo/bar/aPackage.aFunc",
			Line:     123,
		},
	}

	actual1 := instance.formatType()
	assert.ToBeEqual(t, "aPackage.aFunc:123", actual1)
}

func Test_callerImpl_formatType_withoutType_simplified_withoutLine(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingDetail = CallerReportingDetailSimplified
		}),
		frame: &runtime.Frame{
			Function: "github.com/foo/bar/aPackage.aFunc",
			Line:     0,
		},
	}

	actual1 := instance.formatType()
	assert.ToBeEqual(t, "aPackage.aFunc", actual1)
}

func Test_callerImpl_GetFrame(t *testing.T) {
	givenFrame := &runtime.Frame{
		Function: "github.com/foo/bar/aPackage.aFunc",
		Line:     0,
	}
	instance := callerImpl{
		frame: givenFrame,
	}

	actual := instance.GetFrame()

	assert.ToBeEqual(t, *givenFrame, actual)
}

func Test_callerImpl_Get_prefersType_withType(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingType = CallerReportingTypePrefersType
		}),
		frame: &runtime.Frame{
			Function: "github.com/foo/bar/aPackage.aType.aFunc",
			File:     "github.com/foo/bar/aFile.go",
			Line:     123,
		},
	}

	actual1 := instance.Get()
	assert.ToBeEqual(t, "aPackage.aType.aFunc:123", actual1)
}

func Test_callerImpl_Get_prefersType_withoutType(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingType = CallerReportingTypePrefersType
		}),
		frame: &runtime.Frame{
			Function: "",
			File:     "github.com/foo/bar/aFile.go",
			Line:     123,
		},
	}

	actual1 := instance.Get()
	assert.ToBeEqual(t, "aFile.go:123", actual1)
}

func Test_callerImpl_Get_prefersFile_withFile(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingType = CallerReportingTypePrefersFile
		}),
		frame: &runtime.Frame{
			Function: "github.com/foo/bar/aPackage.aType.aFunc",
			File:     "github.com/foo/bar/aFile.go",
			Line:     123,
		},
	}

	actual1 := instance.Get()
	assert.ToBeEqual(t, "aFile.go:123", actual1)
}

func Test_callerImpl_Get_prefersFile_withoutFile(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(func(discovery *CallerDiscovery) {
			discovery.ReportingType = CallerReportingTypePrefersFile
		}),
		frame: &runtime.Frame{
			Function: "github.com/foo/bar/aPackage.aType.aFunc",
			File:     "",
			Line:     123,
		},
	}

	actual1 := instance.Get()
	assert.ToBeEqual(t, "aPackage.aType.aFunc:123", actual1)
}

func Test_callerImpl_Get_withoutEverything(t *testing.T) {
	instance := callerImpl{
		discovery: NewCallerDiscovery(),
		frame: &runtime.Frame{
			Function: "",
			File:     "",
			Line:     123,
		},
	}

	actual := instance.Get()
	assert.ToBeNil(t, actual)
}
