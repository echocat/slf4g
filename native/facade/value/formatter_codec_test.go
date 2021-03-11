package value

import (
	"errors"
	"testing"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/native/formatter"
	"github.com/echocat/slf4g/native/hints"
)

func Test_DefaultFormatterCodec_Parse(t *testing.T) {
	cases := []struct {
		given    string
		expected formatter.Formatter
	}{{
		given:    "text",
		expected: formatter.NewText(),
	}, {
		given:    "json",
		expected: formatter.NewJson(),
	}}

	for _, c := range cases {
		t.Run(c.given, func(t *testing.T) {
			actual, actualErr := DefaultFormatterCodec.Parse(c.given)

			assert.ToBeNoError(t, actualErr)
			assert.ToBeEqual(t, c.expected, actual)
		})
	}
}

func Test_MappingFormatterCodec_Parse(t *testing.T) {
	givenFormatter := formatter.Func(func(log.Event, log.Provider, hints.Hints) ([]byte, error) {
		panic("Should never be called.")
	})
	instance := MappingFormatterCodec{
		"foo": func() (formatter.Formatter, error) {
			return givenFormatter, nil
		},
	}

	actual, actualErr := instance.Parse("foo")

	assert.ToBeNoError(t, actualErr)
	assert.ToBeSame(t, givenFormatter, actual)
}

func Test_MappingFormatterCodec_empty(t *testing.T) {
	instance := MappingFormatterCodec{}

	actual, actualErr := instance.Parse("")

	assert.ToBeNoError(t, actualErr)
	assert.ToBeSame(t, formatter.Default, actual)
}

func Test_MappingFormatterCodec_Parse_failing(t *testing.T) {
	givenFormatter := formatter.Func(func(log.Event, log.Provider, hints.Hints) ([]byte, error) {
		panic("Should never be called.")
	})
	instance := MappingFormatterCodec{
		"foo": func() (formatter.Formatter, error) {
			return givenFormatter, nil
		},
	}

	actual, actualErr := instance.Parse("bar")

	assert.ToBeMatching(t, "^unknown log format: bar$", actualErr)
	assert.ToBeNil(t, actual)
}

func Test_MappingFormatterCodec_Format(t *testing.T) {
	instance := MappingFormatterCodec{
		"foo": func() (formatter.Formatter, error) {
			return formatter.NewText(), nil
		},
		"bar": func() (formatter.Formatter, error) {
			return formatter.NewJson(), nil
		},
	}

	actual, actualErr := instance.Format(formatter.NewText())

	assert.ToBeNoError(t, actualErr)
	assert.ToBeEqual(t, "foo", actual)
}

func Test_MappingFormatterCodec_Format_failing(t *testing.T) {
	instance := MappingFormatterCodec{
		"foo": func() (formatter.Formatter, error) {
			return formatter.NewText(), nil
		},
		"bar": func() (formatter.Formatter, error) {
			return formatter.NewJson(), nil
		},
	}

	actual, actualErr := instance.Format(formatter.Noop())

	assert.ToBeMatching(t, "^unknown log formatter: formatter.Func$", actualErr)
	assert.ToBeEqual(t, "", actual)
}

func Test_Noop(t *testing.T) {
	actual := NoopFormatterCodec()

	assert.ToBeSame(t, noopFormatterCodecV, actual)
}

func Test_noopFormatterCodec_Parse(t *testing.T) {
	instance := NoopFormatterCodec()

	actual, actualErr := instance.Parse("foo")

	assert.ToBeMatching(t, "^unknown log format: foo$", actualErr)
	assert.ToBeNil(t, actual)
}

func Test_noopFormatterCodec_Format(t *testing.T) {
	instance := NoopFormatterCodec()

	actual, actualErr := instance.Format(mockFormatterCodecFormatter)

	assert.ToBeMatching(t, "^unknown log formatter: formatter.Func$", actualErr)
	assert.ToBeEqual(t, "", actual)
}

func Test_NewFormatterCodecFacade(t *testing.T) {
	givenDelegate := &mockFormatterCodec{}

	instance := NewFormatterCodecFacade(func() FormatterCodec {
		return givenDelegate
	})

	assert.ToBeSame(t, givenDelegate, instance.(formatterCodecFacade)())
}

func Test_formatterCodecFacade_Parse(t *testing.T) {
	givenDelegate := &mockFormatterCodec{}

	instance := NewFormatterCodecFacade(func() FormatterCodec {
		return givenDelegate
	})

	actual, actualErr := instance.Parse("parseInput")

	assert.ToBeMatching(t, "^expectedAfterParse$", actualErr)
	assert.ToBeSame(t, mockFormatterCodecFormatter, actual)
	assert.ToBeEqual(t, "parseInput", givenDelegate.parseWasCalledWith)
}

func Test_formatterCodecFacade_Format(t *testing.T) {
	givenDelegate := &mockFormatterCodec{}

	instance := NewFormatterCodecFacade(func() FormatterCodec {
		return givenDelegate
	})

	actual, actualErr := instance.Format(mockFormatterCodecFormatter)

	assert.ToBeMatching(t, "^expectedAfterFormat$", actualErr)
	assert.ToBeEqual(t, "expectedAfterFormat", actual)
	assert.ToBeSame(t, mockFormatterCodecFormatter, givenDelegate.formatWasCalledWith)
}

var mockFormatterCodecFormatter = formatter.Func(func(log.Event, log.Provider, hints.Hints) ([]byte, error) {
	panic("should never be called")
})

type mockFormatterCodec struct {
	parseWasCalledWith  string
	formatWasCalledWith formatter.Formatter
}

func (instance *mockFormatterCodec) Parse(plain string) (formatter.Formatter, error) {
	instance.parseWasCalledWith = plain
	return mockFormatterCodecFormatter, errors.New("expectedAfterParse")
}

func (instance *mockFormatterCodec) Format(what formatter.Formatter) (string, error) {
	instance.formatWasCalledWith = what
	return "expectedAfterFormat", errors.New("expectedAfterFormat")
}
