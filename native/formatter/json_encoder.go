package formatter

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"unicode"
)

type jsonEncoder struct {
	encoder *json.Encoder
	writer  *filteringTailingNewLineWriter
}

func (instance *jsonEncoder) WriteByte(c byte) error {
	return instance.writer.WriteByte(c)
}

func (instance *jsonEncoder) WriteByteChecked(c byte) checkedExecution {
	return func() error {
		return instance.WriteByte(c)
	}
}

func (instance *jsonEncoder) WriteBytes(p []byte) error {
	n, err := instance.Write(p)
	if err != nil {
		return err
	}
	if n < len(p) {
		return io.ErrShortWrite
	}
	return nil
}

func (instance *jsonEncoder) WriteBytesChecked(p []byte) checkedExecution {
	return func() error {
		return instance.WriteBytes(p)
	}
}

func (instance *jsonEncoder) Write(p []byte) (n int, err error) {
	return instance.writer.Write(p)
}

func (instance *jsonEncoder) WriteKeyValue(k string, v interface{}) error {
	return executeChecked(
		instance.WriteValueChecked(k),
		instance.WriteByteChecked(':'),
		instance.WriteValueChecked(v),
	)
}

func (instance *jsonEncoder) WriteKeyValueChecked(k string, v interface{}) checkedExecution {
	return func() error {
		return instance.WriteKeyValue(k, v)
	}
}

func (instance *jsonEncoder) WriteValue(v interface{}) error {
	if ve, ok := v.(error); ok {
		v = ve.Error()
	}
	if vs, ok := v.(*string); ok {
		v = *vs
	}
	if vs, ok := v.(string); ok {
		v = strings.TrimRightFunc(vs, unicode.IsSpace)
	}
	return instance.encoder.Encode(v)
}

func (instance *jsonEncoder) WriteValueChecked(v interface{}) checkedExecution {
	return func() error {
		return instance.WriteValue(v)
	}
}

type filteringTailingNewLineWriter struct {
	io.Writer
}

func (instance *filteringTailingNewLineWriter) Write(p []byte) (n int, err error) {
	if len(p) <= 0 {
		return 0, nil
	}
	if p[len(p)-1] == '\n' {
		p = p[:len(p)-1]
	}
	return instance.Writer.Write(p)
}

func (instance *filteringTailingNewLineWriter) WriteByte(c byte) error {
	if v, ok := instance.Writer.(io.ByteWriter); ok {
		return v.WriteByte(c)
	}
	n, err := instance.Write([]byte{c})
	if n != 1 {
		return io.ErrShortWrite
	}
	return err
}

type bufferedJsonEncoder struct {
	jsonEncoder

	buffer bytes.Buffer
}

func newJsonEncoderBuffered() *bufferedJsonEncoder {
	result := new(bufferedJsonEncoder)
	result.jsonEncoder.writer = &filteringTailingNewLineWriter{&result.buffer}
	result.jsonEncoder.encoder = json.NewEncoder(result.jsonEncoder.writer)
	return result
}

func (instance *bufferedJsonEncoder) Bytes() []byte {
	return instance.buffer.Bytes()
}

func (instance *bufferedJsonEncoder) String() string {
	return string(instance.Bytes())
}
