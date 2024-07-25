package testlog

import (
	"testing"

	log "github.com/echocat/slf4g"
)

func TestHook(t *testing.T) {
	provider := Hook(t)

	log.Info("log.Info(..)")

	provider.GetRootLogger().Info("provider.GetRootLogger().Info(..)")
	provider.GetLogger("foo").Info("provider.GetLogger(foo).Info(..)")
}
