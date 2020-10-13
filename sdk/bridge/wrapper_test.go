package sdk

import (
	stdlog "log"
	"os"
	"testing"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/testing/recording"
)

func Test_Configure(t *testing.T) {
	provider := recording.NewProvider()
	defer provider.HookGlobally()()
	defer func() {
		stdlog.SetFlags(stdlog.LstdFlags)
		stdlog.SetPrefix("")
		stdlog.SetOutput(os.Stderr)
	}()
	logger := provider.GetRootLogger()

	ConfigureWith(logger, level.Warn)
	stdlog.Print()
	stdlog.Printf("a%d%s", 2, "c")
	stdlog.Println("a", 3, "c")

	assert.ToBeEqual(t, 3, provider.Len())
	assert.ToBeEqual(t, true, provider.MustContains(
		log.NewEvent(provider, level.Warn, 4).
			With("message", "\n"),
	))
	assert.ToBeEqual(t, true, provider.MustContains(
		log.NewEvent(provider, level.Warn, 4).
			With("message", "a2c\n"),
	))
	assert.ToBeEqual(t, true, provider.MustContains(
		log.NewEvent(provider, level.Warn, 4).
			With("message", "a 3 c\n"),
	))
}

func Test_ConfigureWith(t *testing.T) {
	defer func() {
		stdlog.SetFlags(stdlog.LstdFlags)
		stdlog.SetPrefix("")
		stdlog.SetOutput(os.Stderr)
	}()

	logger := recording.NewLogger()

	ConfigureWith(logger, level.Warn)
	stdlog.Print()
	stdlog.Printf("a%d%s", 2, "c")
	stdlog.Println("a", 3, "c")

	assert.ToBeEqual(t, 3, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(
		log.NewEvent(logger.GetProvider(), level.Warn, 4).
			With("message", "\n"),
	))
	assert.ToBeEqual(t, true, logger.MustContains(
		log.NewEvent(logger.GetProvider(), level.Warn, 4).
			With("message", "a2c\n"),
	))
	assert.ToBeEqual(t, true, logger.MustContains(
		log.NewEvent(logger.GetProvider(), level.Warn, 4).
			With("message", "a 3 c\n"),
	))
}

func Test_NewWrapper(t *testing.T) {
	logger := recording.NewLogger()

	wrapped := NewWrapper(logger, level.Warn)
	wrapped.Print()
	wrapped.Printf("a%d%s", 2, "c")
	wrapped.Println("a", 3, "c")

	assert.ToBeEqual(t, 3, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(
		log.NewEvent(logger.GetProvider(), level.Warn, 4).
			With("message", "\n"),
	))
	assert.ToBeEqual(t, true, logger.MustContains(
		log.NewEvent(logger.GetProvider(), level.Warn, 4).
			With("message", "a2c\n"),
	))
	assert.ToBeEqual(t, true, logger.MustContains(
		log.NewEvent(logger.GetProvider(), level.Warn, 4).
			With("message", "a 3 c\n"),
	))
}
