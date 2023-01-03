package formatter

import (
	"errors"
	"fmt"
	"testing"

	"github.com/echocat/slf4g/native/formatter/encoding"

	"github.com/echocat/slf4g/fields"

	nlevel "github.com/echocat/slf4g/native/level"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/testing/recording"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_NewJson(t *testing.T) {
	instance := NewJson()

	assert.ToBeEqual(t, "", instance.KeyLevel)
	assert.ToBeNil(t, instance.PrintRootLogger)
}

func Test_NewJson_withCustomization(t *testing.T) {
	instance := NewJson(func(json *Json) {
		vFalse := false
		json.KeyLevel = "foo"
		json.PrintRootLogger = &vFalse
	})

	assert.ToBeEqual(t, "foo", instance.KeyLevel)
	assert.ToBeEqual(t, false, *instance.PrintRootLogger)
}

func Test_Json_getLevelKey_explicit(t *testing.T) {
	instance := NewJson(func(json *Json) {
		json.KeyLevel = "foo"
	})

	actual := instance.getLevelKey()

	assert.ToBeEqual(t, "foo", actual)
}

func Test_Json_getLevelKey_default(t *testing.T) {
	instance := NewJson(func(json *Json) {
		json.KeyLevel = ""
	})

	actual := instance.getLevelKey()

	assert.ToBeEqual(t, DefaultKeyLevel, actual)
}

func Test_Json_getPrintRootLogger_explicit(t *testing.T) {
	instance := NewJson(func(json *Json) {
		vTrue := true
		json.PrintRootLogger = &vTrue
	})

	actual := instance.getPrintRootLogger()

	assert.ToBeEqual(t, true, actual)
}

func Test_Json_getPrintRootLogger_default(t *testing.T) {
	instance := NewJson(func(json *Json) {
		json.PrintRootLogger = nil
	})

	actual := instance.getPrintRootLogger()

	assert.ToBeEqual(t, DefaultPrintRootLogger, actual)
}

func Test_Json_getLevelFormatter_explicit(t *testing.T) {
	givenProvider := recording.NewProvider()
	givenFormatter := LevelFunc(func(in level.Level, using log.Provider) (interface{}, error) {
		return fmt.Sprintf("some-%d", in), nil
	})
	instance := NewJson(func(json *Json) {
		json.LevelFormatter = givenFormatter
	})

	actual := instance.getLevelFormatter(givenProvider)

	assert.ToBeSame(t, givenFormatter, actual)
}

func Test_Json_getLevelFormatter_byProvider(t *testing.T) {
	givenProvider := &someProvider{
		Provider: recording.NewProvider(),
		names:    nlevel.NewNames(),
	}
	instance := NewJson(func(json *Json) {
		json.LevelFormatter = nil
	})

	actual := instance.getLevelFormatter(givenProvider)
	actualFormatted, actualErr := actual.FormatLevel(level.Level(666), nil)

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, "myFormattedLevel-666", actualFormatted)
}

func Test_Json_getLevelFormatter_default(t *testing.T) {
	givenProvider := recording.NewProvider()
	instance := NewJson(func(json *Json) {
		json.LevelFormatter = nil
	})

	actual := instance.getLevelFormatter(givenProvider)

	assert.ToBeSame(t, DefaultLevel, actual)
}

func Test_Json_Format(t *testing.T) {
	provider := recording.NewProvider()
	logger := provider.GetRootLogger()
	instance := NewJson(func(json *Json) {
		json.KeySorter = fields.DefaultKeySorter
	})

	cases := []struct {
		name     string
		given    log.Event
		expected string
	}{{
		name: "withStringAndInteger",
		given: logger.NewEvent(level.Info, map[string]interface{}{
			"foo": "foo",
			"bar": 1,
		}),
		expected: `{"level":"INFO","bar":1,"foo":"foo"}
`,
	}, {
		name:  "withLevelOnly",
		given: logger.NewEvent(level.Warn, nil),
		expected: `{"level":"WARN"}
`,
	}, {
		name:     "nilEvent",
		given:    nil,
		expected: ``,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, actualErr := instance.Format(c.given, provider, nil)

			assert.ToBeNil(t, actualErr)
			assert.ToBeEqual(t, c.expected, string(actual))
		})
	}
}

