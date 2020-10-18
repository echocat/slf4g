package log

import (
	"github.com/echocat/slf4g/level"
)

// LoggingWriter is used to capture lines which might contain log event and
// forward them straight to a configured Logger. This is quite useful with
// old/native logging frameworks which does not have generic hooks for log
// frameworks like slf4g.
type LoggingWriter struct {
	// Logger where to log captured events to. If this field is not set this
	// writer will simply do nothing.
	Logger CoreLogger

	// LevelExtractor is used to determine the level of the current written
	// line when reporting to configured Logger. If nil/not configured it will
	// use level.Info.
	LevelExtractor level.LineExtractor

	// SkipFrames is used to create the event with.
	SkipFrames uint16
}

// Write implements io.Writer.
func (instance *LoggingWriter) Write(p []byte) (int, error) {
	if logger := instance.Logger; logger != nil {
		provider := logger.GetProvider()

		lvl, err := instance.levelOf(p)
		if err != nil {
			return 0, err
		}

		event := instance.Logger.NewEvent(lvl, map[string]interface{}{
			provider.GetFieldKeysSpec().GetMessage(): string(p),
		})

		instance.Logger.Log(event, instance.SkipFrames+1)
	}
	return len(p), nil
}

func (instance *LoggingWriter) levelOf(p []byte) (level.Level, error) {
	if v := instance.LevelExtractor; v != nil {
		return v.ExtractLevelFromLine(p)
	}
	return level.Info, nil
}
