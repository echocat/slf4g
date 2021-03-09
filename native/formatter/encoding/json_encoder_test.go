package encoding

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_NewBufferedJsonEncoder(t *testing.T) {
	actual := NewBufferedJsonEncoder().(*bufferedJsonEncoder)

	assert.ToBeNotNil(t, actual)
	assert.ToBeEqual(t, 0, actual.buffer.Len())

	assert.ToBeNil(t, actual.WriteValue("foo"))
	assert.ToBeEqual(t, `"foo"`, actual.buffer.String())
}

func Test_bufferedJsonEncoder_WriteValue(t *testing.T) {
	instance := NewBufferedJsonEncoder().(*bufferedJsonEncoder)

	actualErr := instance.WriteValue(struct{ Foo string }{Foo: "bar"})
	assert.ToBeNil(t, actualErr)

	assert.ToBeEqual(t, `{"Foo":"bar"}`, instance.buffer.String())
}

func Test_bufferedJsonEncoder_WriteValue_string(t *testing.T) {
	instance := NewBufferedJsonEncoder().(*bufferedJsonEncoder)

	actualErr := instance.WriteValue("\tabc\n\t")
	assert.ToBeNil(t, actualErr)

	assert.ToBeEqual(t, `"\tabc"`, instance.buffer.String())
}

func Test_bufferedJsonEncoder_WriteValue_string_as_pointer(t *testing.T) {
	instance := NewBufferedJsonEncoder().(*bufferedJsonEncoder)

	actualErr := instance.WriteValue(pstring("abc"))
	assert.ToBeNil(t, actualErr)

	assert.ToBeEqual(t, `"abc"`, instance.buffer.String())
}

func Test_bufferedJsonEncoder_WriteValue_error(t *testing.T) {
	instance := NewBufferedJsonEncoder().(*bufferedJsonEncoder)

	actualErr := instance.WriteValue(errors.New("abc"))
	assert.ToBeNil(t, actualErr)

	assert.ToBeEqual(t, `"abc"`, instance.buffer.String())
}

func Test_bufferedJsonEncoder_WriteValueChecked(t *testing.T) {
	instance := NewBufferedJsonEncoder().(*bufferedJsonEncoder)

	actualErr := instance.WriteValueChecked(struct{ Foo string }{Foo: "bar"})()
	assert.ToBeNil(t, actualErr)

	assert.ToBeEqual(t, `{"Foo":"bar"}`, instance.buffer.String())
}

func Test_bufferedJsonEncoder_WriteKeyValue(t *testing.T) {
	instance := NewBufferedJsonEncoder().(*bufferedJsonEncoder)

	actualErr := instance.WriteKeyValue("hello", struct{ Foo string }{Foo: "bar"})
	assert.ToBeNil(t, actualErr)

	assert.ToBeEqual(t, `"hello":{"Foo":"bar"}`, instance.buffer.String())
}

func Test_bufferedJsonEncoder_WriteKeyValueChecked(t *testing.T) {
	instance := NewBufferedJsonEncoder().(*bufferedJsonEncoder)

	actualErr := instance.WriteKeyValueChecked("hello", struct{ Foo string }{Foo: "bar"})()
	assert.ToBeNil(t, actualErr)

	assert.ToBeEqual(t, `"hello":{"Foo":"bar"}`, instance.buffer.String())
}

func Test_filteringTailingNewLineWriter_Write(t *testing.T) {
	cases := []struct {
		given    string
		expected string
	}{
		{"<nil>", ""},
		{"", ""},
		{"abc", "abc"},
		{"abc\nhello\nworld", "abc\nhello\nworld"},
		{"\nabc", "\nabc"},
		{"abc\t", "abc\t"},
		{"abc\r", "abc\r"},
		{"abc\n", "abc"},
		{"abc\nhello\nworld\n", "abc\nhello\nworld"},
	}

	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			givenBuffer := new(bytes.Buffer)
			instance := filteringTailingNewLineWriter{givenBuffer}
			var given []byte
			if c.given != "<nil>" {
				given = []byte(c.given)
			}

			actualN, actualErr := instance.Write(given)

			assert.ToBeNil(t, actualErr)
			assert.ToBeEqual(t, len(given), actualN)
			assert.ToBeEqual(t, c.expected, givenBuffer.String())
		})
	}
}

func pstring(v string) *string {
	return &v
}
