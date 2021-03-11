package value

import (
	"os"
	"testing"

	"github.com/echocat/slf4g/native/formatter"

	"github.com/echocat/slf4g/internal/test/assert"

	"github.com/echocat/slf4g/native/consumer"
)

func Test_NewFormatter(t *testing.T) {
	givenConsumer := consumer.NewWriter(os.Stderr, func(writer *consumer.Writer) {
		writer.Formatter = mockFormatterCodecFormatter
	})

	instance := NewFormatter(givenConsumer)

	assert.ToBeNotNil(t, instance)
	assert.ToBeSame(t, givenConsumer, instance.Target)
	assert.ToBeNil(t, instance.Codec)
}

func Test_NewFormatter_customized(t *testing.T) {
	givenConsumer := consumer.NewWriter(os.Stderr, func(writer *consumer.Writer) {
		writer.Formatter = mockFormatterCodecFormatter
	})

	instance := NewFormatter(givenConsumer, func(formatter *Formatter) {
		formatter.Codec = NoopFormatterCodec()
	})

	assert.ToBeNotNil(t, instance)
	assert.ToBeSame(t, givenConsumer, instance.Target)
	assert.ToBeSame(t, NoopFormatterCodec(), instance.Codec)
}

func Test_Formatter_UnmarshalText(t *testing.T) {
	givenConsumer := consumer.NewWriter(os.Stderr, func(writer *consumer.Writer) {
		writer.Formatter = mockFormatterCodecFormatter
	})
	instance := NewFormatter(givenConsumer)

	actualErr := instance.UnmarshalText([]byte("text"))

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, formatter.NewText(), givenConsumer.Formatter)
}

func Test_Formatter_UnmarshalText_failing(t *testing.T) {
	givenConsumer := consumer.NewWriter(os.Stderr, func(writer *consumer.Writer) {
		writer.Formatter = mockFormatterCodecFormatter
	})
	instance := NewFormatter(givenConsumer)

	actualErr := instance.UnmarshalText([]byte("foo"))

	assert.ToBeMatching(t, "^unknown log format: foo$", actualErr)
	assert.ToBeSame(t, mockFormatterCodecFormatter, givenConsumer.Formatter)
}

func Test_Formatter_String(t *testing.T) {
	givenConsumer := consumer.NewWriter(os.Stderr, func(writer *consumer.Writer) {
		writer.Formatter = formatter.NewText()
	})
	instance := NewFormatter(givenConsumer)

	actual := instance.String()

	assert.ToBeEqual(t, "text", actual)
}

func Test_Formatter_String_withEmptyTarget(t *testing.T) {
	instance := NewFormatter(nil)

	actual := instance.String()

	assert.ToBeEqual(t, "", actual)
}

func Test_Formatter_String_failing(t *testing.T) {
	givenConsumer := consumer.NewWriter(os.Stderr, func(writer *consumer.Writer) {
		writer.Formatter = mockFormatterCodecFormatter
	})
	instance := NewFormatter(givenConsumer)

	actual := instance.String()

	assert.ToBeEqual(t, "ERR-unknown log formatter: formatter.Func", actual)
}

func Test_Formatter_Get(t *testing.T) {
	givenConsumer := consumer.NewWriter(os.Stderr, func(writer *consumer.Writer) {
		writer.Formatter = mockFormatterCodecFormatter
	})
	instance := NewFormatter(givenConsumer)

	actual := instance.Get()

	assert.ToBeSame(t, mockFormatterCodecFormatter, actual)
}

func Test_Formatter_getCodec_explicit(t *testing.T) {
	instance := &Formatter{
		Codec: &mockFormatterCodec{},
	}

	actual := instance.getCodec()

	assert.ToBeSame(t, instance.Codec, actual)
}

func Test_Formatter_getCodec_default(t *testing.T) {
	instance := &Formatter{
		Codec: nil,
	}

	actual := instance.getCodec()

	assert.ToBeEqual(t, DefaultFormatterCodec, actual)
}

func Test_Formatter_getCodec_noop(t *testing.T) {
	old := DefaultFormatterCodec
	defer func() {
		DefaultFormatterCodec = old
	}()
	DefaultFormatterCodec = nil

	instance := &Formatter{
		Codec: nil,
	}

	actual := instance.getCodec()

	assert.ToBeEqual(t, NoopFormatterCodec(), actual)
}
