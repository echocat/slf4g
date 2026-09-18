package log

import (
	"bytes"
	"strings"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"

	"github.com/echocat/slf4g/level"
)

func Test_LoggingWriter_Write(t *testing.T) {
	givenLogger := newMockCoreLogger("foo")
	givenLogger.initLoggedEvents()
	givenLogger.level = level.Trace
	givenExtractor := level.LineExtractorFunc(func(in []byte) (level.Level, error) {
		if strings.HasPrefix(string(in), "I") {
			return level.Info, nil
		}
		if strings.HasPrefix(string(in), "E") {
			return level.Error, nil
		}
		panic("not expected")
	})
	messageKey := givenLogger.GetProvider().GetFieldKeysSpec().GetMessage()
	instance := &LoggingWriter{
		Logger:         givenLogger,
		LevelExtractor: givenExtractor,
		SkipFrames:     666,
	}

	actual1Written, actual1Err := instance.Write([]byte("I hello world!"))
	assert.ToBeNil(t, actual1Err)
	assert.ToBeEqual(t, 14, actual1Written)

	actual2Written, actual2Err := instance.Write([]byte("E hello world!"))
	assert.ToBeNil(t, actual2Err)
	assert.ToBeEqual(t, 14, actual2Written)

	assert.ToBeEqual(t, 2, len(*givenLogger.loggedEvents))

	assert.ToBeEqual(t, level.Info, givenLogger.loggedEvent(0).GetLevel())
	actualMessage0, _ := givenLogger.loggedEvent(0).Get(messageKey)
	assert.ToBeEqual(t, "I hello world!", actualMessage0)

	assert.ToBeEqual(t, level.Error, givenLogger.loggedEvent(1).GetLevel())
	actualMessage1, _ := givenLogger.loggedEvent(1).Get(messageKey)
	assert.ToBeEqual(t, "E hello world!", actualMessage1)
}

func Test_LoggingWriter_Write_withoutLevelExtractor(t *testing.T) {
	givenLogger := newMockCoreLogger("foo")
	givenLogger.initLoggedEvents()
	givenLogger.level = level.Trace
	messageKey := givenLogger.GetProvider().GetFieldKeysSpec().GetMessage()
	instance := &LoggingWriter{
		Logger:     givenLogger,
		SkipFrames: 666,
	}

	actual1Written, actual1Err := instance.Write([]byte("W hello world!"))
	assert.ToBeNil(t, actual1Err)
	assert.ToBeEqual(t, 14, actual1Written)

	actual2Written, actual2Err := instance.Write([]byte("E hello world!"))
	assert.ToBeNil(t, actual2Err)
	assert.ToBeEqual(t, 14, actual2Written)

	assert.ToBeEqual(t, 2, len(*givenLogger.loggedEvents))

	assert.ToBeEqual(t, level.Info, givenLogger.loggedEvent(0).GetLevel())
	actualMessage0, _ := givenLogger.loggedEvent(0).Get(messageKey)
	assert.ToBeEqual(t, "W hello world!", actualMessage0)

	assert.ToBeEqual(t, level.Info, givenLogger.loggedEvent(1).GetLevel())
	actualMessage1, _ := givenLogger.loggedEvent(1).Get(messageKey)
	assert.ToBeEqual(t, "E hello world!", actualMessage1)
}

func Test_LoggingWriter_Write_lineExtractorErrors(t *testing.T) {
	givenLogger := newMockCoreLogger("foo")
	givenLogger.initLoggedEvents()
	givenExtractor := level.LineExtractorFunc(func(in []byte) (level.Level, error) {
		return 0, stringError(in)
	})
	instance := &LoggingWriter{Logger: givenLogger, LevelExtractor: givenExtractor}

	actual1Written, actual1Err := instance.Write([]byte("foo"))
	assert.ToBeEqual(t, stringError("foo"), actual1Err)
	assert.ToBeEqual(t, 0, actual1Written)

	actual2Written, actual2Err := instance.Write([]byte("bar"))
	assert.ToBeEqual(t, stringError("bar"), actual2Err)
	assert.ToBeEqual(t, 0, actual2Written)

	assert.ToBeEqual(t, 0, len(*givenLogger.loggedEvents))
}

