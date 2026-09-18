package eventlog_test

import (
	"testing"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/platform/windows/eventlog"
)

var (
	_ log.Provider   = (*eventlog.Provider)(nil)
	_ log.CoreLogger = (*eventlog.CoreLogger)(nil)
)

func TestPackageDoesNotRegisterProvider(t *testing.T) {
	if providers := log.GetAllProviders(); len(providers) != 0 {
		t.Fatalf("expected no registered providers, got %d", len(providers))
	}
}
