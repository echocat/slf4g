package sdk

import (
	sdklog "log"
	"os"
	"testing"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/testing/recording"
)

func ExampleConfigure() {
	// Configures the whole application with to use the ROOT Logger and logs
	// everything to level.Info.
	Configure()
}

func Test_Configure(t *testing.T) {
	provider := recording.NewProvider()
	defer provider.HookGlobally()()
	defer func() {
		sdklog.SetFlags(sdklog.LstdFlags)
		sdklog.SetPrefix("")
		sdklog.SetOutput(os.Stderr)
	}()
	logger := provider.GetRootLogger()

	ConfigureWith(logger, level.Warn)
	sdklog.Print()
	sdklog.Printf("a%d%s", 2, "c")
	sdklog.Println("a", 3, "c")

	assert.ToBeEqual(t, 3, provider.Len())
	assert.ToBeEqual(t, true, provider.MustContains(
		logger.NewEvent(level.Warn, nil).
			With("message", "\n"),
	))
	assert.ToBeEqual(t, true, provider.MustContains(
		logger.NewEvent(level.Warn, nil).
			With("message", "a2c\n"),
	))
	assert.ToBeEqual(t, true, provider.MustContains(
		logger.NewEvent(level.Warn, nil).
			With("message", "a 3 c\n"),
	))
}

func ExampleConfigureWith() {
	// Configures the whole application with to use the logger named "sdk" and
	// logs everything to level.Debug.
	ConfigureWith(log.GetLogger("sdk"), level.Debug)
}

func Test_ConfigureWith(t *testing.T) {
	defer func() {
		sdklog.SetFlags(sdklog.LstdFlags)
		sdklog.SetPrefix("")
		sdklog.SetOutput(os.Stderr)
	}()

	logger := recording.NewLogger()

	ConfigureWith(logger, level.Warn)
	sdklog.Print()
	sdklog.Printf("a%d%s", 2, "c")
	sdklog.Println("a", 3, "c")

	assert.ToBeEqual(t, 3, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(
		logger.NewEvent(level.Warn, nil).
			With("message", "\n"),
	))
	assert.ToBeEqual(t, true, logger.MustContains(
		logger.NewEvent(level.Warn, nil).
			With("message", "a2c\n"),
	))
	assert.ToBeEqual(t, true, logger.MustContains(
		logger.NewEvent(level.Warn, nil).
			With("message", "a 3 c\n"),
	))
}

func ExampleNewWrapper() {
	// Creates a new SDK logger that uses the slf4g logger "sdk" and logs
	// everything to level.Info.
	wrapped := NewWrapper(log.GetLogger("sdk"), level.Info)

	wrapped.Print("foo", "bar")
}

func Test_NewWrapper(t *testing.T) {
	logger := recording.NewLogger()

	wrapped := NewWrapper(logger, level.Warn)
	wrapped.Print()
	wrapped.Printf("a%d%s", 2, "c")
	wrapped.Println("a", 3, "c")

	assert.ToBeEqual(t, 3, logger.Len())
	assert.ToBeEqual(t, true, logger.MustContains(
		logger.NewEvent(level.Warn, nil).
			With("message", "\n"),
	))
	assert.ToBeEqual(t, true, logger.MustContains(
		logger.NewEvent(level.Warn, nil).
			With("message", "a2c\n"),
	))
	assert.ToBeEqual(t, true, logger.MustContains(
		logger.NewEvent(level.Warn, nil).
			With("message", "a 3 c\n"),
	))
}
