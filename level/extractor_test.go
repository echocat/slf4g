package level

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_FixedLevelExtractor(t *testing.T) {
	actual := FixedLevelExtractor(Fatal)

	assert.ToBeEqual(t, fixedLevelExtractor(Fatal), actual)
}

func Test_LineExtractorFunc(t *testing.T) {
	instance := LineExtractorFunc(func(bytes []byte) (Level, error) {
		switch string(bytes) {
		case "info":
			return Info, nil
		case "warn":
			return Warn, nil
		default:
			return 0, errors.New(string(bytes))
		}
	})

	actual1, acutal1Err := instance.ExtractLevelFromLine([]byte("info"))
	assert.ToBeNil(t, acutal1Err)
	assert.ToBeEqual(t, Info, actual1)
	actual2, acutal2Err := instance.ExtractLevelFromLine([]byte("warn"))
	assert.ToBeNil(t, acutal2Err)
	assert.ToBeEqual(t, Warn, actual2)
	actual3, acutal3Err := instance.ExtractLevelFromLine([]byte("foo"))
	assert.ToBeEqual(t, errors.New("foo"), acutal3Err)
	assert.ToBeEqual(t, Level(0), actual3)
}

func Test_fixedLevelExtractor_ExtractLevelFromLine(t *testing.T) {
	instance := fixedLevelExtractor(Fatal)

	actual, actualErr := instance.ExtractLevelFromLine([]byte("INFO"))

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, Fatal, actual)
}
