package formatter

import (
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_newBufferedTextEncoder(t *testing.T) {
	actual := newBufferedTextEncoder()

	assert.ToBeNotNil(t, actual)
	assert.ToBeEqual(t, 0, actual.buffer.Len())
}

func Test_bufferedTextEncoder_WriteByte(t *testing.T) {
	instance := newBufferedTextEncoder()

	actualErr := instance.WriteByte('x')
	assert.ToBeNil(t, actualErr)

	assert.ToBeEqual(t, `x`, instance.buffer.String())
}

func Test_bufferedTextEncoder_WriteByteChecked(t *testing.T) {
	instance := newBufferedTextEncoder()

	actualErr := instance.WriteByteChecked('x')()
	assert.ToBeNil(t, actualErr)

	assert.ToBeEqual(t, `x`, instance.buffer.String())
}

func Test_bufferedTextEncoder_WriteBytes(t *testing.T) {
	instance := newBufferedTextEncoder()

	actualErr := instance.WriteBytes([]byte{'x', 'y'})
	assert.ToBeNil(t, actualErr)

	assert.ToBeEqual(t, `xy`, instance.buffer.String())
}

func Test_bufferedTextEncoder_WriteBytesChecked(t *testing.T) {
	instance := newBufferedTextEncoder()

	actualErr := instance.WriteBytesChecked([]byte{'x', 'y'})()
	assert.ToBeNil(t, actualErr)

	assert.ToBeEqual(t, `xy`, instance.buffer.String())
}

func Test_bufferedTextEncoder_WriteString(t *testing.T) {
	instance := newBufferedTextEncoder()

	actualErr := instance.WriteString("xy")
	assert.ToBeNil(t, actualErr)

	assert.ToBeEqual(t, `xy`, instance.buffer.String())
}

func Test_bufferedTextEncoder_WriteStringChecked(t *testing.T) {
	instance := newBufferedTextEncoder()

	actualErr := instance.WriteStringChecked("xy")()
	assert.ToBeNil(t, actualErr)

	assert.ToBeEqual(t, `xy`, instance.buffer.String())
}

func Test_bufferedTextEncoder_WriteStringPChecked(t *testing.T) {
	instance := newBufferedTextEncoder()
	givenString := "xy"

	actualErr := instance.WriteStringPChecked(&givenString)()
	assert.ToBeNil(t, actualErr)

	assert.ToBeEqual(t, `xy`, instance.buffer.String())
}

func Test_bufferedTextEncoder_Bytes(t *testing.T) {
	instance := newBufferedTextEncoder()

	instance.buffer.WriteString("hello")
	actual := instance.Bytes()

	assert.ToBeEqual(t, `hello`, string(actual))
}

func Test_bufferedTextEncoder_String(t *testing.T) {
	instance := newBufferedTextEncoder()

	instance.buffer.WriteString("hello")
	actual := instance.String()

	assert.ToBeEqual(t, `hello`, actual)
}
