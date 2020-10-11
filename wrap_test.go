package log

import (
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_UnwrapCoreLogger_withNil(t *testing.T) {
	actual := UnwrapCoreLogger(nil)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapCoreLogger_withoutUnwrapMethod(t *testing.T) {
	given := &testCoreLogger{}

	actual := UnwrapCoreLogger(given)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapCoreLogger_withCoreUnwrapMethod(t *testing.T) {
	wrapped := &testCoreLogger{}
	given := &wrappingCoreTestLogger{wrapped}

	actual := UnwrapCoreLogger(given)

	assert.ToBeSame(t, wrapped, actual)
}

func Test_UnwrapCoreLogger_withUnwrapMethod(t *testing.T) {
	wrapped := &testLogger{}
	given := &wrappingTestLogger{wrapped}

	actual := UnwrapCoreLogger(given)

	assert.ToBeSame(t, wrapped, actual)
}

func Test_UnwrapLogger_withNil(t *testing.T) {
	actual := UnwrapLogger(nil)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapLogger_withoutUnwrapMethod(t *testing.T) {
	given := &testLogger{}

	actual := UnwrapLogger(given)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapLogger_withCoreUnwrapMethod(t *testing.T) {
	wrapped := &testLogger{}
	given := &wrappingCoreTestLogger{wrapped}

	actual := UnwrapLogger(given)

	assert.ToBeOfType(t, &loggerImpl{}, actual)
	assert.ToBeEqual(t, wrapped, actual.(*loggerImpl).coreProvider())
}

func Test_UnwrapLogger_withUnwrapMethod(t *testing.T) {
	wrapped := &testLogger{}
	given := &wrappingTestLogger{wrapped}

	actual := UnwrapLogger(given)

	assert.ToBeSame(t, wrapped, actual)
}

func Test_UnwrapProvider_withNil(t *testing.T) {
	actual := UnwrapProvider(nil)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapProvider_withoutUnwrapMethod(t *testing.T) {
	given := &testProvider{}

	actual := UnwrapProvider(given)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapProvider_withUnwrapMethod(t *testing.T) {
	wrapped := &testProvider{}
	given := &wrappingTestProvider{wrapped}

	actual := UnwrapProvider(given)

	assert.ToBeSame(t, wrapped, actual)
}
