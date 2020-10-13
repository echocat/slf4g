package sdk

import (
	"testing"

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
