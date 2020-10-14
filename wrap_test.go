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
	given := newMockCoreLogger("foo")

	actual := UnwrapCoreLogger(given)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapCoreLogger_withCoreUnwrapMethod(t *testing.T) {
	wrapped := newMockCoreLogger("foo")
	given := newWrappingCoreLogger(wrapped)

	actual := UnwrapCoreLogger(given)

	assert.ToBeSame(t, wrapped, actual)
}

func Test_UnwrapCoreLogger_withUnwrapMethod(t *testing.T) {
	wrapped := newMockLogger("foo")
	given := newWrappingLogger(wrapped)

	actual := UnwrapCoreLogger(given)

	assert.ToBeSame(t, wrapped, actual)
}

func Test_UnwrapLogger_withNil(t *testing.T) {
	actual := UnwrapLogger(nil)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapLogger_withoutUnwrapMethod(t *testing.T) {
	given := &mockCoreLogger{}

	actual := UnwrapLogger(given)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapLogger_withCoreUnwrapMethod(t *testing.T) {
	wrapped := newMockLogger("foo")
	given := &wrappingCoreLogger{wrapped}

	actual := UnwrapLogger(given)

	assert.ToBeSame(t, wrapped, actual)
}

func Test_UnwrapLogger_withUnwrapMethod(t *testing.T) {
	wrapped := newMockLogger("foo")
	given := newWrappingLogger(wrapped)

	actual := UnwrapLogger(given)

	assert.ToBeSame(t, wrapped, actual)
}

func Test_UnwrapProvider_withNil(t *testing.T) {
	actual := UnwrapProvider(nil)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapProvider_withoutUnwrapMethod(t *testing.T) {
	given := newMockProvider("test")

	actual := UnwrapProvider(given)

	assert.ToBeNil(t, actual)
}

func Test_UnwrapProvider_withUnwrapMethod(t *testing.T) {
	wrapped := newMockProvider("test")
	given := newWrappingProvider(wrapped)

	actual := UnwrapProvider(given)

	assert.ToBeSame(t, wrapped, actual)
}
