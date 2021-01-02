package formatter

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"unicode"
)

type jsonEncoder interface {
	io.ByteWriter
	WriteByteChecked(c byte) func() error

	WriteKeyValue(k string, v interface{}) error
	WriteKeyValueChecked(k string, v interface{}) func() error

	WriteValue(v interface{}) error
	WriteValueChecked(v interface{}) func() error
}

func newBufferedJsonEncoder() *bufferedJsonEncoder {
	result := new(bufferedJsonEncoder)
	result.encoder = json.NewEncoder(filteringTailingNewLineWriter{&result.buffer})
	return result
}

type bufferedJsonEncoder struct {
	encoder *json.Encoder
	buffer  bytes.Buffer
}

func (instance *bufferedJsonEncoder) WriteByte(c byte) error {
	return instance.buffer.WriteByte(c)
}

func (instance *bufferedJsonEncoder) WriteByteChecked(c byte) func() error {
	return func() error {
		return instance.WriteByte(c)
	}
}

func (instance *bufferedJsonEncoder) WriteKeyValue(k string, v interface{}) error {
	return executeChecked(
		instance.WriteValueChecked(k),
		instance.WriteByteChecked(':'),
		instance.WriteValueChecked(v),
	)
}

func (instance *bufferedJsonEncoder) WriteKeyValueChecked(k string, v interface{}) func() error {
	return func() error {
		return instance.WriteKeyValue(k, v)
	}
}

func (instance *bufferedJsonEncoder) WriteValue(v interface{}) error {
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

func (instance *bufferedJsonEncoder) WriteValueChecked(v interface{}) func() error {
	return func() error {
		return instance.WriteValue(v)
	}
}

func (instance *bufferedJsonEncoder) Bytes() []byte {
	return instance.buffer.Bytes()
}

func (instance *bufferedJsonEncoder) String() string {
	return string(instance.Bytes())
}

type filteringTailingNewLineWriter struct {
	io.Writer
}

func (instance filteringTailingNewLineWriter) Write(p []byte) (int, error) {
	if len(p) <= 0 {
		return 0, nil
	}
	if p[len(p)-1] == '\n' {
		n, err := instance.Writer.Write(p[:len(p)-1])
		return n + 1, err
	}
	return instance.Writer.Write(p)
}
