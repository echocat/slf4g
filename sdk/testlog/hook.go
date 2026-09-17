package testlog

import (
	"testing"

	"github.com/echocat/slf4g/internal/globalprovider"
)

// Hook creates and registers for the given *testing.T, *testing.B or *testing.F
// a new instance of a log.Logger / log.Provider.
//
// The related Provider will be automatically cleanup at the end of the related
// test run (see testing.TB#Cleanup).
//
// customizer can be used to change the behavior of the managed Provider.
//
// The method returns the related Provider instance but while the test run it
// is also available via log.GetProvider().
func Hook(tb testing.TB, customizer ...func(*Provider)) *Provider {
	provider := NewProvider(tb, customizer...)

	tb.Cleanup(globalprovider.Push(provider))

	return provider
}
