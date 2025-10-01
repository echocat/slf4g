package eventlog

import (
	"syscall"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/eventlog"
)

func TestFoo(t *testing.T) {
	//err := eventlog.InstallAsEventCreate("foo-bar", eventlog.Error|eventlog.Warning|eventlog.Info)
	//assert.ToBeNoError(t, err)

	l, err := eventlog.Open("Acronis Scheduler")
	assert.ToBeNoError(t, err)

	defer l.Close()

	ss := []*uint16{syscall.StringToUTF16Ptr("Foo, Bar")}

	assert.ToBeNoError(t, windows.ReportEvent(
		l.Handle,
		windows.EVENTLOG_INFORMATION_TYPE,
		0,
		666,
		0,
		1,
		0,
		&ss[0],
		nil,
	))
}
