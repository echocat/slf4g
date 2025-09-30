package encoding

import (
	"bytes"
	"io"
)

type Buffered interface {
	Bytes() []byte
	String() string
}

// TextEncoder encodes given elements to text format.
type TextEncoder interface {
	io.ByteWriter
	WriteByteChecked(byte) func() error

	WriteBytes([]byte) error
	WriteBytesChecked([]byte) func() error

	WriteString(string) error
	WriteStringChecked(string) func() error
	WriteStringPChecked(*string) func() error
}

type BufferedTextEncoder interface {
	TextEncoder
	Buffered
}

func NewBufferedTextEncoder() BufferedTextEncoder {
	return &bufferedTextEncoder{}
}

type bufferedTextEncoder struct {
	buffer bytes.Buffer
}

func (instance *bufferedTextEncoder) WriteByte(c byte) error {
	return instance.buffer.WriteByte(c)
}

func (instance *bufferedTextEncoder) WriteByteChecked(c byte) func() error {
	return func() error {
		return instance.WriteByte(c)
	}
}

func (instance *bufferedTextEncoder) WriteBytes(p []byte) error {
	n, err := instance.buffer.Write(p)
	if err != nil {
		return err
	}
	if n < len(p) {
		return io.ErrShortWrite
	}
	return nil
}

func (instance *bufferedTextEncoder) WriteBytesChecked(p []byte) func() error {
	return func() error {
		return instance.WriteBytes(p)
	}
}

func (instance *bufferedTextEncoder) WriteString(s string) error {
	_, err := instance.buffer.WriteString(s)
	return err
}

func (instance *bufferedTextEncoder) WriteStringChecked(v string) func() error {
	return instance.WriteStringPChecked(&v)
}

func (instance *bufferedTextEncoder) WriteStringPChecked(v *string) func() error {
	return func() error {
		return instance.WriteString(*v)
	}
}

func (instance *bufferedTextEncoder) Bytes() []byte {
	return instance.buffer.Bytes()
}

func (instance *bufferedTextEncoder) String() string {
	return string(instance.Bytes())
}
