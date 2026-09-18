package coexistence_test

import (
	"testing"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native"
	_ "github.com/echocat/slf4g/platform/windows/eventlog" //nolint:staticcheck // Verifies compatibility with the deprecated package.
)

func TestEventlogAndNativeCanCoexist(t *testing.T) {
	providers := log.GetAllProviders()
	if len(providers) != 1 {
		t.Fatalf("expected one registered provider, got %d", len(providers))
	}
	if providers[0] != native.DefaultProvider {
		t.Fatalf("expected native default provider, got %T", providers[0])
	}
}