func Test_LoggingWriter_Write_withoutLogger(t *testing.T) {
	instance := &LoggingWriter{}

	actual1Written, actual1Err := instance.Write([]byte("foo"))
	assert.ToBeNil(t, actual1Err)
	assert.ToBeEqual(t, 3, actual1Written)

	actual2Written, actual2Err := instance.Write([]byte("foobar"))
	assert.ToBeNil(t, actual2Err)
	assert.ToBeEqual(t, 6, actual2Written)
}

func Test_LoggingWriter_Write_withInterceptor(t *testing.T) {
	givenLogger := newMockCoreLogger("foo")
	givenLogger.initLoggedEvents()
	givenLogger.level = level.Trace
	givenInterceptor := func(in []byte, lvl level.Level) ([]byte, level.Level, error) {
		assert.ToBeEqual(t, level.Warn, lvl)
		if bytes.HasPrefix(in, []byte("I ")) {
			return in[2:], level.Info, nil
		}
		if bytes.HasPrefix(in, []byte("E ")) {
			return in[2:], level.Error, nil
		}
		panic("not expected")
	}
	messageKey := givenLogger.GetProvider().GetFieldKeysSpec().GetMessage()
	instance := &LoggingWriter{
		Logger:         givenLogger,
		LevelExtractor: level.FixedLevelExtractor(level.Warn),
		Interceptor:    givenInterceptor,
		SkipFrames:     666,
	}

	actual1Written, actual1Err := instance.Write([]byte("I hello world!"))
	assert.ToBeNil(t, actual1Err)
	assert.ToBeEqual(t, 14, actual1Written)

	actual2Written, actual2Err := instance.Write([]byte("E hello world!"))
	assert.ToBeNil(t, actual2Err)
	assert.ToBeEqual(t, 14, actual2Written)

	assert.ToBeEqual(t, 2, len(*givenLogger.loggedEvents))

	assert.ToBeEqual(t, level.Info, givenLogger.loggedEvent(0).GetLevel())
	actualMessage0, _ := givenLogger.loggedEvent(0).Get(messageKey)
	assert.ToBeEqual(t, "hello world!", actualMessage0)

	assert.ToBeEqual(t, level.Error, givenLogger.loggedEvent(1).GetLevel())
	actualMessage1, _ := givenLogger.loggedEvent(1).Get(messageKey)
	assert.ToBeEqual(t, "hello world!", actualMessage1)
}

func Test_LoggingWriter_Write_interceptorErrors(t *testing.T) {
	givenLogger := newMockCoreLogger("foo")
	givenLogger.initLoggedEvents()
	givenInterceptor := func(in []byte, lvl level.Level) ([]byte, level.Level, error) {
		assert.ToBeEqual(t, level.Warn, lvl)
		return nil, 0, stringError(in)
	}
	instance := &LoggingWriter{
		Logger:         givenLogger,
		LevelExtractor: level.FixedLevelExtractor(level.Warn),
		Interceptor:    givenInterceptor,
		SkipFrames:     666,
	}

	actual1Written, actual1Err := instance.Write([]byte("foo"))
	assert.ToBeEqual(t, stringError("foo"), actual1Err)
	assert.ToBeEqual(t, 0, actual1Written)

	actual2Written, actual2Err := instance.Write([]byte("bar"))
	assert.ToBeEqual(t, stringError("bar"), actual2Err)
	assert.ToBeEqual(t, 0, actual2Written)

	assert.ToBeEqual(t, 0, len(*givenLogger.loggedEvents))
}