func Test_Json_Format_failing(t *testing.T) {
	givenProvider := recording.NewProvider()
	givenLogger := givenProvider.GetRootLogger()
	givenEvent := givenLogger.NewEvent(level.Warn, nil)
	givenError := errors.New("expected")
	instance := NewJson(func(json *Json) {
		json.LevelFormatter = LevelFunc(func(level.Level, log.Provider) (interface{}, error) {
			return nil, givenError
		})
	})

	actual, actualErr := instance.Format(givenEvent, givenProvider, nil)

	assert.ToBeMatching(t, "cannot format event .+: expected", actualErr)
	assert.ToBeEqual(t, "", string(actual))
}

func Test_Json_encodeLevelChecked(t *testing.T) {
	givenProvider := recording.NewProvider()
	givenLogger := givenProvider.GetRootLogger()
	givenEvent := givenLogger.NewEvent(level.Warn, nil)
	givenEncoder := encoding.NewBufferedJsonEncoder()
	instance := NewJson(func(json *Json) {
		json.LevelFormatter = LevelFunc(func(actualLevel level.Level, actualProvider log.Provider) (interface{}, error) {
			assert.ToBeEqual(t, givenProvider, actualProvider)
			assert.ToBeEqual(t, givenEvent.GetLevel(), actualLevel)
			return 666, nil
		})
		json.KeyLevel = "myKey"
	})

	actualErr := instance.encodeLevelChecked(givenEvent, givenProvider, givenEncoder)()

	assert.ToBeNil(t, actualErr)
	assert.ToBeEqual(t, `"myKey":666`, givenEncoder.String())
}

func Test_Json_encodeLevelChecked_failingOnLevelFormat(t *testing.T) {
	givenProvider := recording.NewProvider()
	givenLogger := givenProvider.GetRootLogger()
	givenEvent := givenLogger.NewEvent(level.Warn, nil)
	givenEncoder := encoding.NewBufferedJsonEncoder()
	givenError := errors.New("expected")
	instance := NewJson(func(json *Json) {
		json.LevelFormatter = LevelFunc(func(actualLevel level.Level, actualProvider log.Provider) (interface{}, error) {
			assert.ToBeEqual(t, givenProvider, actualProvider)
			assert.ToBeEqual(t, givenEvent.GetLevel(), actualLevel)
			return 0, givenError
		})
		json.KeyLevel = "myKey"
	})

	actualErr := instance.encodeLevelChecked(givenEvent, givenProvider, givenEncoder)()

	assert.ToBeSame(t, givenError, actualErr)
	assert.ToBeEqual(t, ``, givenEncoder.String())
}

func Test_Json_encodeLevelChecked_failingOnEncodeValue(t *testing.T) {
	givenProvider := recording.NewProvider()
	givenLogger := givenProvider.GetRootLogger()
	givenEvent := givenLogger.NewEvent(level.Warn, nil)
	givenEncoder := encoding.NewBufferedJsonEncoder()
	givenError := errors.New("expected")
	instance := NewJson(func(json *Json) {
		json.LevelFormatter = LevelFunc(func(actualLevel level.Level, actualProvider log.Provider) (interface{}, error) {
			assert.ToBeEqual(t, givenProvider, actualProvider)
			assert.ToBeEqual(t, givenEvent.GetLevel(), actualLevel)
			return &failingJsonMarshalling{givenError}, nil
		})
		json.KeyLevel = "myKey"
	})

	actualErr := instance.encodeLevelChecked(givenEvent, givenProvider, givenEncoder)()

	assert.ToBeMatching(t, ".+: expected", actualErr)
	assert.ToBeEqual(t, `"myKey":`, givenEncoder.String())
}

