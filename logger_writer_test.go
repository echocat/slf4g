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
	assert.ToBeEqual(t, 12, actual1Written)

	actual2Written, actual2Err := instance.Write([]byte("E hello world!"))
	assert.ToBeNil(t, actual2Err)
	assert.ToBeEqual(t, 12, actual2Written)

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
