package functions

import (
	"errors"
	"fmt"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_IndentMultiline(t *testing.T) {
	cases := []struct {
		given             string
		givenFirstIndent  string
		givenOthersIndent string
		expected          string
	}{{
		given:             "hello\nworld\nfoo\nbar",
		givenFirstIndent:  "1",
		givenOthersIndent: "2",
		expected:          "1hello\n2world\n2foo\n2bar",
	}, {
		given:             "hello\r\nworld\r\nfoo\r\nbar",
		givenFirstIndent:  "1",
		givenOthersIndent: "2",
		expected:          "1hello\n2world\n2foo\n2bar",
	}, {
		given:             "hello\n",
		givenFirstIndent:  "1",
		givenOthersIndent: "2",
		expected:          "1hello\n2",
	}, {
		given:             "hello",
		givenFirstIndent:  "1",
		givenOthersIndent: "2",
		expected:          "1hello",
	}}

	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			actual := IndentMultiline(c.givenFirstIndent, c.givenOthersIndent, c.given)

			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}

func Test_EncodeMultilineWithIndent_failing(t *testing.T) {
	expectedErr := errors.New("expected")
	actualErr := EncodeMultilineWithIndent("1", "2", failingTextEncoder{expectedErr}, "foo")

	assert.ToBeSame(t, expectedErr, actualErr)
}

type failingTextEncoder struct {
	error
}

func (instance failingTextEncoder) WriteByte(c byte) error {
	return instance.error
}

func (instance failingTextEncoder) WriteByteChecked(b byte) func() error {
	return func() error { return instance.error }
}

func (instance failingTextEncoder) WriteBytes(bytes []byte) error {
	return instance.error
}

func (instance failingTextEncoder) WriteBytesChecked(bytes []byte) func() error {
	return func() error { return instance.error }
}

func (instance failingTextEncoder) WriteString(s string) error {
	return instance.error
}

func (instance failingTextEncoder) WriteStringChecked(s string) func() error {
	return func() error { return instance.error }
}

func (instance failingTextEncoder) WriteStringPChecked(s *string) func() error {
	return func() error { return instance.error }
}