func Test_Json_encodeValuesChecked(t *testing.T) {
	givenProvider := recording.NewProvider()
	givenLogger := givenProvider.GetRootLogger()

	cases := []struct {
		name            string
		given           log.Event
		expected        string
		printRootLogger bool
		unsorted        bool
	}{{
		name: "withStringAndInteger",
		given: givenLogger.NewEvent(0, map[string]interface{}{
			"foo": "foo",
			"bar": 1,
		}),
		expected: `,"bar":1,"foo":"foo"`,
	}, {
		name: "withStringAndMap",
		given: givenLogger.NewEvent(0, map[string]interface{}{
			"foo": "foo",
			"bar": map[string]interface{}{
				"hello": "world",
			},
		}),
		expected: `,"bar":{"hello":"world"},"foo":"foo"`,
	}, {
		name: "withStringAndError",
		given: givenLogger.NewEvent(0, map[string]interface{}{
			"foo": "foo",
			"bar": errors.New("anErrorMessage"),
		}),
		expected: `,"bar":"anErrorMessage","foo":"foo"`,
	}, {
		name: "withStringAndStringPointer",
		given: givenLogger.NewEvent(0, map[string]interface{}{
			"foo": "foo",
			"bar": pstring("barAsPointer"),
		}),
		expected: `,"bar":"barAsPointer","foo":"foo"`,
	}, {
		name: "withStringAndLazy",
		given: givenLogger.NewEvent(0, map[string]interface{}{
			"foo": "foo",
			"bar": aLazy("barAsLazy"),
		}),
		expected: `,"bar":"barAsLazy","foo":"foo"`,
	}, {
		name: "withStringAndFilteredRespected",
		given: givenLogger.NewEvent(level.Info, map[string]interface{}{
			"foo": "foo",
			"bar": fields.RequireMaximalLevel(level.Info, "barAsFiltered"),
		}),
		expected: `,"bar":"barAsFiltered","foo":"foo"`,
	}, {
		name: "withStringAndFilteredIgnored",
		given: givenLogger.NewEvent(level.Info, map[string]interface{}{
			"foo": "foo",
			"bar": fields.RequireMaximalLevel(level.Debug, "barAsFiltered"),
		}),
		expected: `,"foo":"foo"`,
	}, {
		name: "withoutExcluded",
		given: givenLogger.NewEvent(level.Info, map[string]interface{}{
			"foo": "foo",
			"bar": fields.Exclude,
		}),
		expected: `,"foo":"foo"`,
	}, {
		name: "withStringAndSomeLogger",
		given: givenLogger.NewEvent(0, map[string]interface{}{
			"foo":    "foo",
			"logger": "aLogger",
		}),
		expected: `,"foo":"foo","logger":"aLogger"`,
	}, {
		name: "withStringAndHiddenRootLogger",
		given: givenLogger.NewEvent(0, map[string]interface{}{
			"foo":    "foo",
			"logger": "ROOT",
		}),
		expected: `,"foo":"foo"`,
	}, {
		name: "withStringAndShowRootLogger",
		given: givenLogger.NewEvent(0, map[string]interface{}{
			"foo":    "foo",
			"logger": "ROOT",
		}),
		printRootLogger: true,
		expected:        `,"foo":"foo","logger":"ROOT"`,
	}, {
		name: "withStringAndLazy",
		given: givenLogger.NewEvent(0, map[string]interface{}{
			"foo": "foo",
			"bar": aLazy("bar"),
		}),
		printRootLogger: true,
		expected:        `,"bar":"bar","foo":"foo"`,
	}, {
		name: "unsorted",
		given: givenLogger.NewEvent(0, map[string]interface{}{
			"foo": "foo",
		}),
		unsorted: true,
		expected: `,"foo":"foo"`,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			givenEncoder := encoding.NewBufferedJsonEncoder()
			instance := NewJson(func(json *Json) {
				if !c.unsorted {
					json.KeySorter = fields.DefaultKeySorter
				}
				json.PrintRootLogger = &c.printRootLogger
			})

			actualErr := instance.encodeValuesChecked(c.given, givenProvider, givenEncoder)()

			assert.ToBeNil(t, actualErr)

			assert.ToBeEqual(t, c.expected, givenEncoder.String())
		})
	}
}

type someProvider struct {
	log.Provider

	names nlevel.Names
}

func (instance *someProvider) GetLevelNames() nlevel.Names {
	return instance
}

func (instance *someProvider) ToName(l level.Level) (string, error) {
	return instance.toName(l), nil
}

func (instance *someProvider) toName(l level.Level) string {
	return fmt.Sprintf("myFormattedLevel-%d", l)
}

func (instance *someProvider) ToLevel(string) (level.Level, error) {
	panic("should never be called")
}

type failingJsonMarshalling struct {
	err error
}

func (instance failingJsonMarshalling) MarshalJSON() ([]byte, error) {
	return nil, instance.err
}
