package level

import (
	"errors"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
)

func Test_GetLevel(t *testing.T) {
	cases := []struct {
		name          string
		given         interface{}
		expectedLevel Level
		expectedOk    bool
	}{{
		name:          "nil",
		given:         nil,
		expectedLevel: 0,
		expectedOk:    false,
	}, {
		name:          "int",
		given:         666,
		expectedLevel: 0,
		expectedOk:    false,
	}, {
		name:          "error",
		given:         errors.New("foo"),
		expectedLevel: 0,
		expectedOk:    false,
	}, {
		name:          "chan",
		given:         make(chan int),
		expectedLevel: 0,
		expectedOk:    false,
	}, {
		name:          "func",
		given:         func() {},
		expectedLevel: 0,
		expectedOk:    false,
	}, {
		name:          "notAware",
		given:         testNotAware{},
		expectedLevel: 0,
		expectedOk:    false,
	}, {
		name:          "notAwarePointer",
		given:         &testNotAware{},
		expectedLevel: 0,
		expectedOk:    false,
	}, {
		name:          "notAwareWrapped",
		given:         genericWrapper{testNotAware{}},
		expectedLevel: 0,
		expectedOk:    false,
	}, {
		name:          "notAwarePointerWrapped",
		given:         genericWrapper{&testNotAware{}},
		expectedLevel: 0,
		expectedOk:    false,
	}, {
		name:          "nilWrapped",
		given:         genericWrapper{nil},
		expectedLevel: 0,
		expectedOk:    false,
	}, {
		name:          "notAwareWrappedPointer",
		given:         &genericWrapper{testNotAware{}},
		expectedLevel: 0,
		expectedOk:    false,
	}, {
		name:          "aware",
		given:         testAware{1001},
		expectedLevel: 1001,
		expectedOk:    true,
	}, {
		name:          "awareWrapped",
		given:         genericWrapper{testAware{1002}},
		expectedLevel: 1002,
		expectedOk:    true,
	}, {
		name:          "awareWrappedPointer",
		given:         &genericWrapper{testAware{1003}},
		expectedLevel: 1003,
		expectedOk:    true,
	}, {
		name:          "awarePointerWrapped",
		given:         genericWrapper{&testAware{1004}},
		expectedLevel: 1004,
		expectedOk:    true,
	}, {
		name:          "awarePointerWrappedPointerWrapped",
		given:         &genericWrapper{&genericWrapper{testAware{1005}}},
		expectedLevel: 1005,
		expectedOk:    true,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actualLevel, actualOk := Get(c.given)
			assert.ToBeEqual(t, c.expectedLevel, actualLevel)
			assert.ToBeEqual(t, c.expectedOk, actualOk)
		})
	}
}

func Test_SetLevel(t *testing.T) {
	cases := []struct {
		name       string
		given      interface{}
		givenLevel Level
		extractor  func(interface{}) Level
		expectedOk bool
	}{{
		name:       "nil",
		given:      nil,
		givenLevel: 0,
		expectedOk: false,
	}, {
		name:       "int",
		given:      666,
		givenLevel: 0,
		expectedOk: false,
	}, {
		name:       "error",
		given:      errors.New("foo"),
		givenLevel: 0,
		expectedOk: false,
	}, {
		name:       "chan",
		given:      make(chan int),
		givenLevel: 0,
		expectedOk: false,
	}, {
		name:       "func",
		given:      func() {},
		givenLevel: 0,
		expectedOk: false,
	}, {
		name:       "notAware",
		given:      testNotAware{},
		givenLevel: 0,
		expectedOk: false,
	}, {
		name:       "notAwarePointer",
		given:      &testNotAware{},
		givenLevel: 0,
		expectedOk: false,
	}, {
		name:       "notAwareWrapped",
		given:      genericWrapper{testNotAware{}},
		givenLevel: 0,
		expectedOk: false,
	}, {
		name:       "notAwarePointerWrapped",
		given:      genericWrapper{&testNotAware{}},
		givenLevel: 0,
		expectedOk: false,
	}, {
		name:       "nilWrapped",
		given:      genericWrapper{nil},
		givenLevel: 0,
		expectedOk: false,
	}, {
		name:       "notAwareWrappedPointer",
		given:      &genericWrapper{testNotAware{}},
		givenLevel: 0,
		expectedOk: false,
	}, {
		name:       "awarePointer",
		given:      &testAware{0},
		extractor:  func(v interface{}) Level { return v.(*testAware).level },
		givenLevel: 1001,
		expectedOk: true,
	}, {
		name:       "awarePointerWrapped",
		given:      genericWrapper{&testAware{}},
		extractor:  func(v interface{}) Level { return v.(genericWrapper).inner.(*testAware).level },
		givenLevel: 1002,
		expectedOk: true,
	}, {
		name:       "awarePointerWrappedPointer",
		given:      &genericWrapper{&testAware{}},
		extractor:  func(v interface{}) Level { return v.(*genericWrapper).inner.(*testAware).level },
		givenLevel: 1003,
		expectedOk: true,
	}, {
		name:       "awarePointerWrappedWrapped",
		given:      genericWrapper{genericWrapper{&testAware{1005}}},
		extractor:  func(v interface{}) Level { return v.(genericWrapper).inner.(genericWrapper).inner.(*testAware).level },
		givenLevel: 1005,
		expectedOk: true,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actualOk := Set(c.given, c.givenLevel)
			assert.ToBeEqual(t, c.expectedOk, actualOk)
			if c.expectedOk {
				assert.ToBeEqual(t, c.givenLevel, c.extractor(c.given))
			}
		})
	}
}

type testNotAware struct{}

type genericWrapper struct {
	inner interface{}
}

func (instance genericWrapper) Unwrap() interface{} {
	return instance.inner
}

type testAware struct {
	level Level
}

func (instance testAware) GetLevel() Level {
	return instance.level
}

func (instance *testAware) SetLevel(v Level) {
	instance.level = v
}