func Test_LoggingWriter_Write_disabledLevel(t *testing.T) {
	givenLogger := newMockCoreLogger("foo")
	givenLogger.initLoggedEvents()
	givenLogger.level = level.Warn
	extractorCalls := 0
	interceptorCalls := 0
	instance := &LoggingWriter{
		Logger: givenLogger,
		LevelExtractor: level.LineExtractorFunc(func([]byte) (level.Level, error) {
			extractorCalls++
			return level.Info, nil
		}),
		Interceptor: func(in []byte, lvl level.Level) ([]byte, level.Level, error) {
			interceptorCalls++
			return in, lvl, nil
		},
	}
	value := make([]byte, 1024*1024)

	actualWritten, actualErr := instance.Write(value)

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, len(value), actualWritten)
	assert.ToBeEqual(t, 1, extractorCalls)
	assert.ToBeEqual(t, 1, interceptorCalls)
	assert.ToBeEqual(t, 0, len(*givenLogger.loggedEvents))
	actualAllocs := testing.AllocsPerRun(100, func() {
		_, _ = instance.Write(value)
	})
	assert.ToBeEqual(t, float64(0), actualAllocs)
}

func Test_LoggingWriter_Write_usesInterceptedLevel(t *testing.T) {
	givenLogger := newMockCoreLogger("foo")
	givenLogger.initLoggedEvents()
	givenLogger.level = level.Warn
	instance := &LoggingWriter{
		Logger: givenLogger,
		LevelExtractor: level.LineExtractorFunc(func(in []byte) (level.Level, error) {
			if bytes.HasPrefix(in, []byte("up")) {
				return level.Debug, nil
			}
			return level.Error, nil
		}),
		Interceptor: func(in []byte, _ level.Level) ([]byte, level.Level, error) {
			if bytes.HasPrefix(in, []byte("up")) {
				return in, level.Error, nil
			}
			return in, level.Debug, nil
		},
	}

	actualUpWritten, actualUpErr := instance.Write([]byte("up"))
	actualDownWritten, actualDownErr := instance.Write([]byte("down"))

	assert.ToBeNoError(t, actualUpErr)
	assert.ToBeEqual(t, 2, actualUpWritten)
	assert.ToBeNoError(t, actualDownErr)
	assert.ToBeEqual(t, 4, actualDownWritten)
	assert.ToBeEqual(t, 1, len(*givenLogger.loggedEvents))
	assert.ToBeEqual(t, level.Error, givenLogger.loggedEvent(0).GetLevel())
}

func Test_LoggingWriter_Write_usesCurrentLevel(t *testing.T) {
	givenLogger := newMockCoreLogger("foo")
	givenLogger.initLoggedEvents()
	givenLogger.level = level.Warn
	instance := &LoggingWriter{
		Logger:         givenLogger,
		LevelExtractor: level.FixedLevelExtractor(level.Info),
	}

	_, actualDisabledErr := instance.Write([]byte("disabled"))
	givenLogger.level = level.Debug
	_, actualEnabledErr := instance.Write([]byte("enabled"))

	assert.ToBeNoError(t, actualDisabledErr)
	assert.ToBeNoError(t, actualEnabledErr)
	assert.ToBeEqual(t, 1, len(*givenLogger.loggedEvents))
}

func Test_LoggingWriter_Write_usesLoggerReplacedByInterceptor(t *testing.T) {
	originalLogger := newMockCoreLogger("original")
	originalLogger.initLoggedEvents()
	originalLogger.level = level.Fatal
	replacementLogger := newMockCoreLogger("replacement")
	replacementLogger.initLoggedEvents()
	replacementLogger.level = level.Trace
	instance := &LoggingWriter{
		Logger:         originalLogger,
		LevelExtractor: level.FixedLevelExtractor(level.Info),
	}
	instance.Interceptor = func(in []byte, lvl level.Level) ([]byte, level.Level, error) {
		instance.Logger = replacementLogger
		return in, lvl, nil
	}

	actualWritten, actualErr := instance.Write([]byte("message"))

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, 7, actualWritten)
	assert.ToBeEqual(t, 0, len(*originalLogger.loggedEvents))
	assert.ToBeEqual(t, 1, len(*replacementLogger.loggedEvents))
}
