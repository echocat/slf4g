package formatter

import (
	"fmt"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_NewSimpleTextValue(t *testing.T) {
	instance := NewSimpleTextValue()

	assert.ToBeEqual(t, QuoteTypeMinimal, instance.QuoteType)
}

func Test_NewSimpleTextValue_withCustomization(t *testing.T) {
	givenQuoteType := QuoteType(66)

	instance := NewSimpleTextValue(func(simpleTextValue *SimpleTextValue) {
		simpleTextValue.QuoteType = givenQuoteType
	})

	assert.ToBeEqual(t, givenQuoteType, instance.QuoteType)
}

func Test_SimpleTextValue_FormatTextValue(t *testing.T) {
	cases := []struct {
		given      interface{}
		minimal    string
		normal     string
		everything string
	}{
		{aLazy("abc"), "abc", "\"abc\"", "\"abc\""},
		{aLazy("abc%"), "\"abc%\"", "\"abc%\"", "\"abc%\""},
		{pstring("abc"), "abc", "\"abc\"", "\"abc\""},
		{"abc%", "\"abc%\"", "\"abc%\"", "\"abc%\""},
		{pstring("abc%"), "\"abc%\"", "\"abc%\"", "\"abc%\""},
		{aStringer("abc"), "abc", "\"abc\"", "\"abc\""},
		{paStringer("abc"), "abc", "\"abc\"", "\"abc\""},
		{aStringer("abc%"), "\"abc%\"", "\"abc%\"", "\"abc%\""},
		{paStringer("abc%"), "\"abc%\"", "\"abc%\"", "\"abc%\""},
		{aFormatter("abc"), "abc", "\"abc\"", "\"abc\""},
		{paFormatter("abc"), "abc", "\"abc\"", "\"abc\""},
		{aFormatter("abc%"), "\"abc%\"", "\"abc%\"", "\"abc%\""},
		{paFormatter("abc%"), "\"abc%\"", "\"abc%\"", "\"abc%\""},
		{anError("abc"), "abc", "\"abc\"", "\"abc\""},
		{panError("abc"), "abc", "\"abc\"", "\"abc\""},
		{anError("abc%"), "\"abc%\"", "\"abc%\"", "\"abc%\""},
		{panError("abc%"), "\"abc%\"", "\"abc%\"", "\"abc%\""},
		{1, "1", "1", "\"1\""},
		{1.2, "1.2", "1.2", "\"1.2\""},
		{true, "true", "true", "\"true\""},
	}

	instanceMinimal := NewSimpleTextValue(func(value *SimpleTextValue) {
		value.QuoteType = QuoteTypeMinimal
	})
	instanceNormal := NewSimpleTextValue(func(value *SimpleTextValue) {
		value.QuoteType = QuoteTypeNormal
	})
	instanceEverything := NewSimpleTextValue(func(value *SimpleTextValue) {
		value.QuoteType = QuoteTypeEverything
	})

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d-%v-minimal", i, c.given), func(t *testing.T) {
			actual, actualErr := instanceMinimal.FormatTextValue(c.given, nil)

			assert.ToBeNil(t, actualErr)
			assert.ToBeEqual(t, c.minimal, string(actual))
		})

		t.Run(fmt.Sprintf("%d-%v-normal", i, c.given), func(t *testing.T) {
			actual, actualErr := instanceNormal.FormatTextValue(c.given, nil)

			assert.ToBeNil(t, actualErr)
			assert.ToBeEqual(t, c.normal, string(actual))
		})

		t.Run(fmt.Sprintf("%d-%v-everything", i, c.given), func(t *testing.T) {
			actual, actualErr := instanceEverything.FormatTextValue(c.given, nil)

			assert.ToBeNil(t, actualErr)
			assert.ToBeEqual(t, c.everything, string(actual))
		})
	}
}

func pstring(v string) *string {
	return &v
}

type aLazy string

func (instance aLazy) Get() interface{} {
	return string(instance)
}

type aStringer string

func (instance aStringer) String() string {
	return string(instance)
}

func paStringer(v string) *aStringer {
	va := aStringer(v)
	return &va
}

type aFormatter string

func (instance aFormatter) Format(state fmt.State, c rune) {
	_, _ = state.Write([]byte(instance))
}

func paFormatter(v string) *aFormatter {
	va := aFormatter(v)
	return &va
}

type anError string

func (instance anError) Error() string {
	return string(instance)
}

func panError(v string) *anError {
	va := anError(v)
	return &va
}
