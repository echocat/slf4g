package sdk

import (
	"testing"

	log "github.com/echocat/slf4g"

	"github.com/echocat/slf4g/level"

	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/testing/recording"
)

func Test_NewLogger(t *testing.T) {
	givenCoreLogger := recording.NewCoreLogger()

	actual := NewLogger(givenCoreLogger, level.Warn)

	assert.ToBeEqual(t, &LoggerImpl{
		Delegate:   givenCoreLogger,
		PrintLevel: level.Warn,
	}, actual)
}

func ExampleNewLogger() {
	// Creates a new SDK compatible logger that uses the slf4g logger "sdk" and
	// logs everything to level.Info.
	wrapped := NewLogger(log.GetLogger("sdk"), level.Info)

	wrapped.Print("foo", "bar")
}
