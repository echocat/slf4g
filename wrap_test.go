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
	given := &mockCoreLogger{}

	actual := UnwrapCoreLogger(given)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapCoreLogger_withCoreUnwrapMethod(t *testing.T) {
	wrapped := &mockCoreLogger{}
	given := &wrappingMockCoreLogger{wrapped}

	actual := UnwrapCoreLogger(given)

	assert.ToBeSame(t, wrapped, actual)
}

func Test_UnwrapCoreLogger_withUnwrapMethod(t *testing.T) {
	wrapped := &mockLogger{}
	given := &wrappingMockLogger{wrapped}

	actual := UnwrapCoreLogger(given)

	assert.ToBeSame(t, wrapped, actual)
}

func Test_UnwrapLogger_withNil(t *testing.T) {
	actual := UnwrapLogger(nil)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapLogger_withoutUnwrapMethod(t *testing.T) {
	given := &mockLogger{}

	actual := UnwrapLogger(given)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapLogger_withCoreUnwrapMethod(t *testing.T) {
	wrapped := &mockLogger{}
	given := &wrappingMockCoreLogger{wrapped}

	actual := UnwrapLogger(given)

	assert.ToBeOfType(t, &loggerImpl{}, actual)
	assert.ToBeEqual(t, wrapped, actual.(*loggerImpl).coreProvider())
}

func Test_UnwrapLogger_withUnwrapMethod(t *testing.T) {
	wrapped := &mockLogger{}
	given := &wrappingMockLogger{wrapped}

	actual := UnwrapLogger(given)

	assert.ToBeSame(t, wrapped, actual)
}

func Test_UnwrapProvider_withNil(t *testing.T) {
	actual := UnwrapProvider(nil)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapProvider_withoutUnwrapMethod(t *testing.T) {
	given := &mockProvider{}

	actual := UnwrapProvider(given)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapProvider_withUnwrapMethod(t *testing.T) {
	wrapped := &mockProvider{}
	given := &wrappingTestProvider{wrapped}

	actual := UnwrapProvider(given)

	assert.ToBeSame(t, wrapped, actual)
}
